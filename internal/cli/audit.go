package cli

import (
	"context"
	"fmt"
	"os"

	awspkg "github.com/kaustuvbot/devopsctl/internal/aws"
	dockerpkg "github.com/kaustuvbot/devopsctl/internal/docker"
	gitpkg "github.com/kaustuvbot/devopsctl/internal/git"
	"github.com/kaustuvbot/devopsctl/internal/reporter"
	"github.com/spf13/cobra"
)

var auditCmd = &cobra.Command{
	Use:   "audit",
	Short: "Run audit checks",
	Long:  `Run infrastructure audit checks across different platforms.`,
}

var auditAWSCmd = &cobra.Command{
	Use:   "aws",
	Short: "Audit AWS infrastructure",
	Long:  `Audit AWS IAM, S3, EC2 security groups, and EBS volumes.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		clients, err := awspkg.NewAWSClients(AppConfig.AWS)
		if err != nil {
			return fmt.Errorf("failed to initialize AWS clients: %w", err)
		}

		results, err := awspkg.RunAll(context.Background(), clients, AppConfig.AWS)
		if err != nil {
			fmt.Fprintf(os.Stderr, "warning: some checks encountered errors: %v\n", err)
		}

		// Apply filters: ignore patterns and quiet mode
		results = filterByIgnore(results, AppConfig.Ignore.Checks)
		results = filterBySeverity(results, quiet)

		report := &reporter.Report{Module: "aws", Results: results}

		w, err := resolveWriter(cmd)
		if err != nil {
			return err
		}
		if w != os.Stdout {
			defer w.Close()
		}

		rep := resolveReporter()
		if err := rep.Render(w, report); err != nil {
			return err
		}

		if code := exitCodeForResults(results); code > 0 {
			os.Exit(code)
		}
		return nil
	},
}

var dockerfilePath string
var dockerImage string

var auditDockerCmd = &cobra.Command{
	Use:   "docker",
	Short: "Audit Docker configuration",
	Long:  `Run static checks against a Dockerfile and optionally scan an image with Trivy.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		dockerCfg := AppConfig.Docker
		if dockerfilePath != "" {
			dockerCfg.DockerfilePath = dockerfilePath
		}

		opts := dockerpkg.RunOptions{ImageName: dockerImage}
		results, err := dockerpkg.RunAll(dockerCfg, opts)
		if err != nil {
			return err
		}

		// Apply filters: ignore patterns and quiet mode
		results = filterByIgnore(results, AppConfig.Ignore.Checks)
		results = filterBySeverity(results, quiet)

		report := &reporter.Report{Module: "docker", Results: results}

		w, err := resolveWriter(cmd)
		if err != nil {
			return err
		}
		if w != os.Stdout {
			defer w.Close()
		}

		rep := resolveReporter()
		if err := rep.Render(w, report); err != nil {
			return err
		}

		if code := exitCodeForResults(results); code > 0 {
			os.Exit(code)
		}
		return nil
	},
}

var gitRepoPath string

var auditGitCmd = &cobra.Command{
	Use:   "git",
	Short: "Audit Git repository",
	Long:  `Audit Git repository for hygiene issues: size, stale branches, large files.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		repoPath := gitRepoPath
		if repoPath == "" {
			cwd, err := os.Getwd()
			if err != nil {
				return fmt.Errorf("failed to get current directory: %w", err)
			}
			repoPath = cwd
		}

		runner := gitpkg.NewRunner(repoPath, AppConfig.Git)
		results, err := runner.RunAll(context.Background())
		if err != nil {
			fmt.Fprintf(os.Stderr, "warning: some checks encountered errors: %v\n", err)
		}

		// Apply filters: ignore patterns and quiet mode
		results = filterByIgnore(results, AppConfig.Ignore.Checks)
		results = filterBySeverity(results, quiet)

		report := &reporter.Report{Module: "git", Results: results}

		w, err := resolveWriter(cmd)
		if err != nil {
			return err
		}
		if w != os.Stdout {
			defer w.Close()
		}

		rep := resolveReporter()
		if err := rep.Render(w, report); err != nil {
			return err
		}

		if code := exitCodeForResults(results); code > 0 {
			os.Exit(code)
		}
		return nil
	},
}

func init() {
	auditDockerCmd.Flags().StringVar(&dockerfilePath, "file", "", "path to Dockerfile (overrides config)")
	auditDockerCmd.Flags().StringVar(&dockerImage, "image", "", "container image to scan with Trivy")
	auditGitCmd.Flags().StringVar(&gitRepoPath, "repo", "", "path to Git repository (defaults to current directory)")
	auditCmd.AddCommand(auditAWSCmd)
	auditCmd.AddCommand(auditDockerCmd)
	auditCmd.AddCommand(auditGitCmd)
	rootCmd.AddCommand(auditCmd)
}
