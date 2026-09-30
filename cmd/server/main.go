package main

import (
	"fmt"
	"net/http"
	"os"
	"snippet-vault-go/internal/server"
)

var Version = "dev"

func main() {
	srv := server.New(":8080", os.Stdout)

	fmt.Printf("Standalone Vault Server %s listening on http://localhost:8080\n", Version)
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		fmt.Printf("Server error: %v\n", err)
		os.Exit(1)
	}
}
