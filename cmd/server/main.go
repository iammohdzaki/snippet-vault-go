package main

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"snippet-vault-go/internal/server"
)

var Version = "dev"

func getDBPath() string {
	if env := os.Getenv("SNIPPET_VAULT_DB"); env != "" {
		return env
	}
	if Version == "dev" {
		return "vault-dev.db"
	}
	homeDir, _ := os.UserHomeDir()
	return filepath.Join(homeDir, ".snippet-vault", "vault.db")
}

func main() {
	srv := server.New(":8080", os.Stdout, getDBPath())

	fmt.Printf("Standalone Vault Server %s listening on http://localhost:8080\n", Version)
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		fmt.Printf("Server error: %v\n", err)
		os.Exit(1)
	}
}
