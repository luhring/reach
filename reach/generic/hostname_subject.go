package generic

import (
	"github.com/luhring/reach/reach"
)

const SubjectKindHostname = "Hostname"

func NewHostnameSubject(hostname string, role reach.SubjectRole) (*reach.Subject, error) {
	if !reach.ValidSubjectRole(role) {
		return nil, reach.NewSubjectError(reach.ErrSubjectRoleValidation)
	}

	return &reach.Subject{
		Domain: ResourceDomainGeneric,
		Kind:   SubjectKindHostname,
		ID:     hostname,
		Role:   role,
	}, nil
}
