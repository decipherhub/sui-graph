package sui_graphd

import (
	"fmt"
	"github.com/yourname/sui-graph/config"
	"log"
)

func main() {
	cfg, err := config.LoadConfig(".")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	fmt.Printf("Sui RPC URL: %s\n", cfg.Sui.RPCUrl)
	fmt.Printf("DB DSN: %s\n", cfg.Database.DSN)
}
