package reach

import "net"

type NetworkPoint struct {
	IPAddress net.IP              `json:"ipAddress"`
	Lineage   []ResourceReference `json:"lineage"`
}
