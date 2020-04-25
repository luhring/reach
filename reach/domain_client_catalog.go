package reach

type DomainClientCatalog struct {
	domainClients map[Domain]interface{}
}

func NewDomainClientCatalog() *DomainClientCatalog {
	dc := make(map[Domain]interface{})
	return &DomainClientCatalog{
		domainClients: dc,
	}
}

func (p *DomainClientCatalog) Resolve(domain Domain) interface{} {
	return p.domainClients[domain]
}

func (p *DomainClientCatalog) Store(domain Domain, client interface{}) {
	p.domainClients[domain] = client
}
