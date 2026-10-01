package server_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"todo-api/internal/models"
	"todo-api/internal/server"
)

func newTestServer() *httptest.Server {
	return httptest.NewServer(server.New())
}

func createTodo(t *testing.T, base, title string) models.Todo {
	t.Helper()
	body, _ := json.Marshal(models.NewTodo{Title: title})
	resp, err := http.Post(base+"/todos", "application/json", bytes.NewReader(body))
	if err != nil {
		t.Fatalf("POST /todos: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("POST /todos: got status %d, want 201", resp.StatusCode)
	}
	var out models.Todo
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		t.Fatalf("decode created todo: %v", err)
	}
	return out
}

func TestListTodosEmpty(t *testing.T) {
	srv := newTestServer()
	defer srv.Close()

	resp, err := http.Get(srv.URL + "/todos")
	if err != nil {
		t.Fatalf("GET /todos: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("got status %d, want 200", resp.StatusCode)
	}
	var out models.ListResponse
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		t.Fatalf("decode list response: %v", err)
	}
	if out.Count != 0 {
		t.Errorf("count = %d, want 0", out.Count)
	}
	if len(out.Data) != 0 {
		t.Errorf("data length = %d, want 0", len(out.Data))
	}
	if out.Next != nil {
		t.Errorf("next = %v, want nil", *out.Next)
	}
	if out.Previous != nil {
		t.Errorf("previous = %v, want nil", *out.Previous)
	}
}

func TestListTodosReflectsOpenAndDone(t *testing.T) {
	srv := newTestServer()
	defer srv.Close()

	milk := createTodo(t, srv.URL, "Buy milk")
	dog := createTodo(t, srv.URL, "Walk the dog")

	resp, err := http.Post(fmt.Sprintf("%s/todos/%s/toggle", srv.URL, dog.ID), "", nil)
	if err != nil {
		t.Fatalf("toggle: %v", err)
	}
	resp.Body.Close()

	resp, err = http.Get(srv.URL + "/todos")
	if err != nil {
		t.Fatalf("GET /todos: %v", err)
	}
	defer resp.Body.Close()
	var out models.ListResponse
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if out.Count != 2 {
		t.Fatalf("count = %d, want 2", out.Count)
	}
	byID := map[string]models.Todo{}
	for _, d := range out.Data {
		byID[d.ID] = *d
	}
	if byID[milk.ID].Done {
		t.Errorf("Buy milk should not be done")
	}
	if !byID[dog.ID].Done {
		t.Errorf("Walk the dog should be done")
	}
}

func TestListTodosPagination(t *testing.T) {
	srv := newTestServer()
	defer srv.Close()

	for i := 0; i < 5; i++ {
		createTodo(t, srv.URL, fmt.Sprintf("todo-%d", i))
	}

	resp, err := http.Get(srv.URL + "/todos?limit=2&offset=0")
	if err != nil {
		t.Fatalf("GET /todos: %v", err)
	}
	defer resp.Body.Close()
	var out models.ListResponse
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if out.Count != 5 {
		t.Fatalf("count = %d, want 5", out.Count)
	}
	if len(out.Data) != 2 {
		t.Fatalf("data length = %d, want 2", len(out.Data))
	}
	if out.Previous != nil {
		t.Errorf("previous = %v, want nil at offset 0", *out.Previous)
	}
	if out.Next == nil {
		t.Fatalf("next = nil, want a next page link")
	}
	if want := "/todos?limit=2&offset=2"; *out.Next != want {
		t.Errorf("next = %q, want %q", *out.Next, want)
	}

	resp2, err := http.Get(srv.URL + *out.Next)
	if err != nil {
		t.Fatalf("GET next page: %v", err)
	}
	defer resp2.Body.Close()
	var out2 models.ListResponse
	if err := json.NewDecoder(resp2.Body).Decode(&out2); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if out2.Previous == nil || *out2.Previous != "/todos?limit=2&offset=0" {
		t.Errorf("previous on page 2 = %v, want /todos?limit=2&offset=0", out2.Previous)
	}
}

func TestCreateTodoAddsToList(t *testing.T) {
	srv := newTestServer()
	defer srv.Close()

	createTodo(t, srv.URL, "a")
	createTodo(t, srv.URL, "b")
	created := createTodo(t, srv.URL, "Pay rent")

	if created.Title != "Pay rent" {
		t.Errorf("title = %q, want %q", created.Title, "Pay rent")
	}
	if created.Done {
		t.Errorf("new todo should not be done")
	}

	resp, err := http.Get(srv.URL + "/todos")
	if err != nil {
		t.Fatalf("GET /todos: %v", err)
	}
	defer resp.Body.Close()
	var out models.ListResponse
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if out.Count != 3 {
		t.Fatalf("count = %d, want 3", out.Count)
	}
}

func TestCreateTodoBlankTitleRejected(t *testing.T) {
	srv := newTestServer()
	defer srv.Close()

	createTodo(t, srv.URL, "a")
	createTodo(t, srv.URL, "b")

	for _, title := range []string{"", "   ", "\t\n"} {
		body, _ := json.Marshal(models.NewTodo{Title: title})
		resp, err := http.Post(srv.URL+"/todos", "application/json", bytes.NewReader(body))
		if err != nil {
			t.Fatalf("POST /todos: %v", err)
		}
		if resp.StatusCode != http.StatusBadRequest {
			t.Fatalf("title %q: got status %d, want 400", title, resp.StatusCode)
		}
		var errOut models.ErrorResponse
		if err := json.NewDecoder(resp.Body).Decode(&errOut); err != nil {
			t.Fatalf("decode error body: %v", err)
		}
		resp.Body.Close()
		if errOut.Code != http.StatusBadRequest || errOut.Message == "" {
			t.Errorf("title %q: error body = %+v, want code 400 and a message", title, errOut)
		}
	}

	resp, err := http.Get(srv.URL + "/todos")
	if err != nil {
		t.Fatalf("GET /todos: %v", err)
	}
	defer resp.Body.Close()
	var out models.ListResponse
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if out.Count != 2 {
		t.Fatalf("count = %d, want 2 (blank titles must not be added)", out.Count)
	}
}

func TestToggleTodo(t *testing.T) {
	srv := newTestServer()
	defer srv.Close()

	milk := createTodo(t, srv.URL, "Buy milk")

	resp, err := http.Post(fmt.Sprintf("%s/todos/%s/toggle", srv.URL, milk.ID), "", nil)
	if err != nil {
		t.Fatalf("toggle: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("got status %d, want 200", resp.StatusCode)
	}
	var out models.Todo
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if !out.Done {
		t.Errorf("Buy milk should be marked done after toggle")
	}

	resp2, err := http.Post(fmt.Sprintf("%s/todos/%s/toggle", srv.URL, milk.ID), "", nil)
	if err != nil {
		t.Fatalf("toggle again: %v", err)
	}
	defer resp2.Body.Close()
	var out2 models.Todo
	if err := json.NewDecoder(resp2.Body).Decode(&out2); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if out2.Done {
		t.Errorf("Buy milk should be reopened after second toggle")
	}
}

func TestToggleUnknownIdIs404(t *testing.T) {
	srv := newTestServer()
	defer srv.Close()

	resp, err := http.Post(srv.URL+"/todos/does-not-exist/toggle", "", nil)
	if err != nil {
		t.Fatalf("toggle: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusNotFound {
		t.Fatalf("got status %d, want 404", resp.StatusCode)
	}
	var errOut models.ErrorResponse
	if err := json.NewDecoder(resp.Body).Decode(&errOut); err != nil {
		t.Fatalf("decode error body: %v", err)
	}
	if errOut.Code != http.StatusNotFound || errOut.Message == "" {
		t.Errorf("error body = %+v, want code 404 and a message", errOut)
	}
}

func TestDeleteTodo(t *testing.T) {
	srv := newTestServer()
	defer srv.Close()

	milk := createTodo(t, srv.URL, "Buy milk")
	createTodo(t, srv.URL, "Walk the dog")
	createTodo(t, srv.URL, "Pay rent")

	req, _ := http.NewRequest(http.MethodDelete, fmt.Sprintf("%s/todos/%s", srv.URL, milk.ID), nil)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("DELETE: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusNoContent {
		t.Fatalf("got status %d, want 204", resp.StatusCode)
	}

	listResp, err := http.Get(srv.URL + "/todos")
	if err != nil {
		t.Fatalf("GET /todos: %v", err)
	}
	defer listResp.Body.Close()
	var out models.ListResponse
	if err := json.NewDecoder(listResp.Body).Decode(&out); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if out.Count != 2 {
		t.Fatalf("count = %d, want 2", out.Count)
	}
	for _, d := range out.Data {
		if d.ID == milk.ID {
			t.Errorf("Buy milk should no longer appear in the list")
		}
	}
}

func TestDeleteUnknownIdIs404(t *testing.T) {
	srv := newTestServer()
	defer srv.Close()

	req, _ := http.NewRequest(http.MethodDelete, srv.URL+"/todos/does-not-exist", nil)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("DELETE: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusNotFound {
		t.Fatalf("got status %d, want 404", resp.StatusCode)
	}
	var errOut models.ErrorResponse
	if err := json.NewDecoder(resp.Body).Decode(&errOut); err != nil {
		t.Fatalf("decode error body: %v", err)
	}
	if errOut.Code != http.StatusNotFound || errOut.Message == "" {
		t.Errorf("error body = %+v, want code 404 and a message", errOut)
	}
}
