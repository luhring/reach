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

func (rt RouteTable) GetDependencies(provider ResourceProvider) (*reach.ResourceCollection, error) {
	rc := reach.NewResourceCollection()

	vpc, err := provider.GetVPC(rt.VPCID)
	if err != nil {
		return nil, err
	}
	rc.Put(reach.ResourceReference{
		Domain: ResourceDomainAWS,
		Kind:   ResourceKindVPC,
		ID:     vpc.ID,
	}, vpc.ToResource())

	// TODO: Figure out dependencies from RouteTableRoute (i.e. route targets)

	return rc, nil
}
