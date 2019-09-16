package aws

const ResourceKindRouteTable = "RouteTable"

type RouteTable struct {
	ID     string            `json:"id"`
	VPCID  string            `json:"vpcID"`
	Routes []RouteTableRoute `json:"routes"`
}
