package aws

import (
	"github.com/luhring/reach/reach"
)

// ResolveEC2InstanceSubject looks up an EC2Instance using the given provider and returns it as a new subject.
func ResolveEC2InstanceSubject(identifier string, domains reach.DomainClientResolver) (*reach.Subject, error) {
	resources, err := unpackDomainClient(domains)
	if err != nil {
		return nil, err
	}

	// We'll assume the identifier refers to an EC2 instance, even if it doesn't begin with 'i-'.
	// Later, we might use this string to recognize different kinds of AWS resources.
	ec2InstanceID, err := findEC2InstanceID(identifier, resources)
	if err != nil {
		return nil, err
	}

	return NewEC2InstanceSubject(ec2InstanceID), nil
}
