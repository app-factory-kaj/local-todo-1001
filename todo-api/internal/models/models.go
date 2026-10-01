// Package models holds the request/response/domain types for the todo-api
// service, matching specs/design/components/todo-api/openapi.yaml.
package models

// Todo is the domain entity: a server-generated id, a non-blank title, and
// a done flag that starts false and is flipped by the toggle action.
type Todo struct {
	ID    string `json:"id"`
	Title string `json:"title"`
	Done  bool   `json:"done"`
}

// NewTodo is the create request body.
type NewTodo struct {
	Title string `json:"title"`
}

// ErrorResponse is the shared Error schema returned on every 4xx/5xx.
type ErrorResponse struct {
	Code        int    `json:"code"`
	Message     string `json:"message"`
	Description string `json:"description,omitempty"`
	MoreInfo    string `json:"moreInfo,omitempty"`
}

// ListResponse is the paginated envelope returned by GET /todos.
type ListResponse struct {
	Count    int     `json:"count"`
	Next     *string `json:"next"`
	Previous *string `json:"previous"`
	Data     []*Todo `json:"data"`
}
