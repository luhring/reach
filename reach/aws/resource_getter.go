package aws

type ResourceGetter interface {
	GetAllEC2Instances() ([]EC2Instance, error)
	GetEC2Instance(id string) (*EC2Instance, error)
	GetElasticNetworkInterface(id string) (*ElasticNetworkInterface, error)
	GetNetworkACL(id string) (*NetworkACL, error)
	GetRouteTable(id string) (*RouteTable, error)
	GetSecurityGroup(id string) (*SecurityGroup, error)
	GetSubnet(id string) (*Subnet, error)
	GetVPC(id string) (*VPC, error)
}
