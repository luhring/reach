package api

import (
	"fmt"
	"net"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"

	reachAWS "github.com/luhring/reach/reach/aws"
)

// ElasticNetworkInterface queries the AWS API for an ElasticNetworkInterface matching the given ID.
func (client *DomainClient) ElasticNetworkInterface(id string) (*reachAWS.ElasticNetworkInterface, error) {
	if r := client.cachedResource(reachAWS.ElasticNetworkInterfaceRef(id)); r != nil {
		if v, ok := r.(*reachAWS.ElasticNetworkInterface); ok {
			return v, nil
		}
	}

	input := &ec2.DescribeNetworkInterfacesInput{
		NetworkInterfaceIds: []*string{
			aws.String(id),
		},
	}
	result, err := client.ec2.DescribeNetworkInterfaces(input)
	if err != nil {
		return nil, err
	}

	if err = ensureSingleResult(len(result.NetworkInterfaces), reachAWS.ResourceKindElasticNetworkInterface, id); err != nil {
		return nil, err
	}

	networkInterface := newElasticNetworkInterfaceFromAPI(result.NetworkInterfaces[0])
	client.cacheResource(networkInterface)
	return &networkInterface, nil
}

// ElasticNetworkInterfaceByIP queries the AWS API for the ElasticNetworkInterface that's associated with the specified IP address.
func (client *DomainClient) ElasticNetworkInterfaceByIP(ip net.IP) (*reachAWS.ElasticNetworkInterface, error) {
	filterNames := []string{
		"private-ip-address",
		"association.public-ip",
		"ipv6-addresses.ipv6-address",
	}

	for _, name := range filterNames {
		input := &ec2.DescribeNetworkInterfacesInput{
			Filters: []*ec2.Filter{
				{
					Name:   aws.String(name),
					Values: aws.StringSlice([]string{ip.String()}),
				},
			},
		}

		result, err := client.ec2.DescribeNetworkInterfaces(input)
		if err != nil {
			// TODO: Try to differentiate a "Not Found" error vs. more serious errors
			continue
		}
		if err := ensureSingleResult(len(result.NetworkInterfaces), reachAWS.ResourceKindElasticNetworkInterface, ip.String()); err != nil {
			return nil, err
		}

		eniResult := result.NetworkInterfaces[0]
		eni := newElasticNetworkInterfaceFromAPI(eniResult)
		return &eni, nil
	}

	return nil, fmt.Errorf("unable to find matching elastic network interface for IP (%s), either because no such ENI exists or because a more serious error has occurred (such as with the network connection or with AWS authentication)", ip)
}

func newElasticNetworkInterfaceFromAPI(eni *ec2.NetworkInterface) reachAWS.ElasticNetworkInterface {
	publicIPv4Address := publicIPAddress(eni.Association)
	privateIPv4Addresses := privateIPAddresses(eni.PrivateIpAddresses)
	ipv6Addresses := ipv6Addresses(eni.Ipv6Addresses)

	return reachAWS.ElasticNetworkInterface{
		ID:                   aws.StringValue(eni.NetworkInterfaceId),
		NameTag:              nameTag(eni.TagSet),
		SubnetID:             aws.StringValue(eni.SubnetId),
		VPCID:                aws.StringValue(eni.VpcId),
		SecurityGroupIDs:     securityGroupIDs(eni.Groups),
		PublicIPv4Address:    publicIPv4Address,
		PrivateIPv4Addresses: privateIPv4Addresses,
		IPv6Addresses:        ipv6Addresses,
		SrcDstCheck:          aws.BoolValue(eni.SourceDestCheck),
	}
}

func securityGroupID(identifier *ec2.GroupIdentifier) string {
	if identifier == nil {
		return ""
	}

	return aws.StringValue(identifier.GroupId)
}

func securityGroupIDs(identifiers []*ec2.GroupIdentifier) []string {
	ids := make([]string, len(identifiers))

	for i, identifier := range identifiers {
		ids[i] = securityGroupID(identifier)
	}

	return ids
}

func privateIPAddress(address *ec2.NetworkInterfacePrivateIpAddress) net.IP {
	if address == nil {
		return net.IP{}
	}

	return net.ParseIP(aws.StringValue(address.PrivateIpAddress))
}

func privateIPAddresses(addresses []*ec2.NetworkInterfacePrivateIpAddress) []net.IP {
	ips := make([]net.IP, len(addresses))

	for i, address := range addresses {
		ips[i] = privateIPAddress(address)
	}

	return ips
}

func ipv6Address(address *ec2.NetworkInterfaceIpv6Address) net.IP {
	if address == nil {
		return net.IP{}
	}

	return net.ParseIP(aws.StringValue(address.Ipv6Address))
}

func ipv6Addresses(addresses []*ec2.NetworkInterfaceIpv6Address) []net.IP {
	ips := make([]net.IP, len(addresses))

	for i, address := range addresses {
		ips[i] = ipv6Address(address)
	}

	return ips
}

func publicIPAddress(association *ec2.NetworkInterfaceAssociation) net.IP {
	if association == nil {
		return net.IP{}
	}

	return net.ParseIP(aws.StringValue(association.PublicIp))
}
