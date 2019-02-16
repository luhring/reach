package reach

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/luhring/reach/network"
	"net"
)

type SecurityGroupRule struct {
	Ports    *network.PortRange
	IPRanges []*net.IPNet
	SGRefs   []*SecurityGroupReference
}

type SecurityGroupReference struct {
	ID     string
	Name   string
	UserID string
	VPCID  string
}

func NewSecurityGroupRule(permission *ec2.IpPermission) (*SecurityGroupRule, error) {
	portRange, err := network.NewPortRange(
		aws.StringValue(permission.IpProtocol),
		aws.Int64Value(permission.FromPort),
		aws.Int64Value(permission.ToPort),
	)
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

	return &SecurityGroupRule{
		Ports:    portRange,
		IPRanges: ipRanges,
	}, nil
}

func (rule *SecurityGroupRule) doesApplyToInterface(targetInterface *NetworkInterface) bool {
	// Does rule mention a security group of the target interface?
	for _, sgRef := range rule.SGRefs {
		for _, targetSecurityGroup := range targetInterface.SecurityGroups {
			if sgRef.ID == targetSecurityGroup.ID {
				return true
			}
		}
	}

	// Does rule mention IP range that includes private IP address of target? (not yet supporting VPC peering)
	for _, ipRange := range rule.IPRanges {
		for _, privateIPAddress := range targetInterface.PrivateIPAddresses {
			if ipRange.Contains(privateIPAddress) {
				return true
			}
		}
	}

	return false
}
