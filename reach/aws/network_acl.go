package aws

import "github.com/luhring/reach/reach"

// ResourceKindNetworkACL specifies the unique name for the network ACL kind of resource.
const ResourceKindNetworkACL reach.Kind = "NetworkACL"

// A NetworkACL resource representation.
type NetworkACL struct {
	ID            string
	InboundRules  []NetworkACLRule
	OutboundRules []NetworkACLRule
}

// Resource returns the network ACL converted to a generalized Reach resource.
func (nacl NetworkACL) Resource() reach.Resource {
	return reach.Resource{
		Kind:       ResourceKindNetworkACL,
		Properties: nacl,
	}
}

func (nacl NetworkACL) Ref() reach.UniversalReference {
	return NetworkACLRef(nacl.ID)
}

func NetworkACLRef(id string) reach.UniversalReference {
	return reach.UniversalReference{
		Domain: ResourceDomainAWS,
		Kind:   ResourceKindNetworkACL,
		ID:     id,
	}
}
