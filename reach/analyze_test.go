package reach

import (
	"path"
	"reflect"
	"testing"

	"github.com/luhring/reach/reach/acceptance"
	"github.com/luhring/reach/reach/acceptance/terraform"
	"github.com/luhring/reach/reach/aws"
)

func TestAnalyze(t *testing.T) {
	acceptance.Check(t)

	t.Run("deploy two EC2 instances in same subnet", func(t *testing.T) {
		tf, err := terraform.New(t)
		acceptance.IfErrorFailNow(t, err)
		defer func() {
			acceptance.IfErrorFailNow(t, tf.CleanUp())
		}()

		tfFilesDir := path.Join("acceptance", "data", "tf")
		tfFiles := []string{
			"main.tf",
			"ami_ubuntu.tf",
			"ec2_instance_source_and_destination.tf",
		}

		acceptance.IfErrorFailNow(t, tf.LoadFilesFromDir(tfFilesDir, tfFiles))
		acceptance.IfErrorFailNow(t, tf.Init())
		acceptance.IfErrorFailNow(t, tf.Plan())
		acceptance.IfErrorFailNow(t, tf.Apply())
		defer func() {
			acceptance.IfErrorFailNow(t, tf.Destroy())
		}()

		sourceEC2InstanceID, err := tf.Output("source_id")
		acceptance.IfErrorFailNow(t, err)
		destinationEC2InstanceID, err := tf.Output("destination_id")
		acceptance.IfErrorFailNow(t, err)

		cases := []struct {
			goldenFile    string
			expectedError error
		}{
			{
				goldenFile:    "analysis_subjects_two_ec2_instances.json",
				expectedError: nil,
			},
		}

		for _, tc := range cases {
			t.Run(tc.goldenFile, func(t *testing.T) {
				expectedAnalysisJSON, err := acceptance.ProcessTemplateForSubjectPairForTwoEC2Instances(t, tc.goldenFile, sourceEC2InstanceID, destinationEC2InstanceID)
				acceptance.IfErrorFailNow(t, err)

				// Arrange

				src, err := NewEC2InstanceSubject(sourceEC2InstanceID, SubjectRoleSource)
				acceptance.IfErrorFailNow(t, err)

				dest, err := NewEC2InstanceSubject(destinationEC2InstanceID, SubjectRoleDestination)
				acceptance.IfErrorFailNow(t, err)

				// Act

				subjects := []Subject{
					*src,
					*dest,
				}
				analysis, err := Analyze(subjects, aws.NewResourceStore())

				// Assert

				if !reflect.DeepEqual(tc.expectedError, err) {
					DiffErrorf(t, "err", tc.expectedError, err)
				}

				analysisJSON := analysis.ToJSON()

				if expectedAnalysisJSON != analysis.ToJSON() {
					DiffErrorf(t, "analysis", expectedAnalysisJSON, analysisJSON)
				}
			})
		}
	})
}