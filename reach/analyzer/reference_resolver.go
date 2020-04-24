package analyzer

import "github.com/luhring/reach/reach"

type ReferenceResolver struct {
	resolver reach.DomainClientResolver
}

func NewReferenceResolver(resolver reach.DomainClientResolver) ReferenceResolver {
	return ReferenceResolver{
		resolver: resolver,
	}
}

func (r *ReferenceResolver) Resolve(ref reach.UniversalReference) (reach.Resource, error) {
	ref.R.Domain
}
