package aws

import (
	"fmt"

	"github.com/luhring/reach/reach"
)

func GetENIFromLineage(lineage []reach.ResourceReference, collection *reach.ResourceCollection) (*ElasticNetworkInterface, error) {
	const errPrefix = "unable to get ElasticNetworkInterface from lineage"

	for _, ref := range lineage {
		if ref.Domain == ResourceDomainAWS && ref.Kind == ResourceKindElasticNetworkInterface {
			eniResource := collection.Get(ref)
			if eniResource == nil {
				return nil, fmt.Errorf("%s: no resource found in resource collection for reference: %s", errPrefix, ref)
			}

			eni := eniResource.Properties.(ElasticNetworkInterface)
			return &eni, nil
		}
	}

	return nil, fmt.Errorf("%s: lineage does not contain an ElasticNetworkInterface", errPrefix)
}

func GetENIsFromVector(v reach.NetworkVector, collection *reach.ResourceCollection) (*ElasticNetworkInterface, *ElasticNetworkInterface, error) {
	sourceENI, err := GetENIFromLineage(v.Source.Lineage, collection)
	if err != nil {
		return nil, nil, err
	}

	destinationENI, err := GetENIFromLineage(v.Destination.Lineage, collection)
	if err != nil {
		return nil, nil, err
	}

	return sourceENI, destinationENI, nil
}

func GetEC2InstanceFromLineage(lineage []reach.ResourceReference, collection *reach.ResourceCollection) (*EC2Instance, error) {
	const errPrefix = "unable to get EC2Instance from lineage"

	for _, ref := range lineage {
		if ref.Domain == ResourceDomainAWS && ref.Kind == ResourceKindEC2Instance {
			ec2InstanceResource := collection.Get(ref)
			if ec2InstanceResource == nil {
				return nil, fmt.Errorf("%s: no resource found in resource collection for reference: %s", errPrefix, ref)
			}

			ec2Instance := ec2InstanceResource.Properties.(EC2Instance)
			return &ec2Instance, nil
		}
	}

	return nil, fmt.Errorf("%s: lineage does not contain an EC2Instance", errPrefix)
}
