package generic

import (
	"fmt"
	"net"

	"github.com/luhring/reach/reach"
)

// PointsDiscoverer is the generic domain's implementation of the PointsDiscoverer interface.
type PointsDiscoverer struct {
	resourceCollection *reach.ResourceCollection
}

// NewPointsDiscoverer creates a new generic domain PointsDiscoverer.
func NewPointsDiscoverer(resourceCollection *reach.ResourceCollection) PointsDiscoverer {
	return PointsDiscoverer{
		resourceCollection,
	}
}

// Discover identifies and returns all network points for the given subject.
func (d PointsDiscoverer) Discover(subject reach.Subject) ([]reach.NetworkPoint, error) {
	if subject.Domain != ResourceDomainGeneric {
		return nil, fmt.Errorf("non-generic domain subject passed to generic dmoain network points discoverer (domain: %s)", subject.Domain)
	}

	switch subject.Kind {
	case SubjectKindIPAddress:
		ip := net.ParseIP(subject.ID)
		if ip == nil {
			return nil, fmt.Errorf("subject of IP address kind had an ID (%s) that could not be parsed as an IP address", subject.ID)
		}

		return []reach.NetworkPoint{
			{
				IPAddress: ip,
				Lineage:   nil,
				Factors:   nil,
			},
		}, nil
	case SubjectKindHostname:
		ips, err := net.LookupIP(subject.ID)
		if err != nil {
			return nil, fmt.Errorf("subject of hostname kind had an ID (%s) that could not be resolved to an IP address: %v", subject.ID, err)
		}

		var points []reach.NetworkPoint

		for _, ip := range ips {
			points = append(points, reach.NetworkPoint{
				IPAddress: ip,
				Lineage:   nil,
				Factors:   nil,
			})
		}

		return points, nil
	default:
		return nil, fmt.Errorf("unsupported generic resource kind passed to generic domain network points discoverer (kind: %s)", subject.Kind)
	}
}
