package generic

import (
	"fmt"

	"github.com/luhring/reach/reach"
)

func ResolveIPAddressSubject(identifier string) (*reach.Subject, error) {
	err := CheckIPAddress(identifier)
	if err != nil {
		return nil, fmt.Errorf("unable to resolve IP address subject: %v", err)
	}

	return NewIPAddressSubject(identifier), nil
}
