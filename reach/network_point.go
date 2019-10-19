package reach

import (
	"net"
	"strings"
)

// A NetworkPoint is a point of termination for an analyzed network vector (on either the source or destination side), such that there is no further subdivision of a source or destination possible beyond the network point. For example, the CIDR block "10.0.1.0/24" contains numerous individual IP addresses, and the analysis result might vary depending on which of these individual IP addresses is used in real network traffic. To break this problem down, such that an analysis result is as definitive as possible, each individual IP address must be analyzed, one at a time. Each IP address could be considered a network point, whereas the CIDR block could not be considered a network point.
type NetworkPoint struct {
	IPAddress net.IP
	Lineage   []ResourceReference
	Factors   []Factor
}

func (point NetworkPoint) trafficContents() []TrafficContent {
	var components []TrafficContent

	for _, factor := range point.Factors {
		components = append(components, factor.Traffic)
	}

	return components
}

// String returns the text representation of the NetworkPoint
func (point NetworkPoint) String() string {
	var generations []string

	for i := len(point.Lineage) - 1; i >= 0; i-- {
		generations = append(generations, point.Lineage[i].ID)
	}

	generations = append(generations, point.IPAddress.String())

	return strings.Join(generations, " -> ")
}
