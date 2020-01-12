package aws

import (
	"fmt"

	"github.com/luhring/reach/reach"
)

// PointsDiscoverer is the AWS-specific implementation of the PointsDiscoverer interface.
type PointsDiscoverer struct {
	resourceCollection *reach.ResourceCollection
}

// NewPointsDiscoverer creates a new AWS-specific PointsDiscoverer.
func NewPointsDiscoverer(resourceCollection *reach.ResourceCollection) PointsDiscoverer {
	return PointsDiscoverer{
		resourceCollection,
	}
}

// Discover identifies and returns all network points for the given subject.
func (d PointsDiscoverer) Discover(subject reach.Subject) ([]reach.NetworkPoint, error) {
	if subject.Domain != ResourceDomainAWS {
		return nil, fmt.Errorf("non-AWS domain subject passed to AWS-specific network points discoverer (domain: %s)", subject.Domain)
	}

	switch subject.Kind {
	case SubjectKindEC2Instance:
		ec2InstanceResource := d.resourceCollection.Get(reach.ResourceReference{
			Domain: ResourceDomainAWS,
			Kind:   ResourceKindEC2Instance,
			ID:     subject.ID,
		})

		if ec2InstanceResource == nil {
			return nil, fmt.Errorf("resource collection lookup didn't return a resource for %s subject with ID '%s'", SubjectKindEC2Instance, subject.ID)
		}

		ec2Instance := ec2InstanceResource.Properties.(EC2Instance)
		return ec2Instance.networkPoints(d.resourceCollection), nil
	default:
		return nil, fmt.Errorf("unsupported AWS resource kind passed to AWS-specific network points discoverer (kind: %s)", subject.Kind)
	}
}
