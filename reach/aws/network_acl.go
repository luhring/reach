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

// NetworkACLRef returns a Reference for a NetworkACL with the specified ID.
func NetworkACLRef(id string) reach.Reference {
	return reach.Reference{
		Domain: ResourceDomainAWS,
		Kind:   ResourceKindNetworkACL,
		ID:     id,
	}
}

// Resource returns the NetworkACL converted to a generalized Reach resource.
func (nacl NetworkACL) Resource() reach.Resource {
	return reach.Resource{
		Kind:       ResourceKindNetworkACL,
		Properties: nacl,
	}
}

// Ref returns a Reference for the NetworkACL.
func (nacl NetworkACL) Ref() reach.Reference {
	return NetworkACLRef(nacl.ID)
}
