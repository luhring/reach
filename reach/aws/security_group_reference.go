package aws

import "github.com/luhring/reach/reach"

// ResourceKindSecurityGroupReference specifies the unique name for the security group reference kind of resource.
const ResourceKindSecurityGroupReference reach.Kind = "SecurityGroupReference"

// A SecurityGroupReference resource representation. A SecurityGroupReference is similar to a SecurityGroup, except it intentionally omits any further dependencies, so as to prevent a dependency cycle when security groups have security group rules that refer to security groups.
type SecurityGroupReference struct {
	ID        string
	AccountID string
	NameTag   string
	GroupName string
}

// Resource returns the security group reference converted to a generalized Reach resource.
func (sgRef SecurityGroupReference) Resource() reach.Resource {
	return reach.Resource{
		Kind:       ResourceKindSecurityGroupReference,
		Properties: sgRef,
	}
}
