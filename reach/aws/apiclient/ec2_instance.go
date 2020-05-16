package apiclient

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/ec2"

	reachAWS "github.com/luhring/reach/reach/aws"
	"github.com/luhring/reach/reach/reacherr"
)

// EC2Instance queries the AWS API for an EC2 instance matching the given ID.
func (client *DomainClient) EC2Instance(id string) (*reachAWS.EC2Instance, error) {
	if r := client.cachedResource(reachAWS.EC2InstanceRef(id)); r != nil {
		if v, ok := r.(*reachAWS.EC2Instance); ok {
			return v, nil
		}
	}

	input := &ec2.DescribeInstancesInput{
		InstanceIds: []*string{
			aws.String(id),
		},
	}
	result, err := client.ec2.DescribeInstances(input)
	if err != nil {
		return nil, err
	}

	instances := extractEC2Instances(result.Reservations)

	if len(instances) == 0 {
		return nil, fmt.Errorf("AWS API returned no instances for ID '%s'", id)
	}

	if len(instances) > 1 {
		return nil, fmt.Errorf("AWS API returned more than one instance for ID '%s'", id)
	}

	instance := instances[0]
	client.cacheResource(instance)
	return &instance, nil
}

// AllEC2Instances queries the AWS API for all EC2 instances.
func (client *DomainClient) AllEC2Instances() ([]reachAWS.EC2Instance, error) {
	describeInstancesOutput, err := client.ec2.DescribeInstances(nil)
	if err != nil {
		msg := err.Error()
		if awsErr, ok := err.(awserr.Error); ok {
			msg = awsErrMessage(awsErr)
			return nil, reacherr.New(err, msg)
		}
		return nil, err
	}

	reservations := describeInstancesOutput.Reservations
	instances := extractEC2Instances(reservations)
	for _, i := range instances {
		client.cacheResource(i)
	}
	return instances, nil
}

// EC2InstanceByENI queries the AWS API for the EC2Instance that's associated with the specified ID of an ElasticNetworkInterface.
func (client *DomainClient) EC2InstanceByENI(eniID string) (*reachAWS.EC2Instance, error) {
	// TODO: Intelligently cache the result of this query

	input := &ec2.DescribeInstancesInput{
		Filters: []*ec2.Filter{
			{
				Name:   aws.String("network-interface.network-interface-id"),
				Values: aws.StringSlice([]string{eniID}),
			},
		},
	}

	results, err := client.ec2.DescribeInstances(input)
	if err != nil {
		return nil, fmt.Errorf("unable to get EC2 instance by ENI (ENI ID: %s): %v", eniID, err)
	}

	instances := extractEC2Instances(results.Reservations)
	if err = ensureSingleResult(len(instances), reachAWS.ResourceKindEC2Instance, eniID); err != nil {
		return nil, fmt.Errorf("unable to get EC2 instance by ENI (ENI ID: %s): %v", eniID, err)
	}

	i := instances[0]
	client.cacheResource(i)
	return &i, nil
}

func newEC2InstanceFromAPI(instance *ec2.Instance) reachAWS.EC2Instance {
	return reachAWS.EC2Instance{
		ID:                          aws.StringValue(instance.InstanceId),
		NameTag:                     nameTag(instance.Tags),
		State:                       aws.StringValue(instance.State.Name),
		NetworkInterfaceAttachments: networkInterfaceAttachments(instance),
	}
}

func extractEC2Instances(reservations []*ec2.Reservation) []reachAWS.EC2Instance {
	var instances []reachAWS.EC2Instance

	for _, r := range reservations {
		for _, i := range r.Instances {
			instance := newEC2InstanceFromAPI(i)
			instances = append(instances, instance)
		}
	}

	return instances
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
