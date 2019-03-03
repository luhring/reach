package cmd

import (
	"fmt"
	"github.com/luhring/reach/reach"
	"github.com/spf13/cobra"
	"os"
)

const explainFlag = "explain"
const portFlag = "port"

var shouldExplain bool
var port uint16

var rootCmd = &cobra.Command{
	Use:   "reach",
	Short: "reach examines network reachability issues in AWS",
	Long: `reach examines network reachability issues in AWS
See https://github.com/luhring/reach for documentation.`,
	Args: cobra.MinimumNArgs(2),
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
		}

		analysis := instanceVector.Analyze(filter)
		fmt.Print(analysis.Results())

		if shouldExplain {
			fmt.Println("")
			fmt.Print(analysis.Explanation())
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
	rootCmd.Flags().Uint16Var(&port, portFlag, 0, "restrict analysis to a specified TCP port")
}

func exitWithError(err error) {
	fmt.Println(err)
	os.Exit(1)
}
