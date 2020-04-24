package api

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"

	"github.com/luhring/reach/reach"
)

// DomainClient implements an AWS DomainClient using the AWS API (via the AWS SDK).
type DomainClient struct {
	session *session.Session
	ec2     *ec2.EC2
}

// NewDomainClient returns a reference to a new DomainClient for the AWS API.
func NewDomainClient() *DomainClient {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	})) // TODO: Don't call session.Must —- return error, and don't panic, this is a library after all!

	ec2Client := ec2.New(sess)

	return &DomainClient{
		session: sess,
		ec2:     ec2Client,
	}
}

func nameTag(tags []*ec2.Tag) string {
	if tags != nil && len(tags) > 0 {
		for _, tag := range tags {
			if aws.StringValue(tag.Key) == "Name" {
				return aws.StringValue(tag.Value)
			}
		}
	}

	return ""
}

func ensureSingleResult(resultSetLength int, entity reach.Kind, id string) error {
	if resultSetLength == 0 {
		return fmt.Errorf("AWS API did not return a %s for ID '%s'", entity, id)
	}

	if resultSetLength > 1 {
		return fmt.Errorf("AWS API returned more than one %s for ID '%s'", entity, id)
	}

	return nil
}

func convertAWSIPProtocolStringToProtocol(ipProtocol *string) (reach.Protocol, error) {
	if ipProtocol == nil {
		return 0, errors.New("unexpected nil ipProtocol")
	}

	protocolString := strings.ToLower(aws.StringValue(ipProtocol))

	if p, err := strconv.ParseInt(protocolString, 10, 64); err == nil {
		var protocol = reach.Protocol(p)
		return protocol, nil
	}

	var protocolNumber reach.Protocol

	switch protocolString {
	case "tcp":
		protocolNumber = reach.ProtocolTCP
	case "udp":
		protocolNumber = reach.ProtocolUDP
	case "icmp":
		protocolNumber = reach.ProtocolICMPv4
	case "icmpv6":
		protocolNumber = reach.ProtocolICMPv6
	default:
		return 0, errors.New("unrecognized ipProtocol value")
	}

	return protocolNumber, nil
}
