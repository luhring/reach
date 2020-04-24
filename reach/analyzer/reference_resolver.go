package analyzer

import (
	"fmt"

	"github.com/luhring/reach/reach"
	"github.com/luhring/reach/reach/aws"
	"github.com/luhring/reach/reach/generic"
)

type ReferenceResolver struct {
	clientResolver reach.DomainClientResolver
}

func NewReferenceResolver(clientResolver reach.DomainClientResolver) ReferenceResolver {
	return ReferenceResolver{
		clientResolver: clientResolver,
	}
}

func (r *ReferenceResolver) Resolve(ref reach.UniversalReference) (*reach.Resource, error) {
	r.clientResolver.Resolve(ref.R.Domain) // todo: delete?

	switch ref.R.Domain {
	case aws.ResourceDomainAWS:
		awsResolver, err := aws.NewReferenceResolver(r.clientResolver)
		if err != nil {
			return nil, fmt.Errorf("cannot resolve UniversalReference: %v", err)
		}
		return awsResolver.Resolve(ref)
	case generic.ResourceDomainGeneric:
		genericResolver, err := generic.NewReferenceResolver(r.clientResolver)
		if err != nil {
			return nil, fmt.Errorf("cannot resolve UniversalReference: %v", err)
		}
		return genericResolver.Resolve(ref)
	}

	return nil, fmt.Errorf("root resolver encountered an unexpected domain '%s'", ref.R.Domain)
}
