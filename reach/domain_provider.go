package reach

type DomainProvider interface {
	Domain(domain string) interface{}
}
