package aws

import (
	"fmt"
	"strings"

	"github.com/luhring/reach/reach"
)

const errNewSubjectsPrefix = "unable to create subject"
const errNewSubjectsNilInput = "input was nil"
const errNewSubjectsUnrecognizedIdentifier = "unrecognized identifier format"

func NewSubject(identifier string) (*reach.Subject, error) {
	if strings.HasPrefix(identifier, "i-") {
		ec2InstanceID, err := FindEC2InstanceID(identifier, nil)
		if err != nil {
			return nil, err
		}

		subject, err := NewEC2InstanceSubject(ec2InstanceID, reach.SubjectRoleNone)
		if err != nil {
			return nil, err
		}

		return subject, nil
	} else {
		return nil, fmt.Errorf("%s: %s: '%s'", errNewSubjectsPrefix, errNewSubjectsUnrecognizedIdentifier, identifier)
	}
}

func NewSubjectsAsSources(identifiers ...string) ([]reach.Subject, error) {
	return NewSubjectsWithRole(reach.SubjectRoleSource, identifiers...)
}

func NewSubjectsAsDestinations(identifiers ...string) ([]reach.Subject, error) {
	return NewSubjectsWithRole(reach.SubjectRoleDestination, identifiers...)
}

func NewSubjectsWithRole(role string, identifiers ...string) ([]reach.Subject, error) {
	if identifiers == nil {
		return nil, fmt.Errorf("%s: %s", errNewSubjectsPrefix, errNewSubjectsNilInput)
	}

	subjects := make([]reach.Subject, len(identifiers))

	for _, id := range identifiers {
		if strings.HasPrefix(id, "i-") {
			ec2InstanceID, err := FindEC2InstanceID(id, nil)
			if err != nil {
				return nil, err
			}

			subject, err := NewEC2InstanceSubject(ec2InstanceID, role)
			if err != nil {
				return nil, err
			}

			subjects = append(subjects, *subject)
		} else {
			return nil, fmt.Errorf("%s: %s: '%s'", errNewSubjectsPrefix, errNewSubjectsUnrecognizedIdentifier, id)
		}
	}

	return subjects, nil
}
