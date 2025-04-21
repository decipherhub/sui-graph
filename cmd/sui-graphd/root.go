package sui_graphd

import (
	"fmt"
	"github.com/decipherhub/sui-graph/config"
	"github.com/decipherhub/sui-graph/version"
	"log"
)

var (
	configPath  string
	cfg         *config.ClientConfig
	page, limit int
)

var rootCmd = &cobra.Command{
	Use:     "account-cli",
	Short:   "Client CLI for interacting with account-harvest server",
	Long:    `This CLI lets you query and interact with the account-harvest API server.`,
	Version: fmt.Sprintf("Version: %s\nCommit: %s\nDate: %s", version.AppVersion, version.GitCommit, version.BuildDate),
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		var err error
		cfg, err = config.LoadClientConfig("")
		if err != nil {
			log.Fatal(err)
		}
	},
}

var addCmd = &cobra.Command{
	Use:     "add",
	Aliases: []string{"a"},
	Short:   "",
}

var updateCmd = &cobra.Command{
	Use:     "update",
	Aliases: []string{"u"},
	Short:   "",
}

var deleteCmd = &cobra.Command{
	Use:     "delete",
	Aliases: []string{"d"},
	Short:   "",
}

var queryCmd = &cobra.Command{
	Use:     "query",
	Aliases: []string{"q"},
	Short:   "",
}

var processCmd = &cobra.Command{
	Use:     "process",
	Aliases: []string{"p"},
	Short:   "",
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&configPath, "config", "c", "config.toml", "Path to the configuration file (default: config.toml).")

	rootCmd.AddCommand(addCmd, updateCmd, deleteCmd, queryCmd, processCmd)

	addCmd.AddCommand(addCoinCmd, addEndpointCmd, addNetworkCmd, addWalletCmd)

	updateCmd.AddCommand(updateCoinCmd, updateNetworkCmd, updateWalletCmd)

	deleteCmd.AddCommand(deleteCoinCmd, deleteEndpointCmd, deleteNetworkCmd, deleteWalletCmd)

	queryCmd.Flags().IntVar(&page, "page", 0, "Page number")
	queryCmd.Flags().IntVar(&limit, "limit", 50, "Limit number of results")

	queryCmd.AddCommand(queryCoinCmd, queryEndpointCmd, queryEndpointCmd, queryNetworkCmd, queryReportCmd, queryWalletCmd)

	processCmd.AddCommand(processReportCmd, processTxCmd)

}

func Execute() error {
	return rootCmd.Execute()
}

func main() {
	cfg, err := config.LoadConfig(".")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	fmt.Printf("Sui RPC URL: %s\n", cfg.Sui.RPCUrl)
	fmt.Printf("DB DSN: %s\n", cfg.Database.DSN)
}
