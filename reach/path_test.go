package reach

import (
	"reflect"
	"testing"
)

func TestSegments(t *testing.T) {
	cases := []struct {
		name             string
		path             Path
		segmentsExpected []Path
	}{
		{
			name: "single point",
			path: NewPath(
				point(),
			),
			segmentsExpected: []Path{
				NewPath(point()),
			},
		},
		{
			name: "single segment, two points",
			path: Path{
				Points: []Point{
					point(),
					point(),
				},
				Edges: []Edge{
					{},
				},
			},
			segmentsExpected: []Path{
				{
					Points: []Point{
						point(),
						point(),
					},
					Edges: []Edge{
						{},
					},
				},
			},
		},
		{
			name: "two segments",
			path: Path{
				Points: []Point{
					point(),
					pointDivider(),
					point(),
				},
				Edges: []Edge{
					{},
					{},
				},
			},
			segmentsExpected: []Path{
				{
					Points: []Point{
						point(),
						pointDivider(),
					},
					Edges: []Edge{
						{},
					},
				},
				{
					Points: []Point{
						pointDivider(),
						point(),
					},
					Edges: []Edge{
						{},
					},
				},
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			outcome := tc.path.Segments()
			if !reflect.DeepEqual(outcome, tc.segmentsExpected) {
				t.Errorf("outcome was not expected:\noutcome: %+v\nexpected: %+v\n", outcome, tc.segmentsExpected)
			}
		})
	}
}

func point() Point {
	return Point{
		Ref:            Reference{},
		FactorsForward: nil,
		FactorsReturn:  nil,
		SegmentDivider: false,
	}
}

func pointDivider() Point {
	return Point{
		Ref:            Reference{},
		FactorsForward: nil,
		FactorsReturn:  nil,
		SegmentDivider: true,
	}
}
