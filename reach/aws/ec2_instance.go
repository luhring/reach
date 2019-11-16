package aws

import (
	"fmt"
	"strings"

	"github.com/luhring/reach/reach"
)

// ResourceKindEC2Instance specifies the unique name for the EC2 instance kind of resource.
const ResourceKindEC2Instance = "EC2Instance"

// An EC2Instance resource representation.
type EC2Instance struct {
	ID                          string
	NameTag                     string `json:"NameTag,omitempty"`
	State                       string
	NetworkInterfaceAttachments []NetworkInterfaceAttachment
}

// ToResource returns the EC2 instance converted to a generalized Reach resource.
func (i EC2Instance) ToResource() reach.Resource {
	return reach.Resource{
		Kind:       ResourceKindEC2Instance,
		Properties: i,
	}
}

// ToResourceReference returns a resource reference to uniquely identify the EC2 instance.
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

func (i EC2Instance) elasticNetworkInterfaceIDs() []string {
	var ids []string

	for _, attachment := range i.NetworkInterfaceAttachments {
		ids = append(ids, attachment.ElasticNetworkInterfaceID)
	}

	return ids
}

// Dependencies returns a collection of the EC2 instance's resource dependencies.
func (i EC2Instance) Dependencies(provider ResourceProvider) (*reach.ResourceCollection, error) {
	rc := reach.NewResourceCollection()

	for _, attachment := range i.NetworkInterfaceAttachments {
		attachmentDependencies, err := attachment.Dependencies(provider)
		if err != nil {
			return nil, err
		}
		rc.Merge(attachmentDependencies)
	}

	return rc, nil
}

func (i EC2Instance) networkPoints(rc *reach.ResourceCollection) []reach.NetworkPoint {
	var points []reach.NetworkPoint

	for _, id := range i.elasticNetworkInterfaceIDs() {
		eni := rc.Get(reach.ResourceReference{
			Domain: ResourceDomainAWS,
			Kind:   ResourceKindElasticNetworkInterface,
			ID:     id,
		}).Properties.(ElasticNetworkInterface)
		eniNetworkPoints := eni.getNetworkPoints(i.ToResourceReference())
		points = append(points, eniNetworkPoints...)
	}

	return points
}

// Name returns the instance's ID, and, if available, its name tag value.
func (i EC2Instance) Name() string {
	if name := strings.TrimSpace(i.NameTag); name != "" {
		return fmt.Sprintf("\"%s\" (%s)", name, i.ID)
	}
	return i.ID
}
