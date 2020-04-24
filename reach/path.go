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

func (p Path) Contains(ref InfrastructureReference) bool {
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

func (p Path) Factors() []Factor {
	var factors []Factor

	for _, s := range p.Segments {
		factors = append(factors, s.Factors()...)
	}

	return factors
}

func (p Path) ForwardTraffic() TrafficContent { // TODO: Consider moving this for analyzer to conclude
	var tcs []TrafficContent

	for _, factor := range p.Factors() {
		t := factor.Traffic
		tcs = append(tcs, t)
	}

	result, err := NewTrafficContentFromIntersectingMultiple(tcs)
	if err != nil {
		panic(err) // TODO: Don't panic
	}

	return result
}
