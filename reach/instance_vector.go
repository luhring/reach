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

func (instanceVector *InstanceVector) Analyze() Analysis {
	var explanation Explanation

	explanation.AddLineFormat(
		"source instance: %v",
		ansi.Color(instanceVector.Source.LongName(), "default+b"),
	)
	explanation.AddLineFormat(
		"destination instance: %v",
		ansi.Color(instanceVector.Destination.LongName(), "default+b"),
	)
	explanation.AddBlankLine()

	doStatesAllowTraffic, statesExplanation := instanceVector.analyzeInstanceStates()
	explanation.Append(statesExplanation)

	explanation.AddBlankLine()
	explanation.AddLine("source and destination network interface pairings:")

	interfaceVectors := instanceVector.createInterfaceVectors()
	if interfaceVectors == nil {
		var lackOfInterfaceVectors Explanation
		lackOfInterfaceVectors.AddLine(ansi.Color("one or both instances are missing a network interface", "red"))
		explanation.Subsume(lackOfInterfaceVectors)

		return newAnalysisWithNoTrafficAllowances(explanation)
	}

	var allowedTraffic []*network.TrafficAllowance

	for _, v := range interfaceVectors {
		var vectorExplanation Explanation

		vectorExplanation.Append(v.explainSourceAndDestination())
		vectorExplanation.AddBlankLine()

		// Security groups

		reachablePortsViaSecurityGroups, sgExplanation := v.analyzeSecurityGroups()

		if len(reachablePortsViaSecurityGroups) >= 1 {
			allowedTraffic = append(allowedTraffic, reachablePortsViaSecurityGroups...)
		}

		vectorExplanation.Append(sgExplanation)

		// (Other analyses...)

		explanation.Subsume(vectorExplanation)
	}

	allowedTraffic = network.ConsolidateTrafficAllowances(allowedTraffic)

	if doStatesAllowTraffic == false {
		allowedTraffic = []*network.TrafficAllowance{}
	}

	return Analysis{
		allowedTraffic,
		explanation,
	}
}

func (instanceVector *InstanceVector) createInterfaceVectors() []InterfaceVector {
	var interfaceVectors []InterfaceVector = nil

	for _, fromInterface := range instanceVector.Source.NetworkInterfaces {
		for _, toInterface := range instanceVector.Destination.NetworkInterfaces {
			newVector := InterfaceVector{
				Source:      fromInterface,
				Destination: toInterface,
				PortRange:   instanceVector.PortRange,
			}
			interfaceVectors = append(interfaceVectors, newVector)
		}
	}

	return interfaceVectors
}

func (instanceVector *InstanceVector) analyzeInstanceStates() (bool, Explanation) {
	var explanation Explanation

	explanation.AddLineFormat("%v analysis", ansi.Color("instance state", "default+b"))

	isSourceRunning, sourceExplanation := instanceVector.Source.analyzeState("source")
	isDestinationRunning, destinationExplanation := instanceVector.Destination.analyzeState("destination")

	doStatesAllowTraffic := isSourceRunning && isDestinationRunning

	explanation.Subsume(sourceExplanation)
	explanation.Subsume(destinationExplanation)

	return doStatesAllowTraffic, explanation
}
