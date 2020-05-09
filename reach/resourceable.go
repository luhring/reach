package reach

// Resourceable is the interface that wraps the Resource method.
//
// This interface is implemented by pieces of infrastructure that can produce a generalized resource to describe their complete state.
type Resourceable interface {
	Resource() Resource
}
