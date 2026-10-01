// Package handlers implements the HTTP operations of
// specs/design/components/todo-api/openapi.yaml against an in-memory store.
// This API is fully public: no security scheme, no identity to read.
package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"todo-api/internal/models"
	"todo-api/internal/store"
)

const (
	defaultLimit = 20
	maxLimit     = 100
)

// Handlers wires the HTTP layer to the in-memory store.
type Handlers struct {
	store *store.Store
}

// New returns handlers backed by the given store.
func New(s *store.Store) *Handlers {
	return &Handlers{store: s}
}

func writeJSON(w http.ResponseWriter, status int, body interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(body)
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, models.ErrorResponse{Code: status, Message: message})
}

// List handles GET /todos — a paginated list, envelope {count, next, previous, data}.
func (h *Handlers) List(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()

	limit := defaultLimit
	if v := q.Get("limit"); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			limit = n
		}
	}
	if limit <= 0 {
		limit = defaultLimit
	}
	if limit > maxLimit {
		limit = maxLimit
	}

	offset := 0
	if v := q.Get("offset"); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			offset = n
		}
	}
	if offset < 0 {
		offset = 0
	}

	page, total := h.store.List(limit, offset)

	var next, previous *string
	if offset+limit < total {
		n := fmt.Sprintf("/todos?limit=%d&offset=%d", limit, offset+limit)
		next = &n
	}
	if offset > 0 {
		p := offset - limit
		if p < 0 {
			p = 0
		}
		pv := fmt.Sprintf("/todos?limit=%d&offset=%d", limit, p)
		previous = &pv
	}

	writeJSON(w, http.StatusOK, models.ListResponse{
		Count:    total,
		Next:     next,
		Previous: previous,
		Data:     page,
	})
}

// Create handles POST /todos — rejects a blank/whitespace-only title with 400.
func (h *Handlers) Create(w http.ResponseWriter, r *http.Request) {
	var in models.NewTodo
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		writeError(w, http.StatusBadRequest, "request body must be valid JSON")
		return
	}
	if strings.TrimSpace(in.Title) == "" {
		writeError(w, http.StatusBadRequest, "title must not be blank")
		return
	}

	t := h.store.Create(in.Title)
	writeJSON(w, http.StatusCreated, t)
}

// Toggle handles POST /todos/{todoId}/toggle — flips done; 404 on unknown id.
func (h *Handlers) Toggle(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("todoId")

	t, ok := h.store.Toggle(id)
	if !ok {
		writeError(w, http.StatusNotFound, "no todo with that id")
		return
	}
	writeJSON(w, http.StatusOK, t)
}

// Delete handles DELETE /todos/{todoId} — removes it; 404 on unknown id.
func (h *Handlers) Delete(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("todoId")

	if !h.store.Delete(id) {
		writeError(w, http.StatusNotFound, "no todo with that id")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
