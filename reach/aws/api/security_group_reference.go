package api

import (
	reachAWS "github.com/luhring/reach/reach/aws"
)

func (provider *ResourceProvider) GetSecurityGroupReference(id, accountID string) (*reachAWS.SecurityGroupReference, error) {
	// TODO: Incorporate account ID in search.
	// In the meantime, this will be a known bug, where other accounts are not considered.

	sg, err := provider.GetSecurityGroup(id)
	if err != nil {
		return nil, err
	}

	return &reachAWS.SecurityGroupReference{
		ID:        sg.ID,
		AccountID: "",
		NameTag:   sg.NameTag,
		GroupName: sg.GroupName,
	}, nil
}
