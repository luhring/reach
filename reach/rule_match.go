package reach

type RuleMatch interface {
	explain(observedDescriptor string) Explanation
	getRule() *SecurityGroupRule
}
