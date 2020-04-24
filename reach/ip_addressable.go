package reach

import "net"

type IPAddressable interface {
	IPs(resolver DomainClientResolver) ([]net.IP, error)
	InterfaceIPs(resolver DomainClientResolver) ([]net.IP, error)
}
