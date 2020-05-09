package reach

// DomainClientResolver is the interface that wraps the Resolve method.
type DomainClientResolver interface {
	Resolve(domain Domain) interface{}
}
