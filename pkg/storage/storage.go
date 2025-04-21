package storage

import (
	"context"
	"github.com/decipherhub/sui-graph/types"
	"time"
)

type Storage interface {
	SaveGraphCheckpoint(ctx context.Context, checkpoint int64, nodes []types.GraphNodeDTO, edges []types.GraphEdgeDTO, ts time.Time) error
	GetStoredCheckpoint(ctx context.Context, checkpoint int64) (*types.GraphResponseDTO, error)
}
