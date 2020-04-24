package aws

import "github.com/luhring/reach/reach"

type VPCRouter struct {
	VPC VPC
}

func NewVPCRouter(_ DomainClient) (*VPCRouter, error) {
	panic("implement me!")
}

func (r VPCRouter) Ref() reach.UniversalReference {
	return reach.UniversalReference{
		Implicit: true,
		R:        r.VPC.ResourceReference(),
	}
}
