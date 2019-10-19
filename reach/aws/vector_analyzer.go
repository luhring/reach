package aws

import (
	"fmt"

	"github.com/luhring/reach/reach"
)

// VectorAnalyzer is the AWS-specific implementation of the VectorAnalyzer interface.
type VectorAnalyzer struct {
	resourceCollection *reach.ResourceCollection
}

// NewVectorAnalyzer creates a new AWS-specific VectorAnalyzer.
func NewVectorAnalyzer(resourceCollection *reach.ResourceCollection) VectorAnalyzer {
	return VectorAnalyzer{
		resourceCollection,
	}
}

func (analyzer VectorAnalyzer) factorsForPerspective(p reach.Perspective) ([]reach.Factor, error) {
	var factors []reach.Factor

	for _, resourceRef := range p.Self.Lineage {
		if resourceRef.Domain == ResourceDomainAWS {
			if resourceRef.Kind == ResourceKindEC2Instance {
				ec2Instance := analyzer.resourceCollection.Get(resourceRef).Properties.(EC2Instance)

				factors = append(factors, ec2Instance.newInstanceStateFactor())
			}

			if resourceRef.Kind == ResourceKindElasticNetworkInterface {
				eni := analyzer.resourceCollection.Get(resourceRef).Properties.(ElasticNetworkInterface)
				targetENI := ElasticNetworkInterfaceFromNetworkPoint(p.Other, analyzer.resourceCollection)

				var awsP perspective
				if p.SelfRole == reach.SubjectRoleSource {
					awsP = newPerspectiveSourceOriented()
				} else {
					awsP = newPerspectiveDestinationOriented()
				}

				securityGroupRulesFactor, err := eni.newSecurityGroupRulesFactor(
					analyzer.resourceCollection,
					p,
					awsP,
					targetENI,
				)
				if err != nil {
					return nil, err
				}

				factors = append(factors, *securityGroupRulesFactor)

				if !sameVPC(&eni, targetENI) {
					return nil, fmt.Errorf("error: reach is not yet able to analyze EC2 instances in different VPCs, but that's coming soon! (VPCs: %s, %s)", eni.VPCID, targetENI.VPCID)
				}

				if !sameSubnet(&eni, targetENI) {
					return nil, fmt.Errorf("error: reach is not yet able to analyze EC2 instances in different subnets, but that's coming soon! (subnets: %s, %s)", eni.SubnetID, targetENI.SubnetID)
				}
			}
		}
	}

	return factors, nil
}

// Factors calculates the analysis factors for the given network vector.
func (analyzer VectorAnalyzer) Factors(v reach.NetworkVector) ([]reach.Factor, reach.NetworkVector, error) {
	var factors []reach.Factor

	sourcePerspective := v.SourcePerspective()
	sourceFactors, err := analyzer.factorsForPerspective(sourcePerspective)
	if err != nil {
		return nil, reach.NetworkVector{}, err
	}

	destinationPerspective := v.DestinationPerspective()
	destinationFactors, err := analyzer.factorsForPerspective(destinationPerspective)
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
