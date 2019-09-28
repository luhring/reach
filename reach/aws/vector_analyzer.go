package aws

import "github.com/luhring/reach/reach"

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
					p.getSecurityGroupRules,
					p.other.IPAddress,
					targetENI,
				)
				if err != nil {
					return nil, err
				}

				factors = append(factors, *securityGroupRulesFactor)
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