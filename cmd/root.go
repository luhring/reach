package cmd

import (
	"fmt"
	"github.com/luhring/reach/reach"
	"github.com/spf13/cobra"
	"os"
)

var rootCmd = &cobra.Command{
	Use:   "reach",
	Short: "reach examines network reachability issues in AWS",
	Long: `reach examines network reachability issues in AWS
See https://github.com/luhring/reach for documentation.`,
	Args: cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		analyzer := reach.NewAnalyzer()

		vector, err := analyzer.CreateInstanceVector(args[0], args[1])
		if err != nil {
			exitWithError(err)
		}

		analyzer.Analyze(vector)
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		exitWithError(err)
	}
}

func exitWithError(err error) {
	fmt.Println(err)
	os.Exit(1)
}