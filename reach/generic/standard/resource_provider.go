package standard

// ResourceProvider implements a generic domain resource provider using the standard discovery mechanisms, such as DNS lookup provided from the Go standard library.
type ResourceProvider struct {
}

// NewResourceProvider returns a reference to a new ResourceProvider for the standard discovery mechanisms for generic domain resources.
func NewResourceProvider() *ResourceProvider {
	return &ResourceProvider{}
}
