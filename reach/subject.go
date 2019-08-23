package reach

import "fmt"

const (
	roleSource               = "source"
	roleDestination          = "destination"
	ec2InstanceSubjectKind   = "ec2Instance"
	errSubjectPrefix         = "subject creation error"
	errSubjectRoleValidation = "subject role must be 'source' or 'destination'"
	errSubjectIDValidation   = "id must be a non-empty string"
)

type subject struct {
	Kind       string      `json:"kind"`
	Properties interface{} `json:"properties"`
	Role       string      `json:"role"`
}

func NewEC2InstanceSubject(id, role string) (*subject, error) {
	if role != roleSource && role != roleDestination {
		return nil, newSubjectError(errSubjectRoleValidation)
	}

	if len(id) < 1 {
		return nil, newSubjectError(errSubjectIDValidation)
	}

	props := ec2InstanceSubjectProperties{
		ID: id,
	}

	return &subject{
		Kind:       ec2InstanceSubjectKind,
		Properties: props,
		Role:       role,
	}, nil
}

func newSubjectError(details string) error {
	return fmt.Errorf("%s: %s", errSubjectPrefix, details)
}
