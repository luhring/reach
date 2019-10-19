package aws

import (
	"fmt"
	"net"
	"strings"

	"github.com/luhring/reach/reach"
)

// ResourceKindElasticNetworkInterface specifies the unique name for the elastic network interface kind of resource.
const ResourceKindElasticNetworkInterface = "ElasticNetworkInterface"

// An ElasticNetworkInterface resource representation.
type ElasticNetworkInterface struct {
	ID                   string
	NameTag              string `json:"NameTag,omitempty"`
	SubnetID             string
	VPCID                string
	SecurityGroupIDs     []string
	PublicIPv4Address    net.IP   `json:"PublicIPv4Address,omitempty"`
	PrivateIPv4Addresses []net.IP `json:"PrivateIPv4Addresses,omitempty"`
	IPv6Addresses        []net.IP `json:"IPv6Addresses,omitempty"`
}

// ElasticNetworkInterfaceFromNetworkPoint extracts the ElasticNetworkInterface from the lineage of the specified network point.
func ElasticNetworkInterfaceFromNetworkPoint(point reach.NetworkPoint, rc *reach.ResourceCollection) *ElasticNetworkInterface {
	for _, ancestor := range point.Lineage { // assumes there will only be one ENI among the ancestors
		if ancestor.Domain == ResourceDomainAWS && ancestor.Kind == ResourceKindElasticNetworkInterface {
			eni := rc.Get(ancestor).Properties.(ElasticNetworkInterface)
			return &eni
		}
	}

	return nil
}

// ToResource returns the elastic network interface converted to a generalized Reach resource.
func (eni ElasticNetworkInterface) ToResource() reach.Resource {
	return reach.Resource{
		Kind:       ResourceKindElasticNetworkInterface,
		Properties: eni,
	}
}

// ToResourceReference returns a resource reference to uniquely identify the elastic network interface.
func (eni ElasticNetworkInterface) ToResourceReference() reach.ResourceReference {
	return reach.ResourceReference{
		Domain: ResourceDomainAWS,
		Kind:   ResourceKindElasticNetworkInterface,
		ID:     eni.ID,
	}
}

// Dependencies returns a collection of the elastic network interface's resource dependencies.
func (eni ElasticNetworkInterface) Dependencies(provider ResourceProvider) (*reach.ResourceCollection, error) {
	rc := reach.NewResourceCollection()

	subnet, err := provider.GetSubnet(eni.SubnetID)
	if err != nil {
		return nil, err
	}
	rc.Put(reach.ResourceReference{
		Domain: ResourceDomainAWS,
		Kind:   ResourceKindSubnet,
		ID:     subnet.ID,
	}, subnet.ToResource())

	vpc, err := provider.GetVPC(eni.VPCID)
	if err != nil {
		return nil, err
	}
	rc.Put(reach.ResourceReference{
		Domain: ResourceDomainAWS,
		Kind:   ResourceKindVPC,
		ID:     vpc.ID,
	}, vpc.ToResource())

	for _, sgID := range eni.SecurityGroupIDs {
		sg, err := provider.GetSecurityGroup(sgID)
		if err != nil {
			return nil, err
		}
		rc.Put(reach.ResourceReference{
			Domain: ResourceDomainAWS,
			Kind:   ResourceKindSecurityGroup,
			ID:     sg.ID,
		}, sg.ToResource())

		sgDependencies, err := sg.Dependencies(provider)
		if err != nil {
			return nil, err
		}
		rc.Merge(sgDependencies)
	}

	return rc, nil
}

func (eni ElasticNetworkInterface) getNetworkPoints(parent reach.ResourceReference) []reach.NetworkPoint {
	var networkPoints []reach.NetworkPoint

	lineage := []reach.ResourceReference{
		eni.ToResourceReference(),
		parent,
	}

	for _, privateIPv4Address := range eni.PrivateIPv4Addresses {
		point := reach.NetworkPoint{
			IPAddress: privateIPv4Address,
			Lineage:   lineage,
		}

		networkPoints = append(networkPoints, point)
	}

	if !eni.PublicIPv4Address.Equal(nil) {
		networkPoints = append(networkPoints, reach.NetworkPoint{
			IPAddress: eni.PublicIPv4Address,
			Lineage:   lineage,
		})
	}

	for _, ipv6Address := range eni.IPv6Addresses {
		point := reach.NetworkPoint{
			IPAddress: ipv6Address,
			Lineage:   lineage,
		}

		networkPoints = append(networkPoints, point)
	}

	return networkPoints
}

// Name returns the elastic network interface's ID, and, if available, its name tag value.
func (eni ElasticNetworkInterface) Name() string {
	if name := strings.TrimSpace(eni.NameTag); name != "" {
		return fmt.Sprintf("\"%s\" (%s)", name, eni.ID)
	}
	return eni.ID
}
