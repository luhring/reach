package reach

// A Segment is a series of Points in a Path for which it is known that no NAT/PAT has occurred between any two consecutive points in the series. By definition, once NAT/PAT occurs in a Path, a new Segment begins.
type Segment struct {
	Points []Point
	Edges  []Edge
}

func (s Segment) FactorsForward() []Factor {
	var result []Factor

	for _, point := range s.Points {
		result = append(result, point.FactorsForward...)
	}

	return result
}

func (s Segment) FactorsReturn() []Factor {
	var result []Factor

	for _, point := range s.Points {
		result = append(result, point.FactorsReturn...)
	}

	return result
}

func (s Segment) Contains(ref Reference) bool {
	for _, point := range s.Points {
		if point.Ref.Equal(ref) {
			return true
		}
	}

	return false
}

func (s Segment) Add(edge Edge, point Point) Segment {
	return Segment{
		Edges:  append(s.Edges, edge),
		Points: append(s.Points, point),
	}
}
