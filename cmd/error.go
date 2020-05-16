package cmd

import (
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/luhring/reach/reach/reacherr"
)

func exitWithError(err error) { // DEPRECATED
	_, _ = fmt.Fprintf(os.Stderr, "%v\n", err)
	os.Exit(1)
}

func fatal(messages ...string) {
	_, _ = fmt.Fprintf(os.Stderr, "%s\n", strings.Join(messages, "\n\n"))
	os.Exit(1)
}

func handleError(err error) {
	if reachErr, ok := err.(reacherr.ReachErr); ok {
		handleReachError(reachErr)
	} else {
		handleUnexpectedError(err)
	}
}

func handleReachError(reachErr reacherr.ReachErr) {
	errs := reachErr.Unwrap()
	sort.Slice(errs, func(i, j int) bool {
		return true // reverse sort
	})

	for _, e := range errs {
		logger.Error(e.Error())
	}

	originalError := errs[len(errs)-1]
	fatal(originalError.Error())
}

func handleUnexpectedError(err error) {
	logger.Error(err.Error())

	msg := fmt.Sprintf(`*** An unexpected error occurred... ***
%+v

%s
`, err, callToAction)
	fatal(msg)
}

var callToAction = `*** It looks like you've found a bug in Reach! ***

If you're feeling particularly generous with your time, you can help us out by:

1) Re-running your command with the "-v" flag. This exposes all log output.
2) Submitting an issue to help us track down the problem, and including as much information as possible.

	` + githubURL + `/issues

Thank you! You're awesome.`
