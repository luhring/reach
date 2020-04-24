package cmd

// ConfigurationDomainProvider is intended to be implemented based on configuration ingested outside of the reach package boundary.
type ConfigurationDomainProvider struct {
	providers map[string]interface{}
}

func (p *ConfigurationDomainProvider) Domain(domain string) interface{} {
	return p.providers[domain]
}

func (p *ConfigurationDomainProvider) Load(domain string, provider interface{}) {
	p.providers[domain] = provider
}
