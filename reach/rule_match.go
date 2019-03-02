package reach

type RuleMatch interface {
	Explain(observedDescriptor string) Explanation
	GetRule() *SecurityGroupRule
}
