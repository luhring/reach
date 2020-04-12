package reach

type Traceable interface {
	Visitable(alreadyVisited bool) bool
	Ref() InfrastructureReference
	Segments() bool
	ForwardEdges(prev *IPTuple, provider InfrastructureGetter) ([]PathEdge, error)
	Factors() []Factor
}
