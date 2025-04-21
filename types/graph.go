package types

type GraphNodeDTO struct {
	ID          string            `json:"id"`
	Label       string            `json:"label"`
	TimestampMs int64             `json:"timestampMs"`
	Sender      string            `json:"sender"`
	Objects     []string          `json:"objects"`
	Type        string            `json:"type"`
	Status      string            `json:"status"`
	GasUsed     int64             `json:"gasUsed"`
	Tags        map[string]string `json:"tags,omitempty"`
	Highlight   bool              `json:"highlight"`
}

type GraphEdgeDTO struct {
	Source         string `json:"source"`
	Target         string `json:"target"`
	Object         string `json:"object"`
	DependencyType string `json:"dependencyType"` // mutate / read / delete
}

type GraphResponseDTO struct {
	Checkpoint     int64          `json:"checkpoint"`
	TimestampMs    int64          `json:"timestampMs"`
	NodeCount      int            `json:"nodeCount"`
	EdgeCount      int            `json:"edgeCount"`
	ElapsedQueryMs int64          `json:"elapsedQueryMs"`
	Nodes          []GraphNodeDTO `json:"nodes"`
	Edges          []GraphEdgeDTO `json:"edges"`
}
