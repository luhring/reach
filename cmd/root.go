package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

var rootCmd = &cobra.Command{
	Use:   "reach",
	Short: "reach examines network reachability issues in AWS",
	Long: `reach examines network reachability issues in AWS
See https://github.com/luhring/reach for documentation.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Test!")
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

