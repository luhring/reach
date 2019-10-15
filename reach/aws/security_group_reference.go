package aws

import "github.com/luhring/reach/reach"

const ResourceKindSecurityGroupReference = "SecurityGroupReference"

type SecurityGroupReference struct {
	ID        string
	AccountID string
	NameTag   string
	GroupName string
}

func (sgRef SecurityGroupReference) ToResource() reach.Resource {
	return reach.Resource{
		Kind:       ResourceKindSecurityGroupReference,
		Properties: sgRef,
	}
}
