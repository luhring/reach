package cmd

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/luhring/reach/reach/analyzer"
	"github.com/luhring/reach/reach/aws"
	"github.com/luhring/reach/reach/aws/api"
	"github.com/luhring/reach/reach/explainer"
)

const explainFlag = "explain"
const vectorsFlag = "vectors"
const assertReachableFlag = "assert-reachable"
const assertNotReachableFlag = "assert-not-reachable"

var explain bool
var showVectors bool
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
			log.Fatal(err)
		}
		source.SetRoleToSource()

		destination, err := aws.NewSubject(destinationIdentifier, provider)
		if err != nil {
			log.Fatal(err)
		}
		destination.SetRoleToDestination()

		a := analyzer.New()
		analysis, err := a.Analyze(source, destination)
		if err != nil {
			log.Fatal(err)
		}

		mergedTraffic, err := analysis.MergedTraffic()
		if err != nil {
			log.Fatal(err)
		}

		if explain {
			ex := explainer.New(*analysis)
			fmt.Print(ex.Explain())
		} else if showVectors {
			var vectorOutputs []string

			for _, v := range analysis.NetworkVectors {
				output := ""
				output += v.String()

				vectorOutputs = append(vectorOutputs, output)
			}

			fmt.Print(strings.Join(vectorOutputs, "\n"))
		} else {
			fmt.Print("network traffic allowed from source to destination:" + "\n")
			fmt.Print(mergedTraffic.ColorStringWithSymbols())

			if len(analysis.NetworkVectors) > 1 {
				printMergedResultsWarning()
			}
		}

		const canReach = "source is able to reach destination"
		const cannotReach = "source is unable to reach destination"

		if assertReachable {
			if mergedTraffic.None() {
				exitWithFailedAssertion(cannotReach)
			} else {
				exitWithSuccessfulAssertion(canReach)
			}
		}

		if assertNotReachable {
			if mergedTraffic.None() {
				exitWithSuccessfulAssertion(cannotReach)
			} else {
				exitWithFailedAssertion(canReach)
			}
		}
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		exitWithError(err)
	}
}

func init() {
	rootCmd.Flags().BoolVar(&explain, explainFlag, false, "explain how the configuration was analyzed")
	rootCmd.Flags().BoolVar(&showVectors, vectorsFlag, false, "show allowed traffic in terms of network vectors")
	rootCmd.Flags().BoolVar(&assertReachable, assertReachableFlag, false, "exit non-zero if no traffic is allowed from source to destination")
	rootCmd.Flags().BoolVar(&assertNotReachable, assertNotReachableFlag, false, "exit non-zero if any traffic can reach destination from source")
}

func printMergedResultsWarning() {
	const mergedResultsWarning = "IMPORTANT: Reach detected more than one network path between the source and destination. Reach calls these paths \"network vectors\". The analysis result shown above is the merging of all network vectors' analysis results. The impact that infrastructure configuration has on actual network reachability might vary based on the way hosts are configured to use their network interfaces, and Reach is unable to access any configuration internal to a host. To see the network reachability across individual network vectors, run the command again with '--vectors'.\n\n"
	_, _ = fmt.Fprint(os.Stderr, "\n"+mergedResultsWarning)
}
