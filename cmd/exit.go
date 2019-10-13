package cmd

import (
	"fmt"
	"github.com/mgutz/ansi"
	"os"
)

func exitWithError(err error) {
	_, _ = fmt.Fprintf(os.Stderr, "%v\n", err)

	os.Exit(1)
}

func exitWithFailedAssertion(text string) {
	failedMessage := ansi.Color("assertion failed:", "red+b")
	secondaryMessage := ansi.Color(text, "red")
	_, _ = fmt.Fprintf(os.Stderr, "%v %v\n", failedMessage, secondaryMessage)

	os.Exit(2)
}

func exitWithSuccessfulAssertion(text string) {
	succeededMessage := ansi.Color("assertion succeeded:", "green+b")
	secondaryMessage := ansi.Color(text, "green")
	_, _ = fmt.Fprintf(os.Stderr, "%v %v\n", succeededMessage, secondaryMessage)

	os.Exit(0)
}
