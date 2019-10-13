package explainer

import (
	"fmt"
	"log"
	"strings"

	"github.com/luhring/reach/reach"
	"github.com/luhring/reach/reach/aws"
	"github.com/luhring/reach/reach/helper"
)

type Explainer struct {
	analysis reach.Analysis
}

func New(analysis reach.Analysis) *Explainer {
	return &Explainer{
		analysis: analysis,
	}
}

func (ex *Explainer) Explain() string {
	var outputItems []string
	for _, v := range ex.analysis.NetworkVectors {
		outputItems = append(outputItems, ex.ExplainNetworkVector(v))
	}

	output := ""
	output += strings.Join(outputItems, "\n")

	return output
}

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
	results := fmt.Sprintf("%s\n%s", helper.Bold("network traffic allowed from source to destination:"), v.Traffic.ColorStringWithSymbols())
	outputSections = append(outputSections, results)

	return strings.Join(outputSections, "\n")
}

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

func (ex *Explainer) ExplainNetworkPoint(point reach.NetworkPoint, p reach.Perspective) string {
	if aws.IsUsedByNetworkPoint(point) {
		awsEx := aws.NewExplainer(ex.analysis)
		return awsEx.NetworkPoint(point, p)
	}

	return fmt.Sprintf("unable to explain analysis for network point with IP address '%s'", point.IPAddress)
}

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
