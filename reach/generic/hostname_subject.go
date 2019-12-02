package generic

import (
	"github.com/luhring/reach/reach"
)

const SubjectKindHostname = "Hostname"

func NewHostnameSubject(hostname string) *reach.Subject {
	return &reach.Subject{
		Domain: ResourceDomainGeneric,
		Kind:   SubjectKindHostname,
		ID:     hostname,
		Role:   reach.SubjectRoleNone,
	}
}
