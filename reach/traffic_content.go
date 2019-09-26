package reach

import "github.com/luhring/reach/reach/set"

type TrafficContent struct {
	IPProtocol Protocol    `json:"ipProtocol"`
	PortSet    set.PortSet `json:"portSet"`
	ICMPSet    set.ICMPSet `json:"icmpSet"`
}

func NewTrafficContentForAllTraffic() TrafficContent {
	return TrafficContent{
		IPProtocol: ProtocolAll,
	}
}
