package aws

import "github.com/luhring/reach/reach"

// ResourceKindInternetGateway specifies the unique name for the Internet gateway kind of resource.
const ResourceKindInternetGateway reach.Kind = "InternetGateway"

// An InternetGateway resource representation.
type InternetGateway struct {
	ID    string
	VPCID string
}

// Resource returns the Internet gateway converted to a generalized Reach resource.
func (igw InternetGateway) Resource() reach.Resource {
	return reach.Resource{
		Kind:       ResourceKindInternetGateway,
		Properties: igw,
	}
}
