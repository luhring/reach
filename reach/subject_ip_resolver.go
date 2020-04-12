package reach

import "net"

type SubjectIPResolver interface {
	Resolve(role SubjectRole) []net.IP
}
