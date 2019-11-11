package cmd

import (
	"fmt"
	"os"

	"github.com/mgutz/ansi"

	"github.com/luhring/reach/reach"
)

const canReach = "source is able to reach destination"
const cannotReach = "source is unable to reach destination"

func doAssertReachable(analysis reach.Analysis) {
	if passesAssertReachable(analysis) {
		exitWithSuccessfulAssertion(canReach)
	} else {
		exitWithFailedAssertion(cannotReach)
	}
}

func doAssertNotReachable(analysis reach.Analysis) {
	if passesAssertNotReachable(analysis) {
		exitWithSuccessfulAssertion(cannotReach)
	} else {
		exitWithFailedAssertion(canReach)
	}
}

func passesAssertReachable(analysis reach.Analysis) bool {
	return !forwardTraffic.None()
}

func passesAssertNotReachable(analysis reach.Analysis) bool {
	return forwardTraffic.None()
}

func exitWithFailedAssertion(text string) {
	failedMessage := ansi.Color("assertion failed:", "red+b")
	secondaryMessage := ansi.Color(text, "red")
	_, _ = fmt.Fprintf(os.Stderr, "\n%v %v\n", failedMessage, secondaryMessage)

	os.Exit(2)
}

func exitWithSuccessfulAssertion(text string) {
	succeededMessage := ansi.Color("assertion succeeded:", "green+b")
	secondaryMessage := ansi.Color(text, "green")
	_, _ = fmt.Fprintf(os.Stderr, "\n%v %v\n", succeededMessage, secondaryMessage)

	os.Exit(0)
}
