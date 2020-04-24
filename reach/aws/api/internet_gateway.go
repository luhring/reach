package api

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"

	reachAWS "github.com/luhring/reach/reach/aws"
)

// InternetGateway queries the AWS API for an Internet gateway matching the given ID.
func (client *DomainClient) InternetGateway(id string) (*reachAWS.InternetGateway, error) {
	input := &ec2.DescribeInternetGatewaysInput{
		InternetGatewayIds: []*string{aws.String(id)},
	}
	result, err := client.ec2.DescribeInternetGateways(input)
	if err != nil {
		return nil, err
	}

	if err = ensureSingleResult(len(result.InternetGateways), reachAWS.ResourceKindInternetGateway, id); err != nil {
		return nil, err
	}

	internetGateway := result.InternetGateways[0]

	vpcID := vpcIDFromInternetGateway(internetGateway)

	return &reachAWS.InternetGateway{
		ID:    id,
		VPCID: vpcID,
	}, nil
}

func vpcIDFromInternetGateway(igw *ec2.InternetGateway) string {
	return aws.StringValue(igw.Attachments[0].VpcId)
}
