package aws

import (
	"errors"
	"fmt"
	"net"
	"strings"

	"github.com/luhring/reach/reach"
)

// ResourceKindElasticNetworkInterface specifies the unique name for the elastic network interface kind of resource.
const ResourceKindElasticNetworkInterface = "ElasticNetworkInterface"

// An ElasticNetworkInterface resource representation.
type ElasticNetworkInterface struct {
	ID                   string
	NameTag              string `json:"NameTag,omitempty"`
	SubnetID             string
	VPCID                string
	SecurityGroupIDs     []string
	PublicIPv4Address    net.IP   `json:"PublicIPv4Address,omitempty"`
	PrivateIPv4Addresses []net.IP `json:"PrivateIPv4Addresses,omitempty"`
	IPv6Addresses        []net.IP `json:"IPv6Addresses,omitempty"`
	SrcDstCheck          bool
}

// ElasticNetworkInterfaceFromNetworkPoint extracts the ElasticNetworkInterface from the lineage of the specified network point.
func ElasticNetworkInterfaceFromNetworkPoint(point reach.NetworkPoint, rc *reach.ResourceCollection) *ElasticNetworkInterface {
	for _, ancestor := range point.Lineage { // assumes there will only be one ENI among the ancestors
		if ancestor.Domain == ResourceDomainAWS && ancestor.Kind == ResourceKindElasticNetworkInterface {
			eni := rc.Get(ancestor).Properties.(ElasticNetworkInterface)
			return &eni
		}
	}

	return nil
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

// Dependencies returns a collection of the elastic network interface's resource dependencies.
func (eni ElasticNetworkInterface) Dependencies(provider ResourceGetter) (*reach.ResourceCollection, error) {
	rc := reach.NewResourceCollection()

	subnet, err := provider.Subnet(eni.SubnetID)
	if err != nil {
		return nil, err
	}
	rc.Put(reach.ResourceReference{
		Domain: ResourceDomainAWS,
		Kind:   ResourceKindSubnet,
		ID:     subnet.ID,
	}, subnet.ToResource())

	subnetDependencies, err := subnet.Dependencies(provider)
	if err != nil {
		return nil, err
	}
	rc.Merge(subnetDependencies)

	vpc, err := provider.VPC(eni.VPCID)
	if err != nil {
		return nil, err
	}
	rc.Put(reach.ResourceReference{
		Domain: ResourceDomainAWS,
		Kind:   ResourceKindVPC,
		ID:     vpc.ID,
	}, vpc.ToResource())

	for _, sgID := range eni.SecurityGroupIDs {
		sg, err := provider.SecurityGroup(sgID)
		if err != nil {
			return nil, err
		}
		rc.Put(reach.ResourceReference{
			Domain: ResourceDomainAWS,
			Kind:   ResourceKindSecurityGroup,
			ID:     sg.ID,
		}, sg.ToResource())

		sgDependencies, err := sg.Dependencies(provider)
		if err != nil {
			return nil, err
		}
		rc.Merge(sgDependencies)
	}

	return rc, nil
}

func (eni ElasticNetworkInterface) Visitable(_ bool) bool {
	return true
}

func (eni ElasticNetworkInterface) Ref() reach.InfrastructureReference {
	return reach.InfrastructureReference{
		R: eni.ResourceReference(),
	}
}

func (eni ElasticNetworkInterface) Segments() bool {
	return false
}

func (eni ElasticNetworkInterface) ForwardEdges(
	previousEdge *reach.Edge,
	domains reach.DomainProvider,
	_ []net.IP,
) ([]reach.Edge, error) {
	err := eni.checkNilPreviousEdge(previousEdge)
	if err != nil {
		return nil, fmt.Errorf("unable to generate forward edges: %v", err)
	}

	// Elastic Network Interfaces don't mutate the IP tuple
	tuple := previousEdge.Tuple

	resources, err := unpackResourceGetter(domains)
	if err != nil {
		return nil, fmt.Errorf("unable to get resources: %v", err)
	}

	switch eni.flow(tuple, previousEdge.ConnectsInterface) {
	case reach.FlowOutbound:
		return eni.handleEdgeForVPCRouter(tuple, resources)
	case reach.FlowInbound:
		return eni.handleEdgeForEC2Instance(tuple, resources)
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

func (eni ElasticNetworkInterface) handleEdgeForVPCRouter(lastTuple reach.IPTuple, resources ResourceGetter) ([]reach.Edge, error) {
	router, err := eni.connectedVPCRouter(resources)
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

func (eni ElasticNetworkInterface) handleEdgeForEC2Instance(lastTuple reach.IPTuple, resources ResourceGetter) ([]reach.Edge, error) {
	ec2, err := eni.connectedEC2Instance(resources)
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

func (eni ElasticNetworkInterface) FactorsForward(
	previousEdge *reach.Edge,
	domains reach.DomainProvider,
) ([]reach.Factor, error) {
	err := eni.checkNilPreviousEdge(previousEdge)
	if err != nil {
		return nil, fmt.Errorf("unable to generate forward edges: %v", err)
	}

	resources, err := unpackResourceGetter(domains)
	if err != nil {
		return nil, fmt.Errorf("unable to get resources: %v", err)
	}

	var factors []reach.Factor

	sgRulesFactor, err := eni.securityGroupRulesFactor(resources, *previousEdge)
	if err != nil {
		return nil, fmt.Errorf("unable to determine security group rules factors: %v", err)
	}
	factors = append(factors, *sgRulesFactor)

	return factors, nil
}

func (eni ElasticNetworkInterface) securityGroups(resources ResourceGetter) ([]SecurityGroup, error) {
	sgs := make([]SecurityGroup, len(eni.SecurityGroupIDs))
	for _, id := range eni.SecurityGroupIDs {
		sg, err := resources.SecurityGroup(id)
		if err != nil {
			return nil, fmt.Errorf("unable to get security group (id: %s): %v", id, err)
		}
		sgs = append(sgs, *sg)
	}
	return sgs, nil
}

func (eni ElasticNetworkInterface) IPs(_ reach.DomainProvider) ([]net.IP, error) {
	var ips []net.IP

	ips = append(ips, eni.ownedIPs()...)
	ips = append(ips, eni.PublicIPv4Address)

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

func (eni ElasticNetworkInterface) connectedEC2Instance(resources ResourceGetter) (*EC2Instance, error) {
	ec2, err := resources.EC2InstanceByENI(eni.ID)
	if err != nil {
		return nil, fmt.Errorf("unable to get ref of connected EC2 instance: %v", err)
	}
	return ec2, nil
}

func (eni ElasticNetworkInterface) connectedVPCRouter(resources ResourceGetter) (*VPCRouter, error) {
	router, err := NewVPCRouter(resources)
	if err != nil {
		return nil, fmt.Errorf("unable to get VPC router: %v", err)
	}
	return router, nil
}
