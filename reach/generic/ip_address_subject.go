package generic

import (
	"github.com/luhring/reach/reach"
)

// SubjectKindIPAddress specifies the unique name for the IPAddress kind of subject.
const SubjectKindIPAddress = "IPAddress"

// NewIPAddressSubject returns a pointer to a new instance of a reach.Subject for the specified IPAddress.
func NewIPAddressSubject(address string) *reach.Subject {
	return &reach.Subject{
		Domain: ResourceDomainGeneric,
		Kind:   SubjectKindIPAddress,
		ID:     address,
		Role:   reach.SubjectRoleNone,
	}
}
