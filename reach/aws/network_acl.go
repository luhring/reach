package aws

import "github.com/luhring/reach/reach"

const ResourceKindNetworkACL = "NetworkACL"

type NetworkACL struct {
	ID            string
	InboundRules  []NetworkACLRule
	OutboundRules []NetworkACLRule
}

func (nacl NetworkACL) ToResource() reach.Resource {
	return reach.Resource{
		Kind:       ResourceKindNetworkACL,
		Properties: nacl,
	}
}
