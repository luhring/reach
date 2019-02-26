package reach

import (
	"fmt"
	"github.com/logrusorgru/aurora"
)

const indent = "  "

type Explanation struct {
	lines []ExplanationLine
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

func (e *Explanation) AddLineFormatWithEffect(effect func(arg interface{}) aurora.Value, format string, a ...interface{}) {
	text := effect(fmt.Sprintf(format, a...))
	e.AddLineFormat("%v", text)
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
