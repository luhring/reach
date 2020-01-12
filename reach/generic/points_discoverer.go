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
		return nil, fmt.Errorf("non-generic domain subject passed to generic domain network points discoverer (domain: %s)", subject.Domain)
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
		hostnameResource := d.resourceCollection.Get(reach.ResourceReference{
			Domain: ResourceDomainGeneric,
			Kind:   ResourceKindHostname,
			ID:     subject.ID,
		})

		if hostnameResource == nil {
			return nil, fmt.Errorf("resource collection lookup didn't return a resource for %s subject with ID '%s'", SubjectKindHostname, subject.ID)
		}

		h := hostnameResource.Properties.(Hostname)
		return h.networkPoints(), nil
	default:
		return nil, fmt.Errorf("unsupported generic resource kind passed to generic domain network points discoverer (kind: %s)", subject.Kind)
	}
}

func IsNetworkPointGeneric(point reach.NetworkPoint) bool {
	return point.Lineage == nil || len(point.Lineage) == 0 || point.Lineage[0].Domain == ResourceDomainGeneric
}
