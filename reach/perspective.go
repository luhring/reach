package reach

// A Perspective provides a reference to one direction of a network vector with knowledge of which network point is currently being analyzed ("self") and which network point is the "other" or "target" network point, such that properties of the "other" network point can be used when determining of it applies to the analysis of the "self" network point.
type Perspective struct {
	Self      NetworkPoint
	Other     NetworkPoint
	SelfRole  SubjectRole
	OtherRole SubjectRole
}
