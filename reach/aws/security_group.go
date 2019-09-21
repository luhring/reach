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

func (sg SecurityGroup) GetDependencies(provider ResourceProvider) (map[string]map[string]map[string]reach.Resource, error) {
	resources := make(map[string]map[string]map[string]reach.Resource)

	vpc, err := provider.GetVPC(sg.VPCID)
	if err != nil {
		return nil, err
	}
	resources = reach.EnsureResourcePathExists(resources, ResourceDomainAWS, ResourceKindVPC)
	resources[ResourceDomainAWS][ResourceKindVPC][vpc.ID] = vpc.ToResource()

	allRules := append(sg.InboundRules, sg.OutboundRules...)

	for _, rule := range allRules {
		// TODO: sg ref IDs shouldn't be strings, they should be pointers, and this check should be for nil not ""

		if sgRefID := rule.TargetSecurityGroupReferenceID; sgRefID != "" {
			sgRef, err := provider.GetSecurityGroupReference(sgRefID, rule.TargetSecurityGroupReferenceAccountID)
			if err != nil {
				return nil, err
			}
			resources = reach.EnsureResourcePathExists(resources, ResourceDomainAWS, ResourceKindSecurityGroupReference)
			resources[ResourceDomainAWS][ResourceKindSecurityGroupReference][sgRef.ID] = sgRef.ToResource()
		}
	}

	return resources, nil
}
