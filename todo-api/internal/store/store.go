// Package store is the entire persistence layer for todo-api: an in-process,
// mutex-guarded slice of Todo items. No database, no file persistence — data
// loss on restart is expected and correct, per the component's design.json.
package store

import (
	"strconv"
	"sync"

	"todo-api/internal/models"
)

// Store is a mutex-guarded, in-memory list of todos.
type Store struct {
	mu     sync.Mutex
	items  []*models.Todo
	nextID int
}

// New returns an empty store.
func New() *Store {
	return &Store{}
}

// List returns the page of todos starting at offset, up to limit items, and
// the total count of all todos (not just the page). Callers are expected to
// have already clamped limit/offset to their contracted bounds.
func (s *Store) List(limit, offset int) (page []*models.Todo, total int) {
	s.mu.Lock()
	defer s.mu.Unlock()

	total = len(s.items)
	if offset >= total {
		return []*models.Todo{}, total
	}
	end := offset + limit
	if end > total {
		end = total
	}
	page = make([]*models.Todo, end-offset)
	copy(page, s.items[offset:end])
	return page, total
}

// Create appends a new todo with a server-generated id and done=false.
func (s *Store) Create(title string) *models.Todo {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.nextID++
	t := &models.Todo{
		ID:    strconv.Itoa(s.nextID),
		Title: title,
		Done:  false,
	}
	s.items = append(s.items, t)
	return t
}

// Toggle flips the done state of the todo with the given id. It reports
// whether a todo with that id was found.
func (s *Store) Toggle(id string) (*models.Todo, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, t := range s.items {
		if t.ID == id {
			t.Done = !t.Done
			return t, true
		}
	}
	return nil, false
}

// Delete removes the todo with the given id. It reports whether a todo with
// that id was found.
func (s *Store) Delete(id string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	for i, t := range s.items {
		if t.ID == id {
			s.items = append(s.items[:i], s.items[i+1:]...)
			return true
		}
	}
	return false
}
