package aws

import (
	"fmt"

	"github.com/luhring/reach/reach"
)

type ReferenceResolver struct {
	client DomainClient
}

func NewReferenceResolver(clientResolver reach.DomainClientResolver) (*ReferenceResolver, error) {
	client, err := unpackDomainClient(clientResolver)
	if err != nil {
		return nil, fmt.Errorf("cannot create new AWS ReferenceResolver: %v", err)
	}

	return &ReferenceResolver{client: client}, nil
}

func (r *ReferenceResolver) Resolve(ref reach.UniversalReference) (*reach.Resource, error) {
	if ref.Domain != ResourceDomainAWS {
		return nil, fmt.Errorf("%s resolver cannot resolve references for domain '%s'", ResourceDomainAWS, ref.Domain)
	}

	switch ref.Kind {
	case ResourceKindEC2Instance:
		ec2Instance, err := r.client.EC2Instance(ref.ID)
		if err != nil {
			return nil, err
		}
		resource := ec2Instance.Resource()
		return &resource, nil
	case ResourceKindElasticNetworkInterface:
		eni, err := r.client.ElasticNetworkInterface(ref.ID)
		if err != nil {
			return nil, err
		}
		resource := eni.Resource()
		return &resource, nil
	case ResourceKindInternetGateway:
		igw, err := r.client.InternetGateway(ref.ID)
		if err != nil {
			return nil, err
		}
		resource := igw.Resource()
		return &resource, nil
	case ResourceKindNATGateway:
		natgw, err := r.client.NATGateway(ref.ID)
		if err != nil {
			return nil, err
		}
		resource := natgw.Resource()
		return &resource, nil
	case ResourceKindNetworkACL:
		nacl, err := r.client.NetworkACL(ref.ID)
		if err != nil {
			return nil, err
		}
		resource := nacl.Resource()
		return &resource, nil
	case ResourceKindRouteTable:
		rt, err := r.client.RouteTable(ref.ID)
		if err != nil {
			return nil, err
		}
		resource := rt.Resource()
		return &resource, nil
	case ResourceKindSecurityGroup:
		sg, err := r.client.SecurityGroup(ref.ID)
		if err != nil {
			return nil, err
		}
		resource := sg.Resource()
		return &resource, nil
	case ResourceKindSecurityGroupReference:
		sgRef, err := r.client.SecurityGroupReference(ref.ID, "") // TODO: Handle accountID
		if err != nil {
			return nil, err
		}
		resource := sgRef.Resource()
		return &resource, nil
	case ResourceKindSubnet:
		subnet, err := r.client.Subnet(ref.ID)
		if err != nil {
			return nil, err
		}
		resource := subnet.Resource()
		return &resource, nil
	case ResourceKindVPC:
		vpc, err := r.client.VPC(ref.ID)
		if err != nil {
			return nil, err
		}
		resource := vpc.Resource()
		return &resource, nil
	case ResourceKindVPCRouter:
		vpcRouter, err := NewVPCRouter(r.client, ref.ID)
		if err != nil {
			return nil, err
		}
		resource := vpcRouter.Resource()
		return &resource, nil
	}

	return nil, fmt.Errorf("%s resolver encountered an unexpected resource kind '%s'", ResourceDomainAWS, ref.Kind)
}
