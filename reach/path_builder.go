package reach

// PathBuilder providers helper functionality for building a path, one point at a time, with the ability to signal the start of new path segments as needed.
type PathBuilder struct {
	path Path
}

func NewPathBuilder() *PathBuilder {
	s := Segment{}
	path := Path{
		Segments: []Segment{s},
	}

	return &PathBuilder{
		path: path,
	}
}

func ResumePathBuilding(path Path) *PathBuilder {
	return &PathBuilder{
		path: path,
	}
}

// Add adds a new point to the current path segment of the path being built.
func (p *PathBuilder) Add(point Point) {
	i := p.currentSegment()
	p.path.Segments[i].Points = append(p.path.Segments[i].Points, point)
}

// AddSegment tells the PathBuilder to begin a new path segment, such as because some kind of NAT is occurring at this point. Calls to Add will now add points to this new segment.
func (p *PathBuilder) AddSegment() {
	s := Segment{}
	p.path.Segments = append(p.path.Segments, s)
}

// Path returns the current state of the path being built by the PathBuilder.
func (p *PathBuilder) Path() Path {
	return p.path
}

func (p *PathBuilder) currentSegment() int {
	return len(p.path.Segments) - 1
}
