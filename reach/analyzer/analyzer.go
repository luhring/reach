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
	providers map[string]interface{}
}

// New creates a new Analyzer that has a new resource collection.
func New(providers map[string]interface{}) *Analyzer {
	return &Analyzer{
		providers: providers,
	}
}

// Analyze performs a full analysis of allowed network traffic among the specified subjects.
func (a *Analyzer) Analyze(subjects ...*reach.Subject) (*reach.Analysis, error) {
	rc, err := a.buildResourceCollection(subjects) // TODO: Need to catch and prevent cycles (e.g. with route table rules having instance targets).
	if err != nil {
		return nil, err
	}

	// I think this becomes a network path constructor
	var vectorDiscoverer reach.VectorDiscoverer = NewVectorDiscoverer(rc)

	networkVectors, err := vectorDiscoverer.Discover(subjects)
	if err != nil {
		return nil, err
	}

	processedNetworkVectors := make([]reach.NetworkVector, len(networkVectors))

	// TODO: Eventually, this dependency wiring should depend on a passed in config.
	// VectorAnalyzer is no more. Vectors become network paths.
	var vectorAnalyzer reach.VectorAnalyzer = aws.NewVectorAnalyzer(rc)

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

	return reach.NewAnalysis(subjects, rc, processedNetworkVectors), nil
}

func (a *Analyzer) buildResourceCollection(subjects []*reach.Subject) (*reach.ResourceCollection, error) { // TODO: Allow passing any number of providers of various domains
	rc := reach.NewResourceCollection()

	for _, subject := range subjects {
		if subject.Role != reach.SubjectRoleNone && subject.Domain != generic.ResourceDomainGeneric { // For the generic domain, there are no resources to obtain.
			switch subject.Domain {
			case aws.ResourceDomainAWS:
				provider := a.providers[aws.ResourceDomainAWS].(aws.ResourceProvider)

				switch subject.Kind { // An argument could be made that this logic should be pushed into the 'aws' package...
				case aws.SubjectKindEC2Instance:
					id := subject.ID

					ec2Instance, err := provider.EC2Instance(id)
					if err != nil {
						log.Fatalf("couldn't get EC2 instance resource: %v", err)
					}
					rc.Put(reach.ResourceReference{
						Domain: aws.ResourceDomainAWS,
						Kind:   aws.ResourceKindEC2Instance,
						ID:     ec2Instance.ID,
					}, ec2Instance.ToResource())

					dependencies, err := ec2Instance.Dependencies(provider)
					if err != nil {
						return nil, err
					}
					rc.Merge(dependencies)
				default:
					return nil, fmt.Errorf("unsupported subject kind: '%s'", subject.Kind)
				}
			case generic.ResourceDomainGeneric:
				provider := a.providers[generic.ResourceDomainGeneric].(generic.ResourceProvider)

				switch subject.Kind {
				case generic.ResourceKindHostname:
					h, err := provider.Hostname(subject.ID)
					if err != nil {
						log.Fatalf("couldn't get hostname resource: %v", err)
					}

					rc.Put(reach.ResourceReference{
						Domain: generic.ResourceDomainGeneric,
						Kind:   generic.ResourceKindHostname,
						ID:     h.Name,
					}, h.ToResource())
				case generic.SubjectKindIPAddress:
					// This is a special case. We don't create resources for IP addresses, but this is a recognized subject kind.
				default:
					return nil, fmt.Errorf("unsupported subject kind: '%s'", subject.Kind)
				}
			default:
				return nil, fmt.Errorf("unsupported subject domain: '%s'", subject.Domain)
			}
		}
	}

	return rc, nil
}
