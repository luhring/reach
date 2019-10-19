package reach

import "encoding/json"

// A ResourceCollection is a structure used to store any number of Resources, across potentially multiple "domains" (e.g. AWS, GCP, Azure) and kinds (e.g. EC2 instance, subnet, etc.).
type ResourceCollection struct {
	collection map[string]map[string]map[string]Resource
}

// NewResourceCollection returns a reference to a new, empty ResourceCollection.
func NewResourceCollection() *ResourceCollection {
	collection := make(map[string]map[string]map[string]Resource)

	return &ResourceCollection{
		collection: collection,
	}
}

// Put adds a new Resource to the ResourceCollection.
func (rc *ResourceCollection) Put(ref ResourceReference, resource Resource) {
	rc.ensureResourcePathExists(ref.Domain, ref.Kind)

	other := NewResourceCollection()
	other.ensureResourcePathExists(ref.Domain, ref.Kind)
	other.collection[ref.Domain][ref.Kind][ref.ID] = resource

	rc.Merge(other)
}

// Get retrieves a Resource from the ResourceCollection.
func (rc *ResourceCollection) Get(ref ResourceReference) *Resource {
	if _, exists := rc.collection[ref.Domain]; !exists {
		return nil
	}

	if _, exists := rc.collection[ref.Domain][ref.Kind]; !exists {
		return nil
	}

	if resource, exists := rc.collection[ref.Domain][ref.Kind][ref.ID]; exists {
		return &resource
	}
	return nil
}

// Merge safely merges two ResourceCollections such that any unique resource from either collection is represented in the merged collection. For any case where both collections contain a resource for a given domain, kind, and resource ID, the "other" (input parameter) resource will overwrite the corresponding resource in the first collection.
func (rc *ResourceCollection) Merge(other *ResourceCollection) {
	for resourceDomain, resourceKinds := range rc.collection { // e.g. for AWS
		if _, exists := other.collection[resourceDomain]; !exists { // only A has AWS
			rc.collection[resourceDomain] = resourceKinds
		} else { // both have AWS
			rc.ensureResourcePathExists(resourceDomain, "")

			for resourceKind, resources := range rc.collection[resourceDomain] { // e.g. for EC2 instances
				if _, exists := other.collection[resourceDomain][resourceKind]; !exists { // only A has any EC2 instances
					rc.collection[resourceDomain][resourceKind] = resources
				} else { // both have some EC2 instances
					rc.ensureResourcePathExists(resourceDomain, resourceKind)

					for id, resource := range rc.collection[resourceDomain][resourceKind] { // e.g. for EC2 instance with ID i-abc123def456
						rc.collection[resourceDomain][resourceKind][id] = resource
					}

					for id, resource := range other.collection[resourceDomain][resourceKind] {
						rc.collection[resourceDomain][resourceKind][id] = resource
					}
				}
			}

			for resourceKind, resources := range other.collection[resourceDomain] { // e.g. for security groups
				if _, exists := rc.collection[resourceDomain][resourceKind]; !exists { // only B has any security groups
					rc.collection[resourceDomain][resourceKind] = resources
				}
			}
		}
	}

	for resourceDomain, resourceKinds := range other.collection { // e.g. for GCP
		if _, exists := rc.collection[resourceDomain]; !exists { // only B has GCP
			rc.collection[resourceDomain] = resourceKinds
		}
	}
}

// MarshalJSON returns the JSON representation of the ResourceCollection.
func (rc *ResourceCollection) MarshalJSON() ([]byte, error) {
	return json.Marshal(rc.collection)
}

func (rc *ResourceCollection) ensureResourcePathExists(domain, kind string) {
	if domain == "" {
		return
	}

	if _, exists := rc.collection[domain]; !exists {
		rc.collection[domain] = newResourceDomainMap()
	}

	if kind == "" {
		return
	}

	if _, exists := rc.collection[domain][kind]; !exists {
		rc.collection[domain][kind] = newResourceKindMap()
	}

	return
}

func newResourceDomainMap() map[string]map[string]Resource {
	return make(map[string]map[string]Resource)
}

func newResourceKindMap() map[string]Resource {
	return make(map[string]Resource)
}
