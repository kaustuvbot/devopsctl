package cli

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/kaustuvprajapati/devopsctl/internal/doctor"
	"github.com/kaustuvprajapati/devopsctl/internal/reporter"
	"github.com/spf13/cobra"
)

// awsModule wraps AWS checks as a doctor.Module
type awsModule struct{}

func (m *awsModule) Name() string { return "aws" }

func (m *awsModule) Run(ctx context.Context) ([]reporter.CheckResult, error) {
	// Will be implemented in CLI after AWS clients are available
	return []reporter.CheckResult{}, nil
}

// dockerModule wraps Docker checks as a doctor.Module
type dockerModule struct{}

func (m *dockerModule) Name() string { return "docker" }

func (m *dockerModule) Run(ctx context.Context) ([]reporter.CheckResult, error) {
	// Will be implemented in CLI after Docker runner is available
	return []reporter.CheckResult{}, nil
}

// terraformModule wraps Terraform checks as a doctor.Module
type terraformModule struct{}

func (m *terraformModule) Name() string { return "terraform" }

func (m *terraformModule) Run(ctx context.Context) ([]reporter.CheckResult, error) {
	// Will be implemented in CLI after Terraform runner is available
	return []reporter.CheckResult{}, nil
}

// gitModule wraps Git checks as a doctor.Module
type gitModule struct{}

func (m *gitModule) Name() string { return "git" }

func (m *gitModule) Run(ctx context.Context) ([]reporter.CheckResult, error) {
	// Will be implemented in CLI after Git runner is available
	return []reporter.CheckResult{}, nil
}

var doctorCmd = &cobra.Command{
	Use:   "doctor",
	Short: "Run all checks and generate health report",
	Long: `Run all available audit and validation checks, aggregate results,
and generate a comprehensive health report.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		engine := doctor.NewEngine()

		// Register all modules
		if err := engine.Register(&awsModule{}); err != nil {
			return fmt.Errorf("failed to register aws module: %w", err)
		}
		if err := engine.Register(&dockerModule{}); err != nil {
			return fmt.Errorf("failed to register docker module: %w", err)
		}
		if err := engine.Register(&terraformModule{}); err != nil {
			return fmt.Errorf("failed to register terraform module: %w", err)
		}
		if err := engine.Register(&gitModule{}); err != nil {
			return fmt.Errorf("failed to register git module: %w", err)
		}

		// Run all modules
		reports, err := engine.RunAll(context.Background())
		if err != nil {
			fmt.Fprintf(os.Stderr, "warning: some modules encountered errors: %v\n", err)
		}

		// Compute summary
		summary := doctor.ComputeSummary(reports)

		// Output results
		w, err := resolveWriter(cmd)
		if err != nil {
			return err
		}
		if w != os.Stdout {
			defer w.Close()
		}

		rep := resolveReporter()

		// Render each module's results
		for _, r := range reports {
			report := &reporter.Report{Module: r.Module, Results: r.Results}
			if err := rep.Render(w, report); err != nil {
				return err
			}
			if r.Error != "" {
				fmt.Fprintf(w, "  [ERROR] %s\n", r.Error)
			}
		}

		// Print summary for JSON output
		if jsonOutput {
			summaryJSON, _ := json.MarshalIndent(summary, "", "  ")
			fmt.Fprintf(w, "\n%s\n", summaryJSON)
		}

		// Exit with appropriate code
		if code := doctor.ExitCode(reports); code > 0 {
			os.Exit(code)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(doctorCmd)
}
