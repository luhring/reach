package cmd

import (
	"fmt"
	"os"

	"github.com/luhring/reach/reach/reacherr"
)

func exitWithError(err error) {
	_, _ = fmt.Fprintf(os.Stderr, "%v\n", err)

	os.Exit(1)
}

func fatal(message string) {
	_, _ = fmt.Fprintf(os.Stderr, "%s\n", message)
	os.Exit(1)
}

func handleError(err error) {
	// TODO: provide the user with info to include in the bug report

	msg := "An unexpected error occurred; please open a new issue to report this: " + githubURL
	if _, ok := err.(reacherr.ReachErr); ok {
		msg = err.Error()
	}
	fatal(msg)
}
