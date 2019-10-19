package aws

import "github.com/luhring/reach/reach"

// ResourceKindSubnet specifies the unique name for the subnet kind of resource.
const ResourceKindSubnet = "Subnet"

// A Subnet resource representation.
type Subnet struct {
	ID    string
	VPCID string
}

// ToResource returns the subnet converted to a generalized Reach resource.
func (s Subnet) ToResource() reach.Resource {
	return reach.Resource{
		Kind:       ResourceKindSubnet,
		Properties: s,
	}
}

// Dependencies returns a collection of the subnet's resource dependencies.
func (s Subnet) Dependencies(provider ResourceProvider) (*reach.ResourceCollection, error) {
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
