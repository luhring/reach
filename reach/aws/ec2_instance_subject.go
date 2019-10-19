package aws

import "github.com/luhring/reach/reach"

// SubjectKindEC2Instance specifies the unique name for the EC2 instance kind of subject.
const SubjectKindEC2Instance = "EC2Instance"

// NewEC2InstanceSubject returns a new subject for the specified EC2 instance.
func NewEC2InstanceSubject(id string, role reach.SubjectRole) (*reach.Subject, error) {
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
