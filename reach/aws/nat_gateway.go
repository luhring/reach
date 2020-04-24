package aws

import (
	"net"

	"github.com/luhring/reach/reach"
)

// ResourceKindNATGateway specifies the unique name for the NAT gateway kind of resource.
const ResourceKindNATGateway = "NATGateway"

// A NATGateway resource representation.
type NATGateway struct {
	ID        string
	SubnetID  string
	VPCID     string
	PrivateIP net.IP
	PublicIP  net.IP
}

// ToResource returns the NAT gateway converted to a generalized Reach resource.
func (ngw NATGateway) ToResource() reach.Resource {
	return reach.Resource{
		Kind:       ResourceKindNATGateway,
		Properties: ngw,
	}
}

// Dependencies returns a collection of the NAT gateway's resource dependencies.
func (ngw NATGateway) Dependencies(provider DomainClient) (*reach.ResourceCollection, error) {
	rc := reach.NewResourceCollection()

	subnet, err := provider.Subnet(ngw.SubnetID)
	if err != nil {
		return nil, err
	}
	rc.Put(reach.ResourceReference{
		Domain: ResourceDomainAWS,
		Kind:   ResourceKindSubnet,
		ID:     ngw.SubnetID,
	}, subnet.ToResource())

	vpc, err := provider.VPC(ngw.VPCID)
	if err != nil {
		return nil, err
	}
	rc.Put(reach.ResourceReference{
		Domain: ResourceDomainAWS,
		Kind:   ResourceKindVPC,
		ID:     ngw.VPCID,
	}, vpc.ToResource())

	return rc, nil
}
