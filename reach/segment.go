package reach

// A Segment is a series of points in a path for which it is known that no NAT/PAT has occurred between any two consecutive points in the series. By definition, once NAT/PAT occurs in a path, a new segment begins.
type Segment struct {
	Points []Point
	Edges  []Edge
}

// FactorsForward returns a slice of all of the forward-bound factors that exist for each point along the segment.
func (s Segment) FactorsForward() []Factor {
	var result []Factor

	for _, point := range s.Points {
		result = append(result, point.FactorsForward...)
	}

	return result
}

// FactorsReturn returns a slice of all of the return-bound factors that exist for each point along the segment.
func (s Segment) FactorsReturn() []Factor {
	var result []Factor

	for _, point := range s.Points {
		result = append(result, point.FactorsReturn...)
	}

	return result
}

// Contains returns a bool that indicates whether the segment contains a point for the specified reference.
func (s Segment) Contains(ref Reference) bool {
	for _, point := range s.Points {
		if point.Ref.Equal(ref) {
			return true
		}
	}

	return false
}

// Add returns the segment, having been updated to include a new edge and a new point.
func (s Segment) Add(edge Edge, point Point) Segment {
	return Segment{
		Edges:  append(s.Edges, edge),
		Points: append(s.Points, point),
	}
}
