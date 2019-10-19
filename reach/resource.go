package reach

// A Resource is a generic representation of any kind of resource from an infrastructure provider (e.g. AWS). The kind-specific properties can be provided via a kind-specific struct used for the Properties field. Then, given the Kind value, a consumer can assert the kind-specific type when reading the Properties field. Examples of a Resource include an EC2 instance, an AWS VPC, etc.
type Resource struct {
	Kind       string
	Properties interface{}
}
