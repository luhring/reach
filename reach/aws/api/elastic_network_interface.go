package api

import (
	"net"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"

	reachAWS "github.com/luhring/reach/reach/aws"
)

func (provider *ResourceProvider) GetElasticNetworkInterface(id string) (*reachAWS.ElasticNetworkInterface, error) {
	input := &ec2.DescribeNetworkInterfacesInput{
		NetworkInterfaceIds: []*string{
			aws.String(id),
		},
	}
	result, err := provider.ec2.DescribeNetworkInterfaces(input)
	if err != nil {
		return nil, err
	}

	if err = ensureSingleResult(len(result.NetworkInterfaces), "elastic network interface", id); err != nil {
		return nil, err
	}

	networkInterface := newElasticNetworkInterfaceFromAPI(result.NetworkInterfaces[0])
	return &networkInterface, nil
}

func newElasticNetworkInterfaceFromAPI(eni *ec2.NetworkInterface) reachAWS.ElasticNetworkInterface {
	publicIPv4Address := getPublicIPAddress(eni.Association)
	privateIPv4Addresses := getPrivateIPAddresses(eni.PrivateIpAddresses)
	ipv6Addresses := getIPv6Addresses(eni.Ipv6Addresses)

	return reachAWS.ElasticNetworkInterface{
		ID:                   aws.StringValue(eni.NetworkInterfaceId),
		NameTag:              getNameTag(eni.TagSet),
		SubnetID:             aws.StringValue(eni.SubnetId),
		VPCID:                aws.StringValue(eni.VpcId),
		SecurityGroupIDs:     getSecurityGroupIDs(eni.Groups),
		PublicIPv4Address:    publicIPv4Address,
		PrivateIPv4Addresses: privateIPv4Addresses,
		IPv6Addresses:        ipv6Addresses,
	}
}

func getSecurityGroupID(identifier *ec2.GroupIdentifier) string {
	if identifier == nil {
		return ""
	}

	return aws.StringValue(identifier.GroupId)
}

func getSecurityGroupIDs(identifiers []*ec2.GroupIdentifier) []string {
	ids := make([]string, len(identifiers))

	for i, identifier := range identifiers {
		ids[i] = getSecurityGroupID(identifier)
	}

	return ids
}

func getPrivateIPAddress(address *ec2.NetworkInterfacePrivateIpAddress) net.IP {
	if address == nil {
		return net.IP{}
	}

	return net.ParseIP(aws.StringValue(address.PrivateIpAddress))
}

func getPrivateIPAddresses(addresses []*ec2.NetworkInterfacePrivateIpAddress) []net.IP {
	ips := make([]net.IP, len(addresses))

	for i, address := range addresses {
		ips[i] = getPrivateIPAddress(address)
	}

	return ips
}

func getIPv6Address(address *ec2.NetworkInterfaceIpv6Address) net.IP {
	if address == nil {
		return net.IP{}
	}

	return net.ParseIP(aws.StringValue(address.Ipv6Address))
}

func getIPv6Addresses(addresses []*ec2.NetworkInterfaceIpv6Address) []net.IP {
	ips := make([]net.IP, len(addresses))

	for i, address := range addresses {
		ips[i] = getIPv6Address(address)
	}

	return ips
}

func getPublicIPAddress(association *ec2.NetworkInterfaceAssociation) net.IP {
	if association == nil {
		return net.IP{}
	}

	return net.ParseIP(aws.StringValue(association.PublicIp))
}
