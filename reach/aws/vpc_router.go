package aws

import "github.com/luhring/reach/reach"

type VPCRouter struct {
	VPC VPC
}

func NewVPCRouter(resources ResourceGetter) (*VPCRouter, error) {

}

func (r VPCRouter) Ref() reach.UniversalReference {
	return reach.UniversalReference{
		Implicit: true,
		R:        r.VPC.ResourceRef(),
	}
}
