package analyzer

import (
	"net"

	"github.com/luhring/reach/reach"
)

type tracerJob struct {
	partial        reach.PartialPath
	destinationIPs []net.IP
}
