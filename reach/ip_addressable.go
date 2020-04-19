package reach

import "net"

type IPAddressable interface {
	IPs(domains DomainProvider) ([]net.IP, error)
	InterfaceIPs(domains DomainProvider) ([]net.IP, error)
}
