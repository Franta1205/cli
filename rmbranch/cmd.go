package rmbranch

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

func CreateRmBranchCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "rm_branch",
		Short:   "Remove git branches",
		Long:    "A tool to safely remove Git branches with additional sfety checks",
		Example: "rm_branch *, rm_branch feat/new-feature",
		RunE:    rmBranch,
	}
	return cmd
}

func rmBranch(cmd *cobra.Command, args []string) error {
	color.Yellow("Starting branch removal")

	if !isGitRepository() {
		color.Red("error: not a git repository")
		return fmt.Errorf("current directory is not a git repository")
	}

	if len(args) < 1 {
		color.Red("error: no branch specified")
		if err := cmd.Help(); err != nil {
			color.Red("error: %v", err)
		}
		return fmt.Errorf("no branch specified")
	}

	_, err := getCurrentBranch()

	if err != nil {
		color.Red("error: Failed to get current branch: %v", err)
	}

	if args[0] == "all" {
		color.Yellow("Deleting all branches")
	}
	return nil
}

func isGitRepository() bool {
	cmd := exec.Command("git", "rev-parse", "--is-inside-work-tree")
	err := cmd.Run()
	return err == nil
}

func getCurrentBranch() (string, error) {
	cmd := exec.Command("git", "symbolic-ref", "--short", "HEAD")
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(output)), nil
}
