package aws

import (
	"fmt"
	"net"

	"github.com/luhring/reach/reach"
)

// The DomainClient interface wraps all of the necessary methods for accessing AWS-specific resources.
type DomainClient interface {
	AllEC2Instances() ([]EC2Instance, error)
	EC2Instance(id string) (*EC2Instance, error)
	EC2InstanceByENI(eniID string) (*EC2Instance, error)
	ElasticNetworkInterface(id string) (*ElasticNetworkInterface, error)
	ElasticNetworkInterfaceByIP(ip net.IP) (*ElasticNetworkInterface, error)
	InternetGateway(id string) (*InternetGateway, error)
	NATGateway(id string) (*NATGateway, error)
	NetworkACL(id string) (*NetworkACL, error)
	ResolveSecurityGroupReference(sgID string) ([]ElasticNetworkInterface, error)
	RouteTable(id string) (*RouteTable, error)
	RouteTableForGateway(id string) (*RouteTable, error)
	SecurityGroup(id string) (*SecurityGroup, error)
	SecurityGroupReference(id, accountID string) (*SecurityGroupReference, error)
	Subnet(id string) (*Subnet, error)
	SubnetsByVPC(id string) ([]Subnet, error)
	VPC(id string) (*VPC, error)
}

func unpackDomainClient(resolver reach.DomainClientResolver) (DomainClient, error) {
	d := resolver.Resolve(ResourceDomainAWS)
	if d == nil {
		return nil, fmt.Errorf("DomainClientResolver has no entry for domain '%s'", ResourceDomainAWS)
	}
	domainClient, ok := d.(DomainClient)
	if !ok {
		return nil, fmt.Errorf("DomainClient interface not implemented correctly for domain '%s'", ResourceDomainAWS)
	}
	return domainClient, nil
}
