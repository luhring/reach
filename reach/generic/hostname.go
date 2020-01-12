package generic

import (
	"net"

	"github.com/luhring/reach/reach"
)

const ResourceKindHostname = "Hostname"

type Hostname struct {
	Name        string
	IPAddresses []net.IP
}

func NewHostname(name string, ipAddresses []net.IP) Hostname {
	return Hostname{
		Name:        name,
		IPAddresses: ipAddresses,
	}
}

func (h Hostname) ToResource() reach.Resource {
	return reach.Resource{
		Kind:       ResourceKindHostname,
		Properties: h,
	}
}

func (h Hostname) ToResourceReference() reach.ResourceReference {
	return reach.ResourceReference{
		Domain: ResourceDomainGeneric,
		Kind:   ResourceKindHostname,
		ID:     h.Name,
	}
}

func (h Hostname) networkPoints() []reach.NetworkPoint {
	var points []reach.NetworkPoint

	for _, ip := range h.IPAddresses {
		points = append(points, reach.NetworkPoint{
			IPAddress: ip,
			Lineage: []reach.ResourceReference{
				h.ToResourceReference(),
			},
			Factors: nil,
		})
	}

	return points
}
