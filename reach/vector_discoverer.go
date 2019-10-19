package reach

// A VectorDiscoverer can return all network vectors that exist between specified subjects.
type VectorDiscoverer interface {
	Discover([]*Subject) ([]NetworkVector, error)
}
