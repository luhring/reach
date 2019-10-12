package aws

import (
	"fmt"
	"log"
	"sort"
	"strings"

	"github.com/luhring/reach/reach"
	"github.com/luhring/reach/reach/helper"
)

const formatResourceMissing = "unable to explain analysis for network point: resource missing from collection: %s"

type Explainer struct {
	analysis reach.Analysis
}

func NewExplainer(analysis reach.Analysis) *Explainer {
	return &Explainer{
		analysis: analysis,
	}
}

func (ex *Explainer) NetworkPoint(point reach.NetworkPoint) string {
	var outputItems []string

	if instanceStateFactor, _ := GetInstanceStateFactor(point.Factors); instanceStateFactor != nil {
		outputItems = append(outputItems, ex.InstanceState(*instanceStateFactor))
	}

	if securityGroupRulesFactor, _ := GetSecurityGroupRulesFactor(point.Factors); securityGroupRulesFactor != nil {
		outputItems = append(outputItems, ex.SecurityGroupRules(*securityGroupRulesFactor))
	}

	return strings.Join(outputItems, "\n")
}

func (ex *Explainer) InstanceState(factor reach.Factor) string {
	var outputItems []string

	ec2instanceRef := ex.analysis.Resources.Get(factor.Resource)
	if ec2instanceRef == nil {
		return fmt.Sprintf(formatResourceMissing, factor.Resource)
	}

	ec2Instance := ec2instanceRef.Properties.(EC2Instance)
	outputItems = append(outputItems, helper.Bold("instance state:"))
	outputItems = append(outputItems, helper.Indent(fmt.Sprintf("\"%s\"", ec2Instance.State), 2))
	outputItems = append(outputItems, "")
	outputItems = append(outputItems, helper.Indent("traffic allowed by instance state:", 2))
	outputItems = append(outputItems, helper.Indent(factor.Traffic.ColorString(), 4))

	return strings.Join(outputItems, "\n")
}

func (ex *Explainer) SecurityGroupRules(factor reach.Factor) string {
	var outputItems []string
	header := fmt.Sprintf("%s (only showing rules that match analysis):", helper.Bold("security group rules"))
	outputItems = append(outputItems, header)

	props := factor.Properties.(SecurityGroupRulesFactor)

	var bodyItems []string

	if rules := props.ComponentRules; len(rules) == 0 {
		bodyItems = append(bodyItems, "no rules that match analysis\n")
	} else {
		var ruleViewModels []RuleExplanationViewModel

		for _, rule := range rules {
			sgRef := ex.analysis.Resources.Get(rule.SecurityGroup)
			if sgRef == nil {
				log.Fatalf(formatResourceMissing, rule.SecurityGroup)
			}

			sg := sgRef.Properties.(SecurityGroup)
			originalRule, err := sg.GetRule(rule.RuleDirection, rule.RuleIndex)
			if err != nil {
				log.Fatalf(err.Error())
			}

			ruleViewModel := RuleExplanationViewModel{
				securityGroupName: sg.Name(),
				ruleMatchText:     fmt.Sprintf("rule #%d (matches %s: %s)", rule.RuleIndex+1, rule.Match.Basis, rule.Match.Value),
				allowedTraffic:    originalRule.TrafficContent.String(),
			}

			ruleViewModels = append(ruleViewModels, ruleViewModel)
		}

		sort.Slice(ruleViewModels, func(i, j int) bool {
			return sort.StringsAreSorted([]string{
				ruleViewModels[i].securityGroupName,
				ruleViewModels[j].securityGroupName,
			})
		})

		var addedSecurityGroupNames []string

		rulesContent := ""

		for _, ruleViewModel := range ruleViewModels {
			securityGroupNameIsNew := true

			for _, addedName := range addedSecurityGroupNames {
				if ruleViewModel.securityGroupName == addedName {
					securityGroupNameIsNew = false
					break
				}
			}

			if securityGroupNameIsNew {
				if rulesContent != "" {
					bodyItems = append(bodyItems, rulesContent)
				}

				bodyItems = append(bodyItems, fmt.Sprintf("%s:", ruleViewModel.securityGroupName))
				addedSecurityGroupNames = append(addedSecurityGroupNames, ruleViewModel.securityGroupName)
				rulesContent = ""
			}

			rulesContent += helper.Indent(ruleViewModel.String(), 2)
		}

		if rulesContent != "" {
			bodyItems = append(bodyItems, rulesContent)
		}
	}

	bodyItems = append(bodyItems, "traffic allowed by security group rules:")
	bodyItems = append(bodyItems, helper.Indent(factor.Traffic.ColorString(), 2))

	body := strings.Join(bodyItems, "\n")
	outputItems = append(outputItems, helper.Indent(body, 2))

	return strings.Join(outputItems, "\n")
}

type SecurityGroupViewModel struct {
	securityGroupName string
	rules             string
}

func (vm SecurityGroupViewModel) String() string {
	securityGroupLine := vm.securityGroupName
	rulesLines := helper.Indent(vm.rules, 2)

	lines := []string{
		securityGroupLine,
		rulesLines,
	}

	return strings.Join(lines, "\n")
}

type RuleExplanationViewModel struct {
	securityGroupName string
	ruleMatchText     string
	allowedTraffic    string
}

func (vm RuleExplanationViewModel) String() string {
	ruleMatchTextLine := fmt.Sprintf("- %s", vm.ruleMatchText)
	allowedTrafficLines := vm.allowedTraffic

	lines := []string{
		ruleMatchTextLine,
		helper.Indent(allowedTrafficLines, 4),
	}

	return strings.Join(lines, "\n")
}

func (ex Explainer) CheckBothInAWS(v reach.NetworkVector) bool {
	return IsUsedByNetworkPoint(v.Source) && IsUsedByNetworkPoint(v.Destination)
}

func (ex Explainer) CheckBothInSameVPC(v reach.NetworkVector) bool {
	sourceENI, destinationENI, err := GetENIsFromVector(v, ex.analysis.Resources)
	if err != nil {
		return false
	}

	return sourceENI.VPCID == destinationENI.VPCID
}

func (ex Explainer) CheckBothInSameSubnet(v reach.NetworkVector) bool {
	sourceENI, destinationENI, err := GetENIsFromVector(v, ex.analysis.Resources)
	if err != nil {
		return false
	}

	return sourceENI.SubnetID == destinationENI.SubnetID
}
