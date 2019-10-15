package aws

import "github.com/luhring/reach/reach"

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
