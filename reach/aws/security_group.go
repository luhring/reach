package aws

import "github.com/luhring/reach/reach"

const ResourceKindSecurityGroup = "SecurityGroup"

type SecurityGroup struct {
	ID            string              `json:"id"`
	NameTag       string              `json:"nameTag"`
	GroupName     string              `json:"groupName"`
	VPCID         string              `json:"vpcID"`
	InboundRules  []SecurityGroupRule `json:"inboundRules"`
	OutboundRules []SecurityGroupRule `json:"outboundRules"`
}

func (sg SecurityGroup) ToResource() reach.Resource {
	return reach.Resource{
		Kind:       ResourceKindSecurityGroup,
		Properties: sg,
	}
}

func (sg SecurityGroup) GetDependencies(provider ResourceProvider) ([]reach.Resource, error) {
	var resources []reach.Resource = nil

	vpc, err := provider.GetVPC(sg.VPCID)
	if err != nil {
		return nil, err
	}
	resources = append(resources, vpc.ToResource())

	allRules := append(sg.InboundRules, sg.OutboundRules...)

	for _, rule := range allRules {
		// TODO: sg ref IDs shouldn't be strings, they should be pointers, and this check should be for nil not ""

		if sgRefID := rule.TargetSecurityGroupReferenceID; sgRefID != "" {
			sgRef, err := provider.GetSecurityGroupReference(sgRefID, rule.TargetSecurityGroupReferenceAccountID)
			if err != nil {
				return nil, err
			}
			resources = append(resources, sgRef.ToResource())
		}
	}

	return resources, nil
}
