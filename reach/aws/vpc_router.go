package aws

import (
	"errors"
	"fmt"
	"net"

	"github.com/luhring/reach/reach"
)

const ResourceKindVPCRouter reach.Kind = "VPCRouter"

type VPCRouter struct {
	VPC VPC
}

func NewVPCRouter(client DomainClient, id string) (*VPCRouter, error) {
	vpc, err := client.VPC(id)
	if err != nil {
		return nil, fmt.Errorf("unable to get VPC: %v", err)
	}

	return &VPCRouter{VPC: *vpc}, nil
}

func (r VPCRouter) Resource() reach.Resource {
	return reach.Resource{
		Kind:       ResourceKindVPCRouter,
		Properties: r,
	}
}

// ———— Implementing Traceable ————

func (r VPCRouter) Ref() reach.UniversalReference {
	return reach.UniversalReference{
		Domain: ResourceDomainAWS,
		Kind:   ResourceKindVPCRouter,
		ID:     r.VPC.ID,
	}
}

func (r VPCRouter) Visitable(_ bool) bool {
	return true
}

func (r VPCRouter) Segments() bool {
	return false
}

func (r VPCRouter) EdgesForward(resolver reach.DomainClientResolver, previousEdge *reach.Edge, previousRef *reach.UniversalReference, _ []net.IP) ([]reach.Edge, error) {
	err := r.checkNilPreviousEdge(previousEdge)
	if err != nil {
		return nil, fmt.Errorf("unablee to generate forward edges: %v", err)
	}

	// VPC Routers don't mutate the IP tuple
	tuple := previousEdge.Tuple

	client, err := unpackDomainClient(resolver)
	if err != nil {
		return nil, fmt.Errorf("unable to get client: %v", err)
	}

	if r.trafficStaysWithinVPC(tuple) {
		eni, err := client.ElasticNetworkInterfaceByIP(tuple.Dst)
		if err != nil {
			return nil, fmt.Errorf("unable to determine forward edges for VPC router: %v", err)
		}
		return r.newEdges(tuple, eni.Ref()), nil
	}

	rt, err := r.routeTable(client, tuple, *previousRef)
	if err != nil {
		return nil, fmt.Errorf("unable to route traffic: %v", err)
	}
	target, err := rt.routeTarget(tuple.Dst)
	if err != nil {
		return nil, fmt.Errorf("unable to determine next edge for traffic (%s): %v", tuple, err)
	}
	return r.newEdges(tuple, target.Ref()), nil
}

func (r VPCRouter) FactorsForward(resolver reach.DomainClientResolver, previousEdge *reach.Edge) ([]reach.Factor, error) {
	tuple := previousEdge.Tuple
	client, err := unpackDomainClient(resolver)
	if err != nil {
		return nil, fmt.Errorf("unable to get client: %v", err)
	}

	srcSubnet, srcSubnetExists, err := r.VPC.subnetThatContains(client, tuple.Src)
	if err != nil {
		return nil, fmt.Errorf("unable to determine if traffic stays within a subnet: %v", err)
	}
	dstSubnet, dstSubnetExists, err := r.VPC.subnetThatContains(client, tuple.Dst)
	if err != nil {
		return nil, fmt.Errorf("unable to determine if traffic stays within a subnet: %v", err)
	}
	if srcSubnetExists && dstSubnetExists && srcSubnet.equal(*dstSubnet) {
		// Same subnet —— no factors to return!
		return nil, nil
	}

	var factors []reach.Factor

	if srcSubnetExists {
		factor, err := r.networkACLRulesFactor(client, *srcSubnet, networkACLRuleDirectionOutbound, tuple)
		if err != nil {
			return nil, fmt.Errorf("unable to get network ACL rules for src subnet: %v", err)
		}
		factors = append(factors, *factor)
	}
	if dstSubnetExists {
		factor, err := r.networkACLRulesFactor(client, *dstSubnet, networkACLRuleDirectionInbound, tuple)
		if err != nil {
			return nil, fmt.Errorf("unable to get network ACL rules for dst subnet: %v", err)
		}
		factors = append(factors, *factor)
	}

	return factors, nil
}

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

func (r VPCRouter) newEdges(tuple reach.IPTuple, ref reach.UniversalReference) []reach.Edge {
	edge := reach.Edge{
		Tuple:             tuple,
		EndRef:            ref,
		ConnectsInterface: false,
	}
	return []reach.Edge{edge}
}

func (r VPCRouter) routeTable(client DomainClient, tuple reach.IPTuple, previousRef reach.UniversalReference) (*RouteTable, error) {
	if r.VPC.contains(tuple.Src) {
		srcSubnet, exists, err := r.VPC.subnetThatContains(client, tuple.Src)
		if err != nil {
			return nil, fmt.Errorf("VPC router cannot find originating subnet for tuple (%s): %v", tuple, err)
		}
		if !exists {
			return nil, fmt.Errorf("unable to find src's subnet (tuple: %s)", tuple)
		}
		subnetRouteTable, err := client.RouteTable(srcSubnet.RouteTableID)
		if err != nil {
			return nil, fmt.Errorf("unable to get routes for traffic (%s) in subnet (%s): %v", tuple, srcSubnet.ID, err)
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
				return nil, fmt.Errorf("could not load Internet Gateway (the previous point) to get route table: %v", err)
			}
			igwRouteTable, err := client.RouteTableForGateway(igw.ID)
			if err != nil {
				return nil, fmt.Errorf("could not load Route Table for Internet Gateway (id: %s): %v", igw.ID, err)
			}
			return igwRouteTable, nil
		}

		return nil, fmt.Errorf("VPC router is unable to find route table for traffic (%s), src infrastructure (%s) is either unrecognized or not yet supported by Reach", tuple, previousRef.Kind)
	}

	//noinspection GoErrorStringFormat
	return nil, fmt.Errorf("Somehow the VPC Router received traffic from infrastructure not within AWS. Please report this as a bug, and include this information... tuple: %s; previousRef: %v", tuple, previousRef)
}

func (r VPCRouter) trafficStaysWithinVPC(tuple reach.IPTuple) bool {
	return r.VPC.contains(tuple.Src) && r.VPC.contains(tuple.Dst)
}
