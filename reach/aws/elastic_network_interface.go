package aws

import (
	"net"

	"github.com/luhring/reach/reach"
)

const ResourceKindElasticNetworkInterface = "ElasticNetworkInterface"

type ElasticNetworkInterface struct {
	ID                   string   `json:"id"`
	NameTag              string   `json:"nameTag"`
	SubnetID             string   `json:"subnetID"`
	VPCID                string   `json:"vpcID"`
	SecurityGroupIDs     []string `json:"securityGroupIDs"`
	PublicIPv4Address    net.IP   `json:"publicIPv4Address"`
	PrivateIPv4Addresses []net.IP `json:"privateIPv4Addresses"`
	IPv6Addresses        []net.IP `json:"ipv6Addresses"`
}

func (eni ElasticNetworkInterface) ToResource() reach.Resource {
	return reach.Resource{
		Kind:       ResourceKindElasticNetworkInterface,
		Properties: eni,
	}
}

func (eni ElasticNetworkInterface) GetDependencies(provider ResourceProvider) (map[string]map[string]map[string]reach.Resource, error) {
	resources := make(map[string]map[string]map[string]reach.Resource)

	subnet, err := provider.GetSubnet(eni.SubnetID)
	if err != nil {
		return nil, err
	}
	resources = reach.EnsureResourcePathExists(resources, ResourceDomainAWS, ResourceKindSubnet)
	resources[ResourceDomainAWS][ResourceKindSubnet][subnet.ID] = subnet.ToResource()

	vpc, err := provider.GetVPC(eni.VPCID)
	if err != nil {
		return nil, err
	}
	resources = reach.EnsureResourcePathExists(resources, ResourceDomainAWS, ResourceKindVPC)
	resources[ResourceDomainAWS][ResourceKindVPC][vpc.ID] = vpc.ToResource()

	for _, sgID := range eni.SecurityGroupIDs {
		sg, err := provider.GetSecurityGroup(sgID)
		if err != nil {
			return nil, err
		}
		resources = reach.EnsureResourcePathExists(resources, ResourceDomainAWS, ResourceKindSecurityGroup)
		resources[ResourceDomainAWS][ResourceKindSecurityGroup][sg.ID] = sg.ToResource()

		sgDependencies, err := sg.GetDependencies(provider)
		if err != nil {
			return nil, err
		}
		resources = reach.MergeResources(resources, sgDependencies)
	}

	return resources, nil
}
