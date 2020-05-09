package aws

// RouteTableRouteState describes the state of a RouteTableRoute.
type RouteTableRouteState int

// These are the possible values for a RouteTableRouteState.
const (
	RouteStateActive RouteTableRouteState = iota
	RouteStateBlackhole
	RouteStateUnknown
)
