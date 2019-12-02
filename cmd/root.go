package cmd

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/luhring/reach/reach"
	"github.com/luhring/reach/reach/analyzer"
	"github.com/luhring/reach/reach/aws"
	"github.com/luhring/reach/reach/aws/api"
	"github.com/luhring/reach/reach/explainer"
	"github.com/luhring/reach/reach/generic"
)

const explainFlag = "explain"
const vectorsFlag = "vectors"
const jsonFlag = "json"
const assertReachableFlag = "assert-reachable"
const assertNotReachableFlag = "assert-not-reachable"

var explain bool
var showVectors bool
var outputJSON bool
var assertReachable bool
var assertNotReachable bool

var rootCmd = &cobra.Command{
	Use:   "reach",
	Short: "reach examines network reachability issues in AWS",
	Long: `reach examines network reachability issues in AWS
See https://github.com/luhring/reach for documentation.`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 2 {
			return errors.New("requires at least two arguments")
		}

		if assertReachable && assertNotReachable {
			return errors.New("cannot assert both reachable and not reachable at the same time")
		}

		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		sourceIdentifier := args[0]
		destinationIdentifier := args[1]

		var awsResourceProvider aws.ResourceProvider = api.NewResourceProvider()

		// Not sure yet if I like this, but I want to be able to package up a collection of resource providers across arbitrary domains.
		// This relies on type assertions downstream in the code, of course.
		providers := map[string]interface{}{
			aws.ResourceDomainAWS: awsResourceProvider,
		}

		source, err := resolveSubject(sourceIdentifier, os.Stderr, providers)
		if err != nil {
			exitWithError(err)
		}
		source.SetRoleToSource()

		destination, err := resolveSubject(destinationIdentifier, os.Stderr, providers)
		if err != nil {
			exitWithError(err)
		}
		destination.SetRoleToDestination()

		if !outputJSON && !explain && !showVectors {
			fmt.Printf("source: %s\ndestination: %s\n\n", source.ID, destination.ID)
		}

		a := analyzer.New()
		analysis, err := a.Analyze(source, destination)
		if err != nil {
			exitWithError(err)
		}

		mergedTraffic, err := analysis.MergedTraffic()
		if err != nil {
			exitWithError(err)
		}

		if outputJSON {
			fmt.Println(analysis.ToJSON())
		} else if explain {
			ex := explainer.New(*analysis)
			fmt.Print(ex.Explain())
		} else if showVectors {
			var vectorOutputs []string

			for _, v := range analysis.NetworkVectors {
				vectorOutputs = append(vectorOutputs, v.String())
			}

			fmt.Print(strings.Join(vectorOutputs, "\n"))
		} else {
			fmt.Print("network traffic allowed from source to destination:" + "\n")
			fmt.Print(mergedTraffic.ColorStringWithSymbols())

			if len(analysis.NetworkVectors) > 1 { // handling this case with care; this view isn't optimized for multi-vector output!
				printMergedResultsWarning()
				warnIfAnyVectorHasRestrictedReturnTraffic(analysis.NetworkVectors)
			} else {
				// calculate merged return traffic
				mergedReturnTraffic, err := analysis.MergedReturnTraffic()
				if err != nil {
					exitWithError(err)
				}

				restrictedProtocols := mergedTraffic.ProtocolsWithRestrictedReturnPath(mergedReturnTraffic)
				if len(restrictedProtocols) > 0 {
					found, warnings := explainer.WarningsFromRestrictedReturnPath(restrictedProtocols)
					if found {
						fmt.Print("\n" + warnings + "\n")
					}
				}
			}
		}

		if assertReachable {
			doAssertReachable(*analysis)
		}

		if assertNotReachable {
			doAssertNotReachable(*analysis)
		}
	},
}

// Execute runs the root command
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		exitWithError(err)
	}
}

func init() {
	rootCmd.Flags().BoolVar(&explain, explainFlag, false, "explain how the configuration was analyzed")
	rootCmd.Flags().BoolVar(&showVectors, vectorsFlag, false, "show allowed traffic in terms of network vectors")
	rootCmd.Flags().BoolVar(&outputJSON, jsonFlag, false, "output full analysis as JSON (overrides other display flags)")
	rootCmd.Flags().BoolVar(&assertReachable, assertReachableFlag, false, "exit non-zero if no traffic is allowed from source to destination")
	rootCmd.Flags().BoolVar(&assertNotReachable, assertNotReachableFlag, false, "exit non-zero if any traffic can reach destination from source")
}

func resolveSubject(identifier string, progressWriter io.Writer, resourceProviders map[string]interface{}) (*reach.Subject, error) {
	identifierSegments := strings.SplitN(identifier, ":", 2)
	if identifierSegments == nil || len(identifierSegments) < 2 { // implicit resolution (subject type was not specified)
		// 1. Try IP address format.
		err := generic.CheckIPAddress(identifier)
		if err == nil {
			_, _ = fmt.Fprintf(progressWriter, "'%s' is being interpreted as an IP address\n", identifier)
			return generic.NewIPAddressSubject(identifier), nil
		}

		// 2. Try hostname format.
		err = generic.CheckHostname(identifier)
		if err == nil {
			_, _ = fmt.Fprintf(progressWriter, "'%s' is being interpreted as a hostname\n", identifier)
			return generic.NewHostnameSubject(identifier), nil
		}

		// 3. Try EC2 fuzzy matching.
		awsResourceProvider := resourceProviders[aws.ResourceDomainAWS].(aws.ResourceProvider)
		return aws.ResolveEC2InstanceSubject(identifier, awsResourceProvider)
	} else { // explicit resolution (subject type was specified)
		prefix := identifierSegments[0]
		qualifiedIdentifier := identifierSegments[1]

		switch prefix {
		case "ip":
			return generic.ResolveIPAddressSubject(qualifiedIdentifier)
		case "host":
			return generic.ResolveHostnameSubject(qualifiedIdentifier)
		case "ec2":
			awsResourceProvider := resourceProviders[aws.ResourceDomainAWS].(aws.ResourceProvider)
			return aws.ResolveEC2InstanceSubject(qualifiedIdentifier, awsResourceProvider)
		default:
			return nil, fmt.Errorf("unable to resolve subject with identifier '%s' because subject prefix '%s' is not recognized", qualifiedIdentifier, prefix)
		}
	}
}
