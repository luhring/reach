package apiclient

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"

	reachAWS "github.com/luhring/reach/reach/aws"
)

// SecurityGroupReference queries the AWS API for a security group matching the given ID, but returns a security group reference representation instead of the full security group representation.
func (client *DomainClient) SecurityGroupReference(id, accountID string) (*reachAWS.SecurityGroupReference, error) {
	// TODO: Incorporate account ID in search.
	// In the meantime, this will be a known bug, where other accounts are not considered.

	sg, err := client.SecurityGroup(id)
	if err != nil {
		return nil, err
	}

	return &reachAWS.SecurityGroupReference{
		ID:        sg.ID,
		AccountID: "",
		NameTag:   sg.NameTag,
		GroupName: sg.GroupName,
	}, nil
}

// ResolveSecurityGroupReference queries the AWS API to determine the set of ElasticNetworkInterfaces to which the specified SecurityGroup is attached.
func (client *DomainClient) ResolveSecurityGroupReference(sgID string) ([]reachAWS.ElasticNetworkInterface, error) {
	input := &ec2.DescribeNetworkInterfacesInput{
		Filters: []*ec2.Filter{
			{
				Name:   aws.String("group-id"),
				Values: aws.StringSlice([]string{sgID}),
			},
		},
	}

	var enis []reachAWS.ElasticNetworkInterface

	results, err := client.ec2.DescribeNetworkInterfaces(input)
	if err != nil {
		return nil, fmt.Errorf("unable to get network interfaces by security group ID: %v", err)
	}
	for _, resultENI := range results.NetworkInterfaces {
		eni := newElasticNetworkInterfaceFromAPI(resultENI)
		enis = append(enis, eni)
	}

	// TODO: Check other resources besides ENIs?

	return enis, nil
}
