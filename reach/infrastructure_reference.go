package reach

import "fmt"

type InfrastructureReference struct {
	Implicit bool // Is infrastructure implied by referenced resource instead of being the resource itself
	R        ResourceReference
}

func NewInfrastructureReference(domain, kind, id string, implicit bool) InfrastructureReference {
	return InfrastructureReference{
		R: ResourceReference{
			Domain: domain,
			Kind:   kind,
			ID:     id,
		},
		Implicit: implicit,
	}
}

func (i InfrastructureReference) Equal(other InfrastructureReference) bool {
	if i.Implicit != other.Implicit {
		return false
	}

	return i.R.Equal(other.R)
}

func (i InfrastructureReference) String() string {
	var implicitSuffix string
	if i.Implicit {
		implicitSuffix = " (implicit infrastructure)"
	}

	return fmt.Sprintf("%s%s", i.R, implicitSuffix)
}
