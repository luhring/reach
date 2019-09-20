package api

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

type ResourceProvider struct {
	session *session.Session
	ec2     *ec2.EC2
}

func NewResourceProvider() *ResourceProvider {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	})) // TODO: Don't call session.Must —- return error, and don't panic, this is a lib after all!

	ec2Client := ec2.New(sess)

	return &ResourceProvider{
		session: sess,
		ec2:     ec2Client,
	}
}

func getNameTag(tags []*ec2.Tag) string {
	if tags != nil && len(tags) > 0 {
		for _, tag := range tags {
			if aws.StringValue(tag.Key) == "Name" {
				return aws.StringValue(tag.Value)
			}
		}
	}

	return ""
}

func ensureSingleResult(resultSetLength int, entity, id string) error {
	if resultSetLength == 0 {
		return fmt.Errorf("AWS API did not return a %s for ID '%s'", entity, id)
	}

	if resultSetLength > 1 {
		return fmt.Errorf("AWS API returned more than one %s for ID '%s'", entity, id)
	}

	return nil
}
