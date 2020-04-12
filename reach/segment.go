package reach

// A Segment is a series of Points in a Path for which it is known that no NAT/PAT has occurred between any two consecutive points in the series. By definition, once NAT/PAT occurs in a Path, a new Segment begins.
type Segment struct {
	Points []Point
}

func (s Segment) Factors() []Factor {
	var result []Factor

	for _, point := range s.Points {
		result = append(result, point.Factors...)
	}

	return result
}

func (s Segment) Contains(pt Point) bool {
	for _, point := range s.Points {
		if point.Ref.Matches(pt.Ref) { // TODO: (Later) Consider a more intelligent loop detection system that leverages tuples
			return true
		}
	}

	return false
}
