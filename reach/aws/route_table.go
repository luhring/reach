package aws

import "github.com/luhring/reach/reach"

const ResourceKindRouteTable = "RouteTable"

type RouteTable struct {
	ID     string            `json:"id"`
	VPCID  string            `json:"vpcID"`
	Routes []RouteTableRoute `json:"routes"`
}

func (rt RouteTable) ToResource() reach.Resource {
	return reach.Resource{
		Kind:       ResourceKindRouteTable,
		Properties: rt,
	}
}

func (rt RouteTable) GetDependencies(provider ResourceProvider) (map[string]map[string]map[string]reach.Resource, error) {
	resources := make(map[string]map[string]map[string]reach.Resource)

	vpc, err := provider.GetVPC(rt.VPCID)
	if err != nil {
		return nil, err
	}
	resources = reach.EnsureResourcePathExists(resources, ResourceDomainAWS, ResourceKindVPC)
	resources[ResourceDomainAWS][ResourceKindVPC][vpc.ID] = vpc.ToResource()

	// TODO: Figure out dependencies from RouteTableRoute (i.e. route targets)

	return resources, nil
}
