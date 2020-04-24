package reach

type DomainClientResolver interface {
	Resolve(domain Domain) interface{}
}
