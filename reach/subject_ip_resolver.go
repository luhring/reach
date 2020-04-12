package reach

import "net"

type SubjectIPResolver interface {
	Resolve(role SubjectRole, provider InfrastructureGetter) ([]net.IP, error)
}
