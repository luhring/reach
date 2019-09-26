package aws

import "github.com/luhring/reach/reach"

const FactorKindInstanceState = "InstanceState"

func (i EC2Instance) NewInstanceStateFactor() reach.Factor {
	var tc reach.TrafficContent

	if i.isRunning() {
		tc = reach.NewTrafficContentForAllTraffic()
	} else {
		tc = reach.NewTrafficContentForNoTraffic()
	}

	return reach.Factor{
		Kind:     FactorKindInstanceState,
		Resource: i.ToResourceReference(),
		Traffic:  tc,
	}
}
