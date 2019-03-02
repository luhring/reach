package reach

import (
	"github.com/luhring/reach/network"
	"github.com/mgutz/ansi"
)

type InstanceVector struct {
	Source      *EC2Instance
	Destination *EC2Instance
	PortRange   *network.PortRange
}

func (v *InstanceVector) analyzeInstanceStates() (bool, Explanation) {
	var explanation Explanation

	explanation.AddLineFormat("%v analysis", ansi.Color("instance state", "default+b"))

	isSourceRunning, sourceExplanation := v.Source.analyzeState("source")
	isDestinationRunning, destinationExplanation := v.Destination.analyzeState("destination")

	doStatesAllowTraffic := isSourceRunning && isDestinationRunning

	explanation.Subsume(sourceExplanation)
	explanation.Subsume(destinationExplanation)

	return doStatesAllowTraffic, explanation
}
