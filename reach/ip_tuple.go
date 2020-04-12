package reach

import "net"

type IPTuple struct {
	Src net.IP
	Dst net.IP
}
