package reach

import "fmt"

// Reference uniquely identifies an infrastructure resource used by Reach. It specifies the resource's Domain (e.g. AWS), Kind (e.g. EC2 instance), and ID (e.g. "i-0136d3233f0ef1924").
type Reference struct {
	Domain Domain
	Kind   Kind
	ID     string
}

// Equal returns a bool to indicate whether or not two UniversalReferences are equivalent.
func (r Reference) Equal(other Reference) bool {
	return r.Domain == other.Domain && r.Kind == other.Kind && r.ID == other.ID
}

// String returns the string representation of the Reference.
func (r Reference) String() string {
	return fmt.Sprintf("%s->%s->%s", r.Domain, r.Kind, r.ID)
}
