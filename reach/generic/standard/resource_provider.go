package standard

// DomainClient implements a generic domain resource provider using the standard discovery mechanisms, such as DNS lookup provided from the Go standard library.
type DomainClient struct {
}

// NewDomainClient returns a reference to a new DomainClient for the standard discovery mechanisms for generic domain resources.
func NewDomainClient() *DomainClient {
	return &DomainClient{}
}
