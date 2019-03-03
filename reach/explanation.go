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
		explanation.AddLine(line)
	}

	return explanation
}

func (e *Explanation) AddLine(text string) {
	e.AddLineWithIndents(0, text)
}

func (e *Explanation) AddLineWithIndents(indents int, text string) {
	l := NewExplanationLine(indents, text)
	e.lines = append(e.lines, l)
}

func (e *Explanation) AddLineFormat(format string, a ...interface{}) {
	e.AddLine(fmt.Sprintf(format, a...))
}

func (e *Explanation) AddLineFormatWithIndents(indents int, format string, a ...interface{}) {
	e.AddLineWithIndents(indents, fmt.Sprintf(format, a...))
}

func (e *Explanation) AddBlankLine() {
	e.AddLine("")
}

func (e *Explanation) Append(explanation Explanation) {
	for _, line := range explanation.lines {
		e.AddLineWithIndents(line.indents, line.text)
	}
}

func (e *Explanation) Subsume(explanation Explanation) {
	for _, line := range explanation.lines {
		e.AddLineWithIndents(line.indents+1, line.text)
	}
}

func (e *Explanation) Render() string {
	var output string

	for _, line := range e.lines {
		r := line.Render()
		output += r + "\n"
	}

	return output
}

type ExplanationLine struct {
	indents int
	text    string
}

func NewExplanationLine(indents int, text string) ExplanationLine {
	return ExplanationLine{
		indents,
		text,
	}
}

func (l *ExplanationLine) Render() string {
	var lineIndentation string

	for i := 0; i < l.indents; i++ {
		lineIndentation += indent
	}

	return lineIndentation + l.text
}
