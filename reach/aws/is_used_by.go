package aws

import "github.com/luhring/reach/reach"

// IsUsedByNetworkPoint returns a boolean indicating whether or not the specified network point contains an AWS-specific kind of resource.
func IsUsedByNetworkPoint(point reach.NetworkPoint) bool {
	return containsAWSResource(point.Lineage)
}

func containsAWSResource(refs []reach.ResourceReference) bool {
	for _, ref := range refs {
		if ref.Domain == ResourceDomainAWS {
			return true
		}
	}

	return false
}
