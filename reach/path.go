package reach

// Path is the series of points, and the connections in between, from one point in a network to another.
type Path struct {
	Segments []Segment
}

// NewPath returns a new, initialized path.
func NewPath(firstPoint Point) Path {
	s := Segment{
		Points: []Point{firstPoint},
	}
	path := Path{
		Segments: []Segment{s},
	}
	return path
}

// Zero returns a bool that indicates whether the path contains no inner items.
func (p Path) Zero() bool {
	return len(p.Segments) == 0
}

// LastPoint returns the final point in the path.
func (p Path) LastPoint() Point {
	lastSegment := p.Segments[len(p.Segments)-1]
	lastPoint := lastSegment.Points[len(lastSegment.Points)-1]
	return lastPoint
}

// LastEdge returns the final most edge in the path.
func (p Path) LastEdge() Edge {
	lastSegment := p.Segments[len(p.Segments)-1]
	lastEdge := lastSegment.Edges[len(lastSegment.Edges)-1]
	return lastEdge
}

// Contains returns a bool that indicates whether the path contains a point for the specified reference.
//
// This is useful when determining if a given piece of infrastructure is a part of the network path.
func (p Path) Contains(ref Reference) bool {
	for _, s := range p.Segments {
		if s.Contains(ref) {
			return true
		}
	}

	return false
}

// Add returns the path, having been updated to include a new edge and a new point.
func (p Path) Add(edge Edge, point Point, newSegment bool) Path {
	if newSegment {
		p.Segments = append(p.Segments, Segment{})
	}

	lastSegmentIndex := len(p.Segments) - 1
	p.Segments[lastSegmentIndex] = p.Segments[lastSegmentIndex].Add(edge, point)

	return p
}

// FactorsForward returns a slice of all of the forward-bound factors that exist for each point along the path.
func (p Path) FactorsForward() []Factor {
	var factors []Factor

	for _, s := range p.Segments {
		factors = append(factors, s.FactorsForward()...)
	}

	return factors
}

// FactorsReturn returns a slice of all of the return-bound factors that exist for each point along the path.
func (p Path) FactorsReturn() []Factor {
	var factors []Factor

	for _, s := range p.Segments {
		factors = append(factors, s.FactorsReturn()...)
	}

	return factors
}

// TrafficForward returns the traffic that is allowed to travel forward along the entire network path.
func (p Path) TrafficForward() (TrafficContent, error) {
	tcs := TrafficFromFactors(p.FactorsForward())
	result, err := NewTrafficContentFromIntersectingMultiple(tcs)
	if err != nil {
		return TrafficContent{}, err
	}

	return result, nil
}

// TrafficReturn returns the traffic that is allowed to return from the last point to the first point along the network path.
func (p Path) TrafficReturn() (TrafficContent, error) {
	// TODO: This method shouldn't exist — return traffic should not be intersected across multiple segments!
	tcs := TrafficFromFactors(p.FactorsReturn())
	result, err := NewTrafficContentFromIntersectingMultiple(tcs)
	if err != nil {
		return TrafficContent{}, err
	}

	return result, nil
}
