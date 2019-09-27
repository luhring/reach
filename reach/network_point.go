package reach

import "net"

type NetworkPoint struct {
	IPAddress net.IP              `json:"ipAddress"`
	Lineage   []ResourceReference `json:"lineage"` // TODO: More idiomatic approach via DAG
	Factors   []Factor            `json:"factors"`
}

func (p NetworkPoint) TrafficComponents() []TrafficContent {
	var components []TrafficContent

	for _, factor := range p.Factors {
		components = append(components, factor.Traffic)
	}

	return components
}
