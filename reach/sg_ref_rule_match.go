package reach

type SGRefRuleMatch struct {
	Rule  *SecurityGroupRule
	SGRef *SecurityGroupReference
}

func (m *SGRefRuleMatch) Explain(observedDescriptor string) Explanation {
	var explanation Explanation

	explanation.AddLineFormat("security group (%v)", m.SGRef.Name)

	return explanation
}

func (m *SGRefRuleMatch) GetRule() *SecurityGroupRule {
	return m.Rule
}
