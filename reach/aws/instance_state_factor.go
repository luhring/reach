package aws

import (
	"fmt"

	"github.com/luhring/reach/reach"
)

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

func (i EC2Instance) ExplainInstanceState() string {
	output := fmt.Sprintf("EC2 instance state is '%s'", i.State)

	return output
}
