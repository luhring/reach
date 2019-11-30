package reach

// A PointsDiscoverer can identify and return all network points that exist for a given subject.
type PointsDiscoverer interface {
	Discover(subject Subject) ([]NetworkPoint, error)
}
