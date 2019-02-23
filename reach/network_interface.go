package reach

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"log"
	"net"
)

type NetworkInterface struct {
	ID                 string
	Name               string
	PrivateIPAddresses []net.IP
	PublicIPAddress    net.IP
	SecurityGroups     []*SecurityGroup
	SubnetID           string
	VPCID              string
}

func NewNetworkInterface(networkInterface *ec2.InstanceNetworkInterface, findSecurityGroup func(id string) (*SecurityGroup, error)) (*NetworkInterface, error) {
	privateIPAddresses := make([]net.IP, len(networkInterface.PrivateIpAddresses))

	for i, address := range networkInterface.PrivateIpAddresses {
		privateIPAddresses[i] = NewIP(address)
	}

	securityGroupIDs := make([]string, len(networkInterface.Groups))

	for i, group := range networkInterface.Groups {
		securityGroupIDs[i] = aws.StringValue(group.GroupId)
	}

	securityGroups := make([]*SecurityGroup, len(networkInterface.Groups))

	for i, groupIdentifier := range networkInterface.Groups {
		id := aws.StringValue(groupIdentifier.GroupId)

		group, err := findSecurityGroup(id)
		if err != nil {
			return nil, fmt.Errorf("unable to find security group for new network interface: %v", err)
		}

		securityGroups[i] = group
	}

	var publicIPAddress net.IP

	if assoc := networkInterface.Association; assoc != nil {
		if pubIP := assoc.PublicIp; pubIP != nil {
			publicIPAddress = net.ParseIP(aws.StringValue(pubIP))
		}
	}

	return &NetworkInterface{
		ID:                 aws.StringValue(networkInterface.NetworkInterfaceId),
		Name:               aws.StringValue(networkInterface.NetworkInterfaceId),
		PrivateIPAddresses: privateIPAddresses,
		PublicIPAddress:    publicIPAddress,
		SecurityGroups:     securityGroups,
		SubnetID:           aws.StringValue(networkInterface.SubnetId),
		VPCID:              aws.StringValue(networkInterface.VpcId),
	}, nil
}

func NewIP(address *ec2.InstancePrivateIpAddress) net.IP {
	if address.PrivateIpAddress == nil {
		return nil
	}

	ip := net.ParseIP(aws.StringValue(address.PrivateIpAddress))

	if ip == nil {
		log.Printf("unable to parse IP address '%s'\n", aws.StringValue(address.PrivateIpAddress))
	}

	return ip
}
