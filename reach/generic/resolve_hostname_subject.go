package generic

import (
	"fmt"

	"github.com/luhring/reach/reach"
)

func ResolveHostnameSubject(identifier string) (*reach.Subject, error) {
	err := CheckHostname(identifier)
	if err != nil {
		return nil, fmt.Errorf("unable to resolve hostname subject: %v", err)
	}

	return NewHostnameSubject(identifier), nil
}
