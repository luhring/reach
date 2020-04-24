package cmd

import (
	"fmt"
	"io"
	"strings"

	"github.com/luhring/reach/reach"
	"github.com/luhring/reach/reach/aws"
	"github.com/luhring/reach/reach/generic"
)

func resolveSubject(input string, progressWriter io.Writer, domains reach.DomainProvider) (*reach.Subject, error) {
	q := getQualifiedSubject(input)

	if q != nil {
		return resolveSubjectExplicitly(*q, domains)
	} else {
		return resolveSubjectImplicitly(input, progressWriter, domains)
	}
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

func resolveSubjectImplicitly(input string, progressWriter io.Writer, domains reach.DomainProvider) (*reach.Subject, error) {
	// 1. Try IP address format.
	err := generic.CheckIPAddress(input)
	if err == nil {
		_, _ = fmt.Fprintf(progressWriter, "'%s' is being interpreted as an IP address\n", input)
		return generic.NewIPAddressSubject(input), nil
	}

	// 2. Try hostname format.
	err = generic.CheckHostname(input)
	if err == nil {
		_, _ = fmt.Fprintf(progressWriter, "'%s' is being interpreted as a hostname\n", input)
		return generic.NewHostnameSubject(input), nil
	}

	// 3. Try EC2 fuzzy matching.
	return aws.ResolveEC2InstanceSubject(input, domains)
}

func resolveSubjectExplicitly(qualifiedSubject qualifiedSubject, domains reach.DomainProvider) (*reach.Subject, error) {
	switch qualifiedSubject.typePrefix {
	case "ip":
		return generic.ResolveIPAddressSubject(qualifiedSubject.identifier)
	case "host":
		return generic.ResolveHostnameSubject(qualifiedSubject.identifier)
	case "ec2":
		return aws.ResolveEC2InstanceSubject(qualifiedSubject.identifier, domains)
	default:
		return nil, fmt.Errorf("unable to resolve subject with identifier '%s' because subject type prefix '%s' is not recognized", qualifiedSubject.identifier, qualifiedSubject.typePrefix)
	}
}
