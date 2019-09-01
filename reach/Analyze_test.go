package reach

import (
	"reflect"
	"testing"

	acc "github.com/luhring/reach/reach/acceptance"
	"github.com/luhring/reach/reach/acceptance/terraform"
)

func TestAnalyze(t *testing.T) {
	acc.Check(t)

	tf, err := terraform.New(t)
	acc.FailNowIfError(t, err)
	defer func() {
		err := tf.CleanUp()
		if err != nil {
			t.Fatalf("error during cleanup: %v", err)
		}
	}()

	files := []string{
		"main.tf",
		"ami_ubuntu.tf",
		"ec2_instance_source_and_destination.tf",
	}

	filePaths := acc.GetPaths(files...)

	err = tf.Load(filePaths...)
	acc.FailNowIfError(t, err)

	err = tf.Init()
	acc.FailNowIfError(t, err)

	err = tf.PlanAndApply()
	acc.FailNowIfError(t, err)

	defer func() {
		err := tf.Destroy()
		acc.FailNowIfError(t, err)
	}()

	t.Run("subjects only", func(t *testing.T) {
		sourceID, err := tf.Output("source_id")
		acc.FailNowIfError(t, err)

		destinationID, err := tf.Output("destination_id")
		acc.FailNowIfError(t, err)

		data := &acc.TwoSubjects{
			SourceID:      sourceID,
			DestinationID: destinationID,
		}

		templateName := "two_subjects.json"
		expectedAnalysisJSON, err := acc.RenderTemplate(t, templateName, data)
		if err != nil {
			t.Errorf("couldn't complete render of template '%s': %v", templateName, err)
		}

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

		if expectedAnalysisJSON != analysis.ToJSON() {
			diffErrorf(t, "analysis", expectedAnalysisJSON, analysisJSON)
		}
	})
}
