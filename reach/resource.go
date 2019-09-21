package reach

type Resource struct {
	Kind       string      `json:"kind"`
	Properties interface{} `json:"properties"`
}

func MergeResources(a, b map[string]map[string]map[string]Resource) map[string]map[string]map[string]Resource {
	output := make(map[string]map[string]map[string]Resource)

	for resourceDomain, resourceKinds := range a { // e.g. for AWS
		if _, exists := b[resourceDomain]; !exists { // only A has AWS
			output[resourceDomain] = resourceKinds
		} else { // both have AWS
			output = EnsureResourcePathExists(output, resourceDomain, "")

			for resourceKind, resources := range a[resourceDomain] { // e.g. for EC2 instances
				if _, exists := b[resourceDomain][resourceKind]; !exists { // only A has any EC2 instances
					output[resourceDomain][resourceKind] = resources
				} else { // both have some EC2 instances
					output = EnsureResourcePathExists(output, resourceDomain, resourceKind)

					for id, resource := range a[resourceDomain][resourceKind] { // e.g. for EC2 instance with ID i-abc123def456
						output[resourceDomain][resourceKind][id] = resource
					}

					for id, resource := range b[resourceDomain][resourceKind] {
						output[resourceDomain][resourceKind][id] = resource
					}
				}
			}

			for resourceKind, resources := range b[resourceDomain] { // e.g. for security groups
				if _, exists := a[resourceDomain][resourceKind]; !exists { // only B has any security groups
					output[resourceDomain][resourceKind] = resources
				}
			}
		}
	}

	for resourceDomain, resourceKinds := range b { // e.g. for GCP
		if _, exists := a[resourceDomain]; !exists { // only B has GCP
			output[resourceDomain] = resourceKinds
		}
	}

	return output
}

func newResourceDomainMap() map[string]map[string]Resource {
	return make(map[string]map[string]Resource)
}

func newResourceKindMap() map[string]Resource {
	return make(map[string]Resource)
}

func EnsureResourcePathExists(resources map[string]map[string]map[string]Resource, domain, kind string) map[string]map[string]map[string]Resource {
	if domain == "" {
		return resources
	}

	if _, exists := resources[domain]; !exists {
		resources[domain] = newResourceDomainMap()
	}

	if kind == "" {
		return resources
	}

	if _, exists := resources[domain][kind]; !exists {
		resources[domain][kind] = newResourceKindMap()
	}

	return resources
}
