package api

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"

	reachAWS "github.com/luhring/reach/reach/aws"
)

func (provider *ResourceProvider) GetSubnet(id string) (*reachAWS.Subnet, error) {
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

	subnet := newSubnetFromAPI(result.Subnets[0])
	return &subnet, nil
}

func newSubnetFromAPI(subnet *ec2.Subnet) reachAWS.Subnet {
	return reachAWS.Subnet{
		ID:    aws.StringValue(subnet.SubnetId),
		VPCID: aws.StringValue(subnet.VpcId),
	}
}
