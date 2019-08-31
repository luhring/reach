package reach

import (
	"reflect"
	"testing"

	acc "github.com/luhring/reach/reach/acceptance"
	"github.com/luhring/reach/reach/acceptance/terraform"
)

func TestAnalyze(t *testing.T) {
	acc.Check(t)

	tf := terraform.New(t, true)
	defer tf.CleanUp()

	data := []string{
		"main.tf",
		"ami_ubuntu.tf",
		"ec2_instance_source_and_destination.tf",
	}

	files := acc.DataPaths(data...)

	tf.Load(files...)
	tf.ThoroughApply()
	defer tf.Destroy()

	t.Run("subjects only", func(t *testing.T) {
		sourceID := tf.Output("source_id")
		destinationID := tf.Output("destination_id")
		data := &acc.TwoSubjects{
			SourceID:      sourceID,
			DestinationID: destinationID,
		}

		templateName := "two_subjects.json"
		expectedJSON, err := acc.RenderTemplate(t, templateName, data)
		if err != nil {
			t.Errorf("couldn't complete render of template '%s': %v", templateName, err)
		}

		t.Logf("expected JSON:\n\n%s\n\n", expectedJSON) // TODO: Remove this line

		var expectedError error // (nil)

		const setupFailure = "unable to setup Analyze test"

		src, err := NewEC2InstanceSubject(sourceID, RoleSource)
		if err != nil {
			t.Fatalf("%s: %v", setupFailure, err)
		}

		dest, err := NewEC2InstanceSubject(destinationID, RoleDestination)
		if err != nil {
			t.Fatalf("%s: %v", setupFailure, err)
		}

		analysis, err := Analyze(src, dest)
		if !reflect.DeepEqual(expectedError, err) {
			diffErrorf(t, "err", expectedError, err)
		}

		analysisJSON := analysis.ToJSON()

		t.Logf("analysis JSON:\n\n%s\n\n", analysisJSON) // TODO: Remove this line

		if expectedJSON != analysis.ToJSON() {
			diffErrorf(t, "analysis", expectedJSON, analysisJSON)
		}
	})
}
