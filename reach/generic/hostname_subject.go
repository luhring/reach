package generic

import (
	"github.com/luhring/reach/reach"
)

// SubjectKindHostname specifies the unique name for the Hostname kind of subject.
const SubjectKindHostname = "Hostname"

// NewHostnameSubject returns a pointer to a new instance of a reach.Subject for the specified Hostname.
func NewHostnameSubject(hostname string) *reach.Subject {
	return &reach.Subject{
		Domain: ResourceDomainGeneric,
		Kind:   SubjectKindHostname,
		ID:     hostname,
		Role:   reach.SubjectRoleNone,
	}
}
