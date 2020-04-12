package aws

import "github.com/luhring/reach/reach"

// A NetworkInterfaceAttachment resource representation.
type NetworkInterfaceAttachment struct {
	ID                        string
	ElasticNetworkInterfaceID string
	DeviceIndex               int64 // e.g. 0 for "eth0"
}

// Dependencies returns a collection of the network interface attachment's resource dependencies.
func (attachment NetworkInterfaceAttachment) Dependencies(provider ResourceProvider) (*reach.ResourceCollection, error) {
	rc := reach.NewResourceCollection()

	eni, err := provider.ElasticNetworkInterface(attachment.ElasticNetworkInterfaceID)
	if err != nil {
		return nil, err
	}
	rc.Put(reach.ResourceReference{
		Domain: ResourceDomainAWS,
		Kind:   ResourceKindElasticNetworkInterface,
		ID:     eni.ID,
	}, eni.Resource())

	eniDependencies, err := eni.Dependencies(provider)
	if err != nil {
		return nil, err
	}
	rc.Merge(eniDependencies)

	return rc, nil
}
