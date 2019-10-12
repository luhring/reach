package aws

import (
	"fmt"

	"github.com/luhring/reach/reach"
)

type VectorAnalyzer struct {
	resourceCollection *reach.ResourceCollection
}

func NewVectorAnalyzer(resourceCollection *reach.ResourceCollection) VectorAnalyzer {
	return VectorAnalyzer{
		resourceCollection,
	}
}

func (analyzer VectorAnalyzer) FactorsForPerspective(p AnalysisPerspective) ([]reach.Factor, error) {
	var factors []reach.Factor

	for _, resourceRef := range p.self.Lineage {
		if resourceRef.Domain == ResourceDomainAWS {
			if resourceRef.Kind == ResourceKindEC2Instance {
				ec2Instance := analyzer.resourceCollection.Get(resourceRef).Properties.(EC2Instance)

				factors = append(factors, ec2Instance.NewInstanceStateFactor())
			}

			if resourceRef.Kind == ResourceKindElasticNetworkInterface {
				eni := analyzer.resourceCollection.Get(resourceRef).Properties.(ElasticNetworkInterface)
				targetENI := ElasticNetworkInterfaceFromNetworkPoint(p.other, analyzer.resourceCollection)

				securityGroupRulesFactor, err := eni.NewSecurityGroupRulesFactor(
					analyzer.resourceCollection,
					p,
					targetENI,
				)
				if err != nil {
					return nil, err
				}

				factors = append(factors, *securityGroupRulesFactor)

				if !sameSubnet(&eni, targetENI) {
					return nil, fmt.Errorf("unable to analyze without two EC2 instances existing in the same subnet (source subnet: %s, destination subnet: %s)", eni.SubnetID, targetENI.SubnetID)
				}
			}
		}
	}

	return factors, nil
}

func (analyzer VectorAnalyzer) Factors(v reach.NetworkVector) ([]reach.Factor, reach.NetworkVector, error) {
	var factors []reach.Factor

	sourcePerspective := NewAnalysisPerspectiveSourceOriented(v)
	sourceFactors, err := analyzer.FactorsForPerspective(sourcePerspective)
	if err != nil {
		return nil, reach.NetworkVector{}, err
	}

	destinationPerspective := NewAnalysisPerspectiveDestinationOriented(v)
	destinationFactors, err := analyzer.FactorsForPerspective(destinationPerspective)
	if err != nil {
		return nil, reach.NetworkVector{}, err
	}

	factors = append(factors, sourceFactors...)
	factors = append(factors, destinationFactors...)

	v.Source.Factors = sourceFactors
	v.Destination.Factors = destinationFactors

	return factors, v, nil
}

func sameSubnet(first, second *ElasticNetworkInterface) bool {
	if first == nil || second == nil {
		return false
	}

	return first.SubnetID == second.SubnetID
}

func sameVPC(first, second *ElasticNetworkInterface) bool {
	if first == nil || second == nil {
		return false
	}

	return first.VPCID == second.VPCID
}
