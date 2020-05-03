package api

import (
	"errors"
	"fmt"
	"net"
	"strconv"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"

	"github.com/luhring/reach/reach"
	reachAWS "github.com/luhring/reach/reach/aws"
)

// DomainClient implements an AWS DomainClient using the AWS API (via the AWS SDK).
type DomainClient struct {
	session *session.Session
	ec2     *ec2.EC2
	cache   reach.Cache
}

// NewDomainClient returns a reference to a new DomainClient for the AWS API.
func NewDomainClient(cache reach.Cache) *DomainClient {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	})) // TODO: Don't call session.Must —- return error, and don't panic, this is a library after all!

	ec2Client := ec2.New(sess)

	return &DomainClient{
		session: sess,
		ec2:     ec2Client,
		cache:   cache,
	}
}

func (client *DomainClient) cacheResource(r reach.Referable) {
	client.cache.Put(r.Ref().String(), r)
}

func (client *DomainClient) cachedResource(ref reach.UniversalReference) interface{} {
	return client.cache.Get(ref.String())
}

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

func (client *DomainClient) RouteTableForGateway(id string) (*reachAWS.RouteTable, error) {
	panic("implement me")
}

func (client *DomainClient) SubnetsByVPC(id string) ([]reachAWS.Subnet, error) {
	input := &ec2.DescribeSubnetsInput{
		Filters: []*ec2.Filter{
			{
				Name:   aws.String("vpc-id"),
				Values: aws.StringSlice([]string{id}),
			},
		},
	}
	results, err := client.ec2.DescribeSubnets(input)
	if err != nil {
		return nil, err
	}

	var subnets []reachAWS.Subnet
	for _, s := range results.Subnets {
		networkACLID, err := client.networkACLIDFromSubnetID(aws.StringValue(s.SubnetId))
		if err != nil {
			return nil, err
		}

		routeTableID, err := client.routeTableIDFromSubnetID(aws.StringValue(s.SubnetId))
		if err != nil {
			return nil, err
		}

		subnet := newSubnetFromAPI(s, networkACLID, routeTableID)
		subnets = append(subnets, subnet)
	}

	return subnets, nil
}

func nameTag(tags []*ec2.Tag) string {
	if tags != nil && len(tags) > 0 {
		for _, tag := range tags {
			if aws.StringValue(tag.Key) == "Name" {
				return aws.StringValue(tag.Value)
			}
		}
	}

	return ""
}

func ensureSingleResult(resultSetLength int, entity reach.Kind, id string) error {
	if resultSetLength == 0 {
		return fmt.Errorf("AWS API did not return a %s for ID '%s'", entity, id)
	}

	if resultSetLength > 1 {
		return fmt.Errorf("AWS API returned more than one %s for ID '%s'", entity, id)
	}

	return nil
}

func convertAWSIPProtocolStringToProtocol(ipProtocol *string) (reach.Protocol, error) {
	if ipProtocol == nil {
		return 0, errors.New("unexpected nil ipProtocol")
	}

	protocolString := strings.ToLower(aws.StringValue(ipProtocol))

	if p, err := strconv.ParseInt(protocolString, 10, 64); err == nil {
		var protocol = reach.Protocol(p)
		return protocol, nil
	}

	var protocolNumber reach.Protocol

	switch protocolString {
	case "tcp":
		protocolNumber = reach.ProtocolTCP
	case "udp":
		protocolNumber = reach.ProtocolUDP
	case "icmp":
		protocolNumber = reach.ProtocolICMPv4
	case "icmpv6":
		protocolNumber = reach.ProtocolICMPv6
	default:
		return 0, errors.New("unrecognized ipProtocol value")
	}

	return protocolNumber, nil
}
