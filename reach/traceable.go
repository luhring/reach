package reach

type Traceable interface {
	Visitable(alreadyVisited bool) bool
	Ref() InfrastructureReference
	Segments() bool
	ForwardEdges(previousEdge Edge, domains DomainProvider) ([]Edge, error)
	Factors() []Factor
}
