package reach

type Path struct {
	Segments []Segment
}

func NewPath() Path {
	s := Segment{}
	path := Path{
		Segments: []Segment{s},
	}
	return path
}

func (p Path) LastPoint() Point {
	lastSegment := p.Segments[len(p.Segments)-1]
	lastPoint := lastSegment.Points[len(lastSegment.Points)-1]
	return lastPoint
}

func (p Path) Factors() []Factor {
	var factors []Factor

	for _, s := range p.Segments {
		factors = append(factors, s.Factors()...)
	}

	return factors
}

func (p Path) ForwardTraffic() TrafficContent {
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
