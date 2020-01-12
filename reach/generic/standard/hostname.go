package standard

import (
	"fmt"
	"net"

	"github.com/luhring/reach/reach/generic"
)

// Hostname uses the standard library method of resolving a hostname given a specified DNS name.
func (provider *ResourceProvider) Hostname(name string) (*generic.Hostname, error) {
	ips, err := net.LookupIP(name)
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve hostname resource for name '%s': %v", name, err)
	}

	return &generic.Hostname{
		Name:        name,
		IPAddresses: ips,
	}, nil
}
