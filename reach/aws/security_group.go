package aws

import (
	"fmt"

	"github.com/luhring/reach/reach"
)

// ResourceKindSecurityGroup specifies the unique name for the security group kind of resource.
const ResourceKindSecurityGroup reach.Kind = "SecurityGroup"

// A SecurityGroup resource representation.
type SecurityGroup struct {
	ID            string
	NameTag       string
	GroupName     string
	VPCID         string
	InboundRules  []SecurityGroupRule
	OutboundRules []SecurityGroupRule
}

// Resource returns the security group converted to a generalized Reach resource.
func (sg SecurityGroup) Resource() reach.Resource {
	return reach.Resource{
		Kind:       ResourceKindSecurityGroup,
		Properties: sg,
	}
}

// ResourceReference returns a resource reference to uniquely identify the security group.
func (sg SecurityGroup) ResourceReference() reach.ResourceReference {
	return reach.ResourceReference{
		Domain: ResourceDomainAWS,
		Kind:   ResourceKindSecurityGroup,
		ID:     sg.ID,
	}
}

// Name returns the security group's ID, and, if available, its name tag value (or group name).
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

func (sg SecurityGroup) rule(direction securityGroupRuleDirection, ruleIndex int) (*SecurityGroupRule, error) {
	errNotFound := fmt.Errorf("rule not found for direction '%s' and index '%d'", direction, ruleIndex)

	var rules []SecurityGroupRule

	switch direction {
	case securityGroupRuleDirectionInbound:
		rules = sg.InboundRules
	case securityGroupRuleDirectionOutbound:
		rules = sg.OutboundRules
	default:
		return nil, errNotFound
	}

	if ruleIndex < 0 || ruleIndex >= len(rules) {
		return nil, errNotFound
	}

	return &rules[ruleIndex], nil
}
