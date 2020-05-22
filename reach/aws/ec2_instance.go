package aws

import (
	"fmt"
	"net"
	"strings"

	"github.com/luhring/reach/reach"
)

// ResourceKindEC2Instance specifies the unique name for the EC2Instance kind of resource.
const ResourceKindEC2Instance reach.Kind = "EC2Instance"

var _ reach.Traceable = (*EC2Instance)(nil)
var _ reach.IPAddressable = (*EC2Instance)(nil)

// An EC2Instance resource representation.
type EC2Instance struct {
	ID                          string
	NameTag                     string `json:"NameTag,omitempty"`
	State                       string
	NetworkInterfaceAttachments []NetworkInterfaceAttachment
}

// EC2InstanceRef returns a Reference for an EC2Instance with the specified ID.
func EC2InstanceRef(id string) reach.Reference {
	return reach.Reference{
		Domain: ResourceDomainAWS,
		Kind:   ResourceKindEC2Instance,
		ID:     id,
	}
}

// Resource returns the EC2Instance converted to a generalized Reach resource.
func (i EC2Instance) Resource() reach.Resource {
	return reach.Resource{
		Kind:       ResourceKindEC2Instance,
		Properties: i,
	}
}

// Name returns the instance's ID, and, if available, its name tag value.
func (i EC2Instance) Name() string {
	if name := strings.TrimSpace(i.NameTag); name != "" {
		return fmt.Sprintf("\"%s\" (%s)", name, i.ID)
	}
	return i.ID
}

// ———— Implementing Traceable ————

// Ref returns a Reference for the EC2Instance.
func (i EC2Instance) Ref() reach.Reference {
	return EC2InstanceRef(i.ID)
}

// Visitable returns a boolean to indicate whether a tracer is allowed to add this resource to the path it's currently constructing.
//
// In most cases, visiting an EC2 instance multiple times indicates that a network path cycle is occurring, which should be considered fatal for the trace. Network setups that legitimately leverage repeat visits to an EC2Instance are not yet supported by Reach.
func (i EC2Instance) Visitable(alreadyVisited bool) bool {
	return alreadyVisited == false
}

// Segments returns a boolean to indicate whether a tracer should create a new path segment at this point in the path.
//
// Segments always returns false. Although it's certainly possible for an EC2Instance to perform NAT, such network setups are not yet supported by Reach.
func (i EC2Instance) Segments() bool {
	return false // TODO: If this resource can ever perform NAT, this answer would change.
}

// EdgesForward returns the set of all possible edges forward given this point in a path that a tracer is constructing. EdgesForward returns an empty slice of edges if there are no further points for the specified network traffic to travel as it attempts to reach its intended network destination.
func (i EC2Instance) EdgesForward(resolver reach.DomainClientResolver, leftEdge *reach.Edge, _ *reach.Reference, destinationIPs []net.IP) ([]reach.Edge, error) {
	// TODO: Use leftPointRef for more intelligent detection of incoming traffic's origin

	var tuples []reach.IPTuple
	if leftEdge == nil { // This is the first point in the path.
		t, err := i.firstPointTuples(resolver, destinationIPs)
		if err != nil {
			return nil, err
		}
		tuples = t
	} else {
		// Note: If the EC2 instance is changing the IP tuple from a previous tuple state, we don't have visibility into that change, so we'll have to assume no change.
		tuples = []reach.IPTuple{leftEdge.Tuple}
	}

	client, err := unpackDomainClient(resolver)
	if err != nil {
		return nil, err
	}

	enis, err := i.elasticNetworkInterfaces(client)
	if err != nil {
		return nil, err
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

// FactorsForward returns a set of factors that impact the traffic traveling through this point in the direction of source to destination.
func (i EC2Instance) FactorsForward(_ reach.DomainClientResolver, _ *reach.Edge) ([]reach.Factor, error) {
	f := i.newInstanceStateFactor()
	return []reach.Factor{f}, nil
}

// FactorsReturn returns a set of factors that impact the traffic traveling through this point in the direction of destination to source.
func (i EC2Instance) FactorsReturn(_ reach.DomainClientResolver, _ *reach.Edge) ([]reach.Factor, error) {
	f := i.newInstanceStateFactor()
	return []reach.Factor{f}, nil
}

// ———— Implementing IPAddressable ————

// IPs returns the set of IP addresses associated with this resource. This includes both IP addresses known directly by this resource's network interfaces and IP addresses (such as public IPv4 addresses) that are translated (i.e. NAT) to an address associated with this resource's network interface.
func (i EC2Instance) IPs(resolver reach.DomainClientResolver) ([]net.IP, error) {
	client, err := unpackDomainClient(resolver)
	if err != nil {
		return nil, err
	}

	enis, err := i.elasticNetworkInterfaces(client)
	if err != nil {
		return nil, err
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

// InterfaceIPs returns the set of IP addresses that are directly associated with this resource's network interface. Addresses used to reach this resource via address translation are not included.
func (i EC2Instance) InterfaceIPs(resolver reach.DomainClientResolver) ([]net.IP, error) {
	client, err := unpackDomainClient(resolver)
	if err != nil {
		return nil, err
	}

	enis, err := i.elasticNetworkInterfaces(client)
	if err != nil {
		return nil, err
	}

	var ips []net.IP
	for _, eni := range enis {
		ips = append(ips, eni.ownedIPs()...)
	}

	return ips, nil
}

// ———— Supporting methods ————

func (i EC2Instance) firstPointTuples(resolver reach.DomainClientResolver, destinationIPs []net.IP) ([]reach.IPTuple, error) {
	// If the traffic originates from this EC2 instance, any of its owned IP addresses could be used as source.
	// (Technically, any IP address could be used as source, as long as src/dst check is off, but currently we have no way to inform the Tracer about scenarios like this.)

	srcIPs, err := i.InterfaceIPs(resolver)
	if err != nil {
		return nil, err
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

// FactorKindInstanceState specifies the unique name for the EC2 instance state of factor.
const FactorKindInstanceState = "InstanceState"

func (i EC2Instance) newInstanceStateFactor() reach.Factor {
	var traffic reach.TrafficContent

	if i.isRunning() {
		traffic = reach.NewTrafficContentForAllTraffic()
	} else {
		traffic = reach.NewTrafficContentForNoTraffic()
	}

	return reach.Factor{
		Kind:     FactorKindInstanceState,
		Resource: i.Ref(),
		Traffic:  traffic,
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
	var enis []ElasticNetworkInterface

	for _, id := range eniIDs {
		eni, err := client.ElasticNetworkInterface(id)
		if err != nil {
			return nil, err
		}

		enis = append(enis, *eni)
	}

	return enis, nil
}
