package cmd

import (
	"fmt"
	"strings"

	"github.com/luhring/reach/reach"
	"github.com/luhring/reach/reach/aws"
	"github.com/luhring/reach/reach/generic"
	"github.com/luhring/reach/reach/reacherr"
)

func resolveSubject(input string, domains reach.DomainClientResolver) (*reach.Subject, error) {
	q := getQualifiedSubject(input)

	if q != nil {
		return resolveSubjectExplicitly(*q, domains)
	}

	return resolveSubjectImplicitly(input, domains)
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
	logger.Info("resolving subject implicitly", "input", input)
	// 1. Try IP address format.
	err := generic.CheckIPAddress(input)
	if err == nil {
		logger.Info("subject resolution input appears to be an IP address", "input", input)
		return generic.NewIPAddressSubject(input), nil
	}

	// 2. Try hostname format.
	err = generic.CheckHostname(input)
	if err == nil {
		logger.Info("subject resolution input appears to be a hostname", "input", input)
		return generic.NewHostnameSubject(input), nil
	}

	// 3. Try EC2 fuzzy matching.
	return aws.ResolveEC2InstanceSubject(input, domains)
}

func resolveSubjectExplicitly(qualifiedSubject qualifiedSubject, domains reach.DomainClientResolver) (*reach.Subject, error) {
	logger.Info("resolving subject explicitly", "qualifiedSubject", fmt.Sprintf("%+v", qualifiedSubject))
	switch qualifiedSubject.typePrefix {
	case "ip":
		return generic.ResolveIPAddressSubject(qualifiedSubject.identifier)
	case "host":
		return generic.ResolveHostnameSubject(qualifiedSubject.identifier)
	case "ec2":
		return aws.ResolveEC2InstanceSubject(qualifiedSubject.identifier, domains)
	default:
		return nil, reacherr.New(nil, "unable to resolve subject with identifier '%s' because subject type '%s' is not recognized", qualifiedSubject.identifier, qualifiedSubject.typePrefix)
	}
}
