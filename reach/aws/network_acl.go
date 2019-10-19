package aws

import "github.com/luhring/reach/reach"

// ResourceKindNetworkACL specifies the unique name for the network ACL kind of resource.
const ResourceKindNetworkACL = "NetworkACL"

// A NetworkACL resource representation.
type NetworkACL struct {
	ID            string
	InboundRules  []NetworkACLRule
	OutboundRules []NetworkACLRule
}

// ToResource returns the network ACL converted to a generalized Reach resource.
func (nacl NetworkACL) ToResource() reach.Resource {
	return reach.Resource{
		Kind:       ResourceKindNetworkACL,
		Properties: nacl,
	}
}
