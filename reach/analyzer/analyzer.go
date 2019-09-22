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
			switch subject.Domain {
			case aws.ResourceDomainAWS:
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
			default:
				return nil, fmt.Errorf("unsupported subject domain: '%s'", subject.Domain)
			}
		}
	}

	var sourceNetworkPoints []reach.NetworkPoint
	var destinationNetworkPoints []reach.NetworkPoint

	for _, subject := range subjects {
		if subject.Role == reach.SubjectRoleSource {
			switch subject.Domain {
			case aws.ResourceDomainAWS:
				switch subject.Kind {
				case aws.SubjectKindEC2Instance:
					ec2Instance := resources[aws.ResourceDomainAWS][aws.ResourceKindEC2Instance][subject.ID].Properties.(aws.EC2Instance)
					sourceNetworkPoints = append(sourceNetworkPoints, ec2Instance.GetNetworkPoints(resources)...)
				}
			}
		} else if subject.Role == reach.SubjectRoleDestination {
			switch subject.Domain {
			case aws.ResourceDomainAWS:
				switch subject.Kind {
				case aws.SubjectKindEC2Instance:
					ec2Instance := resources[aws.ResourceDomainAWS][aws.ResourceKindEC2Instance][subject.ID].Properties.(aws.EC2Instance)
					destinationNetworkPoints = append(destinationNetworkPoints, ec2Instance.GetNetworkPoints(resources)...)
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

	for i, vector := range networkVectors {
		vector = processFactors(vector, resources)
		networkVectors[i] = vector
	}

	return &reach.Analysis{
		Subjects:       subjects,
		Resources:      resources,
		NetworkVectors: networkVectors,
	}, nil
}

func processFactors(vector reach.NetworkVector, resources map[string]map[string]map[string]reach.Resource) reach.NetworkVector {
	for _, ref := range vector.Source.Lineage {
		if ref.Domain == aws.ResourceDomainAWS && ref.Kind == aws.SubjectKindEC2Instance {
			ec2Instance := resources[ref.Domain][ref.Kind][ref.ID].Properties.(aws.EC2Instance)
			vector.Source.Factors = append(vector.Source.Factors, ec2Instance.NewInstanceStateFactor())
		}
	}

	for _, ref := range vector.Destination.Lineage {
		if ref.Domain == aws.ResourceDomainAWS && ref.Kind == aws.SubjectKindEC2Instance {
			ec2Instance := resources[ref.Domain][ref.Kind][ref.ID].Properties.(aws.EC2Instance)
			vector.Destination.Factors = append(vector.Destination.Factors, ec2Instance.NewInstanceStateFactor())
		}
	}

	return vector
}
