package reach

import (
	"fmt"
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
	tf.Init()
	tf.Plan()
	tf.Apply()
	defer tf.Destroy()

	cases := []struct {
		srcID            string
		destID           string
		expectedAnalysis *NewAnalysis
		expectedError    error
	}{
		{
			"i-0a93117c7575b6d54",
			"i-0136d3233f0ef1924",
			nil,
			nil,
		},
	}

	for _, tc := range cases {
		t.Run(fmt.Sprintf("src %s dest %s", tc.srcID, tc.destID), func(t *testing.T) {
			const setupFailure = "unable to setup Analyze test"

			src, err := NewEC2InstanceSubject(tc.srcID, RoleSource)
			if err != nil {
				t.Fatalf("%s: %v", setupFailure, err)
			}

			dest, err := NewEC2InstanceSubject(tc.destID, RoleDestination)
			if err != nil {
				t.Fatalf("%s: %v", setupFailure, err)
			}

			analysis, err := Analyze(src, dest)

			if !reflect.DeepEqual(tc.expectedAnalysis, analysis) {
				diffErrorf(t, "analysis", tc.expectedAnalysis, analysis)
			}

			if !reflect.DeepEqual(tc.expectedError, err) {
				diffErrorf(t, "err", tc.expectedError, err)
			}
		})
	}
}
