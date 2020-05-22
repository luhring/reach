package aws

// NetworkACLRuleDirection specifies the direction of traffic to which a network ACL rule applies
type NetworkACLRuleDirection string

// Possible values for the NetworkACLRuleDirection type
const (
	NetworkACLRuleDirectionInbound  NetworkACLRuleDirection = "inbound"
	NetworkACLRuleDirectionOutbound NetworkACLRuleDirection = "outbound"
)
