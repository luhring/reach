package reach

import (
	"fmt"
)

const indent = "  "

type Explanation struct {
	lines []ExplanationLine
}

func newExplanation(lines ...string) Explanation {
	var explanation Explanation

	for _, line := range lines {
		explanation.addLine(line)
	}

	return explanation
}

func (e *Explanation) addLine(text string) {
	e.addLineWithIndents(0, text)
}

func (e *Explanation) addLineWithIndents(indents int, text string) {
	l := newExplanationLine(indents, text)
	e.lines = append(e.lines, l)
}

func (e *Explanation) addLineFormat(format string, a ...interface{}) {
	e.addLine(fmt.Sprintf(format, a...))
}

func (e *Explanation) addLineFormatWithIndents(indents int, format string, a ...interface{}) {
	e.addLineWithIndents(indents, fmt.Sprintf(format, a...))
}

func (e *Explanation) addBlankLine() {
	e.addLine("")
}

func (e *Explanation) append(explanation Explanation) {
	for _, line := range explanation.lines {
		e.addLineWithIndents(line.indents, line.text)
	}
}

func (e *Explanation) subsume(explanation Explanation) {
	for _, line := range explanation.lines {
		e.addLineWithIndents(line.indents+1, line.text)
	}
}

func (e *Explanation) render() string {
	var output string

	for _, line := range e.lines {
		r := line.render()
		output += r + "\n"
	}

	return output
}

type ExplanationLine struct {
	indents int
	text    string
}

func newExplanationLine(indents int, text string) ExplanationLine {
	return ExplanationLine{
		indents,
		text,
	}
}

func (l *ExplanationLine) render() string {
	var lineIndentation string

	for i := 0; i < l.indents; i++ {
		lineIndentation += indent
	}

	return lineIndentation + l.text
}
