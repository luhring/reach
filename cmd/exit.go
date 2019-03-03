package cmd

import (
	"fmt"
	"github.com/mgutz/ansi"
	"os"
)

func exitWithError(err error) {
	fmt.Println(err)

	os.Exit(1)
}

func exitWithFailedAssertion(text string) {
	failedMessage := ansi.Color("assertion failed:", "red+b")
	secondaryMessage := ansi.Color(text, "red")
	fmt.Printf("%v %v\n", failedMessage, secondaryMessage)

	os.Exit(2)
}

func exitWithSuccessfulAssertion(text string) {
	succeededMessage := ansi.Color("assertion succeeded:", "green+b")
	secondaryMessage := ansi.Color(text, "green")
	fmt.Printf("%v %v\n", succeededMessage, secondaryMessage)

	os.Exit(0)
}
