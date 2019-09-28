package reach

import (
	"net"
)

type NetworkPoint struct {
	IPAddress net.IP
	Lineage   []ResourceReference
	Factors   []Factor
}

func (point NetworkPoint) TrafficComponents() []TrafficContent {
	var components []TrafficContent

	for _, factor := range point.Factors {
		components = append(components, factor.Traffic)
	}

	return components
}
