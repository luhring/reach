package aws

import "github.com/luhring/reach/reach"

// ResourceKindRouteTable specifies the unique name for the route table kind of resource.
const ResourceKindRouteTable = "RouteTable"

// A RouteTable resource representation.
type RouteTable struct {
	ID     string
	VPCID  string
	Routes []RouteTableRoute
}

// ToResource returns the route table converted to a generalized Reach resource.
func (rt RouteTable) ToResource() reach.Resource {
	return reach.Resource{
		Kind:       ResourceKindRouteTable,
		Properties: rt,
	}
}

// Dependencies returns a collection of the route table's resource dependencies.
func (rt RouteTable) Dependencies(provider ResourceProvider) (*reach.ResourceCollection, error) {
	rc := reach.NewResourceCollection()

	vpc, err := provider.VPC(rt.VPCID)
	if err != nil {
		return nil, err
	}
	rc.Put(reach.ResourceReference{
		Domain: ResourceDomainAWS,
		Kind:   ResourceKindVPC,
		ID:     vpc.ID,
	}, vpc.ToResource())

	for _, route := range rt.Routes {
		switch route.Target.Type {
		case RouteTargetTypeInternetGateway:
			igwID := route.Target.ID

			igw, err := provider.InternetGateway(igwID)
			if err != nil {
				return nil, err
			}
			rc.Put(reach.ResourceReference{
				Domain: ResourceDomainAWS,
				Kind:   ResourceKindInternetGateway,
				ID:     igwID,
			}, igw.ToResource())
		case RouteTargetTypeNATGateway:
			ngwID := route.Target.ID

			ngw, err := provider.NATGateway(ngwID)
			if err != nil {
				return nil, err
			}
			rc.Put(reach.ResourceReference{
				Domain: ResourceDomainAWS,
				Kind:   ResourceKindNATGateway,
				ID:     ngwID,
			}, ngw.ToResource())
		default:
			// target type that's not yet supported, ignore for now
		}
	}

	return rc, nil
}
