package aws

import "github.com/luhring/reach/reach"

const ResourceKindNetworkACL = "NetworkACL"

type NetworkACL struct {
	ID            string           `json:"id"`
	InboundRules  []NetworkACLRule `json:"inboundRules"`
	OutboundRules []NetworkACLRule `json:"outboundRules"`
}

func (nacl NetworkACL) ToResource() reach.Resource {
	return reach.Resource{
		Kind:       ResourceKindNetworkACL,
		Properties: nacl,
	}
}
