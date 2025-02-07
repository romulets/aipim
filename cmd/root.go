package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "aipim",
	Short: "AIPIM is a tool for quickly developing and modifying Painless scripts in Ingest Pipelines",
	Run:   run,
}

func run(cmd *cobra.Command, args []string) {
	cmd.Help()
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
