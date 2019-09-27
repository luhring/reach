package aws

import (
	"net"

	"github.com/luhring/reach/reach"
)

const ResourceKindElasticNetworkInterface = "ElasticNetworkInterface"

type ElasticNetworkInterface struct {
	ID                   string
	NameTag              string
	SubnetID             string
	VPCID                string
	SecurityGroupIDs     []string
	PublicIPv4Address    net.IP   `json:"PublicIPv4Address,omitempty"`
	PrivateIPv4Addresses []net.IP `json:"PrivateIPv4Addresses,omitempty"`
	IPv6Addresses        []net.IP `json:"IPv6Addresses,omitempty"`
}

func (eni ElasticNetworkInterface) ToResource() reach.Resource {
	return reach.Resource{
		Kind:       ResourceKindElasticNetworkInterface,
		Properties: eni,
	}
}

func (eni ElasticNetworkInterface) ToResourceReference() reach.ResourceReference {
	return reach.ResourceReference{
		Domain: ResourceDomainAWS,
		Kind:   ResourceKindElasticNetworkInterface,
		ID:     eni.ID,
	}
}

func (eni ElasticNetworkInterface) GetDependencies(provider ResourceProvider) (*reach.ResourceCollection, error) {
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

		sgDependencies, err := sg.GetDependencies(provider)
		if err != nil {
			return nil, err
		}
		rc.Merge(sgDependencies)
	}

	return rc, nil
}

func (eni ElasticNetworkInterface) GetNetworkPoints(parent reach.ResourceReference) []reach.NetworkPoint {
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
