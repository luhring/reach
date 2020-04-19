package aws

import "github.com/luhring/reach/reach"

// ResourceKindInternetGateway specifies the unique name for the Internet gateway kind of resource.
const ResourceKindInternetGateway = "InternetGateway"

// An InternetGateway resource representation.
type InternetGateway struct {
	ID    string
	VPCID string
}

// ToResource returns the Internet gateway converted to a generalized Reach resource.
func (igw InternetGateway) ToResource() reach.Resource {
	return reach.Resource{
		Kind:       ResourceKindInternetGateway,
		Properties: igw,
	}
}

// Dependencies returns a collection of the Internet gateway's resource dependencies.
func (igw InternetGateway) Dependencies(provider ResourceGetter) (*reach.ResourceCollection, error) {
	rc := reach.NewResourceCollection()

	vpc, err := provider.VPC(igw.VPCID)
	if err != nil {
		return nil, err
	}
	rc.Put(reach.ResourceReference{
		Domain: ResourceDomainAWS,
		Kind:   ResourceKindVPC,
		ID:     igw.VPCID,
	}, vpc.ToResource())

	return rc, nil
}
