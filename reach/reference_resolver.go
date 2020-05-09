package reach

// ReferenceResolver is the interface that wraps the Resolve method.
//
// This interface should be implemented once per domain package, and by the analyzer package, to allow domains to provide complete resources when given just a reference to the resource.
type ReferenceResolver interface {
	Resolve(ref Reference) (*Resource, error)
}
