package generic

import (
	"net"

	"github.com/luhring/reach/reach"
)

const SubjectKindIPAddress = "IPAddress"

func NewIPAddressSubject(addr net.IP, role reach.SubjectRole) (*reach.Subject, error) {
	if !reach.ValidSubjectRole(role) {
		return nil, reach.NewSubjectError(reach.ErrSubjectRoleValidation)
	}

	return &reach.Subject{
		Domain: ResourceDomainGeneric,
		Kind:   SubjectKindIPAddress,
		ID:     addr.String(),
		Role:   role,
	}, nil
}
