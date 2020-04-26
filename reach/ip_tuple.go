package reach

import (
	"fmt"
	"net"
)

type IPTuple struct {
	Src net.IP
	Dst net.IP
}

func (t IPTuple) String() string {
	return fmt.Sprintf("[src: %s, dst: %s]", t.Src, t.Dst)
}
