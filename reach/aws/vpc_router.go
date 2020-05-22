package aws

import (
	"errors"
	"fmt"
	"net"

	"github.com/luhring/reach/reach"
	"github.com/luhring/reach/reach/reacherr"
)

// ResourceKindVPCRouter specifies the unique name for the VPCRouter kind of resource.
const ResourceKindVPCRouter reach.Kind = "VPCRouter"

var _ reach.Traceable = (*VPCRouter)(nil)

// VPCRouter represents the router implicitly present within each AWS VPC.
type VPCRouter struct {
	VPC VPC
}

// NewVPCRouter returns a pointer to a new instance of a VPCRouter. NewVPCRouter returns an error if it is unable to gather enough information required to represent a VPC's router.
func NewVPCRouter(client DomainClient, id string) (*VPCRouter, error) {
	vpc, err := client.VPC(id)
	if err != nil {
		return nil, err
	}

	return &VPCRouter{VPC: *vpc}, nil
}

// VPCRouterRef returns a Reference for a VPCRouter with the specified ID.
func VPCRouterRef(id string) reach.Reference {
	return reach.Reference{
		Domain: ResourceDomainAWS,
		Kind:   ResourceKindVPCRouter,
		ID:     id,
	}
}

// Resource returns the VPCRouter converted to a generalized Reach resource.
func (r VPCRouter) Resource() reach.Resource {
	return reach.Resource{
		Kind:       ResourceKindVPCRouter,
		Properties: r,
	}
}

// ———— Implementing Traceable ————

// Ref returns a Reference for the VPCRouter.
func (r VPCRouter) Ref() reach.Reference {
	return VPCRouterRef(r.VPC.ID)
}

// Visitable returns a boolean to indicate whether a tracer is allowed to add this resource to the path it's currently constructing.
//
// The Visitable method for VPCRouter always returns true because there is no limit to the number of times a tracer can visit a VPCRouter.
func (r VPCRouter) Visitable(_ bool) bool {
	return true
}

// Segments returns a boolean to indicate whether a tracer should create a new path segment at this point in the path.
//
// The Segments method for VPCRouter always returns false because VPC routers never perform NAT.
func (r VPCRouter) Segments() bool {
	return false
}

// EdgesForward returns the set of all possible edges forward given this point in a path that a tracer is constructing. EdgesForward returns an empty slice of edges if there are no further points for the specified network traffic to travel as it attempts to reach its intended network destination.
func (r VPCRouter) EdgesForward(resolver reach.DomainClientResolver, previousEdge *reach.Edge, previousRef *reach.Reference, _ []net.IP) ([]reach.Edge, error) {
	err := r.checkNilPreviousEdge(previousEdge)
	if err != nil {
		return nil, err
	}

	// VPC Routers don't mutate the IP tuple
	tuple := previousEdge.Tuple

	client, err := unpackDomainClient(resolver)
	if err != nil {
		return nil, err
	}

	if r.trafficStaysWithinVPC(tuple) {
		eni, err := client.ElasticNetworkInterfaceByIP(tuple.Dst)
		if err != nil {
			return nil, err
		}
		return r.newEdges(tuple, eni.Ref()), nil
	}

	rt, err := r.routeTable(client, tuple, *previousRef)
	if err != nil {
		return nil, err
	}
	target, err := rt.routeTarget(tuple.Dst)
	if err != nil {
		return nil, fmt.Errorf("unable to determine next edge for traffic (%s): %v", tuple, err)
	}
	return r.newEdges(tuple, target.Ref()), nil
}

// FactorsForward returns a set of factors that impact the traffic traveling through this point in the direction of source to destination.
func (r VPCRouter) FactorsForward(resolver reach.DomainClientResolver, previousEdge *reach.Edge) ([]reach.Factor, error) {
	tuple := previousEdge.Tuple
	client, err := unpackDomainClient(resolver)
	if err != nil {
		return nil, err
	}

	srcSubnet, srcSubnetExists, err := r.VPC.subnetThatContains(client, tuple.Src)
	if err != nil {
		return nil, err
	}
	dstSubnet, dstSubnetExists, err := r.VPC.subnetThatContains(client, tuple.Dst)
	if err != nil {
		return nil, err
	}
	if srcSubnetExists && dstSubnetExists && srcSubnet.equal(*dstSubnet) {
		// Same subnet —— no factors to return!
		return nil, nil
	}

	var factors []reach.Factor

	if srcSubnetExists {
		factor, err := r.networkACLRulesFactor(client, *srcSubnet, networkACLRuleDirectionOutbound, tuple)
		if err != nil {
			return nil, err
		}
		factors = append(factors, *factor)
	}
	if dstSubnetExists {
		factor, err := r.networkACLRulesFactor(client, *dstSubnet, networkACLRuleDirectionInbound, tuple)
		if err != nil {
			return nil, err
		}
		factors = append(factors, *factor)
	}

	return factors, nil
}

// FactorsReturn returns a set of factors that impact the traffic traveling through this point in the direction of destination to source.
func (r VPCRouter) FactorsReturn(resolver reach.DomainClientResolver, nextEdge *reach.Edge) ([]reach.Factor, error) {
	panic("implement me")
}

// ———— Supporting methods ————

func (r VPCRouter) checkNilPreviousEdge(previousEdge *reach.Edge) error {
	if previousEdge == nil {
		return errors.New("reach does not support a VPC router being the first point in a path")
	}
	return nil
}

func (r VPCRouter) newEdges(tuple reach.IPTuple, ref reach.Reference) []reach.Edge {
	edge := reach.Edge{
		Tuple:             tuple,
		EndRef:            ref,
		ConnectsInterface: false,
	}
	return []reach.Edge{edge}
}

func (r VPCRouter) routeTable(client DomainClient, tuple reach.IPTuple, previousRef reach.Reference) (*RouteTable, error) {
	if r.VPC.contains(tuple.Src) {
		srcSubnet, exists, err := r.VPC.subnetThatContains(client, tuple.Src)
		if err != nil {
			return nil, err
		}
		if !exists {
			return nil, fmt.Errorf("unable to find src's subnet (tuple: %s)", tuple)
		}
		subnetRouteTable, err := client.RouteTable(srcSubnet.RouteTableID)
		if err != nil {
			return nil, err
		}
		return subnetRouteTable, nil
	}

	// Traffic originates from outside of VPC.
	// Check edge association (see https://docs.aws.amazon.com/vpc/latest/userguide/VPC_Route_Tables.html#RouteTables)
	// And then gateway route table (see https://docs.aws.amazon.com/vpc/latest/userguide/VPC_Route_Tables.html#gateway-route-table)

	// Figure out where the traffic came from. (IGW or VGW —— If neither, we have a problem.)
	// Get gateway route table for that gateway, and use it to determine next ref.

	if previousRef.Domain == ResourceDomainAWS {
		switch previousRef.Kind {
		case ResourceKindInternetGateway:
			igw, err := client.InternetGateway(previousRef.ID)
			if err != nil {
				return nil, err
			}
			igwRouteTable, err := client.RouteTableForGateway(igw.ID)
			if err != nil {
				return nil, err
			}
			return igwRouteTable, nil
		}

		return nil, reacherr.New(nil, "VPC router is unable to find route table for traffic (%s), src infrastructure (%s) is either unrecognized or not yet supported by Reach", tuple, previousRef.Kind)
	}

	//noinspection GoErrorStringFormat
	return nil, fmt.Errorf("Somehow the VPC Router received traffic from infrastructure not within AWS. tuple: %s; previousRef: %v", tuple, previousRef)
}

func (r VPCRouter) trafficStaysWithinVPC(tuple reach.IPTuple) bool {
	return r.VPC.contains(tuple.Src) && r.VPC.contains(tuple.Dst)
}
