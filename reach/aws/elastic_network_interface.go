package aws

import (
	"net"

	"github.com/luhring/reach/reach"
)

const ResourceKindElasticNetworkInterface = "ElasticNetworkInterface"

type ElasticNetworkInterface struct {
	ID                   string   `json:"id"`
	NameTag              string   `json:"nameTag"`
	SubnetID             string   `json:"subnetID"`
	VPCID                string   `json:"vpcID"`
	SecurityGroupIDs     []string `json:"securityGroupIDs"`
	PublicIPv4Address    net.IP   `json:"publicIPv4Address"`
	PrivateIPv4Addresses []net.IP `json:"privateIPv4Addresses"`
	IPv6Addresses        []net.IP `json:"ipv6Addresses"`
}

func (eni ElasticNetworkInterface) ToResource() reach.Resource {
	return reach.Resource{
		Kind:       ResourceKindElasticNetworkInterface,
		Properties: eni,
	}
}

func (eni ElasticNetworkInterface) GetDependencies(provider ResourceProvider) ([]reach.Resource, error) {
	var resources []reach.Resource = nil

	subnet, err := provider.GetSubnet(eni.SubnetID)
	if err != nil {
		return nil, err
	}
	resources = append(resources, subnet.ToResource())

	vpc, err := provider.GetVPC(eni.VPCID)
	if err != nil {
		return nil, err
	}
	resources = append(resources, vpc.ToResource())

	for _, sgID := range eni.SecurityGroupIDs {
		sg, err := provider.GetSecurityGroup(sgID)
		if err != nil {
			return nil, err
		}
		resources = append(resources, sg.ToResource())

		sgDependencies, err := sg.GetDependencies(provider)
		if err != nil {
			return nil, err
		}
		resources = append(resources, sgDependencies...)
	}

	return resources, nil
}
