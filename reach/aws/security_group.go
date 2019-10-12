package aws

import (
	"fmt"

	"github.com/luhring/reach/reach"
)

const ResourceKindSecurityGroup = "SecurityGroup"

type SecurityGroup struct {
	ID            string
	NameTag       string
	GroupName     string
	VPCID         string
	InboundRules  []SecurityGroupRule
	OutboundRules []SecurityGroupRule
}

func (sg SecurityGroup) ToResource() reach.Resource {
	return reach.Resource{
		Kind:       ResourceKindSecurityGroup,
		Properties: sg,
	}
}

func (sg SecurityGroup) GetDependencies(provider ResourceProvider) (*reach.ResourceCollection, error) {
	rc := reach.NewResourceCollection()

	vpc, err := provider.GetVPC(sg.VPCID)
	if err != nil {
		return nil, err
	}
	rc.Put(reach.ResourceReference{
		Domain: ResourceDomainAWS,
		Kind:   ResourceKindVPC,
		ID:     vpc.ID,
	}, vpc.ToResource())

	allRules := append(sg.InboundRules, sg.OutboundRules...)

	for _, rule := range allRules {
		// TODO: sg ref IDs shouldn't be strings, they should be pointers, and this check should be for nil not ""

		if sgRefID := rule.TargetSecurityGroupReferenceID; sgRefID != "" {
			sgRef, err := provider.GetSecurityGroupReference(sgRefID, rule.TargetSecurityGroupReferenceAccountID)
			if err != nil {
				return nil, err
			}
			rc.Put(reach.ResourceReference{
				Domain: ResourceDomainAWS,
				Kind:   ResourceKindSecurityGroupReference,
				ID:     sgRef.ID,
			}, sgRef.ToResource())
		}
	}

	return rc, nil
}

func (sg SecurityGroup) Name() string {
	var name string

	if sg.NameTag != "" {
		name = sg.NameTag
	} else if sg.GroupName != "" {
		name = sg.GroupName
	}

	if name != "" {
		return fmt.Sprintf("%s (%s)", name, sg.ID)
	}

	return sg.ID
}

func (sg SecurityGroup) GetRule(direction SecurityGroupRuleDirection, ruleIndex int) (*SecurityGroupRule, error) {
	errNotFound := fmt.Errorf("rule not found for direction '%s' and index '%d'", direction, ruleIndex)

	var rules []SecurityGroupRule

	switch direction {
	case SecurityGroupRuleDirectionInbound:
		rules = sg.InboundRules
	case SecurityGroupRuleDirectionOutbound:
		rules = sg.OutboundRules
	default:
		return nil, errNotFound
	}

	if ruleIndex < 0 || ruleIndex >= len(rules) {
		return nil, errNotFound
	}

	return &rules[ruleIndex], nil
}
