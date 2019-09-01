package acceptance

import (
	"flag"
	"fmt"
	"testing"
)

var acceptance = flag.Bool("acceptance", false, "perform full acceptance testing")

func Check(t *testing.T) {
	t.Helper()

	if !*acceptance {
		t.Skip("not running acceptance tests")
	}
}

func IfErrorFailNow(t *testing.T, err error) {
	t.Helper()

	if err != nil {
		fmt.Print("\n\nFAILING NOW!\n\nWriting error to test log...\n\n")
		t.Fatal(err)
	}
}
