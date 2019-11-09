package aws

import (
	"github.com/luhring/reach/reach"
)

// FactorKindInstanceState specifies the unique name for the EC2 instance state of factor.
const FactorKindInstanceState = "InstanceState"

func (i EC2Instance) newInstanceStateFactor() reach.Factor {
	var tc reach.TrafficContent

	if i.isRunning() {
		tc = reach.NewTrafficContentForAllTraffic()
	} else {
		tc = reach.NewTrafficContentForNoTraffic()
	}

	return reach.Factor{
		Kind:          FactorKindInstanceState,
		Resource:      i.ToResourceReference(),
		Traffic:       tc,
		ReturnTraffic: reach.NewTrafficContentForAllTraffic(),
	}
}
