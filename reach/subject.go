package reach

import (
	"fmt"
)

type SubjectRole string

const (
	SubjectRoleNone          SubjectRole = "none"
	SubjectRoleSource        SubjectRole = "source"
	SubjectRoleDestination   SubjectRole = "destination"
	ErrSubjectPrefix                     = "subject creation error"
	ErrSubjectRoleValidation             = "subject role must be 'source' or 'destination'"
	ErrSubjectIDValidation               = "id must be a non-empty string"
)

type Subject struct {
	Domain string
	Kind   string
	ID     string
	Role   SubjectRole
}

func (s *Subject) SetRoleToSource() {
	s.SetRole(SubjectRoleSource)
}

func (s *Subject) SetRoleToDestination() {
	s.SetRole(SubjectRoleDestination)
}

func (s *Subject) SetRole(role SubjectRole) {
	if ValidSubjectRole(role) {
		s.Role = role
	}
}

func ValidSubjectRole(role SubjectRole) bool {
	return role == SubjectRoleNone || role == SubjectRoleSource || role == SubjectRoleDestination
}

func NewSubjectError(details string) error {
	return fmt.Errorf("%s: %s", ErrSubjectPrefix, details)
}
