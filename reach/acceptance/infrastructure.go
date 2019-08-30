package acceptance

import (
	"fmt"
	"testing"

	"github.com/luhring/reach/reach/acceptance/terraform"
)

func Deploy(t *testing.T) func() {
	t.Helper()

	tf := terraform.New(t, true)
	defer tf.CleanUp()

	// insert required .tf files into working directory
	files := []string{
		"acceptance/data/main.tf",
	}

	tf.Load(files...)
	tf.Init()

	// return callback to tear down infrastructure
	return func() {
		fmt.Println("destroyed!")
	}
}
