package apiclient

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"

	"github.com/luhring/reach/reach"
	reachAWS "github.com/luhring/reach/reach/aws"
	"github.com/luhring/reach/reach/reacherr"
	"github.com/luhring/reach/reach/traffic"
)

var _ reachAWS.DomainClient = (*DomainClient)(nil)

// DomainClient implements an AWS DomainClient using the AWS API (via the AWS SDK).
type DomainClient struct {
	session *session.Session
	ec2     *ec2.EC2
	cache   reach.Cache
}

// NewDomainClient returns a reference to a new DomainClient for the AWS API.
func NewDomainClient(cache reach.Cache) (*DomainClient, error) {
	sess, err := session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	})
	if err != nil {
		msg := "unable to start an AWS SDK session"
		if awsErr, ok := err.(awserr.Error); ok {
			msg += ": " + awsErr.Message()
			return nil, reacherr.New(err, msg)
		}

		return nil, err
	}

	ec2Client := ec2.New(sess)

	return &DomainClient{
		session: sess,
		ec2:     ec2Client,
		cache:   cache,
	}, nil
}

func (client *DomainClient) cacheResource(r reach.Referable) {
	client.cache.Put(r.Ref().String(), r)
}

func (client *DomainClient) cachedResource(ref reach.Reference) interface{} {
	return client.cache.Get(ref.String())
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
		err := reacherr.New(nil, "no %s resources found for ID '%s'", entity, id)
		return err
	}

	if resultSetLength > 1 {
		err := reacherr.New(nil, "more than one %s resources found for ID '%s'", entity, id)
		return err
	}

	return nil
}

func convertAWSIPProtocolStringToProtocol(ipProtocol *string) (traffic.Protocol, error) {
	if ipProtocol == nil {
		return 0, errors.New("unexpected nil ipProtocol")
	}

	protocolString := strings.ToLower(aws.StringValue(ipProtocol))

	if p, err := strconv.ParseInt(protocolString, 10, 64); err == nil {
		var protocol = traffic.Protocol(p)
		return protocol, nil
	}

	var protocolNumber traffic.Protocol

	switch protocolString {
	case "tcp":
		protocolNumber = traffic.ProtocolTCP
	case "udp":
		protocolNumber = traffic.ProtocolUDP
	case "icmp":
		protocolNumber = traffic.ProtocolICMPv4
	case "icmpv6":
		protocolNumber = traffic.ProtocolICMPv6
	default:
		return 0, fmt.Errorf("unrecognized ipProtocol value: %v", ipProtocol)
	}

	return protocolNumber, nil
}
