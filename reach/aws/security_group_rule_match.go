package aws

type securityGroupRuleMatch struct {
	Basis       securityGroupRuleMatchBasis
	Requirement interface{}
	Value       interface{}
}
