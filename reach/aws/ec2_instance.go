package aws

import "github.com/luhring/reach/reach"

const ResourceKindEC2Instance = "EC2Instance"

type EC2Instance struct {
	ID                          string                       `json:"id"`
	NameTag                     string                       `json:"nameTag"`
	State                       string                       `json:"state"`
	NetworkInterfaceAttachments []NetworkInterfaceAttachment `json:"networkInterfaceAttachments"`
}

func (i EC2Instance) ToResource() reach.Resource {
	return reach.Resource{
		Kind:       ResourceKindEC2Instance,
		Properties: i,
	}
}

func (i EC2Instance) getUniqueID() string {
	return i.ID
}

func (i EC2Instance) GetDependencies(provider ResourceProvider) (map[string]map[string]map[string]reach.Resource, error) {
	resources := make(map[string]map[string]map[string]reach.Resource)

	for _, attachment := range i.NetworkInterfaceAttachments {
		attachmentDependencies, err := getDependenciesForNetworkInterfaceAttachment(attachment, provider)
		if err != nil {
			return nil, err
		}
		resources = reach.MergeResources(resources, attachmentDependencies)
	}

	return resources, nil
}

func getDependenciesForNetworkInterfaceAttachment(attachment NetworkInterfaceAttachment, provider ResourceProvider) (map[string]map[string]map[string]reach.Resource, error) {
	resources := make(map[string]map[string]map[string]reach.Resource)

	eni, err := provider.GetElasticNetworkInterface(attachment.ElasticNetworkInterfaceID)
	if err != nil {
		return nil, err
	}
	resources = reach.EnsureResourcePathExists(resources, ResourceDomainAWS, ResourceKindElasticNetworkInterface)
	resources[ResourceDomainAWS][ResourceKindElasticNetworkInterface][eni.ID] = eni.ToResource()

	eniDependencies, err := eni.GetDependencies(provider)
	if err != nil {
		return nil, err
	}
	resources = reach.MergeResources(resources, eniDependencies)

	return resources, nil
}
