package reach

import "github.com/luhring/reach/reach/traffic"

// A Point represents a "hop" along the Path that network traffic travels. A Point can refer to resource (such as an AWS Elastic Network Interface) or a non-resource (e.g., an AWS VPC's router), which can be implied by a resource (e.g., the AWS VPC).
type Point struct {
	Ref            Reference
	FactorsForward []Factor
	FactorsReturn  []Factor
	SegmentDivider bool
}

// TrafficForward returns the network traffic allowed to travel forward through this point.
func (p Point) TrafficForward() traffic.Content {
	return traffic.Intersect(TrafficFromFactors(p.FactorsForward))
}

// TrafficReturn returns the network traffic allowed to travel backward through this point.
func (p Point) TrafficReturn() traffic.Content {
	return traffic.Intersect(TrafficFromFactors(p.FactorsReturn))
}
