package explainer

import (
	"fmt"
	"log"
	"strings"

	"github.com/mgutz/ansi"

	"github.com/luhring/reach/reach"
	"github.com/luhring/reach/reach/aws"
	"github.com/luhring/reach/reach/helper"
)

// An Explainer provides mechanisms to explain the business logic behind analyses to users via natural language.
type Explainer struct {
	analysis reach.Analysis
}

// New returns a reference to a new Explainer.
func New(analysis reach.Analysis) *Explainer {
	return &Explainer{
		analysis: analysis,
	}
}

// Explain returns a natural language representation of the logic used during an analysis to compute the final result.
func (ex *Explainer) Explain() string {
	var outputItems []string
	for _, v := range ex.analysis.NetworkVectors {
		outputItems = append(outputItems, ex.ExplainNetworkVector(v))
	}

	output := ""
	output += strings.Join(outputItems, "\n")

	return output
}

// ExplainNetworkVector returns the part of an analysis explanation that's specific to an individual network vector.
func (ex *Explainer) ExplainNetworkVector(v reach.NetworkVector) string {
	var outputSections []string

	// setting the stage: the source and destination
	var vectorHeader string
	vectorHeader += fmt.Sprintf("%s %s\n", helper.Bold("source:"), ex.NetworkPointName(v.Source))
	vectorHeader += fmt.Sprintf("%s %s\n", helper.Bold("destination:"), ex.NetworkPointName(v.Destination))
	outputSections = append(outputSections, vectorHeader)

	// explain source
	sourceHeader := helper.Bold("source factors:")
	outputSections = append(outputSections, sourceHeader)

	sourceContent := ex.ExplainNetworkPoint(v.Source, v.SourcePerspective())
	outputSections = append(outputSections, helper.Indent(sourceContent, 2))

	// explain destination
	destinationHeader := helper.Bold("destination factors:")
	outputSections = append(outputSections, destinationHeader)

	destinationContent := ex.ExplainNetworkPoint(v.Destination, v.DestinationPerspective())
	outputSections = append(outputSections, helper.Indent(destinationContent, 2))

	// final results
	forwardResults := fmt.Sprintf("%s\n%s", helper.Bold("network traffic allowed from source to destination:"), v.Traffic.ColorStringWithSymbols())
	outputSections = append(outputSections, forwardResults)

	returnResults := fmt.Sprintf("%s\n%s", helper.Bold("network traffic allowed to return from destination to source:"), v.ReturnTraffic.StringWithSymbols())
	outputSections = append(outputSections, returnResults)

	return strings.Join(outputSections, "\n")
}

// ExplainCapabilityChecks returns a report on whether or not Reach's capabilities are sufficient to handle the requested analysis.
func (ex *Explainer) ExplainCapabilityChecks(v reach.NetworkVector) string {
	var outputItems []string
	var checksItems []string

	checksHeader := helper.Bold("analysis capability checks:")
	outputItems = append(outputItems, checksHeader)

	awsEx := aws.NewExplainer(ex.analysis)

	if awsEx.CheckBothInAWS(v) {
		checksItems = append(checksItems, "✓ both source and destination are in AWS")
	} else {
		log.Fatal("source and/or destination is not in AWS, and this is not yet supported")
	}

	if awsEx.CheckBothInSameVPC(v) {
		checksItems = append(checksItems, "✓ both source and destination are in same VPC")
	} else {
		log.Fatal("source and/or destination are not in same VPC, and this is not yet supported")
	}

	if awsEx.CheckBothInSameSubnet(v) {
		checksItems = append(checksItems, "✓ both source and destination are in same subnet")
	} else {
		log.Fatal("source and/or destination are not in same subnet, and this is not yet supported")
	}

	outputItems = append(outputItems, checksItems...)

	return strings.Join(outputItems, "\n")
}

// ExplainNetworkPoint returns the part of an analysis explanation that's specific to an individual network point (within a network vector).
func (ex *Explainer) ExplainNetworkPoint(point reach.NetworkPoint, p reach.Perspective) string {
	if aws.IsUsedByNetworkPoint(point) {
		awsEx := aws.NewExplainer(ex.analysis)
		return awsEx.NetworkPoint(point, p)
	}

	return fmt.Sprintf("unable to explain analysis for network point with IP address '%s'", point.IPAddress)
}

// NetworkPointName returns an understandable string representation of a network point.
func (ex *Explainer) NetworkPointName(point reach.NetworkPoint) string {
	// ignoring errors because it's okay if we can't find a particular kind of AWS resource in the lineage
	eni, _ := aws.GetENIFromLineage(point.Lineage, ex.analysis.Resources)
	ec2Instance, _ := aws.GetEC2InstanceFromLineage(point.Lineage, ex.analysis.Resources)

	output := point.IPAddress.String()

	if eni != nil {
		output = fmt.Sprintf("%s -> %s", eni.Name(), output)

		if ec2Instance != nil {
			output = fmt.Sprintf("%s -> %s", ec2Instance.Name(), output)
		}
	}

	return output
}

// WarningsFromRestrictedReturnPath returns a slice of warning strings based on the input slice of restricted protocols.
func WarningsFromRestrictedReturnPath(restrictedProtocols []reach.RestrictedProtocol) (bool, string) {
	if len(restrictedProtocols) == 0 {
		return false, ""
	}

	var warnings []string

	for _, rp := range restrictedProtocols {
		var warning string

		if rp.Protocol == reach.ProtocolTCP { // We have a specific message based on the knowledge that the protocol is TCP.
			if rp.NoReturnTraffic {
				warning = ansi.Color("All TCP connection attempts will be unsuccessful. No TCP traffic is allowed to return to the source.", "red+b")
			} else {
				warning = ansi.Color("TCP connection attempts might be unsuccessful. TCP traffic is allowed to return to the source only at particular source ports.", "yellow+b")
			}
		} else {
			firstSentence := fmt.Sprintf("%s-based communication might be unsuccessful.", rp.Protocol)

			var secondSentence string

			if rp.NoReturnTraffic {
				secondSentence = fmt.Sprintf("No %s traffic is able to return to the source.", rp.Protocol)
			} else {
				secondSentence = fmt.Sprintf("Some %s traffic is unable to return to the source.", rp.Protocol)
			}

			warning = ansi.Color(fmt.Sprintf("%s %s", firstSentence, secondSentence), "yellow+b")
		}

		warnings = append(warnings, warning)
	}

	return true, "warnings from return traffic obstructions:\n" + strings.Join(warnings, "\n")
}
