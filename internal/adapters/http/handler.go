package http

import (
	"encoding/json"
	"errors"
	"net/http"
	"snippet-vault-go/internal/core"
)

type Handler struct {
	snippetService *core.SnippetService
	linkService    *core.LinkService
}

func NewHandler(snippetService *core.SnippetService, linkService *core.LinkService) *Handler {
	return &Handler{
		snippetService: snippetService,
		linkService:    linkService,
	}
}

// RegisterRoutes binds your endpoints to the Go HTTP multiplexer
func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /snippets", h.handleGetAll)
	mux.HandleFunc("POST /snippets", h.handleCreate)
	mux.HandleFunc("PUT /snippets/{id}", h.handleUpdate)
	mux.HandleFunc("DELETE /snippets/{id}", h.handleDelete)
	
	mux.HandleFunc("GET /links", h.handleGetAllLinks)
	mux.HandleFunc("POST /links", h.handleCreateLink)
	mux.HandleFunc("PUT /links/{id}", h.handleUpdateLink)
	mux.HandleFunc("DELETE /links/{id}", h.handleDeleteLink)
}

func (h *Handler) handleGetAll(w http.ResponseWriter, r *http.Request) {
	snippets, err := h.snippetService.GetAllSnippets()
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

	if err := h.snippetService.CreateSnippet(&snippet); err != nil {
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

	if err := h.snippetService.UpdateSnippet(id, &snippet); err != nil {
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

	if err := h.snippetService.DeleteSnippet(id); err != nil {
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

func (h *Handler) handleGetAllLinks(w http.ResponseWriter, r *http.Request) {
	links, err := h.linkService.GetAllLinks()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(links)
}

func (h *Handler) handleCreateLink(w http.ResponseWriter, r *http.Request) {
	var link core.Link
	if err := json.NewDecoder(r.Body).Decode(&link); err != nil {
		http.Error(w, "Invalid JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.linkService.CreateLink(&link); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(`{"status":"created"}`))
}

func (h *Handler) handleUpdateLink(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	var link core.Link
	if err := json.NewDecoder(r.Body).Decode(&link); err != nil {
		http.Error(w, "Invalid JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.linkService.UpdateLink(id, &link); err != nil {
		if errors.Is(err, core.ErrNotFound) {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(link)
}

func (h *Handler) handleDeleteLink(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	if err := h.linkService.DeleteLink(id); err != nil {
		if errors.Is(err, core.ErrNotFound) {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
