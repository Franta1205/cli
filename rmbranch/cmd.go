package rmbranch

import (
	"errors"
	"fmt"
	"os/exec"
	"strings"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

type cmdPair struct {
	B  []string
	Cb string
}

func newCmdPair(b []string, cb string) *cmdPair {
	return &cmdPair{b, cb}
}

func CreateRmBranchCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "rm_branch",
		Short:   "Remove git branches",
		Long:    "A tool to safely remove Git branches with additional sfety checks",
		Example: "rm_branch all",
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

	cb, err := getCurrentBranch()
	if err != nil {
		color.Red("error: Failed to get current branch: %v", err)
		return fmt.Errorf("error: %v", err)
	}
	c := newCmdPair(args, cb)
	if err := c.deleteBranches(); err != nil {
		color.Red("error: %v", err)
		return fmt.Errorf("error: %v", err)
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
	cb := strings.TrimSpace(string(output))
	color.Yellow("current branch: %v\n", cb)
	return cb, nil
}

func (c *cmdPair) deleteBranches() error {
	if c.B[0] == "all" {
		color.Yellow("Deleting all branches from current repository")
		branches, _ := getAllBranches()
		for _, b := range branches {
			if b == "main" || b == "master" || b == c.Cb {
				color.Yellow("skipping %v", b)
				continue
			}

			color.Yellow("Deleting: %v", b)
			cmd := exec.Command("git", "branch", "-D", b)
			cmd.Run()
			color.Green("%v\t branch deleted", b)
		}
	} else {
		return errors.New("branch deletion is not implemented for specific branches")
	}
	return nil
}

func getAllBranches() ([]string, error) {
	cmd := exec.Command("git", "branch", "--list")
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	lines := strings.Split(strings.TrimSpace(string(output)), "\n")
	branches := make([]string, 0, len(lines))

	for _, line := range lines {
		branch := strings.TrimSpace(line)
		if strings.HasPrefix(branch, "* ") {
			branch = strings.TrimPrefix(branch, "* ")
		}
		branches = append(branches, branch)
	}

	return branches, nil
}
