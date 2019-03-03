package reach

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
)

type SecurityGroupReference struct {
	ID     string
	Name   string
	UserID string
	VPCID  string
}

func newSecurityGroupReference(p *ec2.UserIdGroupPair) *SecurityGroupReference {
	return &SecurityGroupReference{
		ID:     aws.StringValue(p.GroupId),
		Name:   aws.StringValue(p.GroupName),
		UserID: aws.StringValue(p.UserId),
		VPCID:  aws.StringValue(p.VpcId),
	}
}
