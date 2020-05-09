package aws

import (
	"fmt"

	"github.com/luhring/reach/reach"
)

// ResourceKindSecurityGroup specifies the unique name for the security group kind of resource.
const ResourceKindSecurityGroup reach.Kind = "SecurityGroup"

// A SecurityGroup resource representation.
type SecurityGroup struct {
	ID            string
	NameTag       string
	GroupName     string
	VPCID         string
	InboundRules  []SecurityGroupRule
	OutboundRules []SecurityGroupRule
}

// Resource returns the security group converted to a generalized Reach resource.
func (sg SecurityGroup) Resource() reach.Resource {
	return reach.Resource{
		Kind:       ResourceKindSecurityGroup,
		Properties: sg,
	}
}

func (sg SecurityGroup) Ref() reach.Reference {
	return SecurityGroupRef(sg.ID)
}

// Name returns the security group's ID, and, if available, its name tag value (or group name).
func (sg SecurityGroup) Name() string {
	var name string

	if sg.NameTag != "" {
		name = sg.NameTag
	} else if sg.GroupName != "" {
		name = sg.GroupName
	}

	if name != "" {
		return fmt.Sprintf("%s (%s)", name, sg.ID)
	}

	return sg.ID
}

func SecurityGroupRef(id string) reach.Reference {
	return reach.Reference{
		Domain: ResourceDomainAWS,
		Kind:   ResourceKindSecurityGroup,
		ID:     id,
	}
}
