package aws

import (
	"github.com/luhring/reach/reach"
)

// ResolveEC2InstanceSubject looks up an EC2 instance using the given provider and returns it as a new subject.
func ResolveEC2InstanceSubject(identifier string, provider ResourceProvider) (*reach.Subject, error) {
	// We'll assume the identifier refers to an EC2 instance, even if it doesn't begin with 'i-'.
	// Later, we might use this string to recognize different kinds of AWS resources.
	ec2InstanceID, err := FindEC2InstanceID(identifier, provider)
	if err != nil {
		return nil, err
	}

	subject := NewEC2InstanceSubject(ec2InstanceID)

	return subject, nil
}
