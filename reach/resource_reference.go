package reach

import "fmt"

// ResourceReference uniquely identifies a Resource used by Reach. It specifies the resource's Domain (e.g. AWS), Kind (e.g. EC2 instance), and ID (e.g. "i-0136d3233f0ef1924").
type ResourceReference struct {
	Domain string
	Kind   string
	ID     string
}

// String returns the string representation of the ResourceReference.
func (r ResourceReference) String() string {
	return fmt.Sprintf("%s->%s->%s", r.Domain, r.Kind, r.ID)
}

func (r ResourceReference) Equal(other ResourceReference) bool {
	return r.Domain == other.Domain && r.Kind == other.Kind && r.ID == other.ID
}
