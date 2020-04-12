package reach

// A Segment is a series of Points in a Path for which it is known that no NAT/PAT has occurred between any two consecutive points in the series. By definition, once NAT/PAT occurs in a Path, a new Segment begins.
type Segment struct {
	Points []Point
	Tuples []*IPTuple
}

func (s Segment) Factors() []Factor {
	var result []Factor

	for _, point := range s.Points {
		result = append(result, point.Factors...)
	}

	return result
}

func (s Segment) Contains(ref InfrastructureReference) bool {
	for _, point := range s.Points {
		if point.Ref.Equal(ref) {
			return true
		}
	}

	return false
}

func (s Segment) Add(tuple *IPTuple, point Point) Segment {
	return Segment{
		Tuples: append(s.Tuples, tuple),
		Points: append(s.Points, point),
	}
}
