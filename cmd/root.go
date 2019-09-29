package cmd

import (
	"errors"
	"fmt"
	"log"

	"github.com/spf13/cobra"

	"github.com/luhring/reach/reach/analyzer"
	"github.com/luhring/reach/reach/aws"
	"github.com/luhring/reach/reach/aws/api"
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

		fmt.Println(analysis.ToJSON())

		if showVectors {
			for _, v := range analysis.NetworkVectors {
				fmt.Print(v)
				fmt.Println(v.Traffic)
			}
		} else {
			mergedTraffic, err := analysis.MergedTraffic()
			if err != nil {
				log.Fatal(err)
			}

			fmt.Print(mergedTraffic)
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
