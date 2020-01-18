package aws

import (
	"errors"
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

func (analyzer VectorAnalyzer) factorsForPerspective(p reach.Perspective) ([]reach.Factor, error) {
	var factors []reach.Factor

	for _, selfResourceRef := range p.Self.Lineage {
		if selfResourceRef.Domain == ResourceDomainAWS {
			if selfResourceRef.Kind == ResourceKindEC2Instance {
				ec2Instance := analyzer.resourceCollection.Get(selfResourceRef).Properties.(EC2Instance)

				factors = append(factors, ec2Instance.newInstanceStateFactor())
			}

			if selfResourceRef.Kind == ResourceKindElasticNetworkInterface {
				eni := analyzer.resourceCollection.Get(selfResourceRef).Properties.(ElasticNetworkInterface)

				if p.Other.Domain() == ResourceDomainAWS {
					otherENIRef := ElasticNetworkInterfaceFromNetworkPoint(p.Other, analyzer.resourceCollection)
					if otherENIRef == nil {
						return nil, errors.New("unable to find elastic network interface for network point within AWS domain")
					}

					otherENI := *otherENIRef

					// Ensure this is scenario that Reach can analyze
					if !sameVPC(eni, otherENI) {
						return nil, errors.New("error: reach is not yet able to analyze EC2 instances in different VPCs, but that's coming soon")
					}

					var awsP perspective
					if p.SelfRole == reach.SubjectRoleSource {
						awsP = newPerspectiveSourceOriented()
					} else {
						awsP = newPerspectiveDestinationOriented()
					}

					// Evaluate factors
					securityGroupRulesFactor, err := eni.newSecurityGroupRulesFactor(
						analyzer.resourceCollection,
						p,
						awsP,
						otherENI,
					)
					if err != nil {
						return nil, err
					}

					factors = append(factors, *securityGroupRulesFactor)

					if sameSubnet(eni, otherENI) {
						// There's nothing further to evaluate for this ENI
						continue
					}

					// Different subnets, same VPC

					networkACLRulesFactor, err := eni.newNetworkACLRulesFactor(
						analyzer.resourceCollection,
						p,
						awsP,
					)
					if err != nil {
						return nil, err
					}

					factors = append(factors, *networkACLRulesFactor)
				} else {
					// Other point is not within AWS. For now, we'll only support the other point having a public IP address that we can connect to.
					if !p.Other.IPAddressIsInternetAccessible() {
						return nil, fmt.Errorf("encountered network point with IP address ('%s') that is not Internet accessible (this scenario is not yet supported)", p.Other.IPAddress)
					}
				}
			}
		}
	}

	return factors, nil
}

func sameSubnet(first, second ElasticNetworkInterface) bool {
	return first.SubnetID == second.SubnetID
}

func sameVPC(first, second ElasticNetworkInterface) bool {
	return first.VPCID == second.VPCID
}
