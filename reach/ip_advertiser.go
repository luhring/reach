package reach

import "net"

type IPAdvertiser interface {
	IPs(provider InfrastructureGetter) ([]net.IP, error)
}
