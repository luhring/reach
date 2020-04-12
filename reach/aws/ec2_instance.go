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

func (i EC2Instance) Resolve(role reach.SubjectRole, provider reach.InfrastructureGetter) ([]net.IP, error) {
	switch role {
	case reach.SubjectRoleSource:
		ips, err := i.ownedIPs(provider)
		if err != nil {
			return nil, fmt.Errorf("couldn't look up source IPs: %v", err)
		}
		return ips, err
	case reach.SubjectRoleDestination:
		ips, err := i.advertisedIPs(provider)
		if err != nil {
			return nil, fmt.Errorf("couldn't look up destination IPs: %v", err)
		}
		return ips, err
	default:
		return nil, fmt.Errorf("cannot look up IPs for subject role: %s", role)
	}
}

func (i EC2Instance) Visitable(alreadyVisited bool) bool {
	return alreadyVisited == false
}

func (i EC2Instance) Destination(ips []net.IP, provider reach.InfrastructureGetter) bool {
	if len(ips) == 0 {
		return false
	}

	ownedIPs, err := i.ownedIPs(provider)
	if err != nil {
		return false // TODO: Consider a better way to report the error. For now, adding an error return value seems excessive.
	}

	for _, ip := range ips {
		for _, ownedIP := range ownedIPs {
			if ip.Equal(ownedIP) {
				return true
			}
		}
	}

	return false
}

func (i EC2Instance) Segments() bool {
	return false // Note: If this resource can ever perform NAT, this answer would change.
}

func (i EC2Instance) NextTuple(prev *reach.IPTuple) *reach.IPTuple {
	// An EC2 Instance doesn't mutate the tuple. (...unless it can perform NAT.)
	return prev
}

func (i EC2Instance) Next(t *reach.IPTuple, provider reach.InfrastructureGetter) ([]reach.InfrastructureReference, error) {
	var refs []reach.InfrastructureReference

	for _, id := range i.elasticNetworkInterfaceIDs() {
		ref := reach.NewInfrastructureReference(
			ResourceDomainAWS,
			ResourceKindElasticNetworkInterface,
			id,
			false,
		)

		eniResource, err := provider.Get(ref)
		if err != nil {
			return nil, fmt.Errorf("couldn't get ENI (%s): %v", ref, err)
		}
		eni := eniResource.Properties.(ElasticNetworkInterface)

		// Only include ENIs that own the tuple's src IP
		for _, ownedIP := range eni.ownedIPs() {
			if t == nil || t.Src.Equal(ownedIP) {
				refs = append(refs, ref)
				break
			}
		}
	}

	return refs, nil
}

func (i EC2Instance) Factors() []reach.Factor {
	f := i.newInstanceStateFactor()
	return []reach.Factor{f}
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
		Resource:      i.ToResourceReference(),
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

func (i EC2Instance) networkPoints(rc *reach.ResourceCollection) []reach.NetworkPoint {
	var points []reach.NetworkPoint

	for _, id := range i.elasticNetworkInterfaceIDs() {
		eni := rc.Get(reach.ResourceReference{
			Domain: ResourceDomainAWS,
			Kind:   ResourceKindElasticNetworkInterface,
			ID:     id,
		}).Properties.(ElasticNetworkInterface)
		eniNetworkPoints := eni.networkPoints(i.ToResourceReference())
		points = append(points, eniNetworkPoints...)
	}

	return points
}

func (i EC2Instance) ownedIPs(provider reach.InfrastructureGetter) ([]net.IP, error) {
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

func (i EC2Instance) advertisedIPs(provider reach.InfrastructureGetter) ([]net.IP, error) {
	var ips []net.IP

	enis, err := i.elasticNetworkInterfaces(provider)
	if err != nil {
		return nil, fmt.Errorf("couldn't look up ENIs: %v", err)
	}

	for _, eni := range enis {
		ips = append(ips, eni.advertisedIPs()...)
	}

	return ips, nil
}
