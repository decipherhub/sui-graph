package sui_graphd

import (
	"context"
	"fmt"
	"github.com/decipherhub/sui-graph/pkg/fetcher"
	"time"

	"github.com/decipherhub/sui-graph/internal/graph"
	"github.com/spf13/cobra"
)

var (
	checkpointNum int64
)

// indexCmd represents the 'index' subcommand
var indexCmd = &cobra.Command{
	Use:   "index",
	Short: "Build and index a transaction graph for a specific checkpoint",
	Long: `This command fetches all transactions in a given Sui checkpoint,
builds a graph of their execution dependencies, and optionally stores or prints the result.`,
	Run: func(cmd *cobra.Command, args []string) {
		log := cfg.Logger.WithField("checkpoint", checkpointNum)

		if checkpointNum < 0 {
			log.Fatal("Invalid checkpoint: must be >= 0")
		}

		log.Info("Starting checkpoint indexing")

		ctx := context.Background()
		start := time.Now()

		svc := graph.NewService(fetcher.New(cfg.SuiClient, nil))
		nodes, edges, ts, err := svc.BuildGraphForCheckpoint(ctx, checkpointNum)
		if err != nil {
			log.WithError(err).Fatal("Failed to build graph")
		}

		duration := time.Since(start)

		log.WithFields(map[string]interface{}{
			"timestamp":  ts.Format(time.RFC3339),
			"node_count": len(nodes),
			"edge_count": len(edges),
			"elapsed_ms": duration.Milliseconds(),
		}).Info("✅ Successfully indexed checkpoint")

		// Optional debug: print first few nodes
		for i, node := range nodes {
			if i >= 5 {
				break
			}
			fmt.Printf("→ [%s] %s (%s)\n", node.ID, node.Label, node.Sender)
		}
	},
}

func init() {
	rootCmd.AddCommand(indexCmd)

	indexCmd.Flags().Int64VarP(&checkpointNum, "checkpoint", "n", -1, "Checkpoint number to index (required)")
	indexCmd.MarkFlagRequired("checkpoint")
}
