package api

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"

	reachAWS "github.com/luhring/reach/reach/aws"
)

// Subnet queries the AWS API for a subnet matching the given ID.
func (provider *ResourceProvider) Subnet(id string) (*reachAWS.Subnet, error) {
	input := &ec2.DescribeSubnetsInput{
		SubnetIds: []*string{
			aws.String(id),
		},
	}
	result, err := provider.ec2.DescribeSubnets(input)
	if err != nil {
		return nil, err
	}

	if err = ensureSingleResult(len(result.Subnets), "subnet", id); err != nil {
		return nil, err
	}

	awsSubnet := result.Subnets[0]
	networkACLID, err := provider.networkACLIDFromSubnetID(aws.StringValue(awsSubnet.SubnetId))
	if err != nil {
		return nil, err
	}

	subnet := newSubnetFromAPI(result.Subnets[0], networkACLID)
	return &subnet, nil
}

func newSubnetFromAPI(subnet *ec2.Subnet, networkACLID string) reachAWS.Subnet {
	return reachAWS.Subnet{
		ID:           aws.StringValue(subnet.SubnetId),
		NetworkACLID: networkACLID,
		VPCID:        aws.StringValue(subnet.VpcId),
	}
}

func (provider *ResourceProvider) networkACLIDFromSubnetID(id string) (string, error) {
	input := &ec2.DescribeNetworkAclsInput{
		Filters: []*ec2.Filter{
			{
				Name: aws.String("association.subnet-id"),
				Values: []*string{
					aws.String(id),
				},
			},
		},
	}
	result, err := provider.ec2.DescribeNetworkAcls(input)
	if err != nil {
		return "", err
	}

	if err = ensureSingleResult(len(result.NetworkAcls), "network ACL (via subnet)", id); err != nil {
		return "", err
	}

	return aws.StringValue(result.NetworkAcls[0].NetworkAclId), nil
}
