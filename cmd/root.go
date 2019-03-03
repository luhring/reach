package cmd

import (
	"errors"
	"fmt"
	"github.com/luhring/reach/reach"
	"github.com/spf13/cobra"
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
		awsManager := reach.NewAWSManager()

		instanceVector, err := awsManager.CreateInstanceVector(args[0], args[1])
		if err != nil {
			exitWithError(err)
		}

		var filter *reach.TrafficAllowance
		if port == 0 {
			filter = nil
		} else {
			filter = reach.NewTrafficAllowanceForTCPPort(port)
			fmt.Printf("analysis scope: TCP %v\n", port)
		}

		analysis := instanceVector.Analyze(filter)
		fmt.Print(analysis.Results())

		if shouldExplain {
			fmt.Println("")
			fmt.Print(analysis.Explanation())
		}

		if assertReachable {
			if analysis.PassesAssertReachable() {
				exitWithSuccessfulAssertion("specified traffic flow is allowed")
			} else {
				exitWithFailedAssertion("specified traffic flow is not allowed")
			}
		}

		if assertNotReachable {
			if analysis.PassesAssertNotReachable() {
				exitWithSuccessfulAssertion("none of specified traffic flow is allowed")
			} else {
				exitWithFailedAssertion("some or all of specified traffic flow is allowed")
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
	rootCmd.Flags().BoolVar(&shouldExplain, explainFlag, false, "explain how the configuration was analyzed")
	rootCmd.Flags().Uint16VarP(&port, portFlag, portFlagShorthand, 0, "restrict analysis to a specified TCP port")
	rootCmd.Flags().BoolVar(&assertReachable, assertReachableFlag, false, "exit non-zero if no traffic is allowed from source to destination (within analysis scope, if specified)")
	rootCmd.Flags().BoolVar(&assertNotReachable, assertNotReachableFlag, false, "exit non-zero if any traffic can reach destination from source (within analysis scope, if specified)")
}
