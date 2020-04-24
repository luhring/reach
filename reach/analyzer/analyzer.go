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
	infrastructure reach.InfrastructureGetter
	domains        reach.DomainProvider
}

// New creates a new Analyzer.
func New(
	infrastructure reach.InfrastructureGetter,
	domains reach.DomainProvider,
) *Analyzer {
	return &Analyzer{
		infrastructure: infrastructure,
		domains:        domains,
	}
}

// Analyze performs a full analysis of allowed network traffic among the specified subjects.
func (a *Analyzer) Analyze(source, destination reach.Subject) (*reach.Analysis, error) {
	rc, err := a.buildResourceCollection([]*reach.Subject{&source, &destination}) // TODO: Need to catch and prevent cycles (e.g. with route table rules having instance targets).
	if err != nil {
		return nil, err
	}

	var tracer reach.Tracer = NewTracer(a.infrastructure, a.domains)
	paths, err := tracer.Trace(source, destination)
	if err != nil {
		return nil, fmt.Errorf("unable to complete trace: %v", err)
	}

	var analyzedPaths []analyzedPath

	var pa interface{} // path analyzer

	for _, path := range paths {
		ft := pa.ForwardTraffic(path)
		rt := pa.ReturnTraffic(path)

		analyzedPaths = append(analyzedPaths, analyzedPath{
			path:           &path,
			forwardTraffic: ft,
			returnTraffic:  rt,
		})
	}

	return reach.NewAnalysis(
		[]reach.Subject{source, destination},
		rc,
		paths,
	), nil
}

func (a *Analyzer) buildResourceCollection(subjects []*reach.Subject) (*reach.ResourceCollection, error) { // TODO: Allow passing any number of providers of various domains
	rc := reach.NewResourceCollection()

	for _, subject := range subjects {
		if subject.Role != reach.SubjectRoleNone && subject.Domain != generic.ResourceDomainGeneric { // For the generic domain, there are no resources to obtain.
			switch subject.Domain {
			case aws.ResourceDomainAWS:
				provider := a.providers[aws.ResourceDomainAWS].(aws.ResourceGetter)

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
					}, ec2Instance.Resource())

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

func (a *Analyzer) determineSubjectPoints(role reach.SubjectRole, subjects []*reach.Subject, rc *reach.ResourceCollection) ([]reach.SubjectPoint, error) {
	var points []reach.SubjectPoint
	var g reach.SubjectPointsGenerator

	for _, subject := range subjects {
		if subject.Role == role {
			switch subject.Domain {
			case aws.ResourceDomainAWS:
				g = aws.NewSubjectPointsGenerator(rc)
			default:
				return nil, fmt.Errorf("unable to determine subject points for subject with domain '%s'", subject.Domain)
			}

			newPoints, err := g.SubjectPoints(*subject)
			if err != nil {
				return nil, err
			}

			points = append(points, newPoints...)
		}
	}

	return points, nil
}

func (a *Analyzer) analyzePath(p reach.Path) AnalyzedPath {

}
