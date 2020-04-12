package reach

type Traceable interface {
	Visitable(alreadyVisited bool) bool
	Ref() InfrastructureReference
	Segments() bool
	Tuple(prev *IPTuple) *IPTuple
	Next(t *IPTuple, provider InfrastructureGetter) ([]InfrastructureReference, error)
	Factors() []Factor
}
