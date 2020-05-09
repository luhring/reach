package reach

// Referable is the interface that wraps the Ref method. It describes infrastructure that has the ability to generate a reference for itself.
type Referable interface {
	Ref() Reference
}
