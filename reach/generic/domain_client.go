package generic

import (
	"fmt"

	"github.com/luhring/reach/reach"
)

// The DomainClient interface wraps all of the necessary methods for accessing generic domain resources.
type DomainClient interface {
	Hostname(name string) (*Hostname, error)
}

func unpackDomainClient(resolver reach.DomainClientResolver) (DomainClient, error) {
	d := resolver.Resolve(ResourceDomainGeneric)
	if d == nil {
		return nil, fmt.Errorf("DomainClientResolver has no entry for domain '%s'", ResourceDomainGeneric)
	}
	domainClient, ok := d.(DomainClient)
	if !ok {
		return nil, fmt.Errorf("DomainClient interface not implemented correctly for domain '%s'", ResourceDomainGeneric)
	}
	return domainClient, nil
}
