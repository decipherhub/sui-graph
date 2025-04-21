package graph

import (
	"context"
	"time"

	"github.com/decipherhub/sui-graph/pkg/fetcher"
)

// GraphDataFetcher defines an interface for retrieving graph data from the chain or DB.
type GraphDataFetcher interface {
	FetchTransactionsByCheckpoint(ctx context.Context, checkpoint int64) ([]fetcher.TxSummary, error)
	FetchCheckpointTimestamp(ctx context.Context, checkpoint int64) (time.Time, error)
}
