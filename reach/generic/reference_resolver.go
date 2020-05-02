package generic

import (
	"fmt"

	"github.com/luhring/reach/reach"
)

type ReferenceResolver struct {
	client DomainClient
}

func NewReferenceResolver(clientResolver reach.DomainClientResolver) (*ReferenceResolver, error) {
	client, err := unpackDomainClient(clientResolver)
	if err != nil {
		return nil, fmt.Errorf("cannot create new Generic ReferenceResolver: %v", err)
	}

	return &ReferenceResolver{client: client}, nil
}

func (r *ReferenceResolver) Resolve(ref reach.UniversalReference) (*reach.Resource, error) {
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
