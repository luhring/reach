package aws

import "github.com/luhring/reach/reach"

const ResourceKindSubnet = "Subnet"

type Subnet struct {
	ID    string `json:"id"`
	VPCID string `json:"vpcID"`
}

func (s Subnet) ToResource() reach.Resource {
	return reach.Resource{
		Kind:       ResourceKindSubnet,
		Properties: s,
	}
}

func (s Subnet) GetDependencies(provider ResourceProvider) (map[string]map[string]map[string]reach.Resource, error) {
	resources := make(map[string]map[string]map[string]reach.Resource)

	vpc, err := provider.GetVPC(s.VPCID)
	if err != nil {
		return nil, err
	}
	resources = reach.EnsureResourcePathExists(resources, ResourceDomainAWS, ResourceKindVPC)
	resources[ResourceDomainAWS][ResourceKindVPC][vpc.ID] = vpc.ToResource()

	return resources, nil
}
