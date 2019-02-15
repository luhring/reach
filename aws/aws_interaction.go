package aws

import (
	"net"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

// CreateAwsSession ...
func CreateAwsSession() (*session.Session, error) {
	session, errorCreatingSession := session.NewSessionWithOptions(
		session.Options{
			SharedConfigState: session.SharedConfigEnable,
		},
	)

	if errorCreatingSession != nil {
		return nil, errorCreatingSession
	}

	return session, nil
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

// GetAllInstancesUsingEc2Client ...
func GetAllInstancesUsingEc2Client(ec2Client *ec2.EC2) ([]*Instance, error) {
	describeInstancesOutput, errorFromDescribeInstances := ec2Client.DescribeInstances(nil)

	if errorFromDescribeInstances != nil {
		return nil, errorFromDescribeInstances
	}

	allReservations := describeInstancesOutput.Reservations
	return getAllInstancesFromReservations(allReservations), nil
}

func getAllInstancesFromReservations(reservations []*ec2.Reservation) []*Instance {
	var instances []*Instance

	for _, reservation := range reservations {
		for _, instanceInReservation := range reservation.Instances {
			instanceToAdd := &Instance{
				ID:                   aws.StringValue(instanceInReservation.InstanceId),
				NameTag:              getNameTagValueFromTags(instanceInReservation.Tags),
				PrivateIPv4Addresses: aggregatePrivateIPv4Addresses(instanceInReservation.NetworkInterfaces),
				PublicIPv4Addresses:  aggregatePublicIPv4Addresses(instanceInReservation.NetworkInterfaces),
				IPv6Addresses:        aggregrateIPv6Addresses(instanceInReservation.NetworkInterfaces),
				SecurityGroupIDs:     getSecurityGroupIDsForInstance(instanceInReservation),
				State:                aws.StringValue(instanceInReservation.State.Name),
				SubnetID:             aws.StringValue(instanceInReservation.SubnetId),
				VpcID:                aws.StringValue(instanceInReservation.VpcId),
			}
			instances = append(instances, instanceToAdd)
		}
	}

	return instances
}

func aggregatePrivateIPv4Addresses(networkInterfaces []*ec2.InstanceNetworkInterface) []net.IP {
	var privateIPv4Addresses []net.IP

	for _, networkInterface := range networkInterfaces {
		for _, privateIPv4Address := range networkInterface.PrivateIpAddresses {
			newIP := net.ParseIP(*privateIPv4Address.PrivateIpAddress)
			privateIPv4Addresses = append(privateIPv4Addresses, newIP)
		}
	}

	return privateIPv4Addresses
}

func aggregatePublicIPv4Addresses(networkInterfaces []*ec2.InstanceNetworkInterface) []net.IP {
	var publicIPv4Addresses []net.IP

	for _, networkInterface := range networkInterfaces {
		for _, privateIPv4Address := range networkInterface.PrivateIpAddresses {
			if privateIPv4Address.Association != nil {
				newIP := net.ParseIP(*privateIPv4Address.Association.PublicIp)
				publicIPv4Addresses = append(publicIPv4Addresses, newIP)
			}
		}
	}

	return publicIPv4Addresses
}

func aggregrateIPv6Addresses(networkInterfaces []*ec2.InstanceNetworkInterface) []net.IP {
	var ipv6Addresses []net.IP

	for _, networkInterface := range networkInterfaces {
		for _, ipv6Address := range networkInterface.Ipv6Addresses {
			newIP := net.ParseIP(*ipv6Address.Ipv6Address)
			ipv6Addresses = append(ipv6Addresses, newIP)
		}
	}

	return ipv6Addresses
}

func getSecurityGroupIDsForInstance(instance *ec2.Instance) []*string {
	var securityGroupIDs []*string

	for _, networkInterface := range instance.NetworkInterfaces {
		for _, securityGroup := range networkInterface.Groups {
			securityGroupIDs = append(securityGroupIDs, securityGroup.GroupId)
		}
	}

	return securityGroupIDs
}

func getSecurityGroupsForIds(ec2Client *ec2.EC2, securityGroupIds []*string) ([]*ec2.SecurityGroup, error) {
	parameters := &ec2.DescribeSecurityGroupsInput{
		GroupIds: securityGroupIds,
	}

	output, errorDescribingSecurityGroups := ec2Client.DescribeSecurityGroups(parameters)

	if errorDescribingSecurityGroups != nil {
		return nil, errorDescribingSecurityGroups
	}

	return output.SecurityGroups, nil
}
