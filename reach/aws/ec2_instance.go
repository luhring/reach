package aws

import (
	"fmt"
	"net"
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

// Resource returns the EC2 instance converted to a generalized Reach resource.
func (i EC2Instance) Resource() reach.Resource {
	return reach.Resource{
		Kind:       ResourceKindEC2Instance,
		Properties: i,
	}
}

// ResourceReference returns a resource reference to uniquely identify the EC2 instance.
func (i EC2Instance) ResourceReference() reach.ResourceReference {
	return reach.ResourceReference{
		Domain: ResourceDomainAWS,
		Kind:   ResourceKindEC2Instance,
		ID:     i.ID,
	}
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

// Name returns the instance's ID, and, if available, its name tag value.
func (i EC2Instance) Name() string {
	if name := strings.TrimSpace(i.NameTag); name != "" {
		return fmt.Sprintf("\"%s\" (%s)", name, i.ID)
	}
	return i.ID
}

func (i EC2Instance) Visitable(alreadyVisited bool) bool {
	return alreadyVisited == false
}

func (i EC2Instance) Ref() reach.InfrastructureReference {
	return reach.InfrastructureReference{
		R: i.ResourceReference(),
	}
}

func (i EC2Instance) Segments() bool {
	return false // Note: If this resource can ever perform NAT, this answer would change.
}

func (i EC2Instance) ForwardEdges(latestTuple *reach.IPTuple, provider reach.InfrastructureGetter) ([]reach.PathEdge, error) {
	// Note: If the EC2 instance is changing the IP tuple from a previous tuple state, we don't have visibility into that change, so we'll have to assume no change.

	enis, err := i.elasticNetworkInterfaces(provider)
	if err != nil {
		return nil, fmt.Errorf("couldn't get ENIs: %v", err)
	}
	var edges []reach.PathEdge
	for _, eni := range enis {
		edge := reach.PathEdge{
			Tuple: latestTuple,
			Ref:   eni.Ref(),
		}
		edges = append(edges, edge)
	}

	return edges, nil
}

func (i EC2Instance) Factors() []reach.Factor {
	f := i.newInstanceStateFactor()
	return []reach.Factor{f}
}

func (i EC2Instance) IPs(provider reach.InfrastructureGetter) ([]net.IP, error) {
	var ips []net.IP

	enis, err := i.elasticNetworkInterfaces(provider)
	if err != nil {
		return nil, fmt.Errorf("couldn't look up ENIs: %v", err)
	}

	for _, eni := range enis {
		eniIPs, err := eni.IPs(provider)
		if err != nil {
			return nil, fmt.Errorf("couldn't get IPs for ENI (%s): %v", eni.Ref(), err)
		}
		ips = append(ips, eniIPs...)
	}

	return ips, nil
}

func (i EC2Instance) InterfaceIPs(provider reach.InfrastructureGetter) ([]net.IP, error) {
	var ips []net.IP

	enis, err := i.elasticNetworkInterfaces(provider)
	if err != nil {
		return nil, fmt.Errorf("couldn't look up ENIs: %v", err)
	}

	for _, eni := range enis {
		ips = append(ips, eni.ownedIPs()...)
	}

	return ips, nil
}

// FactorKindInstanceState specifies the unique name for the EC2 instance state of factor.
const FactorKindInstanceState = "InstanceState"

func (i EC2Instance) newInstanceStateFactor() reach.Factor {
	var traffic reach.TrafficContent
	var returnTraffic reach.TrafficContent

	if i.isRunning() {
		traffic = reach.NewTrafficContentForAllTraffic()
		returnTraffic = reach.NewTrafficContentForAllTraffic()
	} else {
		traffic = reach.NewTrafficContentForNoTraffic()
		returnTraffic = reach.NewTrafficContentForNoTraffic()
	}

	return reach.Factor{
		Kind:          FactorKindInstanceState,
		Resource:      i.ResourceReference(),
		Traffic:       traffic,
		ReturnTraffic: returnTraffic,
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

func (i EC2Instance) elasticNetworkInterfaces(provider reach.InfrastructureGetter) ([]ElasticNetworkInterface, error) {
	eniIDs := i.elasticNetworkInterfaceIDs()
	enis := make([]ElasticNetworkInterface, len(eniIDs))

	for _, id := range eniIDs {
		ref := reach.NewInfrastructureReference(
			ResourceDomainAWS,
			ResourceKindElasticNetworkInterface,
			id,
			false,
		)
		r, err := provider.Get(ref)
		if err != nil {
			return nil, fmt.Errorf("couldn't get ENI (%s): %v", ref, err)
		}
		eni := r.Properties.(ElasticNetworkInterface)

		enis = append(enis, eni)
	}

	return enis, nil
}
