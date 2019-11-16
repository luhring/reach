package cmd

import (
	"fmt"
	"os"
)

func exitWithError(err error) {
	_, _ = fmt.Fprintf(os.Stderr, "%v\n", err)

	os.Exit(1)
}
