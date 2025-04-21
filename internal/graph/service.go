package graph

import (
	"context"
	"github.com/decipherhub/sui-graph/types"
	"time"

	"github.com/decipherhub/sui-graph/pkg/fetcher"
)

// Service encapsulates graph construction logic with an injected fetcher.
type Service struct {
	DataFetcher GraphDataFetcher
}

// NewService creates a new Service instance with the given data fetcher.
func NewService(fetcher GraphDataFetcher) *Service {
	return &Service{
		DataFetcher: fetcher,
	}
}

// BuildGraphForCheckpoint constructs the graph using the injected data fetcher.
func (s *Service) BuildGraphForCheckpoint(ctx context.Context, checkpoint int64) ([]types.GraphNodeDTO, []types.GraphEdgeDTO, time.Time, error) {
	txs, err := s.DataFetcher.FetchTransactionsByCheckpoint(ctx, checkpoint)
	if err != nil {
		return nil, nil, time.Time{}, err
	}
	if len(txs) == 0 {
		return nil, nil, time.Time{}, nil
	}

	ts, err := s.DataFetcher.FetchCheckpointTimestamp(ctx, checkpoint)
	if err != nil {
		return nil, nil, time.Time{}, err
	}

	nodes := make([]types.GraphNodeDTO, 0, len(txs))
	for _, tx := range txs {
		nodes = append(nodes, types.GraphNodeDTO{
			ID:          tx.Digest,
			Label:       tx.Function,
			TimestampMs: ts.UnixMilli(),
			Sender:      tx.Sender,
			Objects:     tx.InvolvedObjects,
			Type:        tx.Type,
			Status:      tx.Status,
			GasUsed:     tx.GasUsed,
		})
	}

	edges := ComputeObjectDependencyEdges(txs)
	return nodes, edges, ts, nil
}

// ComputeObjectDependencyEdges analyzes a list of transactions and identifies object-level
// dependencies between them. If two or more transactions touch the same object (e.g. mutate),
// we assume there is a directional dependency.
//
// This is a simplified model that assumes time-based ordering in the transaction list.
// You may later enhance this with causal ordering (e.g., from events or explicit inputs).
func ComputeObjectDependencyEdges(txs []fetcher.TxSummary) []types.GraphEdgeDTO {
	objectToTxs := make(map[string][]string)

	// Map each object ID to the list of transactions that touched it
	for _, tx := range txs {
		for _, obj := range tx.InvolvedObjects {
			objectToTxs[obj] = append(objectToTxs[obj], tx.Digest)
		}
	}

	var edges []types.GraphEdgeDTO

	// For each object with multiple writers, create pairwise edges
	for objID, txList := range objectToTxs {
		if len(txList) < 2 {
			continue
		}
		// Naively create edges based on transaction order
		for i := 0; i < len(txList)-1; i++ {
			edges = append(edges, types.GraphEdgeDTO{
				Source: txList[i],
				Target: txList[i+1],
				Object: objID,
			})
		}
	}
	return edges
}
