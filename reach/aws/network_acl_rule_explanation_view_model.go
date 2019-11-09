package aws

import (
	"fmt"

	"github.com/luhring/reach/reach"
	"github.com/luhring/reach/reach/helper"
)

type networkACLRuleExplanationViewModel struct {
	ruleNumber      int64
	allowedTraffic  string
	inclusionReason string
}

func newNetworkACLRuleExplanationViewModel(rule networkACLRulesFactorComponent, p reach.Perspective) networkACLRuleExplanationViewModel {
	inclusionReason := fmt.Sprintf(
		"This rule specifies an IP CIDR block \"%s\" that contains the %s's IP address (%s).",
		rule.Match.Requirement.String(),
		p.OtherRole,
		p.Other.IPAddress,
	)

	return networkACLRuleExplanationViewModel{
		ruleNumber:      rule.RuleNumber,
		allowedTraffic:  rule.Traffic.String(),
		inclusionReason: inclusionReason,
	}
}

func (model networkACLRuleExplanationViewModel) String() string {
	output := fmt.Sprintf("- rule # %d\n", model.ruleNumber)

	allowedTrafficHeader := "network traffic allowed:"
	allowedTrafficSection := fmt.Sprintf("%s\n%s", allowedTrafficHeader, helper.Indent(model.allowedTraffic, 2))
	output += helper.Indent(allowedTrafficSection, 4)

	inclusionReasonHeader := "reason for inclusion:"
	inclusionReasonSection := fmt.Sprintf("%s\n%s\n", inclusionReasonHeader, helper.Indent(model.inclusionReason, 2))
	output += helper.Indent(inclusionReasonSection, 4)

	return output
}

func networkACLRuleComponentsToViewModels(rules []networkACLRulesFactorComponent, p reach.Perspective) []networkACLRuleExplanationViewModel {
	var models []networkACLRuleExplanationViewModel

	for _, rule := range rules {
		model := newNetworkACLRuleExplanationViewModel(rule, p)
		models = append(models, model)
	}

	return models
}
