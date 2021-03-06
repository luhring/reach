package api

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"

	reachAWS "github.com/luhring/reach/reach/aws"
)

// RouteTable queries the AWS API for a route table matching the given ID.
func (provider *ResourceProvider) RouteTable(id string) (*reachAWS.RouteTable, error) {
	input := &ec2.DescribeRouteTablesInput{
		RouteTableIds: []*string{
			aws.String(id),
		},
	}
	result, err := provider.ec2.DescribeRouteTables(input)
	if err != nil {
		return nil, err
	}

	if err = ensureSingleResult(len(result.RouteTables), "security group", id); err != nil {
		return nil, err
	}

	routeTable := newRouteTableFromAPI(result.RouteTables[0])
	return &routeTable, nil
}

func newRouteTableFromAPI(routeTable *ec2.RouteTable) reachAWS.RouteTable {
	routes := []reachAWS.RouteTableRoute{} // TODO: implement

	return reachAWS.RouteTable{
		ID:     aws.StringValue(routeTable.RouteTableId),
		VPCID:  aws.StringValue(routeTable.VpcId),
		Routes: routes,
	}
}

func routeTableRoutes(routes []*ec2.Route) []reachAWS.RouteTableRoute {
	return nil // TODO: implement
}

func routeTableRoute(route *ec2.Route) reachAWS.RouteTableRoute {
	if route == nil {
		return reachAWS.RouteTableRoute{}
	}

	panic("need to finish implementing")
}
