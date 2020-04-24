package api

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"

	reachAWS "github.com/luhring/reach/reach/aws"
)

// EC2Instance queries the AWS API for an EC2 instance matching the given ID.
func (client *DomainClient) EC2Instance(id string) (*reachAWS.EC2Instance, error) {
	input := &ec2.DescribeInstancesInput{
		InstanceIds: []*string{
			aws.String(id),
		},
	}
	result, err := client.ec2.DescribeInstances(input)
	if err != nil {
		return nil, err
	}

	instances, err := extractEC2Instances(result.Reservations)
	if err != nil {
		return nil, err
	}

	if len(instances) == 0 {
		return nil, fmt.Errorf("AWS API returned no instances for ID '%s'", id)
	}

	if len(instances) > 1 {
		return nil, fmt.Errorf("AWS API returned more than one instance for ID '%s'", id)
	}

	instance := instances[0]
	return &instance, nil
}

// AllEC2Instances queries the AWS API for all EC2 instances.
func (client *DomainClient) AllEC2Instances() ([]reachAWS.EC2Instance, error) {
	const errFormat = "unable to get all EC2 instances: %v"

	describeInstancesOutput, err := client.ec2.DescribeInstances(nil)

	if err != nil {
		return nil, fmt.Errorf(errFormat, err)
	}

	reservations := describeInstancesOutput.Reservations
	instances, err := extractEC2Instances(reservations)
	if err != nil {
		return nil, fmt.Errorf(errFormat, err)
	}

	return instances, nil
}

func newEC2InstanceFromAPI(instance *ec2.Instance) reachAWS.EC2Instance {
	return reachAWS.EC2Instance{
		ID:                          aws.StringValue(instance.InstanceId),
		NameTag:                     nameTag(instance.Tags),
		State:                       aws.StringValue(instance.State.Name),
		NetworkInterfaceAttachments: networkInterfaceAttachments(instance),
	}
}

func extractEC2Instances(reservations []*ec2.Reservation) ([]reachAWS.EC2Instance, error) {
	var instances []reachAWS.EC2Instance

	for _, r := range reservations {
		for _, i := range r.Instances {
			instance := newEC2InstanceFromAPI(i)
			instances = append(instances, instance)
		}
	}

	return instances, nil
}

func networkInterfaceAttachments(instance *ec2.Instance) []reachAWS.NetworkInterfaceAttachment {
	var attachments []reachAWS.NetworkInterfaceAttachment

	if instance.NetworkInterfaces != nil && len(instance.NetworkInterfaces) > 0 {
		for _, networkInterface := range instance.NetworkInterfaces {
			attachments = append(attachments, newNetworkInterfaceAttachmentFromAPI(networkInterface))
		}
	}

	return attachments
}

func newNetworkInterfaceAttachmentFromAPI(networkInterface *ec2.InstanceNetworkInterface) reachAWS.NetworkInterfaceAttachment {
	return reachAWS.NetworkInterfaceAttachment{
		ID:                        aws.StringValue(networkInterface.Attachment.AttachmentId),
		ElasticNetworkInterfaceID: aws.StringValue(networkInterface.NetworkInterfaceId),
		DeviceIndex:               aws.Int64Value(networkInterface.Attachment.DeviceIndex),
	}
}
