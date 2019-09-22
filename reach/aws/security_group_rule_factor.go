package aws

type SecurityGroupRuleFactor struct {
	RuleIndex  int
	MatchBasis string // named type?
	MatchValue string // IP address, SG Ref ID, Prefix list name?
}
