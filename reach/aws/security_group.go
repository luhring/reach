package aws

const ResourceKindSecurityGroup = "SecurityGroup"

type SecurityGroup struct {
	ID            string              `json:"id"`
	NameTag       string              `json:"nameTag"`
	GroupName     string              `json:"groupName"`
	VPCID         string              `json:"vpcID"`
	InboundRules  []SecurityGroupRule `json:"inboundRules"`
	OutboundRules []SecurityGroupRule `json:"outboundRules"`
}
