package graph_test

import (
	"context"
	"testing"
	"time"

	"github.com/decipherhub/sui-graph/internal/graph"
	"github.com/decipherhub/sui-graph/pkg/fetcher"
	"github.com/stretchr/testify/assert"
)

// mockFetcher implements the GraphDataFetcher interface for testing.
type mockFetcher struct{}

func (m *mockFetcher) FetchTransactionsByCheckpoint(ctx context.Context, checkpoint int64) ([]fetcher.TxSummary, error) {
	return []fetcher.TxSummary{
		{
			Digest:          "0xtx1",
			Function:        "0x2::coin::transfer",
			Sender:          "0xalice",
			InvolvedObjects: []string{"0xobj1", "0xobj2"},
			Type:            "moveCall",
			Status:          "success",
			GasUsed:         1000,
		},
		{
			Digest:          "0xtx2",
			Function:        "0x2::coin::transfer",
			Sender:          "0xbob",
			InvolvedObjects: []string{"0xobj2", "0xobj3"},
			Type:            "moveCall",
			Status:          "success",
			GasUsed:         900,
		},
	}, nil
}

func (m *mockFetcher) FetchCheckpointTimestamp(ctx context.Context, checkpoint int64) (time.Time, error) {
	return time.UnixMilli(1713723000000), nil
}

func TestBuildGraphForCheckpoint(t *testing.T) {
	ctx := context.Background()
	service := graph.NewService(&mockFetcher{})

	nodes, edges, ts, err := service.BuildGraphForCheckpoint(ctx, 1234)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(nodes), "should return 2 graph nodes")
	assert.Equal(t, 1, len(edges), "should return 1 dependency edge")
	assert.Equal(t, int64(1713723000000), ts.UnixMilli())

	// Check first node fields
	node := nodes[0]
	assert.Equal(t, "0xtx1", node.ID)
	assert.Equal(t, "0x2::coin::transfer", node.Label)
	assert.Equal(t, "0xalice", node.Sender)
	assert.Equal(t, "success", node.Status)

	// Check edge details
	edge := edges[0]
	assert.Equal(t, "0xtx1", edge.Source)
	assert.Equal(t, "0xtx2", edge.Target)
	assert.Equal(t, "0xobj2", edge.Object)
}
