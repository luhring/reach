package aws

import "net"

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
