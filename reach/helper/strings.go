package helper

import (
	"strings"

	"github.com/mgutz/ansi"
)

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

func Bold(text string) string {
	return ansi.Color(text, "default+b")
}
