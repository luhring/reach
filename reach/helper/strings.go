package helper

import (
	"strings"

	"github.com/mgutz/ansi"
)

// PrefixLines adds a given prefix to every line in a given slice of lines and returns the result as a single string value.
func PrefixLines(lines []string, prefix string) string {
	var outputLines []string

	for _, line := range lines {
		if line == "" {
			outputLines = append(outputLines, line)
		} else {
			outputLines = append(outputLines, prefix+line)
		}
	}

	return strings.Join(outputLines, "\n")
}

// Indent adds one or more spaces at the beginning of each line of the given text and returns the result.
func Indent(text string, spaces int) string {
	if spaces < 1 {
		return text
	}

	lines := strings.Split(text, "\n")
	var outputLines []string

	var indentation string
	for i := 0; i < spaces; i++ {
		indentation += " "
	}

	for _, line := range lines {
		// don't indent empty lines
		if line == "" {
			outputLines = append(outputLines, line)
		} else {
			outputLine := indentation + line
			outputLines = append(outputLines, outputLine)
		}
	}

	return strings.Join(outputLines, "\n")
}

// Bold uses ANSI escape characters to return the bolded version of the input text.
func Bold(text string) string {
	return ansi.Color(text, "default+b")
}
