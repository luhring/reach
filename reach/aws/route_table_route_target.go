package aws

import (
	"strings"

	"github.com/luhring/reach/reach"
)

type RouteTableRouteTarget struct {
	Type RouteTargetType
	ID   string
}

func (t RouteTableRouteTarget) Ref() reach.Reference {
	var kind reach.Kind
	id := t.ID

	switch t.Type {
	case RouteTargetTypeInternetGateway:
		kind = ResourceKindInternetGateway
	case RouteTargetTypeNATGateway:
		kind = ResourceKindNATGateway
	default:
		kind = reach.ResourceKindUnknown
		id = "Unknown"
	}

	return reach.Reference{
		Domain: ResourceDomainAWS,
		Kind:   kind,
		ID:     id,
	}
}

type RouteTargetType string

const (
	RouteTargetTypeInternetGateway           RouteTargetType = "InternetGateway"
	RouteTargetTypeNATGateway                RouteTargetType = "NATGateway"
	RouteTargetTypeNATInstance               RouteTargetType = "NATInstance"
	RouteTargetTypeVirtualPrivateGateway     RouteTargetType = "VirtualPrivateGateway"
	RouteTargetTypeLocalGateway              RouteTargetType = "LocalGateway" // TODO: Look up how this works; reference: https://docs.aws.amazon.com/vpc/latest/userguide/route-table-options.html#route-tables-lgw
	RouteTargetTypeVPCPeeringConnection      RouteTargetType = "VPCPeeringConnection"
	RouteTargetTypeGatewayVPCEndpoint        RouteTargetType = "GatewayVPCEndpoint"
	RouteTargetTypeEgressOnlyInternetGateway RouteTargetType = "EgressOnlyInternetGateway"
	RouteTargetTypeTransitGateway            RouteTargetType = "TransitGateway"
	RouteTargetTypeUnknown                   RouteTargetType = "Unknown"
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
