package cmd

import (
	"fmt"
	"os"

	"github.com/mgutz/ansi"

	"github.com/luhring/reach/reach"
)

func doAssertReachable(analysis reach.Analysis) {
	if analysis.PassesAssertReachable() {
		exitSuccessfulAssertion("source is able to reach destination")
	} else {
		exitFailedAssertion("one or more forward or return paths of network traffic is obstructed")
	}
}

func doAssertNotReachable(analysis reach.Analysis) {
	if analysis.PassesAssertNotReachable() {
		exitSuccessfulAssertion("source is unable to reach destination")
	} else {
		exitFailedAssertion("source is able to send network traffic to destination")
	}
}

func exitFailedAssertion(text string) {
	failedMessage := ansi.Color("assertion failed:", "red+b")
	secondaryMessage := ansi.Color(text, "red")
	_, _ = fmt.Fprintf(os.Stderr, "\n%v %v\n", failedMessage, secondaryMessage)

	os.Exit(2)
}

func exitSuccessfulAssertion(text string) {
	succeededMessage := ansi.Color("assertion succeeded:", "green+b")
	secondaryMessage := ansi.Color(text, "green")
	_, _ = fmt.Fprintf(os.Stderr, "\n%v %v\n", succeededMessage, secondaryMessage)

	os.Exit(0)
}
