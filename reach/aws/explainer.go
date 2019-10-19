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

// Explainer explains an analysis with respect to AWS.
type Explainer struct {
	analysis reach.Analysis
}

// NewExplainer creates a new AWS-specific explainer.
func NewExplainer(analysis reach.Analysis) *Explainer {
	return &Explainer{
		analysis: analysis,
	}
}

// NetworkPoint explains the analysis component for the specified network point.
func (ex *Explainer) NetworkPoint(point reach.NetworkPoint, p reach.Perspective) string {
	var outputItems []string

	if instanceStateFactor, _ := getInstanceStateFactor(point.Factors); instanceStateFactor != nil {
		outputItems = append(outputItems, ex.InstanceState(*instanceStateFactor))
	}

	if securityGroupRulesFactor, _ := getSecurityGroupRulesFactor(point.Factors); securityGroupRulesFactor != nil {
		outputItems = append(outputItems, ex.SecurityGroupRules(*securityGroupRulesFactor, p))
	}

	return strings.Join(outputItems, "\n")
}

// InstanceState explains the analysis component for the specified instance state factor.
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
	outputItems = append(outputItems, helper.Indent("network traffic allowed based on instance state:", 2))
	outputItems = append(outputItems, helper.Indent(factor.Traffic.ColorString(), 4))

	return strings.Join(outputItems, "\n")
}

// SecurityGroupRules explains the analysis component for the specified security group rules factor.
func (ex *Explainer) SecurityGroupRules(factor reach.Factor, p reach.Perspective) string {
	var outputItems []string
	header := fmt.Sprintf(
		"%s (including only rules from %s that match %s):",
		helper.Bold("security group rules"),
		p.SelfRole,
		p.OtherRole,
	)
	outputItems = append(outputItems, header)

	props := factor.Properties.(securityGroupRulesFactor)

	var bodyItems []string

	if rules := props.ComponentRules; len(rules) == 0 {
		bodyItems = append(bodyItems, "no rules that apply to analysis\n")
	} else {
		var ruleViewModels []ruleExplanationViewModel

		for _, rule := range rules {
			sgRef := ex.analysis.Resources.Get(rule.SecurityGroup)
			if sgRef == nil {
				log.Fatalf(formatResourceMissing, rule.SecurityGroup)
			}

			sg := sgRef.Properties.(SecurityGroup)
			originalRule, err := sg.rule(rule.RuleDirection, rule.RuleIndex)
			if err != nil {
				log.Fatalf(err.Error())
			}

			var inclusionReason string

			switch rule.Match.Basis {
			case securityGroupRuleMatchBasisSGRef:
				inclusionReason = fmt.Sprintf(
					"This rule specifies a security group \"%s\" that is attached to the %s's network interface.",
					rule.Match.Value,
					p.OtherRole,
				)
			case securityGroupRuleMatchBasisIP:
				inclusionReason = fmt.Sprintf(
					"This rule specifies an IP CIDR block \"%s\" that contains the %s's IP address \"%s\".",
					originalRule.TargetIPNetworks[0], // TODO: This could show a different network than the matched network, which would be wrong. Include this IPNet in the Match struct to ensure we use the right network here.
					p.OtherRole,
					p.Other.IPAddress,
				)
			default:
				inclusionReason = fmt.Sprintf("Unknown reason for inclusion. Match basis is '%s'. Please report this.", rule.Match.Basis)
			}

			ruleViewModel := ruleExplanationViewModel{
				securityGroupName: sg.Name(),
				inclusionReason:   inclusionReason,
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

	bodyItems = append(bodyItems, "network traffic allowed based on security group rules:")
	bodyItems = append(bodyItems, helper.Indent(factor.Traffic.ColorString(), 2))

	body := strings.Join(bodyItems, "\n")
	outputItems = append(outputItems, helper.Indent(body, 2))

	return strings.Join(outputItems, "\n")
}

// CheckBothInAWS returns a boolean indicating whether both network points in a network vector are AWS resources.
func (ex Explainer) CheckBothInAWS(v reach.NetworkVector) bool {
	return IsUsedByNetworkPoint(v.Source) && IsUsedByNetworkPoint(v.Destination)
}

// CheckBothInSameVPC returns a boolean indicating whether both network points in a network vector reside in the same AWS VPC.
func (ex Explainer) CheckBothInSameVPC(v reach.NetworkVector) bool {
	sourceENI, destinationENI, err := GetENIsFromVector(v, ex.analysis.Resources)
	if err != nil {
		return false
	}

	return sourceENI.VPCID == destinationENI.VPCID
}

// CheckBothInSameSubnet returns a boolean indicating whether both network points in a network vector reside in the same AWS subnet.
func (ex Explainer) CheckBothInSameSubnet(v reach.NetworkVector) bool {
	sourceENI, destinationENI, err := GetENIsFromVector(v, ex.analysis.Resources)
	if err != nil {
		return false
	}

	return sourceENI.SubnetID == destinationENI.SubnetID
}

type ruleExplanationViewModel struct {
	securityGroupName string
	allowedTraffic    string
	inclusionReason   string
}

func (vm ruleExplanationViewModel) String() string {
	output := "- rule\n"

	allowedTrafficHeader := "network traffic allowed:"
	allowedTrafficSection := fmt.Sprintf("%s\n%s", allowedTrafficHeader, helper.Indent(vm.allowedTraffic, 2))
	output += helper.Indent(allowedTrafficSection, 4)

	inclusionReasonHeader := "reason for inclusion:"
	inclusionReasonSection := fmt.Sprintf("%s\n%s\n", inclusionReasonHeader, helper.Indent(vm.inclusionReason, 2))
	output += helper.Indent(inclusionReasonSection, 4)

	return output
}
