package aws

import (
	"fmt"

	"github.com/luhring/reach/reach"
)

const ResourceKindVPC = "VPC"

type ResourceStore struct {
	resources                []reach.Resource
	ec2Instances             []EC2Instance
	elasticNetworkInterfaces []ElasticNetworkInterface
	networkACLs              []NetworkACL
	routeTables              []RouteTable
	securityGroups           []SecurityGroup
	subnets                  []Subnet
	vpcs                     []VPC
}

func NewResourceStore() *ResourceStore {
	return &ResourceStore{}
}

func (store *ResourceStore) ExportAll() []reach.Resource {
	// TODO: is this the best way to make this data accessible? need to figure out consumer's use
	var resources []reach.Resource

	for _, instance := range store.ec2Instances {
		resource := reach.Resource{
			Kind:       "ec2Instance",
			Properties: instance,
		}

		resources = append(resources, resource)
	}

	// TODO: export other kinds of resources

	return resources
}

func (store *ResourceStore) Store(kind string, resource interface{}) {
	switch kind {
	case ResourceKindEC2Instance:
		store.ec2Instances = append(store.ec2Instances, resource.(EC2Instance))
	case ResourceKindElasticNetworkInterface:
		store.elasticNetworkInterfaces = append(store.elasticNetworkInterfaces, resource.(ElasticNetworkInterface))
	case ResourceKindNetworkACL:
		store.networkACLs = append(store.networkACLs, resource.(NetworkACL))
	case ResourceKindRouteTable:
		store.routeTables = append(store.routeTables, resource.(RouteTable))
	case ResourceKindSecurityGroup:
		store.securityGroups = append(store.securityGroups, resource.(SecurityGroup))
	case ResourceKindSubnet:
		store.subnets = append(store.subnets, resource.(Subnet))
	case ResourceKindVPC:
		store.vpcs = append(store.vpcs, resource.(VPC))
	default:
		panic(fmt.Sprintf("unrecognized resource kind during save: '%s'", kind))
	}
}

func (store *ResourceStore) Retrieve(kind, id string) (interface{}, error) {
	return nil, nil
	// for _, resource := range store.resources {
	// 	if
	// }
}
