package reach

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
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

func (rule *SecurityGroupRule) doesApplyToInterface(targetInterface *NetworkInterface) (bool, MatchedTarget) {
	// Does rule mention a security group of the target interface?
	for _, sgRef := range rule.SGRefs {
		for _, targetSecurityGroup := range targetInterface.SecurityGroups {
			if sgRef.ID == targetSecurityGroup.ID {
				return true, &MatchedSGRef{
					sgRef,
				}
			}
		}
	}

	// Does rule mention IP range that includes private IP address of target? (not yet supporting VPC peering)
	for _, ipRange := range rule.IPRanges {
		for _, privateIPAddress := range targetInterface.PrivateIPAddresses {
			if ipRange.Contains(privateIPAddress) {
				return true, &MatchedIP{
					ipRange,
					privateIPAddress,
					false,
				}
			}
		}
	}

	return false, nil
}

type MatchedSGRef struct {
	SGRef *SecurityGroupReference
}

func (m *MatchedSGRef) Describe() string {
	return fmt.Sprintf("security group (%v)", m.SGRef.Name)
}

type MatchedIP struct {
	MatchedIPRange   *net.IPNet
	TargetIP         net.IP
	IsTargetIPPublic bool
}

func (m *MatchedIP) Describe() string {
	var p string
	if m.IsTargetIPPublic {
		p = "public"
	} else {
		p = "private"
	}

	return fmt.Sprintf(
		"%s IP address (%v, which is within the specified IP range, %v)",
		p,
		m.TargetIP.String(),
		m.MatchedIPRange.String(),
	)
}

type MatchedTarget interface {
	Describe() string
}
