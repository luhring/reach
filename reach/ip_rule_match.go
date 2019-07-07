package reach

import (
	"fmt"
	"github.com/mgutz/ansi"
	"net"
)

type IPRuleMatch struct {
	Rule             *SecurityGroupRule
	MatchedIPRange   *net.IPNet
	TargetIP         net.IP
	IsTargetIPPublic bool
}

func (m *IPRuleMatch) explain(observedDescriptor string) Explanation {
	var publicOrPrivate string
	if m.IsTargetIPPublic {
		publicOrPrivate = "public"
	} else {
		publicOrPrivate = "private"
	}

	explanation := newExplanation(fmt.Sprintf(
		ansi.Color("- rule: allow %v", "green"),
		ansi.Color(m.Rule.TrafficAllowance.String(), "green+b"),
	))

	explanation.addLineFormatWithIndents(
		1,
		"(This rule handles an IP address range '%v' that includes the %s network interface's %s IP address '%v'.)",
		m.MatchedIPRange.String(),
		observedDescriptor,
		publicOrPrivate,
		m.TargetIP.String(),
	)

	return explanation
}

func (m *IPRuleMatch) getRule() *SecurityGroupRule {
	return m.Rule
}
