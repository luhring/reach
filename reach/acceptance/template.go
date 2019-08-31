package acceptance

import (
	"bytes"
	"fmt"
	"html/template"
	"path"
	"strings"
	"testing"
)

type TwoSubjects struct {
	SourceID      string
	DestinationID string
}

func RenderTemplate(t *testing.T, name string, data interface{}) (string, error) {
	t.Helper()

	tmpl, err := template.New(name).ParseFiles(path.Join("acceptance", "data", "golden", name))
	if err != nil {
		t.Fail()
		return "", fmt.Errorf("error: unable to parse template file '%s': %v", name, err)
	}

	var b bytes.Buffer
	err = tmpl.Execute(&b, data)
	if err != nil {
		return "", fmt.Errorf("unable to execute template: %v", err)
	}

	return strings.TrimSpace(b.String()), nil
}
