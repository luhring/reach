package analyzer

import (
	"fmt"
	"log"

	"github.com/luhring/reach/reach"
	"github.com/luhring/reach/reach/aws"
	"github.com/luhring/reach/reach/aws/api"
)

type Analyzer struct {
}

func New() *Analyzer {
	return &Analyzer{}
}

func (a *Analyzer) Analyze(subjects ...*reach.Subject) (*reach.Analysis, error) {
	// TODO: Eventually, this dependency wiring should depend on a passed in config.
	provider := api.NewResourceProvider()

	resources := make(map[string]map[string]map[string]reach.Resource)

	for _, subject := range subjects {
		if subject.Role != reach.SubjectRoleNone {
			switch subject.Kind {
			case aws.SubjectKindEC2Instance:
				id := subject.ID

				ec2Instance, err := provider.GetEC2Instance(id)
				if err != nil {
					log.Fatalf("couldn't get resource: %v", err)
				}
				resources = reach.EnsureResourcePathExists(resources, aws.ResourceDomainAWS, aws.ResourceKindEC2Instance)
				resources[aws.ResourceDomainAWS][aws.ResourceKindEC2Instance][ec2Instance.ID] = ec2Instance.ToResource()

				dependencies, err := ec2Instance.GetDependencies(provider)
				if err != nil {
					return nil, err
				}
				resources = reach.MergeResources(resources, dependencies)
			default:
				return nil, fmt.Errorf("unsupported subject kind: '%s'", subject.Kind)
			}
		}
	}

	var sourceNetworkPoints []reach.NetworkPoint
	var destinationNetworkPoints []reach.NetworkPoint

	for _, subject := range subjects {
		if subject.Role == reach.SubjectRoleSource {
			switch subject.Kind {
			case aws.SubjectKindEC2Instance:
				ec2Instance := resources[aws.ResourceDomainAWS][aws.ResourceKindEC2Instance][subject.ID].Properties.(aws.EC2Instance)
				sourceNetworkPoints = append(sourceNetworkPoints, ec2Instance.GetNetworkPoints(resources)...)
			}
		} else if subject.Role == reach.SubjectRoleDestination {
			switch subject.Kind {
			case aws.SubjectKindEC2Instance:
				ec2Instance := resources[aws.ResourceDomainAWS][aws.ResourceKindEC2Instance][subject.ID].Properties.(aws.EC2Instance)
				destinationNetworkPoints = append(destinationNetworkPoints, ec2Instance.GetNetworkPoints(resources)...)
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

	// From subjects, generate network points.

	// From all network points, generate network vectors.

	// For each network vector, identify analysis type (e.g. ENI-to-ENI), and generate factors.

	return &reach.Analysis{
		Subjects:       subjects,
		Resources:      resources,
		NetworkVectors: networkVectors,
	}, nil
}
