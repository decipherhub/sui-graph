package server

import (
	"fmt"
	"github.com/decipherhub/sui-graph/server/router"
	"log"
	"net/http"
)

func Serve(port int) {

	r := router.SetupRouter()

	log.Printf("Server is running on port %d...", port)
	if err := http.ListenAndServe(fmt.Sprintf(":%d", port), r); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}

}
