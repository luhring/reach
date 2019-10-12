package aws

import (
	"fmt"
	"strings"
)

func FindEC2InstanceID(searchText string, provider ResourceProvider) (string, error) {
	instances, err := provider.GetAllEC2Instances()
	if err != nil {
		return "", err
	}

	var matchesOnID []int
	var matchesOnName []int

	// discover what matches exist... and an exact match on instance ID can return early.

	for i, instance := range instances {
		if isInstanceID(searchText) {
			if strings.EqualFold(searchText, instance.ID) { // exact match -- instance ID
				// no need to examine any more instances
				return instance.ID, nil
			}

			if strings.HasPrefix(instance.ID, searchText) { // partial match -- instance ID
				matchesOnID = append(matchesOnID, i)
			}
		}

		if strings.HasPrefix(instance.NameTag, searchText) { // partial or exact match -- instance name
			matchesOnName = append(matchesOnName, i)
		}
	}

	// first priority goes to partial match on instance ID

	if matchesOnID != nil {
		if len(matchesOnID) == 1 {
			return instances[matchesOnID[0]].ID, nil
		}

		if len(matchesOnID) >= 2 {
			var ids []string

			for _, matchIdx := range matchesOnID {
				ids = append(ids, instances[matchIdx].ID)
			}

			matches := strings.Join(ids, ", ")
			return "", fmt.Errorf("search text matches multiple EC2 instances' IDs (matches for search text '%s': %s)", searchText, matches)
		}
	}

	// next, we hope for a match against name by only one instance (partial or exact)

	if matchesOnName != nil {
		if len(matchesOnName) == 1 {
			return instances[matchesOnName[0]].ID, nil
		}

		if len(matchesOnName) >= 2 {
			// prepare helpful error text
			var matchedInstances []string

			for _, matchIdx := range matchesOnID {
				name := instances[matchIdx].NameTag
				id := instances[matchIdx].ID

				matchedInstances = append(matchedInstances, fmt.Sprintf("'%s' (%s)", name, id))
			}

			matches := strings.Join(matchedInstances, ", ")
			return "", fmt.Errorf("search text matches multiple EC2 instances' name tags (matches for search text '%s': %s)", searchText, matches)
		}
	}

	return "", fmt.Errorf("search text '%s' did not match the ID or name tag of any EC2 instances", searchText)
}

func isInstanceID(text string) bool {
	const instanceIDPrefix = "i-"
	return len(text) >= 3 && strings.HasPrefix(text, instanceIDPrefix)
}
