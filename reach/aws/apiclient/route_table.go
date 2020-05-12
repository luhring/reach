package apiclient

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"

	reachAWS "github.com/luhring/reach/reach/aws"
)

// RouteTable queries the AWS API for a route table matching the given ID.
func (client *DomainClient) RouteTable(id string) (*reachAWS.RouteTable, error) {
	if r := client.cachedResource(reachAWS.RouteTableRef(id)); r != nil {
		if v, ok := r.(*reachAWS.RouteTable); ok {
			return v, nil
		}
	}

	input := &ec2.DescribeRouteTablesInput{
		RouteTableIds: []*string{
			aws.String(id),
		},
	}
	result, err := client.ec2.DescribeRouteTables(input)
	if err != nil {
		return nil, err
	}

	if err = ensureSingleResult(len(result.RouteTables), "route table", id); err != nil {
		return nil, err
	}

	routeTable, err := newRouteTableFromAPI(result.RouteTables[0])
	if err != nil {
		return nil, err
	}
	client.cacheResource(routeTable)
	return &routeTable, nil
}

// RouteTableForGateway returns the route table for the specified gateway.
func (client *DomainClient) RouteTableForGateway(id string) (*reachAWS.RouteTable, error) {
	panic("implement me")
}

func newRouteTableFromAPI(routeTable *ec2.RouteTable) (reachAWS.RouteTable, error) {
	routes, err := routeTableRoutesFromAPI(routeTable.Routes)
	if err != nil {
		return reachAWS.RouteTable{}, err
	}

	return reachAWS.RouteTable{
		ID:     aws.StringValue(routeTable.RouteTableId),
		VPCID:  aws.StringValue(routeTable.VpcId),
		Routes: routes,
	}, nil
}
