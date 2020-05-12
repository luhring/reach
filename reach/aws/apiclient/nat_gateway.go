package apiclient

import (
	"net"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"

	reachAWS "github.com/luhring/reach/reach/aws"
)

// NATGateway queries the AWS API for a NAT gateway matching the given ID.
func (client *DomainClient) NATGateway(id string) (*reachAWS.NATGateway, error) {
	if r := client.cachedResource(reachAWS.NATGatewayRef(id)); r != nil {
		if v, ok := r.(*reachAWS.NATGateway); ok {
			return v, nil
		}
	}

	input := &ec2.DescribeNatGatewaysInput{
		NatGatewayIds: []*string{aws.String(id)},
	}
	result, err := client.ec2.DescribeNatGateways(input)
	if err != nil {
		return nil, err
	}

	if err = ensureSingleResult(len(result.NatGateways), reachAWS.ResourceKindNATGateway, id); err != nil {
		return nil, err
	}

	natGateway := result.NatGateways[0]

	ngw := reachAWS.NATGateway{
		ID:        id,
		SubnetID:  aws.StringValue(natGateway.SubnetId),
		VPCID:     aws.StringValue(natGateway.VpcId),
		PrivateIP: privateIPForNATGateway(natGateway),
		PublicIP:  publicIPForNATGateway(natGateway),
	}
	client.cacheResource(ngw)
	return &ngw, nil
}

func privateIPForNATGateway(ngw *ec2.NatGateway) net.IP {
	input := ngw.NatGatewayAddresses[0].PrivateIp

	return net.ParseIP(aws.StringValue(input))
}

func publicIPForNATGateway(ngw *ec2.NatGateway) net.IP {
	input := ngw.NatGatewayAddresses[0].PublicIp

	return net.ParseIP(aws.StringValue(input))
}
