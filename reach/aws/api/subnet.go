package api

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"

	reachAWS "github.com/luhring/reach/reach/aws"
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
		return nil, err
	}

	if err = ensureSingleResult(len(result.Subnets), reachAWS.ResourceKindSubnet, id); err != nil {
		return nil, err
	}

	awsSubnet := result.Subnets[0]

	networkACLID, err := client.networkACLIDFromSubnetID(aws.StringValue(awsSubnet.SubnetId))
	if err != nil {
		return nil, err
	}

	routeTableID, err := client.routeTableIDFromSubnetID(aws.StringValue(awsSubnet.SubnetId))
	if err != nil {
		return nil, err
	}

	subnet := newSubnetFromAPI(result.Subnets[0], networkACLID, routeTableID)
	client.cacheResource(subnet)
	return &subnet, nil
}

func newSubnetFromAPI(subnet *ec2.Subnet, networkACLID, routeTableID string) reachAWS.Subnet {
	return reachAWS.Subnet{
		ID:           aws.StringValue(subnet.SubnetId),
		NetworkACLID: networkACLID,
		RouteTableID: routeTableID,
		VPCID:        aws.StringValue(subnet.VpcId),
	}
}

func (client *DomainClient) networkACLIDFromSubnetID(id string) (string, error) {
	input := &ec2.DescribeNetworkAclsInput{
		Filters: generateEC2Filters(id),
	}

	result, err := client.ec2.DescribeNetworkAcls(input)
	if err != nil {
		return "", err
	}

	if err = ensureSingleResult(len(result.NetworkAcls), "network ACL (via subnet)", id); err != nil {
		return "", err
	}

	return aws.StringValue(result.NetworkAcls[0].NetworkAclId), nil
}

func (client *DomainClient) routeTableIDFromSubnetID(id string) (string, error) {
	input := &ec2.DescribeRouteTablesInput{
		Filters: generateEC2Filters(id),
	}

	result, err := client.ec2.DescribeRouteTables(input)
	if err != nil {
		return "", nil
	}

	if err = ensureSingleResult(len(result.RouteTables), "route table (via subnet)", id); err != nil {
		return "", nil
	}

	return aws.StringValue(result.RouteTables[0].RouteTableId), nil
}

func generateEC2Filters(subnetID string) []*ec2.Filter {
	return []*ec2.Filter{
		{
			Name: aws.String("association.subnet-subnetID"),
			Values: []*string{
				aws.String(subnetID),
			},
		},
	}
}
