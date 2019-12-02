package analyzer

import (
	"fmt"
	"log"

	"github.com/luhring/reach/reach"
	"github.com/luhring/reach/reach/aws"
	"github.com/luhring/reach/reach/aws/api"
	"github.com/luhring/reach/reach/generic"
)

// Analyzer performs Reach's central network traffic analysis.
type Analyzer struct {
	resourceCollection *reach.ResourceCollection
}

// New creates a new Analyzer that has a new resource collection.
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

					ec2Instance, err := provider.EC2Instance(id)
					if err != nil {
						log.Fatalf("couldn't get resource: %v", err)
					}
					a.resourceCollection.Put(reach.ResourceReference{
						Domain: aws.ResourceDomainAWS,
						Kind:   aws.ResourceKindEC2Instance,
						ID:     ec2Instance.ID,
					}, ec2Instance.ToResource())

					dependencies, err := ec2Instance.Dependencies(provider)
					if err != nil {
						return err
					}
					a.resourceCollection.Merge(dependencies)
				case generic.ResourceDomainGeneric:
					// No resource to add (but this is a supported subject domain).
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

// Analyze performs a full analysis of allowed network traffic among the specified subjects.
func (a *Analyzer) Analyze(subjects ...*reach.Subject) (*reach.Analysis, error) {
	// TODO: Eventually, this dependency wiring should depend on a passed in config.
	var provider aws.ResourceProvider = api.NewResourceProvider()

	err := a.buildResourceCollection(subjects, provider)
	if err != nil {
		return nil, err
	}

	var vectorDiscoverer reach.VectorDiscoverer = NewVectorDiscoverer(a.resourceCollection) // TODO: Consider: We might not need an interface for this, since we might never have multiple implementations

	networkVectors, err := vectorDiscoverer.Discover(subjects)
	if err != nil {
		return nil, err
	}

	processedNetworkVectors := make([]reach.NetworkVector, len(networkVectors))

	// TODO: Eventually, this dependency wiring should depend on a passed in config.
	var vectorAnalyzer reach.VectorAnalyzer = aws.NewVectorAnalyzer(a.resourceCollection)

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

		returnTrafficContents := reach.ReturnTrafficContentsFromFactors(factors)
		returnTrafficContent, err := reach.NewTrafficContentFromIntersectingMultiple(returnTrafficContents)
		if err != nil {
			return nil, err
		}

		processedVector.Traffic = &trafficContent
		processedVector.ReturnTraffic = &returnTrafficContent

		processedNetworkVectors[i] = processedVector
	}

	return reach.NewAnalysis(subjects, a.resourceCollection, processedNetworkVectors), nil
}
