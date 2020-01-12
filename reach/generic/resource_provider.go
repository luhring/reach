package generic

// The ResourceProvider interface wraps all of the necessary methods for accessing generic domain resources.
type ResourceProvider interface {
	Hostname(name string) (*Hostname, error)
}
