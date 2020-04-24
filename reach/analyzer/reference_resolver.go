package analyzer

import "github.com/luhring/reach/reach"

type ReferenceResolver struct {
	domains reach.DomainProvider
}

func NewReferenceResolver(domains reach.DomainProvider) ReferenceResolver {
	return ReferenceResolver{
		domains: domains,
	}
}

func (r *ReferenceResolver) Resolve(ref reach.UniversalReference) (reach.Resource, error) {
	ref.R.Domain
}
