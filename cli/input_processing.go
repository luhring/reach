package cli

import (
	"fmt"
	"os"
	"strings"

	"github.com/luhring/cnct/aws"
)

// GetInstancePairFromIdentificationInputs ...
func GetInstancePairFromIdentificationInputs(
	instanceIdentificationInputs *[2]string,
	instanceCollection []*aws.Instance,
) *aws.InstancePair {
	return &aws.InstancePair{
		getInstanceThatMatchesInput(
			instanceIdentificationInputs[0],
			instanceCollection,
		),
		getInstanceThatMatchesInput(
			instanceIdentificationInputs[1],
			instanceCollection,
		),
	}
}

// GetInstanceIdentificationInputsFromCli ...
func GetInstanceIdentificationInputsFromCli() [2]string {
	if len(os.Args) != 3 {
		fmt.Println("Invalid number of arguments supplied.\n\nExpected syntax:\ncnct instanceA instanceB")
		os.Exit(1)
	}

	return [2]string{
		os.Args[1],
		os.Args[2],
	}
}

func getInstanceThatMatchesInput(
	input string,
	allInstances []*aws.Instance,
) *aws.Instance {
	var indicesOfInstanceIDSubstringMatches []int
	var indicesOfInstanceNameTagExactMatches []int
	var indicesOfInstanceNameTagSubstringMatches []int
	const noMatchingIndices = -1
	const instanceIDPrefix = "i-"

	matchingIndex := noMatchingIndices

	for index, instance := range allInstances {
		if len(input) >= 3 && strings.EqualFold(input[:2], instanceIDPrefix) {
			if strings.EqualFold(input, instance.ID) {
				// exact match to instance ID (using "i-a1b2c3..." format)
				return instance
			}

			if doesFirstItemMatchBeginningSubstringOfSecondItem(input, instance.ID) {
				indicesOfInstanceIDSubstringMatches =
					append(indicesOfInstanceIDSubstringMatches, index)
			}
		} else if strings.EqualFold(input, instance.NameTag) {
			indicesOfInstanceNameTagExactMatches =
				append(indicesOfInstanceNameTagExactMatches, index)
		} else if doesFirstItemMatchBeginningSubstringOfSecondItem(input, instance.NameTag) {
			indicesOfInstanceNameTagSubstringMatches = append(indicesOfInstanceNameTagSubstringMatches, index)
		}
	}

	countOfInstancesWithIDSubstringMatches := len(indicesOfInstanceIDSubstringMatches)
	countOfInstancesWithNameTagExactMatches := len(indicesOfInstanceNameTagExactMatches)
	countOfInstancesWithNameTagSubstringMatches := len(indicesOfInstanceNameTagSubstringMatches)

	if countOfInstancesWithIDSubstringMatches == 1 {
		matchingIndex = indicesOfInstanceIDSubstringMatches[0]
	} else if countOfInstancesWithNameTagExactMatches == 1 {
		matchingIndex = indicesOfInstanceNameTagExactMatches[0]
	} else if countOfInstancesWithNameTagSubstringMatches == 1 {
		matchingIndex = indicesOfInstanceNameTagSubstringMatches[0]
	}

	if matchingIndex != noMatchingIndices {
		return allInstances[matchingIndex]
	}

	if countOfInstancesWithIDSubstringMatches > 1 {
		exitWithErrorMatchingInputToInstance(
			fmt.Sprintf("The input, '%s', was found the IDs of more than one instance and thus could not be used to uniquely identify an instance.", input),
		)
	}

	if countOfInstancesWithNameTagExactMatches > 1 {
		exitWithErrorMatchingInputToInstance(
			fmt.Sprintf("The input, '%s', exactly matched the 'Name' tags of more than one instance, and thus could not be used to uniquely identify an instance.", input),
		)
	}

	if countOfInstancesWithNameTagSubstringMatches > 1 {
		exitWithErrorMatchingInputToInstance(
			fmt.Sprintf("The input, '%s', was found in the 'Name' tags of more than one instance and thus could not be used to uniquely identify an instance.", input),
		)
	}

	exitWithErrorMatchingInputToInstance(
		fmt.Sprintf("No instances found in EC2 that match the input, '%s'.", input),
	)

	return &aws.Instance{}
}

func exitWithErrorMatchingInputToInstance(errorMessage string) {
	fmt.Println(errorMessage)
	fmt.Println("Unable to determine network accessibility because we were unable to find the instance(s) to which you're referring.")
	os.Exit(1)
}

func doesFirstItemMatchBeginningSubstringOfSecondItem(firstItem string, secondItem string) bool {
	lengthOfFirstItem := len(firstItem)

	if lengthOfFirstItem > len(secondItem) {
		return false
	}

	truncatedSecondItem := secondItem[:lengthOfFirstItem]
	return strings.EqualFold(firstItem, truncatedSecondItem)
}
