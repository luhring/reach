package reach

// DomainClientCatalog contains a set of configured domain-specific DomainClients, each of which is implemented within a domain-specific package and used for accessing information specific to that domain.
type DomainClientCatalog struct {
	domainClients map[Domain]interface{}
}

// NewDomainClientCatalog returns a pointer to a new instance of a DomainClientCatalog.
func NewDomainClientCatalog() *DomainClientCatalog {
	dc := make(map[Domain]interface{})
	return &DomainClientCatalog{
		domainClients: dc,
	}
}

// Resolve returns a domain client for the specific domain.
func (p *DomainClientCatalog) Resolve(domain Domain) interface{} {
	return p.domainClients[domain]
}

// Store loads a domain client into the DomainClientCatalog associated with the specified domain.
func (p *DomainClientCatalog) Store(domain Domain, client interface{}) {
	p.domainClients[domain] = client
}
