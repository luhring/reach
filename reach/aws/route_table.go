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

func (rt RouteTable) GetDependencies(provider ResourceProvider) ([]reach.Resource, error) {
	var resources []reach.Resource = nil

	vpc, err := provider.GetVPC(rt.VPCID)
	if err != nil {
		return nil, err
	}
	resources = append(resources, vpc.ToResource())

	// TODO: Figure out dependencies from RouteTableRoute (i.e. route targets)

	return resources, nil
}
