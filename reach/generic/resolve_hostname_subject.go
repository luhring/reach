package generic

import (
	"fmt"

	"github.com/luhring/reach/reach"
)

// ResolveHostnameSubject looks up a Hostname using the given provider and returns it as a new subject.
func ResolveHostnameSubject(identifier string) (*reach.Subject, error) {
	err := CheckHostname(identifier)
	if err != nil {
		return nil, fmt.Errorf("unable to resolve hostname subject: %v", err)
	}

	return NewHostnameSubject(identifier), nil
}
