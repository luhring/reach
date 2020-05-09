package reach

// A Flow represents a direction of travel for network traffic with respect to a point in a network path.
type Flow int

// The possible values for a Flow.
const (
	FlowUnknown Flow = iota
	FlowOutbound
	FlowInbound
	FlowDropped
)

// String returns the string representation of the Flow.
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
