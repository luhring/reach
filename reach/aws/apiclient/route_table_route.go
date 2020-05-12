package apiclient

import (
	"errors"
	"net"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"

	reachAWS "github.com/luhring/reach/reach/aws"
)

func newRouteTableRouteFromAPI(route *ec2.Route) (reachAWS.RouteTableRoute, error) {
	// figure out target type, then put the route together using the struct from aws package

	destination, err := destinationIPNet(route)
	if err != nil {
		return reachAWS.RouteTableRoute{}, err
	}

	state := routeStateFromAPI(aws.StringValue(route.State))

	target := routeTargetFromAPI(route)

	return reachAWS.RouteTableRoute{
		Destination: destination,
		State:       state,
		Target:      target,
	}, nil
}

func destinationIPNet(route *ec2.Route) (*net.IPNet, error) {
	destinationIPv4 := route.DestinationCidrBlock
	destinationIPv6 := route.DestinationIpv6CidrBlock

	var cidr string

	if destinationIPv4 != nil {
		cidr = aws.StringValue(destinationIPv4)
	} else if destinationIPv6 != nil {
		cidr = aws.StringValue(destinationIPv6)
	} else {
		return nil, errors.New("unable to get destination IPNet from AWS route, route had empty IPv4 destination and empty IPv6 destination")
	}

	_, result, err := net.ParseCIDR(cidr)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func routeTableRoutesFromAPI(routes []*ec2.Route) ([]reachAWS.RouteTableRoute, error) {
	var result []reachAWS.RouteTableRoute

	for _, route := range routes {
		r, err := newRouteTableRouteFromAPI(route)
		if err != nil {
			return nil, err
		}

		result = append(result, r)
	}

	return result, nil
}

func routeStateFromAPI(v string) reachAWS.RouteTableRouteState {
	switch v {
	case ec2.RouteStateActive:
		return reachAWS.RouteStateActive
	case ec2.RouteStateBlackhole:
		return reachAWS.RouteStateBlackhole
	default:
		return reachAWS.RouteStateUnknown
	}
}

func routeTargetFromAPI(route *ec2.Route) reachAWS.RouteTableRouteTarget {
	if route.NatGatewayId != nil {
		return reachAWS.RouteTableRouteTarget{
			Type: reachAWS.RouteTargetTypeNATGateway,
			ID:   aws.StringValue(route.NatGatewayId),
		}
	}

	if route.EgressOnlyInternetGatewayId != nil {
		return reachAWS.RouteTableRouteTarget{
			Type: reachAWS.RouteTargetTypeEgressOnlyInternetGateway,
			ID:   aws.StringValue(route.EgressOnlyInternetGatewayId),
		}
	}

	if route.InstanceId != nil { // TODO: Handle InstanceOwnerId, too
		return reachAWS.RouteTableRouteTarget{
			Type: reachAWS.RouteTargetTypeNATInstance,
			ID:   aws.StringValue(route.InstanceId),
		}
	}

	if route.TransitGatewayId != nil {
		return reachAWS.RouteTableRouteTarget{
			Type: reachAWS.RouteTargetTypeTransitGateway,
			ID:   aws.StringValue(route.TransitGatewayId),
		}
	}

	if route.VpcPeeringConnectionId != nil {
		return reachAWS.RouteTableRouteTarget{
			Type: reachAWS.RouteTargetTypeVPCPeeringConnection,
			ID:   aws.StringValue(route.VpcPeeringConnectionId),
		}
	}

	if route.GatewayId != nil {
		id := aws.StringValue(route.GatewayId)

		return reachAWS.RouteTableRouteTarget{
			Type: reachAWS.RouteTargetTypeFromPrefix(id),
			ID:   id,
		}
	}

	return reachAWS.RouteTableRouteTarget{
		Type: reachAWS.RouteTargetTypeUnknown,
		ID:   "unknown",
	}
}
