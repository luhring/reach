package reach

import "net"

type IPAddressable interface {
	IPs(provider InfrastructureGetter) ([]net.IP, error)
	InterfaceIPs(provider InfrastructureGetter) ([]net.IP, error)
}
