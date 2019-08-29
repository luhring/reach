package reach

import (
	"reflect"
	"testing"
)

func TestNewEC2InstanceSubject(t *testing.T) {
	cases := []struct {
		name            string
		id              string
		role            string
		expectedSubject *subject
		expectedError   error
	}{
		{
			name: "valid input with source role",
			id:   "i-abc123",
			role: RoleSource,
			expectedSubject: &subject{
				Kind:       ec2InstanceSubjectKind,
				Properties: ec2InstanceSubjectProperties{ID: "i-abc123"},
				Role:       RoleSource,
			},
			expectedError: nil,
		},
		{
			name: "valid input with destination role",
			id:   "i-def456",
			role: destination,
			expectedSubject: &subject{
				Kind:       ec2InstanceSubjectKind,
				Properties: ec2InstanceSubjectProperties{ID: "i-def456"},
				Role:       destination,
			},
			expectedError: nil,
		},
		{
			name:            "invalid ID value",
			id:              "",
			role:            RoleSource,
			expectedSubject: nil,
			expectedError:   newSubjectError(errSubjectIDValidation),
		},
		{
			name:            "invalid role value",
			id:              "i-abc123",
			role:            "custom-role",
			expectedSubject: nil,
			expectedError:   newSubjectError(errSubjectRoleValidation),
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			subj, err := NewEC2InstanceSubject(tc.id, tc.role)

			if !reflect.DeepEqual(tc.expectedSubject, subj) {
				diffErrorf(t, "subj", tc.expectedSubject, subj)
			}

			if !reflect.DeepEqual(tc.expectedError, err) {
				diffErrorf(t, "err", tc.expectedError, err)
			}
		})
	}
}
