package analyzer

import (
	"log"
	"testing"

	"github.com/luhring/reach/reach"
	"github.com/luhring/reach/reach/acceptance"
	"github.com/luhring/reach/reach/acceptance/terraform"
	"github.com/luhring/reach/reach/aws"
	"github.com/luhring/reach/reach/set"
)

func TestAnalyze(t *testing.T) {
	acceptance.Check(t)

	cases := []struct {
		name                   string
		files                  []string
		expectedForwardTraffic reach.TrafficContent
		expectedReturnTraffic  reach.TrafficContent
	}{
		{
			"same subnet/all traffic",
			[]string{
				"main.tf",
				"outputs.tf",
				"ami_ubuntu.tf",
				"vpc.tf",
				"subnet_single.tf",
				"ec2_instances_same_subnet_all_traffic.tf",
				"security_group_outbound_allow_all.tf",
				"security_group_inbound_allow_all.tf",
			},
			reach.NewTrafficContentForAllTraffic(),
			reach.NewTrafficContentForAllTraffic(),
		},
		{
			"same subnet/SSH",
			[]string{
				"main.tf",
				"outputs.tf",
				"ami_ubuntu.tf",
				"vpc.tf",
				"subnet_single.tf",
				"ec2_instances_same_subnet_ssh.tf",
				"security_group_outbound_allow_all.tf",
				"security_group_inbound_allow_ssh.tf",
			},
			trafficSSH(),
			reach.NewTrafficContentForAllTraffic(),
		},
		{
			"same subnet/HTTPS via two-way IP match",
			[]string{
				"main.tf",
				"outputs.tf",
				"ami_ubuntu.tf",
				"vpc.tf",
				"subnet_single.tf",
				"ec2_instances_https_via_two-way_sg_ip_match.tf",
				"security_group_outbound_allow_https_to_ip.tf",
				"security_group_inbound_allow_https_from_ip.tf",
			},
			trafficHTTPS(),
			reach.NewTrafficContentForAllTraffic(),
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
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
				tc.files...,
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

			// Analyze

			analyzer := New()
			analysis, err := analyzer.Analyze(source, destination)
			if err != nil {
				t.Fatal(err)
			}

			// Tests

			if forwardTraffic := analysis.NetworkVectors[0].Traffic; forwardTraffic.String() != tc.expectedForwardTraffic.String() { // TODO: consider a better comparison method besides strings
				t.Errorf("expected: %v\nbut was: %v\n", tc.expectedForwardTraffic, forwardTraffic)
			} else {
				log.Print("forward traffic analysis was successful")
			}

			if returnTraffic := analysis.NetworkVectors[0].ReturnTraffic; returnTraffic.String() != tc.expectedReturnTraffic.String() {
				t.Errorf("expected: %v\nbut was: %v\n", tc.expectedReturnTraffic, returnTraffic)
			} else {
				log.Print("return traffic analysis was successful")
			}
		})
	}
}

func trafficSSH() reach.TrafficContent {
	ports, err := set.NewPortSetFromRange(22, 22)
	if err != nil {
		panic(err)
	}

	return reach.NewTrafficContentForPorts(reach.ProtocolTCP, ports)
}

func trafficHTTPS() reach.TrafficContent {
	ports, err := set.NewPortSetFromRange(443, 443)
	if err != nil {
		panic(err)
	}

	return reach.NewTrafficContentForPorts(reach.ProtocolTCP, ports)
}
