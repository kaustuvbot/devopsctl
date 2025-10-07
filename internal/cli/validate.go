package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

var validateCmd = &cobra.Command{
	Use:   "validate",
	Short: "Run validation checks",
	Long:  `Run validation checks against infrastructure code.`,
}

var validateTerraformCmd = &cobra.Command{
	Use:   "terraform",
	Short: "Validate Terraform configuration",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("Terraform validation not yet implemented")
		return nil
	},
}

func init() {
	validateCmd.AddCommand(validateTerraformCmd)
	rootCmd.AddCommand(validateCmd)
}
