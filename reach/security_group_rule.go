package reach

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/logrusorgru/aurora"
	"github.com/luhring/reach/network"
	"net"
)

type SecurityGroupRule struct {
	TrafficAllowance *network.TrafficAllowance
	IPRanges         []*net.IPNet
	SGRefs           []*SecurityGroupReference
}

func NewSecurityGroupRule(permission *ec2.IpPermission) (*SecurityGroupRule, error) {
	trafficAllowance, err := network.NewTrafficAllowanceFromAWS(permission.IpProtocol, permission.FromPort, permission.ToPort)
	if err != nil {
		return nil, fmt.Errorf("unable to ingest security group rule: %v", err)
	}

	ipRanges := make([]*net.IPNet, len(permission.IpRanges))

	for i, r := range permission.IpRanges {
		cidr := aws.StringValue(r.CidrIp)
		_, ipNetwork, err := net.ParseCIDR(cidr)
		if err != nil {
			return nil, fmt.Errorf("unable to ingest security group rule: %v", err)
		}

		ipRanges[i] = ipNetwork
	}

	var sgRefs []*SecurityGroupReference

	if pairs := permission.UserIdGroupPairs; pairs != nil {
		sgRefs := make([]*SecurityGroupReference, len(pairs))

		for i, p := range pairs {
			sgRefs[i] = NewSecurityGroupReference(p)
		}
	}

	return &SecurityGroupRule{
		TrafficAllowance: trafficAllowance,
		IPRanges:         ipRanges,
		SGRefs:           sgRefs,
	}, nil
}

func (rule *SecurityGroupRule) matchWithInterface(targetInterface *NetworkInterface) RuleMatch {
	// Does rule mention a security group of the target interface?
	for _, sgRef := range rule.SGRefs {
		for _, targetSecurityGroup := range targetInterface.SecurityGroups {
			if sgRef.ID == targetSecurityGroup.ID {
				return &SGRefRuleMatch{
					rule,
					sgRef,
				}
			}
		}
	}

	// Does rule mention IP range that includes private IP address of target? (not yet supporting VPC peering)
	for _, ipRange := range rule.IPRanges {
		for _, privateIPAddress := range targetInterface.PrivateIPAddresses {
			if ipRange.Contains(privateIPAddress) {
				return &IPRuleMatch{
					rule,
					ipRange,
					privateIPAddress,
					false,
				}
			}
		}
	}

	return nil
}

type SGRefRuleMatch struct {
	Rule  *SecurityGroupRule
	SGRef *SecurityGroupReference
}

func (m *SGRefRuleMatch) Explain(observedDescriptor string) Explanation {
	var explanation Explanation

	explanation.AddLineFormat("security group (%v)", m.SGRef.Name)

	return explanation
}

type IPRuleMatch struct {
	Rule             *SecurityGroupRule
	MatchedIPRange   *net.IPNet
	TargetIP         net.IP
	IsTargetIPPublic bool
}

func (m *IPRuleMatch) Explain(observedDescriptor string) Explanation {
	var explanation Explanation

	var publicOrPrivate string
	if m.IsTargetIPPublic {
		publicOrPrivate = "public"
	} else {
		publicOrPrivate = "private"
	}

	explanation.AddLineFormatWithEffect(aurora.Green, "- rule: allow %v", aurora.Bold(m.Rule.TrafficAllowance.Describe()))
	explanation.AddLineFormatWithIndents(
		1,
		"(This rule handles an IP address range '%v' that includes the %s network interface's %s IP address '%v'.)",
		m.MatchedIPRange.String(),
		observedDescriptor,
		publicOrPrivate,
		m.TargetIP.String(),
	)

	return explanation
}

type RuleMatch interface {
	Explain(observedDescriptor string) Explanation
}