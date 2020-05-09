package generic

import (
	"fmt"

	"github.com/luhring/reach/reach"
)

// ReferenceResolver is the generic domain's implementation of the interface reach.ReferenceResolver. This ReferenceResolver can resolve only generic references (like hostnames and IP addresses).
type ReferenceResolver struct {
	client DomainClient
}

// NewReferenceResolver returns a pointer a new ReferenceResolver. If NewReferenceResolver is unable to find a generic.DomainClient using the provided clientResolver, it returns an error.
func NewReferenceResolver(clientResolver reach.DomainClientResolver) (*ReferenceResolver, error) {
	client, err := unpackDomainClient(clientResolver)
	if err != nil {
		return nil, fmt.Errorf("cannot create new Generic ReferenceResolver: %v", err)
	}

	return &ReferenceResolver{client: client}, nil
}

// Resolve returns a Resource for the specified Reference. Resolve returns the error if the Reference does not specify the generic domain or if there is an error encountered while querying for the resource itself.
func (r *ReferenceResolver) Resolve(ref reach.Reference) (*reach.Resource, error) {
	if ref.Domain != ResourceDomainGeneric {
		return nil, fmt.Errorf("%s resolver cannot resolve references for domain '%s'", ResourceDomainGeneric, ref.Domain)
	}

	switch ref.Kind {
	case ResourceKindHostname:
		hostname, err := r.client.Hostname(ref.ID)
		if err != nil {
			return nil, err
		}
		resource := hostname.Resource()
		return &resource, nil
	}

	return nil, fmt.Errorf("%s resolver encountered an unexpected resource kind '%s'", ResourceDomainGeneric, ref.Kind)
}
