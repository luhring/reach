package aws

import (
	"errors"
	"fmt"
	"net"
	"strings"

	"github.com/luhring/reach/reach"
)

// ResourceKindElasticNetworkInterface specifies the unique name for the ElasticNetworkInterface kind of resource.
const ResourceKindElasticNetworkInterface reach.Kind = "ElasticNetworkInterface"

// An ElasticNetworkInterface resource representation.
type ElasticNetworkInterface struct {
	ID                   string
	NameTag              string `json:",omitempty"`
	SubnetID             string
	VPCID                string
	SecurityGroupIDs     []string
	PublicIPv4Address    net.IP   `json:",omitempty"`
	PrivateIPv4Addresses []net.IP `json:",omitempty"`
	IPv6Addresses        []net.IP `json:",omitempty"`
	SrcDstCheck          bool
}

// ElasticNetworkInterfaceRef returns a Reference for an ElasticNetworkInterface with the specified ID.
func ElasticNetworkInterfaceRef(id string) reach.Reference {
	return reach.Reference{
		Domain: ResourceDomainAWS,
		Kind:   ResourceKindElasticNetworkInterface,
		ID:     id,
	}
}

// Name returns the ElasticNetworkInterface's ID, and, if available, its name tag value.
func (eni ElasticNetworkInterface) Name() string {
	if name := strings.TrimSpace(eni.NameTag); name != "" {
		return fmt.Sprintf("\"%s\" (%s)", name, eni.ID)
	}
	return eni.ID
}

// Resource returns the ElasticNetworkInterface converted to a generalized Reach resource.
func (eni ElasticNetworkInterface) Resource() reach.Resource {
	return reach.Resource{
		Kind:       ResourceKindElasticNetworkInterface,
		Properties: eni,
	}
}

// ———— Implementing Traceable ————

// Ref returns a Reference for the ElasticNetworkInterface.
func (eni ElasticNetworkInterface) Ref() reach.Reference {
	return ElasticNetworkInterfaceRef(eni.ID)
}

// Visitable returns a boolean to indicate whether a tracer is allowed to add this resource to the path it's currently constructing.
//
// The Visitable method for ElasticNetworkInterface always returns true because there is no limit to the number of times a tracer can visit an ElasticNetworkInterface.
func (eni ElasticNetworkInterface) Visitable(_ bool) bool {
	return true
}

// Segments returns a boolean to indicate whether a tracer should create a new path segment at this point in the path.
//
// The Segments method for ElasticNetworkInterface always returns false because ENIs never perform NAT.
func (eni ElasticNetworkInterface) Segments() bool {
	return false
}

// EdgesForward returns the set of all possible edges forward given this point in a path that a tracer is constructing. EdgesForward returns an empty slice of edges if there are no further points for the specified network traffic to travel as it attempts to reach its intended network destination.
func (eni ElasticNetworkInterface) EdgesForward(resolver reach.DomainClientResolver, previousEdge *reach.Edge, _ *reach.Reference, _ []net.IP) ([]reach.Edge, error) {
	// TODO: Use previousRef for more intelligent detection of incoming traffic's origin

	err := eni.checkNilPreviousEdge(previousEdge)
	if err != nil {
		return nil, fmt.Errorf("unable to generate forward edges: %v", err)
	}

	// Elastic Network Interfaces don't mutate the IP tuple
	tuple := previousEdge.Tuple

	client, err := unpackDomainClient(resolver)
	if err != nil {
		return nil, fmt.Errorf("unable to get client: %v", err)
	}

	switch eni.flow(tuple, previousEdge.ConnectsInterface) {
	case reach.FlowOutbound:
		return eni.handleEdgeForVPCRouter(client, tuple)
	case reach.FlowInbound:
		return eni.handleEdgeForEC2Instance(client, tuple)
	case reach.FlowDropped:
		return nil, nil
	default:
		return nil, fmt.Errorf("cannot determine direction of flow for tuple (%v)", tuple)
	}
}

// FactorsForward returns a set of factors that impact the traffic traveling through this point in the direction of source to destination.
func (eni ElasticNetworkInterface) FactorsForward(resolver reach.DomainClientResolver, previousEdge *reach.Edge) ([]reach.Factor, error) {
	err := eni.checkNilPreviousEdge(previousEdge)
	if err != nil {
		return nil, fmt.Errorf("unable to generate forward edges: %v", err)
	}

	client, err := unpackDomainClient(resolver)
	if err != nil {
		return nil, fmt.Errorf("unable to get client: %v", err)
	}

	var factors []reach.Factor

	sgRulesFactor, err := eni.securityGroupRulesFactor(client, *previousEdge)
	if err != nil {
		return nil, fmt.Errorf("unable to determine security group rules factors: %v", err)
	}
	factors = append(factors, *sgRulesFactor)

	return factors, nil
}

// FactorsReturn returns a set of factors that impact the traffic traveling through this point in the direction of destination to source.
func (eni ElasticNetworkInterface) FactorsReturn(_ reach.DomainClientResolver, _ *reach.Edge) ([]reach.Factor, error) {
	panic("implement me!")
}

// ———— Implementing IPAddressable ————

// IPs returns the set of IP addresses associated with this resource. This includes both IP addresses known directly by this network interfaces and IP addresses (such as public IPv4 addresses) that are translated (i.e. NAT) to an address associated with this network interface.
func (eni ElasticNetworkInterface) IPs(_ reach.DomainClientResolver) ([]net.IP, error) {
	var ips []net.IP

	ips = append(ips, eni.ownedIPs()...)
	if eni.PublicIPv4Address.Equal(nil) == false {
		ips = append(ips, eni.PublicIPv4Address)
	}

	return ips, nil
}

// TODO: Either implement method InterfaceIPs, or modify interface

// ———— Supporting methods ————

func (eni ElasticNetworkInterface) flow(tuple reach.IPTuple, previousEdgeConnectsInterface bool) reach.Flow {
	if eni.owns(tuple.Dst) {
		return reach.FlowInbound
	}

	if eni.owns(tuple.Src) {
		return reach.FlowOutbound
	}

	if eni.SrcDstCheck == false {
		if previousEdgeConnectsInterface {
			return reach.FlowOutbound
		}

		return reach.FlowInbound
	}

	// SrcDstCheck is on, but neither of the IPs in the tuple belongs to this ENI.
	// This traffic would be dropped by the ENI.
	return reach.FlowDropped
}

func (eni ElasticNetworkInterface) checkNilPreviousEdge(previousEdge *reach.Edge) error {
	if previousEdge == nil {
		return errors.New("reach does not support an Elastic Network Interface being the first point in a path")
	}
	return nil
}

func (eni ElasticNetworkInterface) handleEdgeForVPCRouter(client DomainClient, lastTuple reach.IPTuple) ([]reach.Edge, error) {
	router, err := eni.connectedVPCRouter(client)
	if err != nil {
		return nil, fmt.Errorf("cannot produce forward edge: %v", err)
	}
	edge := reach.Edge{
		Tuple:             lastTuple,
		EndRef:            router.Ref(),
		ConnectsInterface: false,
	}
	return []reach.Edge{edge}, nil
}

func (eni ElasticNetworkInterface) handleEdgeForEC2Instance(client DomainClient, lastTuple reach.IPTuple) ([]reach.Edge, error) {
	ec2, err := eni.connectedEC2Instance(client)
	if err != nil {
		return nil, fmt.Errorf("cannot produce forward edge: %v", err)
	}
	edge := reach.Edge{
		Tuple:             lastTuple,
		EndRef:            ec2.Ref(),
		ConnectsInterface: true,
	}
	return []reach.Edge{edge}, nil
}

func (eni ElasticNetworkInterface) securityGroups(client DomainClient) ([]SecurityGroup, error) {
	sgs := make([]SecurityGroup, len(eni.SecurityGroupIDs))
	for _, id := range eni.SecurityGroupIDs {
		sg, err := client.SecurityGroup(id)
		if err != nil {
			return nil, fmt.Errorf("unable to get security group (id: %s): %v", id, err)
		}
		sgs = append(sgs, *sg)
	}
	return sgs, nil
}

func (eni ElasticNetworkInterface) ownedIPs() []net.IP {
	var ips []net.IP

	ips = append(ips, eni.PrivateIPv4Addresses...)
	ips = append(ips, eni.IPv6Addresses...)

	return ips
}

func (eni ElasticNetworkInterface) owns(ip net.IP) bool {
	for _, ownedIP := range eni.ownedIPs() {
		if ownedIP.Equal(ip) {
			return true
		}
	}

	return false
}

func (eni ElasticNetworkInterface) connectedEC2Instance(client DomainClient) (*EC2Instance, error) {
	ec2, err := client.EC2InstanceByENI(eni.ID)
	if err != nil {
		return nil, fmt.Errorf("unable to get ref of connected EC2 instance: %v", err)
	}
	return ec2, nil
}

func (eni ElasticNetworkInterface) connectedVPCRouter(client DomainClient) (*VPCRouter, error) {
	router, err := NewVPCRouter(client, eni.VPCID)
	if err != nil {
		return nil, fmt.Errorf("unable to get connected VPC router: %v", err)
	}
	return router, nil
}
