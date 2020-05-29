package apiclient

import (
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"

	reachAWS "github.com/luhring/reach/reach/aws"
)

func TestDirectionMatches(t *testing.T) {
	inboundRule := ec2.NetworkAclEntry{
		Egress: aws.Bool(false),
	}
	outboundRule := ec2.NetworkAclEntry{
		Egress: aws.Bool(true),
	}
	cases := []struct {
		name      string
		direction reachAWS.NetworkACLRuleDirection
		entry     ec2.NetworkAclEntry
		matches   bool
	}{
		{
			name:      "inbound direction should match inbound rule",
			direction: reachAWS.NetworkACLRuleDirectionInbound,
			entry:     inboundRule,
			matches:   true,
		},
		{
			name:      "outbound direction should not match inbound rule",
			direction: reachAWS.NetworkACLRuleDirectionOutbound,
			entry:     inboundRule,
			matches:   false,
		},
		{
			name:      "inbound direction should not match outbound rule",
			direction: reachAWS.NetworkACLRuleDirectionInbound,
			entry:     outboundRule,
			matches:   false,
		},
		{
			name:      "outbound direction should match outbound rule",
			direction: reachAWS.NetworkACLRuleDirectionOutbound,
			entry:     outboundRule,
			matches:   true,
		},
	}

	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			result := directionMatches(testCase.direction, testCase.entry)
			if result != testCase.matches {
				t.Fail()
			}
		})
	}
}
