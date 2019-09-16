package reach

const (
	all                       = -1 // old
	ProtocolAll               = -1
	ProtocolICMP              = 1
	ProtocolTCP               = 6
	ProtocolUDP               = 17
	ProtocolICMPv6            = 58
	allName                   = "all"
	icmpName                  = "ICMP"
	tcpName                   = "TCP"
	udpName                   = "UDP"
	icmpv6Name                = "ICMPv6"
	ipProtocolNumberForICMP   = 1
	ipProtocolNumberForTCP    = 6
	ipProtocolNumberForUDP    = 17
	ipProtocolNumberForICMPv6 = 58
)

type TrafficContent struct {
	IPProtocol int         `json:"ipProtocol"`
	PortSet    interface{} `json:"portSet"`
	ICMPSet    interface{} `json:"icmpSet"`
}

func NewTrafficContentForAllTraffic() TrafficContent {
	return TrafficContent{
		IPProtocol: ProtocolAll,
	}
}
