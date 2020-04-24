package reach

type Flow int

const (
	FlowUnknown Flow = iota
	FlowOutbound
	FlowInbound
	FlowDropped
)

func (f Flow) String() string {
	switch f {
	case FlowOutbound:
		return "FlowOutbound"
	case FlowInbound:
		return "FlowInbound"
	case FlowDropped:
		return "FlowDropped"
	default:
		return "FlowUnknown"
	}
}
