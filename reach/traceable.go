package reach

type Traceable interface {
	Visitable(alreadyVisited bool) bool
	Ref() InfrastructureReference
	Segments() bool
	ForwardEdges(latestTuple *IPTuple, provider InfrastructureGetter) ([]PathEdge, error)
	Factors() []Factor
}
