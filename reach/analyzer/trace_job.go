package analyzer

import (
	"net"

	"github.com/luhring/reach/reach"
)

type traceJob struct {
	ref            reach.Reference // The ref point to focus on during the trace
	path           reach.Path      // The state of the path being traced
	edge           reach.Edge      // The edge following the furthest point in the path under construction
	sourceRef      reach.Reference // The ref of the source for the original query
	destinationRef reach.Reference // The ref of the destination for the original query
	destinationIPs []net.IP        // The IP addresses associated with the destination for the original query
}
