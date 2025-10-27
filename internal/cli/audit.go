package cli

import (
	"context"
	"fmt"
	"os"

	awspkg "github.com/kaustuvprajapati/devopsctl/internal/aws"
	"github.com/kaustuvprajapati/devopsctl/internal/reporter"
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

var auditDockerCmd = &cobra.Command{
	Use:   "docker",
	Short: "Audit Docker configuration",
	Long:  `Run static checks against a Dockerfile.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("Docker audit not yet implemented")
		return nil
	},
}

var auditGitCmd = &cobra.Command{
	Use:   "git",
	Short: "Audit Git repository",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("Git audit not yet implemented")
		return nil
	},
}

func init() {
	auditDockerCmd.Flags().StringVar(&dockerfilePath, "file", "", "path to Dockerfile (overrides config)")
	auditCmd.AddCommand(auditAWSCmd)
	auditCmd.AddCommand(auditDockerCmd)
	auditCmd.AddCommand(auditGitCmd)
	rootCmd.AddCommand(auditCmd)
}
