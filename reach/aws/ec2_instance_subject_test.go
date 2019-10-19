package aws

import (
	"reflect"
	"testing"

	"github.com/luhring/reach/reach"
)

func TestNewEC2InstanceSubject(t *testing.T) {
	cases := []struct {
		name            string
		id              string
		role            reach.SubjectRole
		expectedSubject *reach.Subject
		expectedError   error
	}{
		{
			name: "valid input with source role",
			id:   "i-abc123",
			role: reach.SubjectRoleSource,
			expectedSubject: &reach.Subject{
				Domain: ResourceDomainAWS,
				Kind:   SubjectKindEC2Instance,
				ID:     "i-abc123",
				Role:   reach.SubjectRoleSource,
			},
			expectedError: nil,
		},
		{
			name: "valid input with destination role",
			id:   "i-def456",
			role: reach.SubjectRoleDestination,
			expectedSubject: &reach.Subject{
				Domain: ResourceDomainAWS,
				Kind:   SubjectKindEC2Instance,
				ID:     "i-def456",
				Role:   reach.SubjectRoleDestination,
			},
			expectedError: nil,
		},
		{
			name:            "invalid ID value",
			id:              "",
			role:            reach.SubjectRoleSource,
			expectedSubject: nil,
			expectedError:   reach.NewSubjectError(reach.ErrSubjectIDValidation),
		},
		{
			name:            "invalid role value",
			id:              "i-abc123",
			role:            "custom-role",
			expectedSubject: nil,
			expectedError:   reach.NewSubjectError(reach.ErrSubjectRoleValidation),
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			subj, err := NewEC2InstanceSubject(tc.id, tc.role)

			if !reflect.DeepEqual(tc.expectedSubject, subj) {
				reach.DiffErrorf(t, "subj", tc.expectedSubject, subj)
			}

			if !reflect.DeepEqual(tc.expectedError, err) {
				reach.DiffErrorf(t, "err", tc.expectedError, err)
			}
		})
	}
}
