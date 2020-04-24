package aws

import (
	"errors"
	"fmt"
	"net"

	"github.com/luhring/reach/reach"
)

// The ResourceGetter interface wraps all of the necessary methods for accessing AWS-specific resources.
type ResourceGetter interface {
	AllEC2Instances() ([]EC2Instance, error)
	EC2Instance(id string) (*EC2Instance, error)
	EC2InstanceByENI(eniID string) (*EC2Instance, error)
	ElasticNetworkInterface(id string) (*ElasticNetworkInterface, error)
	InternetGateway(id string) (*InternetGateway, error)
	NATGateway(id string) (*NATGateway, error)
	NetworkACL(id string) (*NetworkACL, error)
	ResolveSecurityGroupReference(sgID string) ([]net.IPNet, error)
	RouteTable(id string) (*RouteTable, error)
	SecurityGroup(id string) (*SecurityGroup, error)
	SecurityGroupReference(id, accountID string) (*SecurityGroupReference, error)
	Subnet(id string) (*Subnet, error)
	VPC(id string) (*VPC, error)
}

func unpackResourceGetter(domains reach.DomainProvider) (ResourceGetter, error) {
	domainResourceGetter := domains.Domain(ResourceDomainAWS)
	if domainResourceGetter == nil {
		return nil, fmt.Errorf("DomainProvider has no entry for domain '%s'", ResourceDomainAWS)
	}
	resourceGetter, ok := domainResourceGetter.(ResourceGetter)
	if !ok {
		return nil, errors.New("domain ResourceGetter interface not implemented correctly")
	}
	return resourceGetter, nil
}
