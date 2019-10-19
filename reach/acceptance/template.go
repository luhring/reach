package acceptance

import (
	"bytes"
	"fmt"
	"html/template"
	"path"
	"strings"
	"testing"
)

type subjectPairForTwoEC2Instances struct {
	SourceEC2InstanceID      string
	DestinationEC2InstanceID string
}

// processTemplate handles injecting data into specified template file.
func processTemplate(t *testing.T, name string, data interface{}) (string, error) {
	t.Helper()

	tmpl, err := template.New(name).ParseFiles(path.Join("..", "acceptance", "data", "golden", name)) // TODO: Need a better way to coordinate path construction with CWD of test
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

// ProcessTemplateForSubjectPairForTwoEC2Instances calls processTemplate using the subjectPairForTwoEC2Instances data type.
func ProcessTemplateForSubjectPairForTwoEC2Instances(t *testing.T, name, sourceID, destinationID string) (string, error) {
	t.Helper()

	return processTemplate(t, name, &subjectPairForTwoEC2Instances{
		SourceEC2InstanceID:      sourceID,
		DestinationEC2InstanceID: destinationID,
	})
}
