package aws

type securityGroupRuleMatchBasis string

const securityGroupRuleMatchBasisIP securityGroupRuleMatchBasis = "IP"
const securityGroupRuleMatchBasisSGRef securityGroupRuleMatchBasis = "SecurityGroupReference"

// String returns the string representation of a security group rule match.
func (basis securityGroupRuleMatchBasis) String() string {
	switch basis {
	case securityGroupRuleMatchBasisIP:
		return "IP address"
	case securityGroupRuleMatchBasisSGRef:
		return "attached security group"
	default:
		return "[unknown match basis]"
	}
}
