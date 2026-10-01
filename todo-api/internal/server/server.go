// Package server assembles the routed HTTP handler for todo-api, so tests can
// exercise the real routing instead of calling handler funcs directly.
package server

import (
	"net/http"

	"todo-api/internal/handlers"
	"todo-api/internal/store"
)

// New returns the full todo-api HTTP handler wired to a fresh in-memory store.
func New() http.Handler {
	return NewWithStore(store.New())
}

// NewWithStore returns the HTTP handler wired to the given store — tests use
// this to inspect or pre-seed store state.
func NewWithStore(s *store.Store) http.Handler {
	h := handlers.New(s)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /todos", h.List)
	mux.HandleFunc("POST /todos", h.Create)
	mux.HandleFunc("POST /todos/{todoId}/toggle", h.Toggle)
	mux.HandleFunc("DELETE /todos/{todoId}", h.Delete)

	return mux
}
