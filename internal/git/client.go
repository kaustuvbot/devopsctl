package git

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"time"
)

// Client wraps git command execution.
type Client struct {
	repoPath string
	timeout  time.Duration
}

// NewClient creates a new git client for the given repository path.
func NewClient(repoPath string) *Client {
	return &Client{
		repoPath: repoPath,
		timeout:  30 * time.Second,
	}
}

// Run executes a git command and returns the output.
// ctx can be used to cancel the command.
func (c *Client) Run(ctx context.Context, args ...string) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()

	cmd := exec.CommandContext(ctx, "git", args...)
	cmd.Dir = c.repoPath

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("git %s: %w", args[0], err)
	}

	return stdout.String(), nil
}

// RunSimple is a convenience method for running git commands without context.
func (c *Client) RunSimple(args ...string) (string, error) {
	return c.Run(context.Background(), args...)
}

// IsRepo checks if the given path is a git repository.
func (c *Client) IsRepo() bool {
	_, err := c.RunSimple("rev-parse", "--is-inside-work-tree")
	return err == nil
}

// CurrentBranch returns the current branch name.
func (c *Client) CurrentBranch() (string, error) {
	return c.RunSimple("rev-parse", "--abbrev-ref", "HEAD")
}

// BranchList returns all local and remote branches.
func (c *Client) BranchList() (string, error) {
	return c.RunSimple("branch", "-a")
}

// LastCommitDate returns the date of the last commit on a branch.
func (c *Client) LastCommitDate(branch string) (time.Time, error) {
	output, err := c.RunSimple("log", "-1", "--format=%ci", branch)
	if err != nil {
		return time.Time{}, err
	}

	// Parse the date (format: "2025-12-16 10:00:00 +0000")
	layout := "2006-01-02 15:04:05 -0700"
	t, err := time.Parse(layout, output[:len(output)-1]) // trim newline
	if err != nil {
		return time.Time{}, err
	}

	return t, nil
}

// CountObjects returns the count-objects -vH output for repo size.
func (c *Client) CountObjects() (string, error) {
	return c.RunSimple("count-objects", "-vH")
}

// ListFiles lists all tracked files with their sizes.
func (c *Client) ListFiles() (string, error) {
	return c.RunSimple("ls-files", "-s")
}
