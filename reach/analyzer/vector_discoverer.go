package analyzer

import (
	"errors"
	"fmt"

	"github.com/luhring/reach/reach"
	"github.com/luhring/reach/reach/aws"
	"github.com/luhring/reach/reach/generic"
)

// VectorDiscoverer is the domain-agnostic implementation of the VectorDiscoverer interface.
type VectorDiscoverer struct {
	resourceCollection *reach.ResourceCollection
}

func NewVectorDiscoverer(collection *reach.ResourceCollection) VectorDiscoverer {
	return VectorDiscoverer{
		resourceCollection: collection,
	}
}

func (d VectorDiscoverer) Discover(subjects []*reach.Subject) ([]reach.NetworkVector, error) {
	var sourceNetworkPoints []reach.NetworkPoint
	var destinationNetworkPoints []reach.NetworkPoint

	for _, subject := range subjects {
		if subject == nil {
			return nil, errors.New("encountered nil subject during network vector discovery")
		}

		var pointsDiscoverer reach.PointsDiscoverer

		switch subject.Domain {
		case aws.ResourceDomainAWS:
			pointsDiscoverer = aws.NewPointsDiscoverer(d.resourceCollection)
		case generic.ResourceDomainGeneric:
			pointsDiscoverer = generic.NewPointsDiscoverer(d.resourceCollection)
		default:
			return nil, fmt.Errorf("unable to discover points for subject with unrecognized resource domain: %v", subject)
		}

		networkPoints, err := pointsDiscoverer.Discover(*subject)
		if err != nil {
			return nil, fmt.Errorf("error encountered while discovering network vectors: %v", err)
		}

		if subject.Role == reach.SubjectRoleSource {
			sourceNetworkPoints = append(sourceNetworkPoints, networkPoints...)
		} else if subject.Role == reach.SubjectRoleDestination {
			destinationNetworkPoints = append(destinationNetworkPoints, networkPoints...)
		}
	}

	var networkVectors []reach.NetworkVector

	for _, source := range sourceNetworkPoints {
		for _, destination := range destinationNetworkPoints {
			if bothPointsAreGeneric(source, destination) {
				return nil, fmt.Errorf("cannot perform analysis if both points in a network vector are generic (e.g. just an IP address or a hostname)")
			}

			vector, err := reach.NewNetworkVector(source, destination)
			if err != nil {
				return nil, err
			}

			networkVectors = append(networkVectors, vector)
		}
	}

	return networkVectors, nil
}

func bothPointsAreGeneric(source, destination reach.NetworkPoint) bool {
	return generic.IsNetworkPointGeneric(source) && generic.IsNetworkPointGeneric(destination)
}
