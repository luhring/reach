package reach

import (
	"net"
	"strings"
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

func (point NetworkPoint) String() string {
	var generations []string

	for i := len(point.Lineage) - 1; i >= 0; i -= 1 {
		generations = append(generations, point.Lineage[i].ID)
	}

	generations = append(generations, point.IPAddress.String())

	return strings.Join(generations, " -> ")
}
