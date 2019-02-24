package reach

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"strings"
)

const sourceWord = "source"
const destinationWord = "destination"

type EC2Instance struct {
	ID                string
	NameTag           string
	State             string
	NetworkInterfaces []*NetworkInterface
}

func NewEC2Instance(instance *ec2.Instance, findSecurityGroup func(id string) (*SecurityGroup, error)) (*EC2Instance, error) {
	networkInterfaces := make([]*NetworkInterface, len(instance.NetworkInterfaces))

	for i, networkInterface := range instance.NetworkInterfaces {
		newInterface, err := NewNetworkInterface(networkInterface, findSecurityGroup)
		if err != nil {
			return nil, fmt.Errorf("unable to create new EC2Instance object due to network interface error: %v", err)
		}

		networkInterfaces[i] = newInterface
	}

	return &EC2Instance{
		ID:                aws.StringValue(instance.InstanceId),
		NameTag:           getNameTagValueFromTags(instance.Tags),
		State:             aws.StringValue(instance.State.Name),
		NetworkInterfaces: networkInterfaces,
	}, nil
}

func (i *EC2Instance) doesStateAllowAccess() bool {
	const running = "running"
	return strings.EqualFold(i.State, running)
}

func (i *EC2Instance) LongName() string {
	if len(i.NameTag) >= 1 {
		return fmt.Sprintf("\"%v\" (%v)", i.NameTag, i.ID)
	}

	return i.ID
}

func getNameTagValueFromTags(tags []*ec2.Tag) string {
	const keyForNameTag = "Name"

	for _, tag := range tags {
		if *tag.Key == keyForNameTag {
			return aws.StringValue(tag.Value)
		}
	}

	return ""
}