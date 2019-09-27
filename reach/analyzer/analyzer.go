package analyzer

import (
	"fmt"
	"log"

	"github.com/luhring/reach/reach"
	"github.com/luhring/reach/reach/aws"
	"github.com/luhring/reach/reach/aws/api"
)

type Analyzer struct {
	collection *reach.ResourceCollection
}

func New() *Analyzer {
	rc := reach.NewResourceCollection()
	return &Analyzer{
		collection: rc,
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
					a.collection.Put(reach.ResourceReference{
						Domain: aws.ResourceDomainAWS,
						Kind:   aws.ResourceKindEC2Instance,
						ID:     ec2Instance.ID,
					}, ec2Instance.ToResource())

					dependencies, err := ec2Instance.GetDependencies(provider)
					if err != nil {
						return err
					}
					a.collection.Merge(dependencies)
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
					ec2Instance := a.collection.Get(reach.ResourceReference{
						Domain: aws.ResourceDomainAWS,
						Kind:   aws.ResourceKindEC2Instance,
						ID:     subject.ID,
					}).Properties.(aws.EC2Instance)

					sourceNetworkPoints = append(sourceNetworkPoints, ec2Instance.GetNetworkPoints(a.collection)...)
				}
			}
		} else if subject.Role == reach.SubjectRoleDestination {
			switch subject.Domain {
			case aws.ResourceDomainAWS:
				switch subject.Kind {
				case aws.SubjectKindEC2Instance:
					ec2Instance := a.collection.Get(reach.ResourceReference{
						Domain: aws.ResourceDomainAWS,
						Kind:   aws.ResourceKindEC2Instance,
						ID:     subject.ID,
					}).Properties.(aws.EC2Instance)

					destinationNetworkPoints = append(destinationNetworkPoints, ec2Instance.GetNetworkPoints(a.collection)...)
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

func (a *Analyzer) computeFactorsForSingleVector(vector reach.NetworkVector) (reach.NetworkVector, error) {
	for _, ref := range vector.Source.Lineage {
		if ref.Domain == aws.ResourceDomainAWS {
			if ref.Kind == aws.ResourceKindEC2Instance {
				ec2Instance := a.collection.Get(ref).Properties.(aws.EC2Instance)

				vector.Source.Factors = append(vector.Source.Factors, ec2Instance.NewInstanceStateFactor())
			}

			if ref.Kind == aws.ResourceKindElasticNetworkInterface {
				eni := a.collection.Get(ref).Properties.(aws.ElasticNetworkInterface)

				var targetENI aws.ElasticNetworkInterface

				for _, ancestor := range vector.Destination.Lineage {
					if ancestor.Domain == aws.ResourceDomainAWS && ancestor.Kind == aws.ResourceKindElasticNetworkInterface {
						targetENI = a.collection.Get(ancestor).Properties.(aws.ElasticNetworkInterface)
					}
				}

				securityGroupRulesFactor, err := eni.NewSecurityGroupRulesFactor(
					a.collection,
					func(sg aws.SecurityGroup) []aws.SecurityGroupRule {
						return sg.OutboundRules
					},
					vector.Destination.IPAddress, &targetENI,
				)
				if err != nil {
					return reach.NetworkVector{}, err
				}

				vector.Source.Factors = append(vector.Source.Factors, *securityGroupRulesFactor)
			}
		}
	}

	for _, ref := range vector.Destination.Lineage { // TODO: refactor to extract symmetry
		if ref.Domain == aws.ResourceDomainAWS {
			if ref.Kind == aws.SubjectKindEC2Instance {
				ec2Instance := a.collection.Get(ref).Properties.(aws.EC2Instance)
				vector.Destination.Factors = append(vector.Destination.Factors, ec2Instance.NewInstanceStateFactor())
			}

			if ref.Kind == aws.ResourceKindElasticNetworkInterface {
				eni := a.collection.Get(ref).Properties.(aws.ElasticNetworkInterface)

				var targetENI aws.ElasticNetworkInterface

				for _, ancestor := range vector.Source.Lineage {
					if ancestor.Domain == aws.ResourceDomainAWS && ancestor.Kind == aws.ResourceKindElasticNetworkInterface {
						targetENI = a.collection.Get(ancestor).Properties.(aws.ElasticNetworkInterface)
					}
				}

				securityGroupRulesFactor, err := eni.NewSecurityGroupRulesFactor(
					a.collection,
					func(sg aws.SecurityGroup) []aws.SecurityGroupRule {
						return sg.InboundRules
					},
					vector.Destination.IPAddress, &targetENI,
				)
				if err != nil {
					return reach.NetworkVector{}, err
				}

				vector.Destination.Factors = append(vector.Destination.Factors, *securityGroupRulesFactor)
			}
		}
	}

	return vector, nil
}

func (a *Analyzer) computeFactors(vectors []reach.NetworkVector) ([]reach.NetworkVector, error) {
	var result []reach.NetworkVector

	for _, vector := range vectors { // TODO: Consider concurrent processing
		processedVector, err := a.computeFactorsForSingleVector(vector)
		if err != nil {
			return nil, err
		}

		result = append(result, processedVector)
	}

	return result, nil
}

func (a *Analyzer) netTraffic(vectors []reach.NetworkVector) ([]reach.NetworkVector, error) {
	var result []reach.NetworkVector

	for _, vector := range vectors { // TODO: Consider concurrent processing
		traffic, err := vector.NetTraffic()
		if err != nil {
			return nil, err
		}

		vector.Traffic = traffic
		result = append(result, vector)
	}

	return result, nil
}

func (a *Analyzer) Analyze(subjects ...*reach.Subject) (*reach.Analysis, error) {
	// TODO: Eventually, this dependency wiring should depend on a passed in config.
	var provider aws.ResourceProvider = api.NewResourceProvider()

	err := a.buildResourceCollection(subjects, provider)
	if err != nil {
		return nil, err
	}

	networkVectors, err := a.identifyNetworkVectors(subjects)
	if err != nil {
		return nil, err
	}

	networkVectors, err = a.computeFactors(networkVectors)
	if err != nil {
		return nil, err
	}

	networkVectors, err = a.netTraffic(networkVectors)
	if err != nil {
		return nil, err
	}

	return reach.NewAnalysis(subjects, a.collection, networkVectors), nil
}
