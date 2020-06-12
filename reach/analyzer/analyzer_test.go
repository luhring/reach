package analyzer

import (
	"log"
	"reflect"
	"testing"

	"github.com/luhring/reach/reach"
	"github.com/luhring/reach/reach/acceptance"
	"github.com/luhring/reach/reach/acceptance/terraform"
	"github.com/luhring/reach/reach/aws"
	"github.com/luhring/reach/reach/aws/apiclient"
	"github.com/luhring/reach/reach/cache"
	"github.com/luhring/reach/reach/generic"
	"github.com/luhring/reach/reach/generic/standard"
	"github.com/luhring/reach/reach/reachlog"
	"github.com/luhring/reach/reach/set"
)

func TestAnalyze(t *testing.T) {
	acceptance.Check(t)

	type testCase struct {
		name                   string
		files                  []string
		expectedForwardTraffic reach.TrafficContent
		expectedReturnTraffic  reach.TrafficContent
	}

	groupings := []struct {
		name  string
		files []string
		cases []testCase
	}{
		{
			"same subnet",
			[]string{
				"main.tf",
				"outputs.tf",
				"ami_ubuntu.tf",
				"vpc.tf",
				"subnet_single.tf",
			},
			[]testCase{
				{
					"no security group rules",
					[]string{
						"ec2_instances_same_subnet_no_security_group_rules.tf",
						"security_group_no_rules.tf",
					},
					reach.NewTrafficContentForNoTraffic(),
					reach.NewTrafficContentForAllTraffic(),
				},
				{
					"multiple protocols",
					[]string{
						"ec2_instances_same_subnet_multiple_protocols.tf",
						"security_group_no_rules.tf",
						"security_group_outbound_allow_all_udp_to_sg_no_rules.tf",
						"security_group_outbound_allow_esp.tf",
						"security_group_outbound_allow_all_tcp.tf",
						"security_group_inbound_allow_udp_dns_from_sg_no_rules.tf",
						"security_group_inbound_allow_esp.tf",
						"security_group_inbound_allow_ssh.tf",
					},
					trafficAssorted(),
					reach.NewTrafficContentForAllTraffic(),
				},
				{
					"UDP DNS via SG reference",
					[]string{
						"ec2_instances_same_subnet_udp_dns_via_sg_reference.tf",
						"security_group_no_rules.tf",
						"security_group_outbound_allow_all_udp_to_sg_no_rules.tf",
						"security_group_inbound_allow_udp_dns_from_sg_no_rules.tf",
					},
					trafficDNS(),
					reach.NewTrafficContentForAllTraffic(),
				},
				{
					"HTTPS via two-way IP match",
					[]string{
						"ec2_instances_same_subnet_https_via_two-way_sg_ip_match.tf",
						"security_group_outbound_allow_https_to_ip.tf",
						"security_group_inbound_allow_https_from_ip.tf",
					},
					trafficHTTPS(),
					reach.NewTrafficContentForAllTraffic(),
				},
				{
					"SSH",
					[]string{
						"ec2_instances_same_subnet_ssh.tf",
						"security_group_outbound_allow_all.tf",
						"security_group_inbound_allow_ssh.tf",
					},
					trafficSSH(),
					reach.NewTrafficContentForAllTraffic(),
				},
				{
					"all traffic",
					[]string{
						"ec2_instances_same_subnet_all_traffic.tf",
						"security_group_outbound_allow_all.tf",
						"security_group_inbound_allow_all.tf",
					},
					reach.NewTrafficContentForAllTraffic(),
					reach.NewTrafficContentForAllTraffic(),
				},
			},
		},
		{
			"same VPC",
			[]string{
				"main.tf",
				"outputs.tf",
				"ami_ubuntu.tf",
				"vpc.tf",
				"subnet_pair.tf",
			},
			[]testCase{
				{
					"all traffic",
					[]string{
						"network_acl_both_subnets_all_traffic.tf",
						"ec2_instances_same_vpc_all_traffic.tf",
						"security_group_outbound_allow_all.tf",
						"security_group_inbound_allow_all.tf",
					},
					reach.NewTrafficContentForAllTraffic(),
					reach.NewTrafficContentForAllTraffic(),
				},
				{
					"no NACL allow rules",
					[]string{
						"network_acl_both_subnets_no_traffic.tf",
						"ec2_instances_same_vpc_all_traffic.tf",
						"security_group_outbound_allow_all.tf",
						"security_group_inbound_allow_all.tf",
					},
					reach.NewTrafficContentForNoTraffic(),
					reach.NewTrafficContentForNoTraffic(),
				},
				{
					"NACL rules don't match SG rules",
					[]string{
						"network_acl_both_subnets_all_tcp.tf",
						"ec2_instances_same_vpc_all_esp.tf",
						"security_group_outbound_allow_esp.tf",
						"security_group_inbound_allow_all.tf",
					},
					reach.NewTrafficContentForNoTraffic(),
					trafficTCP(), // TODO: Revisit return traffic calculation for this scenario
				},
				{
					"Postgres with tightened rules",
					[]string{
						"network_acl_source_subnet_tightened_postgres.tf",
						"network_acl_destination_subnet_tightened_postgres.tf",
						"ec2_instances_same_vpc_postgres.tf",
						"security_group_no_rules.tf",
						"security_group_outbound_allow_postgres_to_sg_no_rules.tf",
						"security_group_inbound_allow_postgres_from_sg_no_rules.tf",
					},
					trafficPostgres(),
					trafficTCP(),
				},
			},
		},
	}

	for _, g := range groupings {
		t.Run(g.name, func(t *testing.T) {
			for _, tc := range g.cases {
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
						append(g.files, tc.files...)...,
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

					source, err := aws.NewEC2InstanceSubjectWithRole(sourceID, reach.SubjectRoleSource)
					if err != nil {
						t.Fatal(err)
					}
					destination, err := aws.NewEC2InstanceSubjectWithRole(destinationID, reach.SubjectRoleDestination)
					if err != nil {
						t.Fatal(err)
					}

					// Analyze

					catalog := reach.NewDomainClientCatalog()

					c := cache.New()
					awsClient, err := apiclient.NewDomainClient(&c)
					if err != nil {
						t.Fatal(err)
					}

					logger := reachlog.New(reachlog.LevelDebug)

					catalog.Store(aws.ResourceDomainAWS, awsClient)
					catalog.Store(generic.ResourceDomainGeneric, standard.NewDomainClient())
					analyzer := New(catalog, logger)

					log.Print("analyzing...")
					analysis, err := analyzer.Analyze(*source, *destination)
					if err != nil {
						t.Fatal(err)
					}

					// Tests

					log.Print("verifying analysis results...")

					ft := analysis.Paths[0].TrafficForward()

					if ft.String() != tc.expectedForwardTraffic.String() { // TODO: consider a better comparison method besides strings
						t.Errorf("forward traffic -- expected:\n%v\nbut was:\n%v\n", tc.expectedForwardTraffic, ft)
					} else {
						log.Print("✓ forward traffic content is correct")
					}
				})
			}
		})
	}
}

func TestConnectionPredictions(t *testing.T) {
	cases := []struct {
		name        string
		path        reach.Path
		predictions reach.ConnectionPredictionSet
	}{
		{
			name: "single point, all TCP return",
			path: pathWithPoints(
				pointWithReturnTraffic(trafficTCP(), false),
			),
			predictions: connectionPredictionsByProtocol(
				reach.ConnectionPredictionSuccess,
				reach.ConnectionPredictionPossibleFailure,
				reach.ConnectionPredictionFailure,
				reach.ConnectionPredictionFailure,
			),
		},
		{
			name: "single point, some TCP return",
			path: pathWithPoints(
				pointWithReturnTraffic(trafficTCPHighPortsOnly(), false),
			),
			predictions: connectionPredictionsByProtocol(
				reach.ConnectionPredictionPossibleFailure,
				reach.ConnectionPredictionPossibleFailure,
				reach.ConnectionPredictionFailure,
				reach.ConnectionPredictionFailure,
			),
		},
		{
			name: "single point, no return traffic",
			path: pathWithPoints(
				pointWithReturnTraffic(reach.NewTrafficContentForNoTraffic(), false),
			),
			predictions: connectionPredictionsByProtocol(
				reach.ConnectionPredictionFailure,
				reach.ConnectionPredictionPossibleFailure,
				reach.ConnectionPredictionFailure,
				reach.ConnectionPredictionFailure,
			),
		},
		{
			name: "multiple points, no port translation, all traffic",
			path: pathWithPoints(
				pointWithReturnTraffic(reach.NewTrafficContentForAllTraffic(), false),
				pointWithReturnTraffic(reach.NewTrafficContentForAllTraffic(), false),
				pointWithReturnTraffic(reach.NewTrafficContentForAllTraffic(), false),
			),
			predictions: connectionPredictionsByProtocol(
				reach.ConnectionPredictionSuccess,
				reach.ConnectionPredictionSuccess,
				reach.ConnectionPredictionSuccess,
				reach.ConnectionPredictionSuccess,
			),
		},
		{
			name: "multiple points, port translation, all traffic",
			path: pathWithPoints(
				pointWithReturnTraffic(reach.NewTrafficContentForAllTraffic(), false),
				pointWithReturnTraffic(reach.NewTrafficContentForAllTraffic(), true),
				pointWithReturnTraffic(reach.NewTrafficContentForAllTraffic(), false),
			),
			predictions: connectionPredictionsByProtocol(
				reach.ConnectionPredictionSuccess,
				reach.ConnectionPredictionSuccess,
				reach.ConnectionPredictionSuccess,
				reach.ConnectionPredictionSuccess,
			),
		},
		{
			name: "multiple points, no port translation, mix of TCP some and all",
			path: pathWithPoints(
				pointWithReturnTraffic(trafficTCP(), false),
				pointWithReturnTraffic(reach.NewTrafficContentForAllTraffic(), false),
				pointWithReturnTraffic(trafficSSH(), false),
			),
			predictions: connectionPredictionsByProtocol(
				reach.ConnectionPredictionPossibleFailure,
				reach.ConnectionPredictionPossibleFailure,
				reach.ConnectionPredictionFailure,
				reach.ConnectionPredictionFailure,
			),
		},
		{
			name: "multiple points, port translation, mix of TCP some and all",
			path: pathWithPoints(
				pointWithReturnTraffic(reach.NewTrafficContentForAllTraffic(), false),
				pointWithReturnTraffic(reach.NewTrafficContentForAllTraffic(), true),
				pointWithReturnTraffic(trafficTCPHighPortsOnly(), false),
			),
			predictions: connectionPredictionsByProtocol(
				reach.ConnectionPredictionPossibleFailure,
				reach.ConnectionPredictionPossibleFailure,
				reach.ConnectionPredictionFailure,
				reach.ConnectionPredictionFailure,
			),
		},
		{
			name: "multiple points, port translation, mix of TCP none and some",
			path: pathWithPoints(
				pointWithReturnTraffic(reach.NewTrafficContentForNoTraffic(), false),
				pointWithReturnTraffic(reach.NewTrafficContentForAllTraffic(), true),
				pointWithReturnTraffic(trafficTCPHighPortsOnly(), false),
			),
			predictions: connectionPredictionsByProtocol(
				reach.ConnectionPredictionFailure,
				reach.ConnectionPredictionPossibleFailure,
				reach.ConnectionPredictionFailure,
				reach.ConnectionPredictionFailure,
			),
		},
		{
			name: "multiple points, no port translation, mix of mutually exclusive TCP ports",
			path: pathWithPoints(
				point(trafficTCP(), trafficSSH(), false),
				point(trafficTCP(), trafficTCPHighPortsOnly(), false),
			),
			predictions: reach.ConnectionPredictionSet{
				reach.ProtocolTCP: reach.ConnectionPredictionFailure,
			},
		},
		{
			name: "multiple points, port translation, mix of mutually exclusive TCP ports",
			path: pathWithPoints(
				point(trafficTCP(), trafficSSH(), false),
				pointWithReturnTraffic(reach.NewTrafficContentForAllTraffic(), true),
				point(trafficTCP(), trafficTCPHighPortsOnly(), false),
			),
			predictions: reach.ConnectionPredictionSet{
				reach.ProtocolTCP: reach.ConnectionPredictionPossibleFailure,
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			// if tc.name != "multiple points, no port translation, mix of TCP some and all" {
			// 	t.Skip()
			// }
			result := ConnectionPredictions(tc.path)

			if !reflect.DeepEqual(result, tc.predictions) {
				t.Errorf("result did not match expectation\nresult: %v\nexpected: %v\n", result, tc.predictions)
			}
		})
	}
}

func connectionPredictionsByProtocol(tcp, udp, icmpv4, icmpv6 reach.ConnectionPrediction) reach.ConnectionPredictionSet {
	return map[reach.Protocol]reach.ConnectionPrediction{
		reach.ProtocolTCP:    tcp,
		reach.ProtocolUDP:    udp,
		reach.ProtocolICMPv4: icmpv4,
		reach.ProtocolICMPv6: icmpv6,
	}
}

func pointWithReturnTraffic(returnTraffic reach.TrafficContent, translatesPorts bool) reach.Point {
	return point(reach.NewTrafficContentForAllTraffic(), returnTraffic, translatesPorts)
}

func point(forwardTraffic, returnTraffic reach.TrafficContent, translatesPorts bool) reach.Point {
	return reach.Point{
		FactorsForward: []reach.Factor{
			{
				Traffic: forwardTraffic,
			},
		},
		FactorsReturn: []reach.Factor{
			{
				Traffic: returnTraffic,
			},
		},
		SegmentDivider: translatesPorts,
	}
}

func trafficSSH() reach.TrafficContent {
	ports := set.NewPortSetFromRange(22, 22)

	return reach.NewTrafficContentForPorts(reach.ProtocolTCP, ports)
}

func trafficHTTPS() reach.TrafficContent {
	ports := set.NewPortSetFromRange(443, 443)

	return reach.NewTrafficContentForPorts(reach.ProtocolTCP, ports)
}

func trafficDNS() reach.TrafficContent {
	ports := set.NewPortSetFromRange(53, 53)

	return reach.NewTrafficContentForPorts(reach.ProtocolUDP, ports)
}

func trafficESP() reach.TrafficContent {
	return reach.NewTrafficContentForCustomProtocol(50, true)
}

func trafficAssorted() reach.TrafficContent {
	tc, err := reach.NewTrafficContentFromMergingMultiple([]reach.TrafficContent{
		trafficDNS(),
		trafficSSH(),
		trafficESP(),
	})
	if err != nil {
		panic(err)
	}

	return tc
}

func trafficTCP() reach.TrafficContent {
	return reach.NewTrafficContentForPorts(reach.ProtocolTCP, set.NewFullPortSet())
}

func trafficTCPHighPortsOnly() reach.TrafficContent {
	return reach.NewTrafficContentForPorts(reach.ProtocolTCP, set.NewPortSetFromRange(1024, 65535))
}

func trafficPostgres() reach.TrafficContent {
	ports := set.NewPortSetFromRange(5432, 5432)

	return reach.NewTrafficContentForPorts(reach.ProtocolTCP, ports)
}

func pathWithPoints(points ...reach.Point) reach.Path {
	return reach.Path{
		Points: points,
		Edges:  make([]reach.Edge, len(points)-1),
	}
}
