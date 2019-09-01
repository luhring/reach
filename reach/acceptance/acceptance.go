package acceptance

import (
	"flag"
	"fmt"
	"testing"
)

var acceptance = flag.Bool("acceptance", false, "perform full acceptance testing")

// Check to see if the acceptance flag was set, and if not, skip the current test.
func Check(t *testing.T) {
	t.Helper()

	if !*acceptance {
		t.Skip("not running acceptance tests")
	}
}

// IfErrorFailNow determines if the err value contains an error, and if so, calls t.Fatal(err).
func IfErrorFailNow(t *testing.T, err error) {
	t.Helper()

	if err != nil {
		fmt.Print("\n\nFAILING NOW!\n\nWriting error to test log...\n\n")
		t.Fatal(err)
	}
}
