package reach

type DomainClientCatalog struct {
	domainClients map[Domain]interface{}
}

func (p *DomainClientCatalog) Resolve(domain Domain) interface{} {
	return p.domainClients[domain]
}

func (p *DomainClientCatalog) Store(domain Domain, client interface{}) {
	p.domainClients[domain] = client
}
