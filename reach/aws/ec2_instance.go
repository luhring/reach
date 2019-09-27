package aws

import "github.com/luhring/reach/reach"

const ResourceKindEC2Instance = "EC2Instance"

type EC2Instance struct {
	ID                          string
	NameTag                     string
	State                       string
	NetworkInterfaceAttachments []NetworkInterfaceAttachment
}

func (i EC2Instance) ToResource() reach.Resource {
	return reach.Resource{
		Kind:       ResourceKindEC2Instance,
		Properties: i,
	}
}

func (i EC2Instance) ToResourceReference() reach.ResourceReference {
	return reach.ResourceReference{
		Domain: ResourceDomainAWS,
		Kind:   ResourceKindEC2Instance,
		ID:     i.ID,
	}
}

func (i EC2Instance) isRunning() bool {
	return i.State == "running"
}

func (i EC2Instance) getElasticNetworkInterfaceIDs() []string {
	var ids []string

	for _, attachment := range i.NetworkInterfaceAttachments {
		ids = append(ids, attachment.ElasticNetworkInterfaceID)
	}

	return ids
}

func (i EC2Instance) GetDependencies(provider ResourceProvider) (*reach.ResourceCollection, error) {
	rc := reach.NewResourceCollection()

	for _, attachment := range i.NetworkInterfaceAttachments {
		attachmentDependencies, err := getDependenciesForNetworkInterfaceAttachment(attachment, provider)
		if err != nil {
			return nil, err
		}
		rc.Merge(attachmentDependencies)
	}

	return rc, nil
}

func getDependenciesForNetworkInterfaceAttachment(attachment NetworkInterfaceAttachment, provider ResourceProvider) (*reach.ResourceCollection, error) {
	rc := reach.NewResourceCollection()

	eni, err := provider.GetElasticNetworkInterface(attachment.ElasticNetworkInterfaceID)
	if err != nil {
		return nil, err
	}
	rc.Put(reach.ResourceReference{
		Domain: ResourceDomainAWS,
		Kind:   ResourceKindElasticNetworkInterface,
		ID:     eni.ID,
	}, eni.ToResource())

	eniDependencies, err := eni.GetDependencies(provider)
	if err != nil {
		return nil, err
	}
	rc.Merge(eniDependencies)

	return rc, nil
}

func (i EC2Instance) GetNetworkPoints(rc *reach.ResourceCollection) []reach.NetworkPoint {
	var points []reach.NetworkPoint

	for _, id := range i.getElasticNetworkInterfaceIDs() {
		eni := rc.Get(reach.ResourceReference{
			Domain: ResourceDomainAWS,
			Kind:   ResourceKindElasticNetworkInterface,
			ID:     id,
		}).Properties.(ElasticNetworkInterface)
		eniNetworkPoints := eni.GetNetworkPoints(i.ToResourceReference())
		points = append(points, eniNetworkPoints...)
	}

	return points
}
