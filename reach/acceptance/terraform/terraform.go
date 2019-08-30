package terraform

import (
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"testing"
)

type terraform struct {
	t        *testing.T
	execPath string
	tempDir  string
	logging  bool
}

// New creates a new instance of a terraform struct that can perform the operations of a basic Terraform workflow.
// Callers should call CleanUp() when done with the object to ensure there are no lingering side effects after use.
func New(t *testing.T, logging bool) *terraform {
	t.Helper()

	const execName = "terraform"

	// make sure we have access to terraform executable
	execPath, err := exec.LookPath(execName)
	if err != nil {
		t.Fatalf("unable to find %s: %v", execName, err)
	}
	if logging {
		t.Logf("found %s at: %s", execName, execPath)
	}

	// create temporary working directory for all Terraform operations
	tempDir, err := ioutil.TempDir("", "terraform")
	if err != nil {
		t.Fatalf("unable to create temp dir: %v", err)
	}
	if logging {
		t.Logf("created temporary working directory for terraform files: %s", tempDir)
	}

	return &terraform{
		t:        t,
		execPath: execPath,
		tempDir:  tempDir,
		logging:  logging,
	}
}

func (tf *terraform) CleanUp() {
	tf.t.Helper()

	err := os.RemoveAll(tf.tempDir)
	if tf.logging {
		if err != nil {
			tf.t.Logf("unable to clean up temp dir '%s': %v", tf.tempDir, err)
		} else {
			tf.t.Logf("temp dir was successfully removed ('%s')", tf.tempDir)
		}
	}
}

// Load copies files for Terraform to use into the temporary working directory.
// It is recommended that callers pass in absolute file paths, since the terraform object occasionally needs to change the working directory.
func (tf *terraform) Load(files ...string) {
	tf.t.Helper()

	for _, filePath := range files {
		data, err := ioutil.ReadFile(filePath)
		if err != nil {
			tf.t.Fatalf("unable to read file '%s': %v", filePath, err)
		}

		_, filename := path.Split(filePath)
		destFile := path.Join(tf.tempDir, filename)
		err = ioutil.WriteFile(destFile, data, 0644)
		if err != nil {
			tf.t.Fatalf("unable to write file '%s': %v", destFile, err)
		}
		if tf.logging {
			tf.t.Logf("copied '%s' to '%s'", filePath, destFile)
		}
	}
}

// Init calls 'terraform init' to initialize Terraform in the temporary working directory.
func (tf *terraform) Init() {
	tf.t.Helper()

	revert := tf.changeDirToTempDir()
	defer revert()

	// TODO: This should stream back output bytes as it gets them, not wait until the end (due to CombinedOutput())

	_ = tf.exec("init") // TODO: Handle this better
}

func (tf *terraform) Plan() {
	tf.t.Helper()
}

func (tf *terraform) Apply() {
	tf.t.Helper()
}

func (tf *terraform) Destroy() {
	tf.t.Helper()
}

func (tf *terraform) changeDirToTempDir() func() {
	tf.t.Helper()

	originalWorkDir, err := os.Getwd()
	if err != nil {
		tf.t.Fatalf("unable to determine original working directory: %v", err)
	}

	err = os.Chdir(tf.tempDir)
	if err != nil {
		tf.t.Fatalf("unable to change directory to '%s': %v", tf.tempDir, err)
	}

	if tf.logging {
		tf.t.Logf("changed directory from '%s' to '%s'", originalWorkDir, tf.tempDir)
	}

	changeDirToOriginalDir := func() {
		err = os.Chdir(originalWorkDir)
		if err != nil {
			tf.t.Fatalf("unable to change directory to '%s': %v", originalWorkDir, err)
		}
		if tf.logging {
			tf.t.Logf("changed directory back to '%s'", originalWorkDir)
		}
	}

	return changeDirToOriginalDir
}

func (tf *terraform) Version() { // TODO: This should just return a string of the 'terraform version' output, and let the caller do what it wants with that
	tf.t.Helper()

	_ = tf.exec("version") // TODO: Handle this better
}

func (tf *terraform) exec(args ...string) error {
	tf.t.Helper()

	// TODO: return real error if found

	cmd := exec.Command(tf.execPath, args...)
	b, err := cmd.CombinedOutput()
	if err != nil {
		tf.t.Fatalf("encountered error when executing %s: %v", tf.execPath, err)
	}

	if tf.logging {
		_, _ = io.WriteString(os.Stdout, string(b))
	}

	return nil
}
