package reach

const (
	source      = "source"
	destination = "destination"
	outbound    = "outbound"
	inbound     = "inbound"
)

type perspective struct {
	self           string
	selfInterface  *NetworkInterface
	other          string
	otherInterface *NetworkInterface
	direction      string
	rules          func(sg *SecurityGroup) []*SecurityGroupRule
}

func newPerspectiveFromSource(vector *InterfaceVector) perspective {
	return perspective{
		source,
		vector.Source,
		destination,
		vector.Destination,
		outbound,
		func(sg *SecurityGroup) []*SecurityGroupRule { return sg.OutboundRules },
	}
}

func newPerspectiveFromDestination(vector *InterfaceVector) perspective {
	return perspective{
		destination,
		vector.Destination,
		source,
		vector.Source,
		inbound,
		func(sg *SecurityGroup) []*SecurityGroupRule { return sg.InboundRules },
	}
}
