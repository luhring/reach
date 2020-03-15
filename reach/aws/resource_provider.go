package aws

// The ResourceProvider interface wraps all of the necessary methods for accessing AWS-specific resources.
type ResourceProvider interface {
	AllEC2Instances() ([]EC2Instance, error)
	EC2Instance(id string) (*EC2Instance, error)
	ElasticNetworkInterface(id string) (*ElasticNetworkInterface, error)
	InternetGateway(id string) (*InternetGateway, error)
	NATGateway(id string) (*NATGateway, error)
	NetworkACL(id string) (*NetworkACL, error)
	RouteTable(id string) (*RouteTable, error)
	SecurityGroup(id string) (*SecurityGroup, error)
	SecurityGroupReference(id, accountID string) (*SecurityGroupReference, error)
	Subnet(id string) (*Subnet, error)
	VPC(id string) (*VPC, error)
}
