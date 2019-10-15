package reach

type VectorDiscoverer interface {
	Discover([]*Subject) ([]NetworkVector, error)
}
