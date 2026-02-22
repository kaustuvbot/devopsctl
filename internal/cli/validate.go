package cli

import (
	"fmt"
	"os"

	"github.com/kaustuvbot/devopsctl/internal/reporter"
	terraformpkg "github.com/kaustuvbot/devopsctl/internal/terraform"
	"github.com/spf13/cobra"
)

var validateCmd = &cobra.Command{
	Use:   "validate",
	Short: "Run validation checks",
	Long:  `Run validation checks against infrastructure code.`,
}

var terraformDir string

var validateTerraformCmd = &cobra.Command{
	Use:   "terraform",
	Short: "Validate Terraform configuration",
	RunE: func(cmd *cobra.Command, args []string) error {
		workingDir := terraformDir
		if workingDir == "" {
			workingDir = "."
		}

		runner := terraformpkg.NewRunner(workingDir)
		results, err := runner.RunAllChecks()
		if err != nil {
			fmt.Fprintf(os.Stderr, "warning: some checks encountered errors: %v\n", err)
		}

		// Convert terraform.CheckResult to reporter.CheckResult
		reportResults := make([]reporter.CheckResult, len(results))
		for i, r := range results {
			reportResults[i] = reporter.CheckResult{
				CheckName:      r.CheckName,
				Severity:      string(r.Severity),
				ResourceID:     r.ResourceID,
				Message:        r.Message,
				Recommendation: r.Recommendation,
			}
		}

		report := &reporter.Report{Module: "terraform", Results: reportResults}

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

		if code := exitCodeForResults(reportResults); code > 0 {
			os.Exit(code)
		}
		return nil
	},
}

func init() {
	validateTerraformCmd.Flags().StringVar(&terraformDir, "dir", "", "path to Terraform directory (default: current directory)")
	validateCmd.AddCommand(validateTerraformCmd)
	rootCmd.AddCommand(validateCmd)
}
