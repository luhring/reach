package aws

import "github.com/luhring/reach/reach"

const ResourceKindSecurityGroupReference = "SecurityGroupReference"

type SecurityGroupReference struct {
	ID        string `json:"id"`
	AccountID string `json:"accountID"`
	NameTag   string `json:"nameTag"`
	GroupName string `json:"groupName"`
}

func (sgRef SecurityGroupReference) ToResource() reach.Resource {
	return reach.Resource{
		Kind:       ResourceKindSecurityGroupReference,
		Properties: sgRef,
	}
}
