package reach

import "net"

// IPAddressable is the interface that describes the ability to obtain the associated IP addresses from a point in a network.
type IPAddressable interface {
	IPs(resolver DomainClientResolver) ([]net.IP, error)
	InterfaceIPs(resolver DomainClientResolver) ([]net.IP, error)
}
