package reach

import (
	"fmt"
)

const (
	SubjectRoleNone          = "none"
	SubjectRoleSource        = "source"
	SubjectRoleDestination   = "destination"
	ErrSubjectPrefix         = "subject creation error"
	ErrSubjectRoleValidation = "subject role must be 'source' or 'destination'"
	ErrSubjectIDValidation   = "id must be a non-empty string"
)

type Subject struct {
	Kind         string      `json:"kind"`
	Properties   interface{} `json:"properties"`
	Role         string      `json:"role"`
	GetResources func() []Resource
}

func (s *Subject) SetRoleToSource() {
	s.SetRole(SubjectRoleSource)
}

func (s *Subject) SetRoleToDestination() {
	s.SetRole(SubjectRoleDestination)
}

func (s *Subject) SetRole(role string) {
	if ValidSubjectRole(role) {
		s.Role = role
	}
}

func ValidSubjectRole(role string) bool {
	return role == SubjectRoleNone || role == SubjectRoleSource || role == SubjectRoleDestination
}

func NewSubjectError(details string) error {
	return fmt.Errorf("%s: %s", ErrSubjectPrefix, details)
}
