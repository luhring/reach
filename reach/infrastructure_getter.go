package reach

type InfrastructureGetter interface {
	Get(ref InfrastructureReference) (Resource, error)
}
