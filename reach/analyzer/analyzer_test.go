package analyzer

import (
	"testing"

	"github.com/luhring/reach/reach"
	"github.com/luhring/reach/reach/acceptance"
	"github.com/luhring/reach/reach/acceptance/terraform"
	"github.com/luhring/reach/reach/aws"
)

func TestAnalyze(t *testing.T) {
	acceptance.Check(t)

	// Setup (and deferred teardown)
	tf, err := terraform.New(t)
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		err = tf.CleanUp()
		if err != nil {
			t.Fatal(err)
		}
	}()

	err = tf.LoadFilesFromDir(
		"../acceptance/data/tf",
		"main.tf",
		"ami_ubuntu.tf",
		"vpc.tf",
		"subnet_single.tf",
		"ec2_instance_source_and_destination_in_single_subnet.tf",
	)
	if err != nil {
		t.Fatal(err)
	}

	err = tf.Init()
	if err != nil {
		t.Fatal(err)
	}

	err = tf.PlanAndApply()
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		err = tf.Destroy()
		if err != nil {
			t.Fatal(err)
		}
	}()

	// Test

	sourceID, err := tf.Output("source_id")
	if err != nil {
		t.Fatal(err)
	}
	destinationID, err := tf.Output("destination_id")
	if err != nil {
		t.Fatal(err)
	}

	source, err := aws.NewEC2InstanceSubject(sourceID, reach.SubjectRoleSource)
	if err != nil {
		t.Fatal(err)
	}
	destination, err := aws.NewEC2InstanceSubject(destinationID, reach.SubjectRoleDestination)
	if err != nil {
		t.Fatal(err)
	}

	analyzer := New()

	analysis, err := analyzer.Analyze(source, destination)
	if err != nil {
		t.Fatal(err)
	}

	if vectorsCount := len(analysis.NetworkVectors); vectorsCount != 1 {
		t.Errorf("vectorsCount should be 1, but it was %d", vectorsCount)
		t.Errorf(analysis.ToJSON())
	}
}
