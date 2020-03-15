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

	if err = ensureSingleResult(len(result.RouteTables), "route table", id); err != nil {
		return nil, err
	}

	routeTable, err := newRouteTableFromAPI(result.RouteTables[0])
	if err != nil {
		return nil, err
	}

	return &routeTable, nil
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
