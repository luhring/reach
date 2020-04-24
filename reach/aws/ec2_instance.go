package aws

import (
	"fmt"
	"net"
	"strings"

	"github.com/luhring/reach/reach"
)

// ResourceKindEC2Instance specifies the unique name for the EC2 instance kind of resource.
const ResourceKindEC2Instance reach.Kind = "EC2Instance"

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

func (i EC2Instance) Ref() reach.UniversalReference {
	return reach.UniversalReference{
		R: i.ResourceReference(),
	}
}

func (i EC2Instance) Segments() bool {
	return false // TODO: If this resource can ever perform NAT, this answer would change.
}

func (i EC2Instance) ForwardEdges(resolver reach.DomainClientResolver, previousEdge *reach.Edge, destinationIPs []net.IP) ([]reach.Edge, error) {
	var tuples []reach.IPTuple
	if previousEdge == nil { // This is the first point in the path.
		t, err := i.firstPointTuples(resolver, destinationIPs)
		if err != nil {
			return nil, fmt.Errorf("cannot generate tuples for first point: %v", err)
		}
		tuples = t
	} else {
		// Note: If the EC2 instance is changing the IP tuple from a previous tuple state, we don't have visibility into that change, so we'll have to assume no change.
		tuples = []reach.IPTuple{previousEdge.Tuple}
	}

	client, err := unpackDomainClient(resolver)
	if err != nil {
		return nil, fmt.Errorf("unable to get client: %v", err)
	}

	enis, err := i.elasticNetworkInterfaces(client)
	if err != nil {
		return nil, fmt.Errorf("couldn't get ENIs: %v", err)
	}
	var edges []reach.Edge
	for _, eni := range enis {
		for _, tuple := range tuples {
			edge := reach.Edge{
				Tuple:             tuple,
				EndRef:            eni.Ref(),
				ConnectsInterface: true,
			}
			edges = append(edges, edge)
		}
	}

	return edges, nil
}

func (i EC2Instance) firstPointTuples(
	resolver reach.DomainClientResolver,
	destinationIPs []net.IP,
) ([]reach.IPTuple, error) {
	// If the traffic originates from this EC2 instance, any of its owned IP addresses could be used as source.
	// (Technically, any IP address could be used as source, as long as src/dst check is off, but currently we have no way to inform the Tracer about scenarios like this.)

	srcIPs, err := i.InterfaceIPs(resolver)
	if err != nil {
		return nil, fmt.Errorf("cannot determine possible src IPs: %v", err)
	}
	// TODO: We need some mechanism to confirm destination addresses are valid in the context of source's network.
	var tuples []reach.IPTuple
	for _, src := range srcIPs {
		for _, dst := range destinationIPs {
			tuples = append(tuples, reach.IPTuple{
				Src: src,
				Dst: dst,
			})
		}
	}
	return tuples, nil
}

func (i EC2Instance) FactorsForward(
	_ reach.DomainClientResolver,
	_ *reach.Edge,
) ([]reach.Factor, error) {
	f := i.newInstanceStateFactor()
	return []reach.Factor{f}, nil
}

func (i EC2Instance) IPs(resolver reach.DomainClientResolver) ([]net.IP, error) {
	client, err := unpackDomainClient(resolver)
	if err != nil {
		return nil, fmt.Errorf("unable to get client: %v", err)
	}

	enis, err := i.elasticNetworkInterfaces(client)
	if err != nil {
		return nil, fmt.Errorf("couldn't look up ENIs: %v", err)
	}

	var ips []net.IP
	for _, eni := range enis {
		eniIPs, err := eni.IPs(resolver)
		if err != nil {
			return nil, fmt.Errorf("couldn't get IPs for ENI (%s): %v", eni.Ref(), err)
		}
		ips = append(ips, eniIPs...)
	}

	return ips, nil
}

func (i EC2Instance) InterfaceIPs(resolver reach.DomainClientResolver) ([]net.IP, error) {
	client, err := unpackDomainClient(resolver)
	if err != nil {
		return nil, fmt.Errorf("unable to get client: %v", err)
	}

	enis, err := i.elasticNetworkInterfaces(client)
	if err != nil {
		return nil, fmt.Errorf("couldn't look up ENIs: %v", err)
	}

	var ips []net.IP
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

func (i EC2Instance) elasticNetworkInterfaces(client DomainClient) ([]ElasticNetworkInterface, error) {
	eniIDs := i.elasticNetworkInterfaceIDs()
	enis := make([]ElasticNetworkInterface, len(eniIDs))

	for _, id := range eniIDs {
		eni, err := client.ElasticNetworkInterface(id)
		if err != nil {
			return nil, fmt.Errorf("couldn't get ENI (%s): %v", id, err)
		}

		enis = append(enis, *eni)
	}

	return enis, nil
}
