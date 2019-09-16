package aws

const ResourceKindNetworkACL = "NetworkACL"

type NetworkACL struct {
	ID            string           `json:"id"`
	InboundRules  []NetworkACLRule `json:"inboundRules"`
	OutboundRules []NetworkACLRule `json:"outboundRules"`
}
