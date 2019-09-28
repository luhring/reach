package analyzer

import (
	"fmt"
	"log"

	"github.com/luhring/reach/reach"
	"github.com/luhring/reach/reach/aws"
	"github.com/luhring/reach/reach/aws/api"
)

type Analyzer struct {
	resourceCollection *reach.ResourceCollection
}

func New() *Analyzer {
	rc := reach.NewResourceCollection()
	return &Analyzer{
		resourceCollection: rc,
	}
}

func (a *Analyzer) buildResourceCollection(subjects []*reach.Subject, provider aws.ResourceProvider) error { // TODO: Allow passing any number of providers of various domains
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
					a.resourceCollection.Put(reach.ResourceReference{
						Domain: aws.ResourceDomainAWS,
						Kind:   aws.ResourceKindEC2Instance,
						ID:     ec2Instance.ID,
					}, ec2Instance.ToResource())

					dependencies, err := ec2Instance.GetDependencies(provider)
					if err != nil {
						return err
					}
					a.resourceCollection.Merge(dependencies)
				default:
					return fmt.Errorf("unsupported subject kind: '%s'", subject.Kind)
				}
			default:
				return fmt.Errorf("unsupported subject domain: '%s'", subject.Domain)
			}
		}
	}

	return nil
}

func (a *Analyzer) identifyNetworkVectors(subjects []*reach.Subject) ([]reach.NetworkVector, error) {
	var sourceNetworkPoints []reach.NetworkPoint
	var destinationNetworkPoints []reach.NetworkPoint

	for _, subject := range subjects {
		if subject.Role == reach.SubjectRoleSource {
			switch subject.Domain {
			case aws.ResourceDomainAWS:
				switch subject.Kind {
				case aws.SubjectKindEC2Instance:
					ec2Instance := a.resourceCollection.Get(reach.ResourceReference{
						Domain: aws.ResourceDomainAWS,
						Kind:   aws.ResourceKindEC2Instance,
						ID:     subject.ID,
					}).Properties.(aws.EC2Instance)

					sourceNetworkPoints = append(sourceNetworkPoints, ec2Instance.GetNetworkPoints(a.resourceCollection)...)
				}
			}
		} else if subject.Role == reach.SubjectRoleDestination {
			switch subject.Domain {
			case aws.ResourceDomainAWS:
				switch subject.Kind {
				case aws.SubjectKindEC2Instance:
					ec2Instance := a.resourceCollection.Get(reach.ResourceReference{
						Domain: aws.ResourceDomainAWS,
						Kind:   aws.ResourceKindEC2Instance,
						ID:     subject.ID,
					}).Properties.(aws.EC2Instance)

					destinationNetworkPoints = append(destinationNetworkPoints, ec2Instance.GetNetworkPoints(a.resourceCollection)...)
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

func (a *Analyzer) Analyze(subjects ...*reach.Subject) (*reach.Analysis, error) {
	// TODO: Eventually, this dependency wiring should depend on a passed in config.
	var provider aws.ResourceProvider = api.NewResourceProvider()

	err := a.buildResourceCollection(subjects, provider)
	if err != nil {
		return nil, err
	}

	// TODO: Eventually, this dependency wiring should depend on a passed in config.
	var vectorAnalyzer reach.VectorAnalyzer = aws.NewVectorAnalyzer(a.resourceCollection)

	networkVectors, err := a.identifyNetworkVectors(subjects)
	if err != nil {
		return nil, err
	}

	processedNetworkVectors := make([]reach.NetworkVector, len(networkVectors))

	for i, v := range networkVectors {
		factors, processedVector, err := vectorAnalyzer.Factors(v)
		if err != nil {
			return nil, err
		}

		trafficContents := reach.TrafficContentsFromFactors(factors)

		trafficContent, err := reach.NewTrafficContentFromIntersectingMultiple(trafficContents)
		if err != nil {
			return nil, err
		}

		processedVector.Traffic = &trafficContent
		processedNetworkVectors[i] = processedVector
	}

	return reach.NewAnalysis(subjects, a.resourceCollection, processedNetworkVectors), nil
}
