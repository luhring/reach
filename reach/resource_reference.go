package reach

import "fmt"

type ResourceReference struct {
	Domain string
	Kind   string
	ID     string
}

func (r ResourceReference) String() string {
	return fmt.Sprintf("%s->%s->%s", r.Domain, r.Kind, r.ID)
}
