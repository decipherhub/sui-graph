package handler

import (
	"encoding/json"
	"github.com/decipherhub/sui-graph/internal/graph"
	"net/http"
	"strconv"
	"time"
)

func GetGraphByCheckpoint(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query().Get("checkpoint")
	cp, err := strconv.ParseInt(q, 10, 64)
	if err != nil {
		http.Error(w, "invalid checkpoint", http.StatusBadRequest)
		return
	}

	start := time.Now()
	nodes, edges, ts, err := graph.BuildGraphForCheckpoint(r.Context(), cp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	res := graph.GraphResponseDTO{
		Checkpoint:     cp,
		TimestampMs:    ts.UnixMilli(),
		NodeCount:      len(nodes),
		EdgeCount:      len(edges),
		ElapsedQueryMs: time.Since(start).Milliseconds(),
		Nodes:          nodes,
		Edges:          edges,
	}
	json.NewEncoder(w).Encode(res)
}
