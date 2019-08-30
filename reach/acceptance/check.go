package acceptance

import (
	"flag"
	"testing"
)

var acceptance = flag.Bool("acceptance", false, "perform full acceptance testing")

func Check(t *testing.T) {
	t.Helper()

	if !*acceptance {
		t.Skip("not running acceptance tests")
	}
}
