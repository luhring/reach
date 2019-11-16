package aws

import (
	"fmt"

	"github.com/luhring/reach/reach/helper"
)

type securityGroupRuleExplanationViewModel struct {
	securityGroupName string
	allowedTraffic    string
	inclusionReason   string
}

func (model securityGroupRuleExplanationViewModel) String() string {
	output := "- rule\n"

	allowedTrafficHeader := "network traffic allowed:"
	allowedTrafficSection := fmt.Sprintf("%s\n%s", allowedTrafficHeader, helper.Indent(model.allowedTraffic, 2))
	output += helper.Indent(allowedTrafficSection, 4)

	inclusionReasonHeader := "reason for inclusion:"
	inclusionReasonSection := fmt.Sprintf("%s\n%s\n", inclusionReasonHeader, helper.Indent(model.inclusionReason, 2))
	output += helper.Indent(inclusionReasonSection, 4)

	return output
}
