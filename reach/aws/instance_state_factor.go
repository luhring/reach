package aws

import "github.com/luhring/reach/reach"

const FactorKindInstanceState = "InstanceState"

func (i EC2Instance) NewInstanceStateFactor() reach.Factor {
	var set []reach.TrafficContent

	if i.isRunning() {
		// TODO: make universal set
	}

	return reach.Factor{
		Kind:              FactorKindInstanceState,
		Resource:          i.ToResourceReference(),
		TrafficContentSet: set,
	}
}
