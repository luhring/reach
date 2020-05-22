package standard

import (
	"net"

	"github.com/luhring/reach/reach/generic"
	"github.com/luhring/reach/reach/reacherr"
)

// Hostname uses the standard library method of resolving a hostname given a specified DNS name.
func (provider *DomainClient) Hostname(name string) (*generic.Hostname, error) {
	ips, err := net.LookupIP(name)
	if err != nil {
		return nil, reacherr.New(err, "unable to retrieve hostname resource for name '%s': %v", name, err)
	}

	return &generic.Hostname{
		Name:        name,
		IPAddresses: ips,
	}, nil
}
