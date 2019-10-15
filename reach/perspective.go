package reach

type Perspective struct {
	Self      NetworkPoint
	Other     NetworkPoint
	SelfRole  SubjectRole
	OtherRole SubjectRole
}
