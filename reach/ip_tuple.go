package reach

import (
	"fmt"
	"net"
)

// An IPTuple represents the src and dst values found in IP packet metadata.
type IPTuple struct {
	Src net.IP
	Dst net.IP
}

// String returns the string representation of the IPTuple.
func (t IPTuple) String() string {
	return fmt.Sprintf("[src: %s, dst: %s]", t.Src, t.Dst)
}
