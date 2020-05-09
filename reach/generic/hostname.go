package generic

import (
	"net"

	"github.com/luhring/reach/reach"
)

// ResourceKindHostname specifies the unique name for the Hostname kind of resource.
const ResourceKindHostname reach.Kind = "Hostname"

// A Hostname represents a real hostname, used to describe a point in a network for which only the hostname is known.
type Hostname struct {
	Name        string
	IPAddresses []net.IP
}

// Resource returns the Hostname converted to a generalized Reach resource.
func (h Hostname) Resource() reach.Resource {
	return reach.Resource{
		Kind:       ResourceKindHostname,
		Properties: h,
	}
}

// Ref returns a Reference for the Hostname.
func (h Hostname) Ref() reach.Reference {
	return reach.Reference{
		Domain: ResourceDomainGeneric,
		Kind:   ResourceKindHostname,
		ID:     h.Name,
	}
}
