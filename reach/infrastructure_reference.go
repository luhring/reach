package reach

type InfrastructureReference struct {
	Implicit bool // Is infrastructure implied by referenced resource instead of being the resource itself
	R        ResourceReference
}
