package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/luhring/reach/reach/reacherr"
)

func exitWithError(err error) {
	_, _ = fmt.Fprintf(os.Stderr, "%v\n", err)

	os.Exit(1)
}

func fatal(messages ...string) {
	_, _ = fmt.Fprintf(os.Stderr, "%s\n", strings.Join(messages, "\n\n"))
	os.Exit(1)
}

func showErrorToUser(err error) {
	if reachErr, ok := err.(reacherr.ReachErr); ok { // properly-defined Reach error
		errs := reachErr.Unwrap()
		var msgs []string
		for _, e := range errs {
			msgs = append(msgs, e.Error())
		}
		fatal(msgs...)
		return
	}

	callToAction := `*** It looks like you've found a bug in Reach! ***

If you're feeling particularly generous with your time, you can help us out by:

1) Re-running your command with the "-v" flag. This exposes all log output.
2) Submitting an issue to help us track down the problem, and including as much information as possible.

	` + githubURL + `/issues

Thank you! You're awesome.`

	msg := fmt.Sprintf(`*** An unexpected error occurred... ***
%+v

%s
`, err, callToAction)
	fatal(msg)
}
