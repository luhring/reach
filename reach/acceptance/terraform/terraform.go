package terraform

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"
	"strings"
	"testing"
)

var logTF = flag.Bool("log-tf", false, "log output of Terraform helper")

const execName = "terraform"
const terraformPlan = "terraform.plan"

// Terraform is a usable, programmatic representation of the Terraform application.
type Terraform struct {
	t        *testing.T
	execPath string
	tempDir  string
	logging  bool
}

// New creates a new instance of a Terraform struct that can perform the operations of a basic Terraform workflow.
// Callers should call CleanUp when done with the object to ensure there are no lingering side effects.
func New(t *testing.T) (*Terraform, error) {
	t.Helper()

	logging := *logTF

	execPath, err := findExecutable(logging)
	if err != nil {
		return nil, err
	}

	tempDir, err := createTempDir(logging)
	if err != nil {
		return nil, err
	}

	return &Terraform{
		t:        t,
		execPath: execPath,
		tempDir:  tempDir,
		logging:  logging,
	}, nil
}

// CleanUp removes the temporary directory set up for Terraform assets and all contents.
func (tf *Terraform) CleanUp() error {
	tf.t.Helper()

	if tf.logging {
		log.Println("cleaning up...")
	}

	err := os.RemoveAll(tf.tempDir)
	if err != nil {
		return fmt.Errorf("unable to clean up temp dir '%s': %v", tf.tempDir, err)
	}

	if tf.logging {
		log.Printf("temp dir was successfully removed ('%s')", tf.tempDir)
	}

	return nil
}

// LoadFilesFromDir calls LoadFile for all specified files within the specified directory.
func (tf *Terraform) LoadFilesFromDir(dir string, files ...string) error {
	for _, file := range files {
		fullPath := path.Join(dir, file)
		err := tf.LoadFile(fullPath)
		if err != nil {
			return err
		}
	}

	return nil
}

// LoadFile copies a file for Terraform to use into the temporary working directory.
func (tf *Terraform) LoadFile(file string) error {
	tf.t.Helper()

	if tf.logging {
		log.Printf("loading '%s'...", file)
	}

	data, err := ioutil.ReadFile(file)
	if err != nil {
		return fmt.Errorf("unable to read file '%s': %v", file, err)
	}

	_, filename := path.Split(file)
	destFile := path.Join(tf.tempDir, filename)
	err = ioutil.WriteFile(destFile, data, 0644)
	if err != nil {
		return fmt.Errorf("unable to write file '%s': %v", destFile, err)
	}

	return nil
}

// Init calls 'Terraform init' to initialize Terraform in the temporary working directory.
func (tf *Terraform) Init() error {
	tf.t.Helper()

	if tf.logging {
		log.Print("initializing Terraform...")
	}

	err := tf.action("unable to initialize Terraform", "init")
	if err != nil {
		return err
	}

	return nil
}

// Plan calls 'Terraform plan' to create a Terraform plan that can be run by calling Apply.
func (tf *Terraform) Plan() error {
	tf.t.Helper()

	if tf.logging {
		log.Print("creating a Terraform plan...")
	}

	err := tf.action("unable to create plan", "plan", "-out", terraformPlan)
	if err != nil {
		return err
	}

	return nil
}

// Apply calls 'Terraform apply', which applies the plan generated by calling Plan. This means you must first call Plan.
func (tf *Terraform) Apply() error {
	tf.t.Helper()

	if tf.logging {
		log.Print("applying Terraform plan...")
	}

	err := tf.action("unable to apply plan", "apply", terraformPlan)
	if err != nil {
		return err
	}

	return nil
}

// PlanAndApply does the equivalent of calling Plan and Apply.
func (tf *Terraform) PlanAndApply() error {
	tf.t.Helper()

	err := tf.Plan()
	if err != nil {
		return err
	}

	err = tf.Apply()
	if err != nil {
		return err
	}

	return nil
}

// Destroy calls 'Terraform destroy' non-interactively, which tears down all resources referenced within the temp dir.
func (tf *Terraform) Destroy() error {
	tf.t.Helper()

	if tf.logging {
		log.Print("destroying resources...")
	}

	err := tf.action("unable to destroy", "destroy", "-auto-approve")
	if err != nil {
		return err
	}

	return nil
}

// Output calls 'Terraform output' to retrieve a predefined output value from available resources.
func (tf *Terraform) Output(name string) (string, error) {
	tf.t.Helper()

	pop, err := tf.changeToTempDir()
	defer func() {
		err := pop()
		if err != nil {
			// Report error irrespective of logTF setting for visibility
			fmt.Printf("encountered error in deferred call: %v", err)
		}
	}()

	output, err := tf.execForOutput("output", "-no-color", name)
	if err != nil {
		return "", fmt.Errorf("unable to retrieve output '%s': %v", name, err)
	}

	value := strings.TrimSpace(output)

	if tf.logging {
		log.Printf("retrieved output value: %s=%s", name, value)
	}

	return value, nil
}

// Version retrieves the current Terraform version by calling 'Terraform version'.
func (tf *Terraform) Version() (string, error) {
	tf.t.Helper()

	b, err := tf.execForOutput("version")
	if err != nil {
		return "", fmt.Errorf("unable to get version: %v", err)
	}

	return b, nil
}

func (tf *Terraform) action(errMessage string, args ...string) error {
	tf.t.Helper()

	pop, err := tf.changeToTempDir()
	defer func() {
		err := pop()
		if err != nil {
			// Report error irrespective of logTF setting for visibility
			fmt.Printf("encountered error in deferred call: %v", err)
		}
	}()

	err = tf.exec(args...)
	if err != nil {
		return fmt.Errorf("%s: %v", errMessage, err)
	}

	return nil
}

func (tf *Terraform) changeToTempDir() (func() error, error) {
	tf.t.Helper()

	originalWorkDir, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("unable to determine original working directory: %v", err)
	}

	err = os.Chdir(tf.tempDir)
	if err != nil {
		tf.t.Fatalf("unable to change directory to '%s': %v", tf.tempDir, err)
	}

	changeToOriginalDir := func() error {
		err = os.Chdir(originalWorkDir)
		if err != nil {
			return fmt.Errorf("unable to change directory to '%s': %v", originalWorkDir, err)
		}

		return nil
	}

	return changeToOriginalDir, nil
}

func (tf *Terraform) exec(args ...string) error {
	tf.t.Helper()

	cmd := exec.Command(tf.execPath, args...)

	var stdout, stderr io.ReadCloser

	if tf.logging {
		var err error
		stdout, err = cmd.StdoutPipe()
		if err != nil {
			return err
		}
		stderr, err = cmd.StderrPipe()
		if err != nil {
			return err
		}
	}

	err := cmd.Start()
	if err != nil {
		return err
	}

	if tf.logging {
		multi := io.MultiReader(stdout, stderr)
		in := bufio.NewScanner(multi)

		for in.Scan() {
			fmt.Printf("%v\n", in.Text())
		}
		if err := in.Err(); err != nil {
			fmt.Printf("error: %s", err)
		}
	}

	err = cmd.Wait()
	if err != nil {
		return fmt.Errorf("command exited non-zero: %v", err)
	}

	return nil
}

func (tf *Terraform) execForOutput(args ...string) (string, error) {
	tf.t.Helper()

	cmd := exec.Command(tf.execPath, args...)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return "", fmt.Errorf("unable to connect to stdout: %v", err)
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return "", fmt.Errorf("unable to connect to stderr: %v", err)
	}

	err = cmd.Start()
	if err != nil {
		return "", fmt.Errorf("unable to start command: %v", err)
	}

	fromStdout := bufio.NewScanner(stdout)
	fromStderr := bufio.NewScanner(stderr)

	var output string
	for fromStdout.Scan() {
		output += fmt.Sprintf("%s\n", fromStdout.Text())
	}

	var errText string
	for fromStderr.Scan() {
		errText += fmt.Sprintf("%s\n", fromStderr.Text())
	}

	err = cmd.Wait()
	if err != nil {
		return "", fmt.Errorf("command exited non-zero: %v\n\n%s", err, errText)
	}

	return output, nil
}

func findExecutable(logging bool) (string, error) {
	execPath, err := exec.LookPath(execName)
	if err != nil {
		return "", fmt.Errorf("unable to find %s: %v", execName, err)
	}

	if logging {
		log.Printf("found %s at: %s", execName, execPath)
	}

	return execPath, nil
}

func createTempDir(logging bool) (string, error) {
	// create temporary working directory for all Terraform operations
	tempDir, err := ioutil.TempDir("", "Terraform")
	if err != nil {
		return "", fmt.Errorf("unable to create temp dir: %v", err)
	}

	if logging {
		log.Printf("created temporary working directory for Terraform files: %s", tempDir)
	}

	return tempDir, nil
}
