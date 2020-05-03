package reach

type Path struct {
	Segments []Segment
}

func NewPath(firstPoint Point) Path {
	s := Segment{
		Points: []Point{firstPoint},
	}
	path := Path{
		Segments: []Segment{s},
	}
	return path
}

func (p Path) Zero() bool {
	return len(p.Segments) == 0
}

func (p Path) LastPoint() Point {
	lastSegment := p.Segments[len(p.Segments)-1]
	lastPoint := lastSegment.Points[len(lastSegment.Points)-1]
	return lastPoint
}

func (p Path) LastEdge() Edge {
	lastSegment := p.Segments[len(p.Segments)-1]
	lastEdge := lastSegment.Edges[len(lastSegment.Edges)-1]
	return lastEdge
}

func (p Path) Contains(ref UniversalReference) bool {
	for _, s := range p.Segments {
		if s.Contains(ref) {
			return true
		}
	}

	return false
}

func (p Path) Add(edge Edge, point Point, newSegment bool) Path {
	if newSegment {
		p.Segments = append(p.Segments, Segment{})
	}

	lastSegmentIndex := len(p.Segments) - 1
	p.Segments[lastSegmentIndex] = p.Segments[lastSegmentIndex].Add(edge, point)

	return p
}

func (p Path) FactorsForward() []Factor {
	var factors []Factor

	for _, s := range p.Segments {
		factors = append(factors, s.FactorsForward()...)
	}

	return factors
}

func (p Path) FactorsReturn() []Factor {
	var factors []Factor

	for _, s := range p.Segments {
		factors = append(factors, s.FactorsReturn()...)
	}

	return factors
}

func (p Path) TrafficForward() (TrafficContent, error) {
	tcs := TrafficFromFactors(p.FactorsForward())
	result, err := NewTrafficContentFromIntersectingMultiple(tcs)
	if err != nil {
		return TrafficContent{}, err
	}

	return result, nil
}

func (p Path) TrafficReturn() (TrafficContent, error) {
	// TODO: This method shouldn't exist — return traffic should not be intersected across multiple segments!
	tcs := TrafficFromFactors(p.FactorsReturn())
	result, err := NewTrafficContentFromIntersectingMultiple(tcs)
	if err != nil {
		return TrafficContent{}, err
	}

	return result, nil
}
