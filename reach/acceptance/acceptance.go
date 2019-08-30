package acceptance

import (
	"flag"
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

func DataPath(filename string) string {
	dataDir := path.Join("acceptance", "data")
	return path.Join(dataDir, filename)
}

func DataPaths(filenames ...string) []string {
	filePaths := make([]string, len(filenames))

	for i, f := range filenames {
		filePaths[i] = DataPath(f)
	}

	return filePaths
}
