package cmd

import (
	"strings"

	"github.com/luhring/reach/reach"
	"github.com/luhring/reach/reach/aws"
	"github.com/luhring/reach/reach/generic"
	"github.com/luhring/reach/reach/reacherr"
)

func resolveSubject(input string, domains reach.DomainClientResolver) (*reach.Subject, error) {
	q := getQualifiedSubject(input)
	if q != nil {
		subject, err := resolveSubjectExplicitly(*q, domains)
		if err != nil {
			return nil, err
		}
		return subject, nil
	}

	subject, err := resolveSubjectImplicitly(input, domains)
	if err != nil {
		return nil, err
	}
	return subject, nil
}

func getQualifiedSubject(input string) *qualifiedSubject {
	inputSegments := strings.SplitN(input, ":", 2)

	if inputSegments == nil || len(inputSegments) < 2 {
		return nil
	}

	return &qualifiedSubject{
		typePrefix: inputSegments[0],
		identifier: inputSegments[1],
	}
}

type qualifiedSubject struct {
	typePrefix string
	identifier string
}

func resolveSubjectImplicitly(input string, domains reach.DomainClientResolver) (*reach.Subject, error) {
	logger.Debug("resolving subject implicitly", "input", input)

	// 1. Try IP address format.
	err := generic.CheckIPAddress(input)
	if err == nil {
		logger.Debug("subject resolution input appears to be an IP address", "input", input)
		return generic.NewIPAddressSubject(input), nil
	}

	// 2. Try hostname format.
	err = generic.CheckHostname(input)
	if err == nil {
		logger.Debug("subject resolution input appears to be a hostname", "input", input)
		return generic.NewHostnameSubject(input), nil
	}

	// 3. Try EC2 fuzzy matching.
	return aws.ResolveEC2InstanceSubject(input, domains)
}

func resolveSubjectExplicitly(subj qualifiedSubject, domains reach.DomainClientResolver) (*reach.Subject, error) {
	logger.Debug("resolving subject explicitly", "prefix", subj.typePrefix, "identifier", subj.identifier)

	switch subj.typePrefix {
	case "ip":
		return generic.ResolveIPAddressSubject(subj.identifier)
	case "host":
		return generic.ResolveHostnameSubject(subj.identifier)
	case "ec2":
		return aws.ResolveEC2InstanceSubject(subj.identifier, domains)
	default:
		return nil, reacherr.New(nil, "unable to resolve subject because subject type '%s' is not recognized", subj.identifier, subj.typePrefix)
	}
}
