package aws

import "github.com/luhring/reach/reach"

// VectorDiscoverer is the AWS-specific implementation of the VectorDiscoverer interface.
type VectorDiscoverer struct {
	resourceCollection *reach.ResourceCollection
}

// NewVectorDiscoverer creates a new AWS-specific VectorDiscoverer.
func NewVectorDiscoverer(resourceCollection *reach.ResourceCollection) VectorDiscoverer {
	return VectorDiscoverer{
		resourceCollection,
	}
}

// Discover identifies all of the network vectors that could exist between the given subjects.
func (d VectorDiscoverer) Discover(subjects []*reach.Subject) ([]reach.NetworkVector, error) {
	// TODO: Re-evaluate: As non-AWS network points are introduced, we may need to rethink how we divvy up this logic

	var sourceNetworkPoints []reach.NetworkPoint
	var destinationNetworkPoints []reach.NetworkPoint

	for _, subject := range subjects {
		if subject.Role == reach.SubjectRoleSource {
			switch subject.Domain {
			case ResourceDomainAWS:
				switch subject.Kind {
				case SubjectKindEC2Instance:
					ec2Instance := d.resourceCollection.Get(reach.ResourceReference{
						Domain: ResourceDomainAWS,
						Kind:   ResourceKindEC2Instance,
						ID:     subject.ID,
					}).Properties.(EC2Instance)

					sourceNetworkPoints = append(sourceNetworkPoints, ec2Instance.networkPoints(d.resourceCollection)...)
				}
			}
		} else if subject.Role == reach.SubjectRoleDestination {
			switch subject.Domain {
			case ResourceDomainAWS:
				switch subject.Kind {
				case SubjectKindEC2Instance:
					ec2Instance := d.resourceCollection.Get(reach.ResourceReference{
						Domain: ResourceDomainAWS,
						Kind:   ResourceKindEC2Instance,
						ID:     subject.ID,
					}).Properties.(EC2Instance)

					destinationNetworkPoints = append(destinationNetworkPoints, ec2Instance.networkPoints(d.resourceCollection)...)
				}
			}
		}
	}

	var networkVectors []reach.NetworkVector

	for _, source := range sourceNetworkPoints {
		for _, destination := range destinationNetworkPoints {
			vector, err := reach.NewNetworkVector(source, destination)
			if err != nil {
				return nil, err
			}

			networkVectors = append(networkVectors, vector)
		}
	}

	return networkVectors, nil
}
