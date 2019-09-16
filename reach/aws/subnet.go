package aws

const ResourceKindSubnet = "Subnet"

type Subnet struct {
	ID    string `json:"id"`
	VPCID string `json:"vpcID"`
}
