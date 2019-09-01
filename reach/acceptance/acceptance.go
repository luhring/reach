package acceptance

import (
	"flag"
	"fmt"
	"path"
	"testing"
)

var acceptance = flag.Bool("acceptance", false, "perform full acceptance testing")

func Check(t *testing.T) {
	t.Helper()

	if !*acceptance {
		t.Skip("not running acceptance tests")
	}
}

func GetPath(filename string) string {
	dataDir := path.Join("acceptance", "data", "tf")
	return path.Join(dataDir, filename)
}

func GetPaths(filenames ...string) []string {
	filePaths := make([]string, len(filenames))

	for i, f := range filenames {
		filePaths[i] = GetPath(f)
	}

	return filePaths
}

func FailNowIfError(t *testing.T, err error) {
	t.Helper()

	if err != nil {
		fmt.Println("\n\nFAILING NOW!\n\nWriting error to test log...\n\n")
		t.Fatal(err)
	}
}
