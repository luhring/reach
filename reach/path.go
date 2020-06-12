package reach

import "github.com/luhring/reach/reach/traffic"

// Path is the series of points, and the connections in between, from one point in a network to another.
type Path struct {
	Points []Point
	Edges  []Edge
}

// NewPath returns a new, initialized path.
func NewPath(firstPoint Point) Path {
	path := Path{
		Points: []Point{firstPoint},
		Edges:  []Edge{},
	}
	return path
}

// Zero returns a bool that indicates whether the path contains no inner items.
func (p Path) Zero() bool {
	return len(p.Points) == 0
}

// LastPoint returns the final point in the path.
func (p Path) LastPoint() Point {
	lastPoint := p.Points[len(p.Points)-1]
	return lastPoint
}

// LastEdge returns the final most edge in the path.
func (p Path) LastEdge() Edge {
	lastEdge := p.Edges[len(p.Edges)-1]
	return lastEdge
}

// Contains returns a bool that indicates whether the path contains a point for the specified reference.
//
// This is useful when determining if a given piece of infrastructure is a part of the network path.
func (p Path) Contains(ref Reference) bool {
	for _, point := range p.Points {
		if point.Ref.Equal(ref) {
			return true
		}
	}

	return false
}

// Add returns the path, having been updated to include a new edge and a new point.
func (p Path) Add(edge Edge, point Point) Path {
	return Path{
		Points: append(p.Points, point),
		Edges:  append(p.Edges, edge),
	}
}

// FactorsForward returns a slice of all of the forward-bound factors that exist for each point along the path.
func (p Path) FactorsForward() []Factor {
	var factors []Factor

	for _, point := range p.Points {
		factors = append(factors, point.FactorsForward...)
	}

	return factors
}

// FactorsReturn returns a slice of all of the return-bound factors that exist for each point along the path.
func (p Path) FactorsReturn() []Factor {
	var factors []Factor

	for _, point := range p.Points {
		factors = append(factors, point.FactorsReturn...)
	}

	return factors
}

// TrafficForward returns the traffic that is allowed to travel forward along the entire network path.
func (p Path) TrafficForward() traffic.Content {
	return traffic.Intersect(TrafficFromFactors(p.FactorsForward()))
}

// TrafficReturn returns the traffic that is allowed to travel backward along the entire network path. IMPORTANT: This operation only makes sense if this path is not divided into multiple segments.
func (p Path) TrafficReturn() traffic.Content {
	fs := p.FactorsReturn()
	ts := TrafficFromFactors(fs)
	return traffic.Intersect(ts)
}

// MapPoints creates a new version of the path where each point has been transformed by the supplied mapper function.
func (p Path) MapPoints(
	mapper func(point Point, leftEdge, rightEdge *Edge) (Point, error),
) (Path, error) {
	mappedPoints := make([]Point, len(p.Points))
	for i, point := range p.Points {
		var leftEdge *Edge
		if i > 0 {
			leftEdge = &p.Edges[i-1]
		}

		var rightEdge *Edge
		if i < len(p.Points)-1 {
			rightEdge = &p.Edges[i]
		}

		mappedPoint, err := mapper(point, leftEdge, rightEdge)
		if err != nil {
			return Path{}, err
		}
		mappedPoints[i] = mappedPoint
	}

	return Path{
		Points: mappedPoints,
		Edges:  p.Edges,
	}, nil
}

// Segments returns the continuous sets of points where no PAT occurs that comprise the Path.
func (p Path) Segments() []Path {
	var segments []Path
	var currentSegment Path
	points := p.Points

	for i, point := range points {
		if i != 0 {
			currentSegment = currentSegment.Add(p.Edges[i-1], point)
		} else {
			currentSegment = NewPath(point)
		}

		// I.e. if segment is complete
		if len(points) == 1 || point.SegmentDivider || i == len(points)-1 {
			segments = append(segments, currentSegment)
		}

		if point.SegmentDivider {
			currentSegment = NewPath(point)
		}
	}

	return segments
}
