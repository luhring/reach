package apiclient

import (
	"fmt"
	"net"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/ec2"

	reachAWS "github.com/luhring/reach/reach/aws"
	"github.com/luhring/reach/reach/reacherr"
)

// Subnet queries the AWS API for a subnet matching the given ID.
func (client *DomainClient) Subnet(id string) (*reachAWS.Subnet, error) {
	if r := client.cachedResource(reachAWS.SubnetRef(id)); r != nil {
		if v, ok := r.(*reachAWS.Subnet); ok {
			return v, nil
		}
	}

	input := &ec2.DescribeSubnetsInput{
		SubnetIds: []*string{
			aws.String(id),
		},
	}
	result, err := client.ec2.DescribeSubnets(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			return nil, reacherr.New(err, awsErrMessage(aerr))
		}
		return nil, err
	}

	if err = ensureSingleResult(len(result.Subnets), reachAWS.ResourceKindSubnet, id); err != nil {
		return nil, err
	}

	awsSubnet := result.Subnets[0]

	subnet, err := client.newSubnetFromAPI(awsSubnet)
	if err != nil {
		return nil, err
	}

	client.cacheResource(*subnet)
	return subnet, nil
}

// SubnetsByVPC returns the set of Subnets that exist within the specified VPC.
func (client *DomainClient) SubnetsByVPC(id string) ([]reachAWS.Subnet, error) {
	input := &ec2.DescribeSubnetsInput{
		Filters: []*ec2.Filter{
			{
				Name:   aws.String("vpc-id"),
				Values: aws.StringSlice([]string{id}),
			},
		},
	}
	results, err := client.ec2.DescribeSubnets(input)
	if err != nil {
		return nil, err
	}

	var subnets []reachAWS.Subnet
	for _, s := range results.Subnets {
		subnet, err := client.newSubnetFromAPI(s)
		if err != nil {
			return nil, err
		}
		client.cacheResource(*subnet)
		subnets = append(subnets, *subnet)
	}

	return subnets, nil
}

func (client *DomainClient) newSubnetFromAPI(subnet *ec2.Subnet) (*reachAWS.Subnet, error) {
	networkACLID, err := client.networkACLIDFromSubnet(subnet)
	if err != nil {
		return nil, err
	}

	routeTableID := client.routeTableIDFromSubnetID(subnet)

	ipv4CIDR, err := ipv4CIDRFromSubnet(subnet)
	if err != nil {
		return nil, err
	}
	ipv6CIDR, err := ipv6CIDRFromSubnet(subnet)
	if err != nil {
		return nil, err
	}

	return &reachAWS.Subnet{
		ID:           aws.StringValue(subnet.SubnetId),
		NetworkACLID: networkACLID,
		RouteTableID: routeTableID,
		VPCID:        aws.StringValue(subnet.VpcId),
		IPv4CIDR:     *ipv4CIDR,
		IPv6CIDR:     ipv6CIDR,
	}, nil
}

func ipv4CIDRFromSubnet(subnet *ec2.Subnet) (*net.IPNet, error) {
	_, cidr, err := net.ParseCIDR(aws.StringValue(subnet.CidrBlock))
	if err != nil {
		return nil, err
	}
	return cidr, nil
}

func ipv6CIDRFromSubnet(subnet *ec2.Subnet) (*net.IPNet, error) {
	set := subnet.Ipv6CidrBlockAssociationSet
	if set == nil {
		return nil, nil
	}

	if numSets := len(set); numSets != 1 {
		return nil, fmt.Errorf("could not obtain IPv6 CIDR block for subnet, expected response to contain exactly 1 association set (%d sets found)", numSets)
	}

	assoc := set[0]
	_, cidr, err := net.ParseCIDR(aws.StringValue(assoc.Ipv6CidrBlock))
	if err != nil {
		return nil, err
	}
	return cidr, nil
}

func (client *DomainClient) networkACLIDFromSubnet(subnet *ec2.Subnet) (string, error) {
	subnetID := subnet.SubnetId
	input := &ec2.DescribeNetworkAclsInput{
		Filters: []*ec2.Filter{
			{
				Name:   aws.String("association.subnet-id"),
				Values: []*string{subnetID},
			},
		},
	}

	result, err := client.ec2.DescribeNetworkAcls(input)
	if err != nil {
		return "", err
	}

	if err = ensureSingleResult(len(result.NetworkAcls), "network ACL (via subnet)", *subnetID); err != nil {
		return "", err
	}

	return aws.StringValue(result.NetworkAcls[0].NetworkAclId), nil
}

func (client *DomainClient) routeTableIDFromSubnetID(subnet *ec2.Subnet) string {
	subnetID := subnet.SubnetId
	input := &ec2.DescribeRouteTablesInput{
		Filters: []*ec2.Filter{
			{
				Name:   aws.String("association.subnet-id"),
				Values: []*string{subnetID},
			},
		},
	}

	result, err := client.ec2.DescribeRouteTables(input)
	if err != nil {
		return ""
	}

	if err = ensureSingleResult(len(result.RouteTables), "route table (via subnet)", *subnetID); err != nil {
		return ""
	}

	return aws.StringValue(result.RouteTables[0].RouteTableId)
}
