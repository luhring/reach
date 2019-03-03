package reach

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"net"
)

type SecurityGroupRule struct {
	TrafficAllowance *TrafficAllowance
	IPRanges         []*net.IPNet
	SGRefs           []*SecurityGroupReference
}

func newSecurityGroupRule(permission *ec2.IpPermission) (*SecurityGroupRule, error) {
	trafficAllowance, err := newTrafficAllowanceFromAWS(permission.IpProtocol, permission.FromPort, permission.ToPort)
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
			sgRefs[i] = newSecurityGroupReference(p)
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
