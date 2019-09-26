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

	rc := reach.NewResourceCollection()

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
					rc.Put(reach.ResourceReference{
						Domain: aws.ResourceDomainAWS,
						Kind:   aws.ResourceKindEC2Instance,
						ID:     ec2Instance.ID,
					}, ec2Instance.ToResource())

					dependencies, err := ec2Instance.GetDependencies(provider)
					if err != nil {
						return nil, err
					}
					rc.Merge(dependencies)
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
					ec2Instance := rc.Get(reach.ResourceReference{
						Domain: aws.ResourceDomainAWS,
						Kind:   aws.ResourceKindEC2Instance,
						ID:     subject.ID,
					}).Properties.(aws.EC2Instance)

					sourceNetworkPoints = append(sourceNetworkPoints, ec2Instance.GetNetworkPoints(rc)...)
				}
			}
		} else if subject.Role == reach.SubjectRoleDestination {
			switch subject.Domain {
			case aws.ResourceDomainAWS:
				switch subject.Kind {
				case aws.SubjectKindEC2Instance:
					ec2Instance := rc.Get(reach.ResourceReference{
						Domain: aws.ResourceDomainAWS,
						Kind:   aws.ResourceKindEC2Instance,
						ID:     subject.ID,
					}).Properties.(aws.EC2Instance)

					destinationNetworkPoints = append(destinationNetworkPoints, ec2Instance.GetNetworkPoints(rc)...)
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
		vector = processFactors(vector, rc)
		networkVectors[i] = vector
	}

	return &reach.Analysis{
		Subjects:       subjects,
		Resources:      rc,
		NetworkVectors: networkVectors,
	}, nil
}

func processFactors(vector reach.NetworkVector, rc *reach.ResourceCollection) reach.NetworkVector {
	for _, ref := range vector.Source.Lineage {
		if ref.Domain == aws.ResourceDomainAWS && ref.Kind == aws.SubjectKindEC2Instance {
			ec2Instance := rc.Get(ref).Properties.(aws.EC2Instance)
			vector.Source.Factors = append(vector.Source.Factors, ec2Instance.NewInstanceStateFactor())
		}
	}

	for _, ref := range vector.Destination.Lineage {
		if ref.Domain == aws.ResourceDomainAWS && ref.Kind == aws.SubjectKindEC2Instance {
			ec2Instance := rc.Get(ref).Properties.(aws.EC2Instance)
			vector.Destination.Factors = append(vector.Destination.Factors, ec2Instance.NewInstanceStateFactor())
		}
	}

	return vector
}
