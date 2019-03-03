package reach

import "fmt"

type SGRefRuleMatch struct {
	Rule  *SecurityGroupRule
	SGRef *SecurityGroupReference
}

func (m *SGRefRuleMatch) Explain(observedDescriptor string) Explanation {
	explanation := newExplanation(
		fmt.Sprintf("security group (%v)", m.SGRef.Name),
	)

	return explanation
}

func (m *SGRefRuleMatch) GetRule() *SecurityGroupRule {
	return m.Rule
}
