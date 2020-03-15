package aws

import (
	"net"
	"strings"
)

// A RouteTableRoute resource representation.
type RouteTableRoute struct {
	Destination *net.IPNet
	State       RouteState
	Target      RouteTableRouteTarget
}

type RouteTableRouteTarget struct {
	Type RouteTargetType
	ID   string
}

type RouteTargetType int

const (
	RouteTargetTypeInternetGateway RouteTargetType = iota
	RouteTargetTypeNATGateway
	RouteTargetTypeNATInstance
	RouteTargetTypeVirtualPrivateGateway
	RouteTargetTypeLocalGateway // TODO: Look up how this works; reference: https://docs.aws.amazon.com/vpc/latest/userguide/route-table-options.html#route-tables-lgw
	RouteTargetTypeVPCPeeringConnection
	RouteTargetTypeGatewayVPCEndpoint
	RouteTargetTypeEgressOnlyInternetGateway
	RouteTargetTypeTransitGateway
	RouteTargetTypeUnknown
)

func RouteTargetTypeFromPrefix(id string) RouteTargetType {
	prefix := strings.Split(id, "-")[0]

	switch prefix {
	case "igw":
		return RouteTargetTypeInternetGateway
	case "vgw":
		return RouteTargetTypeVirtualPrivateGateway
	case "lgw":
		return RouteTargetTypeLocalGateway
	case "pcx":
		return RouteTargetTypeVPCPeeringConnection
	case "vpce":
		return RouteTargetTypeGatewayVPCEndpoint
	case "eigw":
		return RouteTargetTypeEgressOnlyInternetGateway
	case "tgw":
		return RouteTargetTypeTransitGateway
	default:
		return RouteTargetTypeUnknown
	}
}

type RouteState int

const (
	RouteStateActive RouteState = iota
	RouteStateBlackhole
	RouteStateUnknown
)
