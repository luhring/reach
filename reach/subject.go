package reach

const (
	roleSource      = "source"
	roleDestination = "destination"
)

type subject struct {
	Kind       string      `json:"kind"`
	Properties interface{} `json:"properties"`
	Role       string      `json:"role"`
}

func newSubjectForEC2Instance(id, role string) subject {
	props := ec2InstanceSubjectProperties{
		ID: id,
	}

	return subject{
		Kind:       "ec2InstanceSubject",
		Properties: props,
		Role:       role,
	}
}
