package reach

type Tracer interface {
	Trace(source, destination Subject) ([]Path, error)
}
