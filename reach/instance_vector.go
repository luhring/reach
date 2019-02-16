package reach

import "github.com/luhring/reach/network"

type InstanceVector struct {
	Source      *EC2Instance
	Destination *EC2Instance
	PortRange   *network.PortRange
}
