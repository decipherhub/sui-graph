package sui_graphd

import (
	"fmt"
	"log"

	"github.com/decipherhub/sui-graph/config"
	"github.com/decipherhub/sui-graph/version"
	"github.com/spf13/cobra"
)

var (
	configPath string
	cfg        *config.Config
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "sui-graphd",
	Short: "A tool for building and serving Sui transaction dependency graphs",
	Long: `sui-graphd is a backend service for indexing, analyzing,
and exposing Sui blockchain transaction dependency graphs via a Graph API.`,
	Version: fmt.Sprintf("Version: %s\nCommit: %s\nDate: %s", version.AppVersion, version.GitCommit, version.BuildDate),
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		var err error
		cfg, err = config.LoadConfig(configPath)
		if err != nil {
			log.Fatalf("Failed to load config from %s: %v", configPath, err)
		}
	},
}

// init sets up persistent flags for the CLI
func init() {
	rootCmd.PersistentFlags().StringVarP(&configPath, "config", "c", "config.yaml", "Path to the configuration file (default: config.yaml)")
}

// Execute runs the CLI root command
func Execute() error {
	return rootCmd.Execute()
}
