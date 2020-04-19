package aws

import "github.com/luhring/reach/reach"

type VPCRouter struct {
	VPC VPC
}

func NewVPCRouter(resources ResourceGetter) (*VPCRouter, error) {

}

func (r VPCRouter) Ref() reach.InfrastructureReference {
	return reach.InfrastructureReference{
		Implicit: true,
		R:        r.VPC.ResourceRef(),
	}
}
