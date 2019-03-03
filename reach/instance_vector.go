package reach

import (
	"fmt"
	"github.com/mgutz/ansi"
)

type InstanceVector struct {
	Source      *EC2Instance
	Destination *EC2Instance
}

func (instanceVector *InstanceVector) Analyze(filter *TrafficAllowance) Analysis {
	explanation := newExplanation(
		fmt.Sprintf("source instance: %v", ansi.Color(instanceVector.Source.LongName(), "default+b")),
		fmt.Sprintf("destination instance: %v", ansi.Color(instanceVector.Destination.LongName(), "default+b")),
		"",
	)

	doStatesAllowTraffic, statesExplanation := instanceVector.analyzeInstanceStates()
	explanation.append(statesExplanation)

	explanation.addBlankLine()
	explanation.addLine("source and destination network interface pairings:")

	interfaceVectors := instanceVector.createInterfaceVectors()
	if interfaceVectors == nil {
		lackOfInterfaceVectors := newExplanation(
			ansi.Color("one or both instances are missing a network interface", "red"),
		)
		explanation.subsume(lackOfInterfaceVectors)

		return newAnalysisWithNoTrafficAllowances(explanation)
	}

	var allowedTraffic []*TrafficAllowance

	for _, v := range interfaceVectors {
		var vectorExplanation Explanation

		vectorExplanation.append(v.explainSourceAndDestination())
		vectorExplanation.addBlankLine()

		// Security groups

		reachablePortsViaSecurityGroups, sgExplanation := v.analyzeSecurityGroups(filter)

		if len(reachablePortsViaSecurityGroups) >= 1 {
			allowedTraffic = append(allowedTraffic, reachablePortsViaSecurityGroups...)
		}

		vectorExplanation.append(sgExplanation)

		// (Other analyses...)

		explanation.subsume(vectorExplanation)
	}

	allowedTraffic = consolidateTrafficAllowances(allowedTraffic)

	if doStatesAllowTraffic == false {
		allowedTraffic = []*TrafficAllowance{}
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
			}
			interfaceVectors = append(interfaceVectors, newVector)
		}
	}

	return interfaceVectors
}

func (instanceVector *InstanceVector) analyzeInstanceStates() (bool, Explanation) {
	explanation := newExplanation(
		fmt.Sprintf("%v analysis", ansi.Color("instance state", "default+b")),
	)

	isSourceRunning, sourceExplanation := instanceVector.Source.analyzeState("source")
	isDestinationRunning, destinationExplanation := instanceVector.Destination.analyzeState("destination")

	doStatesAllowTraffic := isSourceRunning && isDestinationRunning

	explanation.subsume(sourceExplanation)
	explanation.subsume(destinationExplanation)

	return doStatesAllowTraffic, explanation
}
