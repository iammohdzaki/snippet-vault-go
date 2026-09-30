package server

import (
	"fmt"
	"net/http"
	"os"
)

func main() {
	srv := New(":8080", os.Stdout)

	fmt.Println("Standalone Vault Server listening on http://localhost:8080")
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		fmt.Printf("Server error: %v\n", err)
		os.Exit(1)
	}
}
