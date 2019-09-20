package aws

import "github.com/luhring/reach/reach"

const ResourceKindSubnet = "Subnet"

type Subnet struct {
	ID    string `json:"id"`
	VPCID string `json:"vpcID"`
}

func (s Subnet) ToResource() reach.Resource {
	return reach.Resource{
		Kind:       ResourceKindSubnet,
		Properties: s,
	}
}

func (s Subnet) GetDependencies(provider ResourceProvider) ([]reach.Resource, error) {
	vpc, err := provider.GetVPC(s.VPCID)
	if err != nil {
		return nil, err
	}

	return []reach.Resource{
		vpc.ToResource(),
	}, nil
}
