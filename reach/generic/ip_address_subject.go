package generic

import (
	"github.com/luhring/reach/reach"
)

const SubjectKindIPAddress = "IPAddress"

func NewIPAddressSubject(address string) *reach.Subject {
	return &reach.Subject{
		Domain: ResourceDomainGeneric,
		Kind:   SubjectKindIPAddress,
		ID:     address,
		Role:   reach.SubjectRoleNone,
	}
}
