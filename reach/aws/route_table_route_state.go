package aws

type RouteTableRouteState int

const (
	RouteStateActive RouteTableRouteState = iota
	RouteStateBlackhole
	RouteStateUnknown
)
