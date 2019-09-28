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
const portFlag = "port"
const portFlagShorthand = "p"
const assertReachableFlag = "assert-reachable"
const assertNotReachableFlag = "assert-not-reachable"

var shouldExplain bool
var port uint16
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

		// traffic, err := analysis.MergeVectorTraffic()
		// if err != nil {
		// 	log.Fatal(err)
		// }
		//
		// fmt.Print(traffic)
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		exitWithError(err)
	}
}

func init() {
	rootCmd.Flags().BoolVar(&shouldExplain, explainFlag, false, "explain how the configuration was analyzed")
	rootCmd.Flags().Uint16VarP(&port, portFlag, portFlagShorthand, 0, "restrict analysis to a specified TCP port")
	rootCmd.Flags().BoolVar(&assertReachable, assertReachableFlag, false, "exit non-zero if no traffic is allowed from source to destination (within analysis scope, if specified)")
	rootCmd.Flags().BoolVar(&assertNotReachable, assertNotReachableFlag, false, "exit non-zero if any traffic can reach destination from source (within analysis scope, if specified)")
}
