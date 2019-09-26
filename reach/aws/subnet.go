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

func (s Subnet) GetDependencies(provider ResourceProvider) (*reach.ResourceCollection, error) {
	rc := reach.NewResourceCollection()

	vpc, err := provider.GetVPC(s.VPCID)
	if err != nil {
		return nil, err
	}
	rc.Put(reach.ResourceReference{
		Domain: ResourceDomainAWS,
		Kind:   ResourceKindVPC,
		ID:     vpc.ID,
	}, vpc.ToResource())

	return rc, nil
}
