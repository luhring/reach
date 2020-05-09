package aws

import (
	"fmt"

	"github.com/luhring/reach/reach"
)

// ReferenceResolver is the AWS-specific implementation of the interface reach.ReferenceResolver. This ReferenceResolver can resolve only AWS references.
type ReferenceResolver struct {
	client DomainClient
}

// NewReferenceResolver returns a pointer a new ReferenceResolver. If NewReferenceResolver is unable to find an aws.DomainClient using the provided clientResolver, it returns an error.
func NewReferenceResolver(clientResolver reach.DomainClientResolver) (*ReferenceResolver, error) {
	client, err := unpackDomainClient(clientResolver)
	if err != nil {
		return nil, fmt.Errorf("cannot create new AWS ReferenceResolver: %v", err)
	}

	return &ReferenceResolver{client: client}, nil
}

// Resolve returns a Resource for the specified Reference. Resolve returns the error if the Reference does not specify the AWS domain or if there is an error encountered while querying for the resource itself.
func (r *ReferenceResolver) Resolve(ref reach.Reference) (*reach.Resource, error) {
	if ref.Domain != ResourceDomainAWS {
		return nil, fmt.Errorf("%s resolver cannot resolve references for domain '%s'", ResourceDomainAWS, ref.Domain)
	}

	var get func(id string) (reach.Resourceable, error)

	switch ref.Kind {
	case ResourceKindEC2Instance:
		get = func(id string) (reach.Resourceable, error) {
			return r.client.EC2Instance(id)
		}
	case ResourceKindElasticNetworkInterface:
		get = func(id string) (reach.Resourceable, error) {
			return r.client.ElasticNetworkInterface(id)
		}
	case ResourceKindInternetGateway:
		get = func(id string) (reach.Resourceable, error) {
			return r.client.InternetGateway(id)
		}
	case ResourceKindNATGateway:
		get = func(id string) (reach.Resourceable, error) {
			return r.client.NATGateway(id)
		}
	case ResourceKindNetworkACL:
		get = func(id string) (reach.Resourceable, error) {
			return r.client.NetworkACL(id)
		}
	case ResourceKindRouteTable:
		get = func(id string) (reach.Resourceable, error) {
			return r.client.RouteTable(id)
		}
	case ResourceKindSecurityGroup:
		get = func(id string) (reach.Resourceable, error) {
			return r.client.SecurityGroup(id)
		}
	case ResourceKindSecurityGroupReference:
		get = func(id string) (reach.Resourceable, error) {
			// TODO: Handle accountID
			return r.client.SecurityGroupReference(id, "")
		}
	case ResourceKindSubnet:
		get = func(id string) (reach.Resourceable, error) {
			return r.client.Subnet(id)
		}
	case ResourceKindVPC:
		get = func(id string) (reach.Resourceable, error) {
			return r.client.VPC(id)
		}
	case ResourceKindVPCRouter:
		get = func(id string) (reach.Resourceable, error) {
			return NewVPCRouter(r.client, id)
		}
	default:
		return nil, fmt.Errorf("%s resolver encountered an unexpected resource kind '%s'", ResourceDomainAWS, ref.Kind)
	}

	result, err := get(ref.ID)
	if err != nil {
		return nil, fmt.Errorf("%s resource resolution failed (ref: %s): %v", ResourceDomainAWS, ref, err)
	}

	resource := result.Resource()
	return &resource, nil
}
