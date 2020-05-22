package cmd

import (
	"errors"
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/luhring/reach/reach"
	"github.com/luhring/reach/reach/analyzer"
	"github.com/luhring/reach/reach/aws"
	"github.com/luhring/reach/reach/aws/apiclient"
	"github.com/luhring/reach/reach/cache"
	"github.com/luhring/reach/reach/generic"
	"github.com/luhring/reach/reach/generic/standard"
	"github.com/luhring/reach/reach/reachlog"
)

const githubURL = "https://github.com/luhring/reach"
const explainFlag = "explain"
const pathsFlag = "paths"
const jsonFlag = "json"
const assertReachableFlag = "assert-reachable"
const assertNotReachableFlag = "assert-not-reachable"

var explain bool
var showPaths bool
var outputJSON bool
var assertReachable bool
var assertNotReachable bool
var verbose bool

var logger = reachlog.New(reachlog.LevelNone)

var rootCmd = &cobra.Command{
	Use:   "reach",
	Short: "reach examines network reachability issues in AWS",
	Long: `reach examines network reachability issues in AWS
See https://github.com/luhring/reach for documentation.`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) != 2 {
			return errors.New("requires exactly two arguments")
		}

		if assertReachable && assertNotReachable {
			return errors.New("cannot assert both reachable and not reachable at the same time")
		}

		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		if verbose {
			logger = reachlog.New(reachlog.LevelDebug)
		}

		sourceInput := args[0]
		destinationInput := args[1]

		catalog := reach.NewDomainClientCatalog()
		c := cache.New()
		awsClient, err := apiclient.NewDomainClient(&c)
		if err != nil {
			handleError(err)
		}
		catalog.Store(aws.ResourceDomainAWS, awsClient)
		catalog.Store(generic.ResourceDomainGeneric, standard.NewDomainClient())

		source, err := resolveSubject(sourceInput, catalog)
		if err != nil {
			handleError(err)
		}
		source.SetRoleToSource()

		destination, err := resolveSubject(destinationInput, catalog)
		if err != nil {
			handleError(err)
		}
		destination.SetRoleToDestination()

		if !(outputJSON || explain || showPaths) {
			fmt.Printf("source: %s\ndestination: %s\n\n", source.ID, destination.ID)
		}

		a := analyzer.New(catalog, logger)
		analysis, err := a.Analyze(*source, *destination)
		if err != nil {
			handleError(err)
		}

		if outputJSON {
			fmt.Println(analysis.ToJSON())
		} else if showPaths {
			var pathDescriptions []string

			for _, p := range analysis.Paths {
				pathDescriptions = append(pathDescriptions, fmt.Sprint(p))
			}

			fmt.Print(strings.Join(pathDescriptions, "\n"))
		} else {
			paths := analysis.Paths
			tcs, err := reach.TrafficContentsFromPaths(paths)
			if err != nil {
				handleError(err)
			}
			mergedTraffic, err := reach.MergeTraffic(tcs...)
			if err != nil {
				handleError(err)
			}

			fmt.Print("network traffic allowed from source to destination:" + "\n")
			fmt.Print(mergedTraffic.ColorStringWithSymbols())

			if len(paths) > 1 { // handling this case with care; this view isn't optimized for multi-vector output!
				printMergedResultsWarning()
				// warnIfAnyVectorHasRestrictedReturnTraffic(paths)
			} else {
				// warnings about return traffic?
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
	rootCmd.Flags().BoolVar(&showPaths, pathsFlag, false, "show allowed traffic in terms of network paths")
	rootCmd.Flags().BoolVar(&outputJSON, jsonFlag, false, "output full analysis as JSON (overrides other display flags)")
	rootCmd.Flags().BoolVar(&assertReachable, assertReachableFlag, false, "exit non-zero if no traffic is allowed from source to destination")
	rootCmd.Flags().BoolVar(&assertNotReachable, assertNotReachableFlag, false, "exit non-zero if any traffic can reach destination from source")
	rootCmd.Flags().BoolVarP(&verbose, "", "v", false, "show verbose output (displays full log output)")
}
