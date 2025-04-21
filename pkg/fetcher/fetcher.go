package fetcher

import (
	"context"
	"fmt"
	"github.com/block-vision/sui-go-sdk/models"
	"strconv"
	"time"

	"github.com/block-vision/sui-go-sdk/sui"
)

type TxSummary struct {
	Digest          string
	Function        string
	Sender          string
	InvolvedObjects []string
	Type            string
	Status          string
	GasUsed         int64
}

// Fetcher implements GraphDataFetcher using Sui JSON-RPC API.
type Fetcher struct {
	api sui.ISuiAPI
	ws  sui.ISuiWebsocketAPI
}

func New(api sui.ISuiAPI, ws sui.ISuiWebsocketAPI) *Fetcher {
	return &Fetcher{
		api: api,
		ws:  ws,
	}
}

func (f *Fetcher) FetchTransactionsByCheckpoint(ctx context.Context, checkpoint int64) ([]TxSummary, error) {
	// Step 1: Get checkpoint details
	cpResp, err := f.api.SuiGetCheckpoint(ctx, models.SuiGetCheckpointRequest{CheckpointID: fmt.Sprintf("%d", checkpoint)})
	if err != nil {
		return nil, fmt.Errorf("fetch checkpoint %d: %w", checkpoint, err)
	}
	if len(cpResp.Transactions) == 0 {
		return nil, nil
	}

	var results []TxSummary
	for _, digest := range cpResp.Transactions {
		tx, err := f.api.SuiGetTransactionBlock(ctx, models.SuiGetTransactionBlockRequest{
			Digest: digest,
			Options: models.SuiTransactionBlockOptions{
				ShowInput:          true,
				ShowEffects:        true,
				ShowEvents:         false,
				ShowObjectChanges:  false,
				ShowBalanceChanges: false,
			},
		})
		if err != nil {
			// Ignore individual tx errors, or log them
			continue
		}
		// Parse basic fields
		digest := tx.Digest
		sender := tx.Transaction.Data.Sender
		txKind := tx.Transaction.Data.Transaction.Kind
		moveCall := extractMoveFunction(tx.Transaction.Data.Transaction)

		var objectIDs []string
		for _, obj := range tx.Transaction.Data.Transaction.Inputs {
			if obj != nil && obj["objectId"] != nil {
				objectIDs = append(objectIDs, obj["objectId"].(string))
			}
		}

		// Set status := status (success or failure)
		status := tx.Effects.Status.Status

		gas := int64(0)
		computationCost, err := strconv.ParseInt(tx.Effects.GasUsed.ComputationCost, 10, 64)
		storageCost, err := strconv.ParseInt(tx.Effects.GasUsed.StorageCost, 10, 64)
		gas = computationCost + storageCost

		results = append(results, TxSummary{
			Digest:          digest,
			Function:        moveCall,
			Sender:          sender,
			InvolvedObjects: objectIDs,
			Type:            txKind,
			Status:          status,
			GasUsed:         gas,
		})
	}
	return results, nil
}

func (f *Fetcher) FetchCheckpointTimestamp(ctx context.Context, checkpoint int64) (time.Time, error) {
	cpResp, err := f.api.SuiGetCheckpoint(ctx, models.SuiGetCheckpointRequest{CheckpointID: fmt.Sprintf("%d", checkpoint)})
	if err != nil {
		return time.Time{}, fmt.Errorf("fetch checkpoint timestamp %d: %w", checkpoint, err)
	}
	ts, err := strconv.ParseInt(cpResp.TimestampMs, 10, 64)
	if err != nil {
		return time.Time{}, fmt.Errorf("fetch checkpoint timestamp %d: %w", checkpoint, err)
	}
	return time.UnixMilli(ts), nil
}

// extractMoveFunction returns the formatted Move call signature: package::module::function
func extractMoveFunction(tx models.SuiTransactionBlockKind) string {
	for _, pTx := range tx.Transactions {
		if enum, ok := pTx.(models.SuiTransactionEnum); ok {
			if enum.MoveCall != nil {
				return fmt.Sprintf("%s::%s::%s", enum.MoveCall.Package, enum.MoveCall.Module, enum.MoveCall.Function)
			}
		}
	}
	return ""
}
