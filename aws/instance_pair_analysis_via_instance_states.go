package aws

import "fmt"

type analysisOfInstanceStates struct {
	instanceNames                       [2]string
	doesAnalysisSuggestThatAccessExists bool
	instanceIsRunningBools              [2]bool
	instanceStateStrings                [2]string
}

func (instancePair *InstancePair) analyzeNetworkAccessViaInstanceStates() analysisOfInstanceStates {
	instanceStateStrings := [2]string{
		instancePair[0].State,
		instancePair[1].State,
	}

	instanceIsRunningBools := [2]bool{
		instancePair[0].isRunning(),
		instancePair[1].isRunning(),
	}

	doesAnalysisSuggestThatAccessExists :=
		instanceIsRunningBools[0] && instanceIsRunningBools[1]

	instanceNames := [2]string{
		instancePair[0].GetFriendlyName(),
		instancePair[1].GetFriendlyName(),
	}

	return analysisOfInstanceStates{
		instanceNames,
		doesAnalysisSuggestThatAccessExists,
		instanceIsRunningBools,
		instanceStateStrings,
	}
}

func (analysis *analysisOfInstanceStates) generateExplanationForLackOfAccessDueToInstanceStates() string {
	const requirementMessage = "Both instances need to be running. The following instance(s) are not currently running:"

	var listOfInstancesNotRunning string

	for instanceIndex, isInstanceRunning := range analysis.instanceIsRunningBools {
		if false == isInstanceRunning {
			listOfInstancesNotRunning += fmt.Sprintf(
				"  - %s (%s)\n",
				analysis.instanceNames[instanceIndex],
				analysis.instanceStateStrings[instanceIndex],
			)
		}
	}

	return requirementMessage + "\n" + listOfInstancesNotRunning
}
