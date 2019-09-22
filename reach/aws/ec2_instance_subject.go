package aws

import "github.com/luhring/reach/reach"

const SubjectKindEC2Instance = "EC2Instance"

func NewEC2InstanceSubject(id, role string) (*reach.Subject, error) {
	if !reach.ValidSubjectRole(role) {
		return nil, reach.NewSubjectError(reach.ErrSubjectRoleValidation)
	}

	if len(id) < 1 {
		return nil, reach.NewSubjectError(reach.ErrSubjectIDValidation)
	}

	return &reach.Subject{
		Domain: ResourceDomainAWS,
		Kind:   SubjectKindEC2Instance,
		ID:     id,
		Role:   role,
	}, nil
}
