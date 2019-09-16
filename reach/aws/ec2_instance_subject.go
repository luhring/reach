package aws

import "github.com/luhring/reach/reach"

const SubjectKindEC2Instance = "EC2Instance"

type EC2InstanceSubject struct {
	ID string `json:"id"`
}

func NewEC2InstanceSubject(id, role string) (*reach.Subject, error) {
	if !reach.ValidSubjectRole(role) {
		return nil, reach.NewSubjectError(reach.ErrSubjectRoleValidation)
	}

	if len(id) < 1 {
		return nil, reach.NewSubjectError(reach.ErrSubjectIDValidation)
	}

	props := EC2InstanceSubject{
		ID: id,
	}

	return &reach.Subject{
		Kind:       SubjectKindEC2Instance,
		Properties: props,
		Role:       role,
	}, nil
}
