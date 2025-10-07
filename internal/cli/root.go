package cli

import (
	"fmt"
	"os"

	"github.com/kaustuvprajapati/devopsctl/internal/config"
	"github.com/spf13/cobra"
)

var (
	cfgFile    string
	jsonOutput bool
	outputFile string

	// AppConfig holds the loaded configuration.
	AppConfig *config.Config
)

var rootCmd = &cobra.Command{
	Use:   "devopsctl",
	Short: "Infrastructure hygiene and DevOps validation toolkit",
	Long: `devopsctl is a lightweight CLI toolkit for infrastructure hygiene
and DevOps validation across AWS, Docker, Terraform, and Git.

Run checks against your infrastructure to identify security issues,
misconfigurations, and maintenance problems.`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		return initConfig()
	},
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

func initConfig() error {
	path := cfgFile
	if path == "" {
		path = config.FindConfigFile()
	}

	if path == "" {
		AppConfig = config.DefaultConfig()
		return nil
	}

	cfg, err := config.Load(path)
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	AppConfig = cfg
	return nil
}
