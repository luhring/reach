package reach

// A Point represents a "hop" along the Path that network traffic travels. A Point can refer to resource (such as an AWS Elastic Network Interface) or a non-resource (e.g., an AWS VPC's router), which can be implied by a resource (e.g., the AWS VPC).
type Point struct {
	Ref            Reference
	FactorsForward []Factor
	FactorsReturn  []Factor
}
