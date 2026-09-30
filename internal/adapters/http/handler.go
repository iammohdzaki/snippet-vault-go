package http

import (
	"encoding/json"
	"errors"
	"net/http"
	"snippet-vault-go/internal/core"
)

type Handler struct {
	service *core.SnippetService
}

func NewHandler(service *core.SnippetService) *Handler {
	return &Handler{service: service}
}

// RegisterRoutes binds your endpoints to the Go HTTP multiplexer
func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /snippets", h.handleGetAll)
	mux.HandleFunc("POST /snippets", h.handleCreate)
	mux.HandleFunc("PUT /snippets/{id}", h.handleUpdate)
	mux.HandleFunc("DELETE /snippets/{id}", h.handleDelete)
}

func (h *Handler) handleGetAll(w http.ResponseWriter, r *http.Request) {
	snippets, err := h.service.GetAllSnippets()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(snippets)
}

func (h *Handler) handleCreate(w http.ResponseWriter, r *http.Request) {
	var snippet core.Snippet
	if err := json.NewDecoder(r.Body).Decode(&snippet); err != nil {
		// Return the exact parsing error so your IDE HTTP client shows it
		http.Error(w, "Invalid JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.service.CreateSnippet(&snippet); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(`{"status":"created"}`))
}

func (h *Handler) handleUpdate(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	var snippet core.Snippet
	if err := json.NewDecoder(r.Body).Decode(&snippet); err != nil {
		http.Error(w, "Invalid JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.service.UpdateSnippet(id, &snippet); err != nil {
		if errors.Is(err, core.ErrNotFound) {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(snippet)
}

func (h *Handler) handleDelete(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	if err := h.service.DeleteSnippet(id); err != nil {
		if errors.Is(err, core.ErrNotFound) {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// 204 No Content is the standard HTTP response for a successful DELETE
	w.WriteHeader(http.StatusNoContent)
}
