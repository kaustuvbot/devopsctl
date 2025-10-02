package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	cfgFile    string
	jsonOutput bool
	outputFile string
)

var rootCmd = &cobra.Command{
	Use:   "devopsctl",
	Short: "Infrastructure hygiene and DevOps validation toolkit",
	Long: `devopsctl is a lightweight CLI toolkit for infrastructure hygiene
and DevOps validation across AWS, Docker, Terraform, and Git.`,
}

// Execute runs the root command.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is .devopsctl.yaml)")
	rootCmd.PersistentFlags().BoolVar(&jsonOutput, "json", false, "output in JSON format")
	rootCmd.PersistentFlags().StringVar(&outputFile, "output", "", "write report to file")
}
