package cmd

import (
	"errors"
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/luhring/reach/reach/analyzer"
	"github.com/luhring/reach/reach/aws"
	"github.com/luhring/reach/reach/aws/api"
	"github.com/luhring/reach/reach/explainer"
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

		var provider aws.ResourceProvider = api.NewResourceProvider()

		source, err := aws.NewSubject(sourceIdentifier, provider)
		if err != nil {
			exitWithError(err)
		}
		source.SetRoleToSource()

		destination, err := aws.NewSubject(destinationIdentifier, provider)
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
