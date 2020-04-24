package cmd

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/luhring/reach/reach"
	"github.com/luhring/reach/reach/analyzer"
	"github.com/luhring/reach/reach/aws"
	"github.com/luhring/reach/reach/aws/api"
	"github.com/luhring/reach/reach/explainer"
	"github.com/luhring/reach/reach/generic"
	"github.com/luhring/reach/reach/generic/standard"
)

const explainFlag = "explain"
const vectorsFlag = "vectors"
const jsonFlag = "json"
const assertReachableFlag = "assert-reachable"
const assertNotReachableFlag = "assert-not-reachable"

var explain bool
var showPaths bool
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
		sourceInput := args[0]
		destinationInput := args[1]

		domains := &ConfigurationDomainProvider{}
		domains.Load(aws.ResourceDomainAWS, api.NewResourceProvider())
		domains.Load(generic.ResourceDomainGeneric, standard.NewResourceProvider())

		source, err := resolveSubject(sourceInput, os.Stderr, domains)
		if err != nil {
			exitWithError(err)
		}
		source.SetRoleToSource()

		destination, err := resolveSubject(destinationInput, os.Stderr, domains)
		if err != nil {
			exitWithError(err)
		}
		destination.SetRoleToDestination()

		if !(outputJSON || explain || showPaths) {
			fmt.Printf("source: %s\ndestination: %s\n\n", source.ID, destination.ID)
		}

		a := analyzer.New(nil, domains) // TODO: Pass in InfrastructureGetter
		analysis, err := a.Analyze(*source, *destination)
		if err != nil {
			exitWithError(err)
		}

		if outputJSON {
			fmt.Println(analysis.ToJSON())
		} else if explain {
			ex := explainer.New(*analysis)
			fmt.Print(ex.Explain())
		} else if showPaths {
			var pathDescriptions []string

			for _, p := range analysis.Paths {
				pathDescriptions = append(pathDescriptions, fmt.Sprint(p))
			}

			fmt.Print(strings.Join(pathDescriptions, "\n"))
		} else {
			paths := analysis.Paths
			tcs := reach.TrafficContentsFromPaths(paths)
			mergedTraffic, err := reach.MergeTraffic(tcs...)
			if err != nil {
				exitWithError(err)
			}

			fmt.Print("network traffic allowed from source to destination:" + "\n")
			fmt.Print(mergedTraffic.ColorStringWithSymbols())

			if len(paths) > 1 { // handling this case with care; this view isn't optimized for multi-vector output!
				printMergedResultsWarning()
				// warnIfAnyVectorHasRestrictedReturnTraffic(paths)
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
	rootCmd.Flags().BoolVar(&showPaths, vectorsFlag, false, "show allowed traffic in terms of network vectors")
	rootCmd.Flags().BoolVar(&outputJSON, jsonFlag, false, "output full analysis as JSON (overrides other display flags)")
	rootCmd.Flags().BoolVar(&assertReachable, assertReachableFlag, false, "exit non-zero if no traffic is allowed from source to destination")
	rootCmd.Flags().BoolVar(&assertNotReachable, assertNotReachableFlag, false, "exit non-zero if any traffic can reach destination from source")
}
