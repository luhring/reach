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

func (h Hostname) Resource() reach.Resource {
	return reach.Resource{
		Kind:       ResourceKindHostname,
		Properties: h,
	}
}

func (h Hostname) Ref() reach.UniversalReference {
	return reach.UniversalReference{
		Domain: ResourceDomainGeneric,
		Kind:   ResourceKindHostname,
		ID:     h.Name,
	}
}
