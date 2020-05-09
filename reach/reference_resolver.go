package reach

type ReferenceResolver interface {
	Resolve(ref Reference) (*Resource, error)
}
