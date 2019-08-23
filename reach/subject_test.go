package reach

import (
	"fmt"
	"reflect"
	"strings"
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
			role: roleSource,
			expectedSubject: &subject{
				Kind:       ec2InstanceSubjectKind,
				Properties: ec2InstanceSubjectProperties{ID: "i-abc123"},
				Role:       roleSource,
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
			role:            roleSource,
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
			subject, err := newEC2InstanceSubject(tc.id, tc.role)

			var problems []string

			if false == reflect.DeepEqual(err, tc.expectedError) {
				problems = append(problems, fmt.Sprintf("expected error to be %v but it was %v", tc.expectedError, err))
			}

			if false == reflect.DeepEqual(subject, tc.expectedSubject) {
				problems = append(problems, fmt.Sprintf("expected subject to be %v but it was %v", tc.expectedSubject, subject))
			}

			if len(problems) > 0 {
				message := strings.Join(problems, ", ")

				t.Errorf("one or more expectations were failed: %s", message)
			}
		})
	}
}
