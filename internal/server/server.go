package server

import (
	"io"
	"log/slog"
	"net/http"
	httpAdapter "snippet-vault-go/internal/adapters/http"
	"snippet-vault-go/internal/adapters/storage"
	"snippet-vault-go/internal/core"
	"time"
)

type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func loggingMiddleware(logger *slog.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		wrapped := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}

		next.ServeHTTP(wrapped, r)

		logger.Info("call",
			slog.String("method", r.Method),
			slog.String("path", r.URL.Path),
			slog.Int("status", wrapped.statusCode),
			slog.Duration("duration", time.Since(start)),
		)
	})
}

// New builds and returns a configured HTTP server without starting it yet.
func New(addr string, logOutput io.Writer, dbPath string) *http.Server {
	logger := slog.New(slog.NewTextHandler(logOutput, nil))

	// Swap MemoryRepo for SQLiteRepo!
	repo, err := storage.NewSQLiteRepo(dbPath)
	if err != nil {
		panic("failed to initialize sqlite database: " + err.Error())
	}
	service := core.NewSnippetService(repo)
	handler := httpAdapter.NewHandler(service)

	mux := http.NewServeMux()
	handler.RegisterRoutes(mux)

	return &http.Server{
		Addr:    addr,
		Handler: loggingMiddleware(logger, mux),
	}
}
