package reach

// Tracer is the interface that wraps the Trace method. Implementers of Tracer provide a mechanism for constructing network paths between specified sources and destinations across one or more domains.
type Tracer interface {
	Trace(source, destination Subject) ([]Path, error)
}
