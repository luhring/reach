package reach

type Protocol int

const (
	ProtocolNone       Protocol = -2
	ProtocolAll        Protocol = -1
	ProtocolICMPv4     Protocol = 1
	ProtocolTCP        Protocol = 6
	ProtocolUDP        Protocol = 17
	ProtocolICMPv6     Protocol = 58
	ProtocolNameAll             = "all"
	ProtocolNameICMP            = "ICMP"
	ProtocolNameTCP             = "TCP"
	ProtocolNameUDP             = "UDP"
	ProtocolNameICMPv6          = "ICMPv6"
)
