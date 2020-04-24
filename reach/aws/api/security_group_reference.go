package api

import (
	reachAWS "github.com/luhring/reach/reach/aws"
)

// SecurityGroupReference queries the AWS API for a security group matching the given ID, but returns a security group reference representation instead of the full security group representation.
func (client *DomainClient) SecurityGroupReference(id, accountID string) (*reachAWS.SecurityGroupReference, error) {
	// TODO: Incorporate account ID in search.
	// In the meantime, this will be a known bug, where other accounts are not considered.

	sg, err := client.SecurityGroup(id)
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
