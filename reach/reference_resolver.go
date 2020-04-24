package reach

type ReferenceResolver interface {
	Resolve(ref UniversalReference) (*Resource, error)
}
