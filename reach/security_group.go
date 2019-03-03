package reach

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
)

type SecurityGroup struct {
	ID            string
	Name          string
	VPCID         string
	InboundRules  []*SecurityGroupRule
	OutboundRules []*SecurityGroupRule
}

func newSecurityGroup(group *ec2.SecurityGroup) (*SecurityGroup, error) {
	inboundRules := make([]*SecurityGroupRule, len(group.IpPermissions))
	for i, r := range group.IpPermissions {
		newRule, err := newSecurityGroupRule(r)
		if err != nil {
			return nil, fmt.Errorf("error: unable to ingest inbound security group rule at index %v: %v", i, err)
		}

		inboundRules[i] = newRule
	}

	outboundRules := make([]*SecurityGroupRule, len(group.IpPermissionsEgress))
	for i, r := range group.IpPermissionsEgress {
		newRule, err := newSecurityGroupRule(r)
		if err != nil {
			return nil, fmt.Errorf("error: unable to ingest outbound security group rule at index %v: %v", i, err)
		}

		outboundRules[i] = newRule
	}

	return &SecurityGroup{
		ID:            aws.StringValue(group.GroupId),
		Name:          aws.StringValue(group.GroupName),
		VPCID:         aws.StringValue(group.VpcId),
		InboundRules:  inboundRules,
		OutboundRules: outboundRules,
	}, nil
}

func (sg *SecurityGroup) longName() string {
	if len(sg.Name) >= 1 {
		return fmt.Sprintf("\"%v\" (%v)", sg.Name, sg.ID)
	}

	return sg.ID
}
