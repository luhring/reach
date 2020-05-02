package reach

import "fmt"

// UniversalReference uniquely identifies a Resource used by Reach. It specifies the resource's Domain (e.g. AWS), Kind (e.g. EC2 instance), and ID (e.g. "i-0136d3233f0ef1924").
type UniversalReference struct {
	Domain Domain
	Kind   Kind
	ID     string
}

// Equal returns a bool to indicate whether or not two UniversalReferences are equivalent.
func (r UniversalReference) Equal(other UniversalReference) bool {
	return r.Domain == other.Domain && r.Kind == other.Kind && r.ID == other.ID
}

// String returns the string representation of the UniversalReference.
func (r UniversalReference) String() string {
	return fmt.Sprintf("%s->%s->%s", r.Domain, r.Kind, r.ID)
}
