package reach

type InfrastructureReference struct {
	Implicit bool // Is infrastructure implied by referenced resource instead of being the resource itself
	R        ResourceReference
}

func (i InfrastructureReference) Matches(other InfrastructureReference) bool {
	if i.Implicit != other.Implicit {
		return false
	}

	return i.R.Matches(other.R)
}
