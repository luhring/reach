package analyzer

import (
	"fmt"
	"log"

	"github.com/luhring/reach/reach"
	"github.com/luhring/reach/reach/aws"
	"github.com/luhring/reach/reach/generic"
)

// Analyzer performs Reach's central network traffic analysis.
type Analyzer struct {
	providers          map[string]interface{}
	resourceCollection *reach.ResourceCollection
}

// New creates a new Analyzer that has a new resource collection.
func New(providers map[string]interface{}) *Analyzer {
	rc := reach.NewResourceCollection()
	return &Analyzer{
		providers:          providers,
		resourceCollection: rc,
	}
}

// Analyze performs a full analysis of allowed network traffic among the specified subjects.
func (a *Analyzer) Analyze(subjects ...*reach.Subject) (*reach.Analysis, error) {
	err := a.buildResourceCollection(subjects) // TODO: Need to catch and prevent cycles (e.g. with route table rules having instance targets).
	if err != nil {
		return nil, err
	}

	var vectorDiscoverer reach.VectorDiscoverer = NewVectorDiscoverer(a.resourceCollection)

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

func (a *Analyzer) buildResourceCollection(subjects []*reach.Subject) error { // TODO: Allow passing any number of providers of various domains
	for _, subject := range subjects {
		if subject.Role != reach.SubjectRoleNone && subject.Domain != generic.ResourceDomainGeneric { // For the generic domain, there are no resources to obtain.
			switch subject.Domain {
			case aws.ResourceDomainAWS:
				provider := a.providers[aws.ResourceDomainAWS].(aws.ResourceProvider)

				switch subject.Kind {
				case aws.SubjectKindEC2Instance:
					id := subject.ID

					ec2Instance, err := provider.EC2Instance(id)
					if err != nil {
						log.Fatalf("couldn't EC2 instance resource: %v", err)
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
				default:
					return fmt.Errorf("unsupported subject kind: '%s'", subject.Kind)
				}
			case generic.ResourceDomainGeneric:
				provider := a.providers[generic.ResourceDomainGeneric].(generic.ResourceProvider)

				switch subject.Kind {
				case generic.ResourceKindHostname:
					h, err := provider.Hostname(subject.ID)
					if err != nil {
						log.Fatalf("couldn't get hostname resource: %v", err)
					}

					a.resourceCollection.Put(reach.ResourceReference{
						Domain: generic.ResourceDomainGeneric,
						Kind:   generic.ResourceKindHostname,
						ID:     h.Name,
					}, h.ToResource())
				case generic.SubjectKindIPAddress:
					// This is a special case. We don't create resources for IP addresses, but this is a recognized subject kind.
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
