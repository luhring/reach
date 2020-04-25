package aws

import (
	"errors"
	"fmt"
	"net"
	"strings"

	"github.com/luhring/reach/reach"
)

// ResourceKindElasticNetworkInterface specifies the unique name for the elastic network interface kind of resource.
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

// Name returns the elastic network interface's ID, and, if available, its name tag value.
func (eni ElasticNetworkInterface) Name() string {
	if name := strings.TrimSpace(eni.NameTag); name != "" {
		return fmt.Sprintf("\"%s\" (%s)", name, eni.ID)
	}
	return eni.ID
}

// Resource returns the elastic network interface converted to a generalized Reach resource.
func (eni ElasticNetworkInterface) Resource() reach.Resource {
	return reach.Resource{
		Kind:       ResourceKindElasticNetworkInterface,
		Properties: eni,
	}
}

// ResourceReference returns a resource reference to uniquely identify the elastic network interface.
func (eni ElasticNetworkInterface) ResourceReference() reach.ResourceReference {
	return reach.ResourceReference{
		Domain: ResourceDomainAWS,
		Kind:   ResourceKindElasticNetworkInterface,
		ID:     eni.ID,
	}
}

func (eni ElasticNetworkInterface) Visitable(_ bool) bool {
	return true
}

func (eni ElasticNetworkInterface) Ref() reach.UniversalReference {
	return reach.UniversalReference{
		R: eni.ResourceReference(),
	}
}

func (eni ElasticNetworkInterface) Segments() bool {
	return false
}

func (eni ElasticNetworkInterface) EdgesForward(resolver reach.DomainClientResolver, previousEdge *reach.Edge, _ []net.IP) ([]reach.Edge, error) {
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
		return errors.New("reach does not currently support an Elastic Network Interface being the first point in a path")
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

func (eni ElasticNetworkInterface) FactorsReturn(resolver reach.DomainClientResolver, nextEdge *reach.Edge) ([]reach.Factor, error) {
	panic("implement me!")
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

func (eni ElasticNetworkInterface) IPs(_ reach.DomainClientResolver) ([]net.IP, error) {
	var ips []net.IP

	ips = append(ips, eni.ownedIPs()...)
	if eni.PublicIPv4Address.Equal(nil) == false {
		ips = append(ips, eni.PublicIPv4Address)
	}

	return ips, nil
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
	router, err := NewVPCRouter(client)
	if err != nil {
		return nil, fmt.Errorf("unable to get VPC router: %v", err)
	}
	return router, nil
}
