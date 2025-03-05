package main

import (
	"cli/rmbranch"
	"fmt"

	"github.com/spf13/cobra"
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "cli",
		Short: "CLI tools for development",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Use a subcommand like 'rm_branch'")
			cmd.Help()
		},
	}

	rootCmd.AddCommand(rmbranch.CreateRmBranchCommand())

	rootCmd.Execute()
}
