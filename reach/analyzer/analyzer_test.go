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
		"security_group_outbound_allow_all.tf",
		"security_group_inbound_allow_all.tf",
		"ec2_instance_source_and_destination_in_single_subnet.tf",
	)
	if err != nil {
		t.Fatal(err)
	}

	err = tf.Init()
	if err != nil {
		t.Fatal(err)
	}

	defer func() {
		err = tf.Destroy() // Putting this before apply so that we're not left with some resources not destroyed after failure from apply step.
		if err != nil {
			t.Fatal(err)
		}
	}()
	err = tf.PlanAndApply()
	if err != nil {
		t.Fatal(err)
	}

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

	// Act

	analyzer := New()
	analysis, err := analyzer.Analyze(source, destination)
	if err != nil {
		t.Fatal(err)
	}

	// Tests

	t.Run("correct count of network vectors", func(t *testing.T) {
		if vectorsCount := len(analysis.NetworkVectors); vectorsCount != 1 {
			t.Fatalf("vectorsCount should be 1, but it was %d", vectorsCount)
		}
	})

	t.Run("all forward traffic allowed", func(t *testing.T) {
		if traffic := analysis.NetworkVectors[0].Traffic; !traffic.All() {
			t.Errorf("expected all traffic allowed, but traffic was: %v", traffic)
		}
	})

	t.Run("all return traffic allowed", func(t *testing.T) {
		if returnTraffic := analysis.NetworkVectors[0].ReturnTraffic; !returnTraffic.All() {
			t.Errorf("expected all returnTraffic allowed, but returnTraffic was: %v", returnTraffic)
		}
	})
}
