package reach

import (
	"fmt"
)

// SubjectRole specifies the role the subject plays in an analysis -- i.e. that this subject is the "source" or the "destination".
type SubjectRole string

// Allowed values for SubjectRole.
const (
	SubjectRoleNone        SubjectRole = "none"
	SubjectRoleSource      SubjectRole = "source"
	SubjectRoleDestination SubjectRole = "destination"
)

// Common errors for the Subject type.
const (
	ErrSubjectPrefix         = "subject creation error"
	ErrSubjectRoleValidation = "subject role must be 'source' or 'destination'"
	ErrSubjectIDValidation   = "id must be a non-empty string"
)

// A Subject is an entity about which a network traffic question is being asked. Reach analyses are conducted between "source" subjects and "destination" subjects. For example, when asking about network traffic allowed between instance A and instance B, instances A and B are the "subjects" of the analysis.
type Subject struct {
	Domain Domain
	Kind   Kind
	ID     string
	Role   SubjectRole
}

func (s Subject) Ref() UniversalReference {
	return UniversalReference{
		Domain: s.Domain,
		Kind:   s.Kind,
		ID:     s.ID,
	}
}

// SetRoleToSource sets the subject's role to "source".
func (s *Subject) SetRoleToSource() {
	s.setRole(SubjectRoleSource)
}

// SetRoleToDestination sets the subject's role to "destination".
func (s *Subject) SetRoleToDestination() {
	s.setRole(SubjectRoleDestination)
}

// ValidSubjectRole returns a boolean indicating whether or not the specified subject role is valid.
func ValidSubjectRole(role SubjectRole) bool {
	return role == SubjectRoleNone || role == SubjectRoleSource || role == SubjectRoleDestination
}

// NewSubjectError generates a new error related to a subject operation.
func NewSubjectError(details string) error {
	return fmt.Errorf("%s: %s", ErrSubjectPrefix, details)
}

func (s *Subject) setRole(role SubjectRole) {
	if ValidSubjectRole(role) {
		s.Role = role
	}
}
