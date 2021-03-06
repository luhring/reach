package api

import (
	"net"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"

	reachAWS "github.com/luhring/reach/reach/aws"
)

// VPC queries the AWS API for a VPC matching the given ID.
func (provider *ResourceProvider) VPC(id string) (*reachAWS.VPC, error) {
	input := &ec2.DescribeVpcsInput{
		VpcIds: []*string{
			aws.String(id),
		},
	}
	result, err := provider.ec2.DescribeVpcs(input)
	if err != nil {
		return nil, err
	}

	if err = ensureSingleResult(len(result.Vpcs), "VPC", id); err != nil {
		return nil, err
	}

	vpc := newVPCFromAPI(result.Vpcs[0])
	return &vpc, nil
}

func newVPCFromAPI(vpc *ec2.Vpc) reachAWS.VPC {
	ipv4CIDRs := cidrs(vpc.CidrBlockAssociationSet)
	ipv6CIDRs := ipv6CIDRs(vpc.Ipv6CidrBlockAssociationSet)

	return reachAWS.VPC{
		ID:        aws.StringValue(vpc.VpcId),
		IPv4CIDRs: ipv4CIDRs,
		IPv6CIDRs: ipv6CIDRs,
	}
}

func cidrs(associationSet []*ec2.VpcCidrBlockAssociation) []net.IPNet {
	cidrs := make([]net.IPNet, len(associationSet))

	for i, association := range associationSet {
		cidrs[i] = cidr(association)
	}

	return cidrs
}

func cidr(association *ec2.VpcCidrBlockAssociation) net.IPNet {
	if association == nil {
		return net.IPNet{}
	}

	_, cidr, err := net.ParseCIDR(aws.StringValue(association.CidrBlock))
	if err != nil {
		return net.IPNet{}
	}

	return *cidr
}

func ipv6CIDRs(associationSet []*ec2.VpcIpv6CidrBlockAssociation) []net.IPNet {
	cidrs := make([]net.IPNet, len(associationSet))

	for i, association := range associationSet {
		cidrs[i] = ipv6CIDR(association)
	}

	return cidrs
}

func ipv6CIDR(association *ec2.VpcIpv6CidrBlockAssociation) net.IPNet {
	if association == nil {
		return net.IPNet{}
	}

	_, cidr, err := net.ParseCIDR(aws.StringValue(association.Ipv6CidrBlock))
	if err != nil {
		return net.IPNet{}
	}

	return *cidr
}
