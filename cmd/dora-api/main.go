package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/dora/dora/internal/configuration"
)

func main() {
	fmt.Printf("DORA REST API Server %s\n", configuration.Version)
	
	http.HandleFunc("/api/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintln(w, `{"status":"ok","version":"`+configuration.Version+`"}`)
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Printf("Listening on :%s...\n", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		fmt.Fprintf(os.Stderr, "API server failed: %v\n", err)
		os.Exit(1)
	}
}
