import { http, HttpResponse } from "msw";
import type { components } from "../src/generated/todo-api";

type Todo = components["schemas"]["Todo"];
type ApiError = components["schemas"]["Error"];

// Seed matches both wireframes.dsl's TodoList table ("Buy milk", "Walk the
// dog") and todos.feature's story-1 scenario (Buy milk not done, Walk the dog
// done), so the walk and the acceptance scenarios see the same starting list.
//
// State lives in module scope, not on a server: a full page load (reload, a
// typed URL, a link that leaves the SPA) re-runs this module and resets to
// this seed. Only in-app navigation — the only kind this single-screen app
// has — carries a change forward.
let todos: Todo[] = [
  { id: "1", title: "Buy milk", done: false },
  { id: "2", title: "Walk the dog", done: true },
];
let nextId = 3;

function errorBody(code: number, message: string): ApiError {
  return { code, message };
}

export const handlers = [
  // No todos:read / scope check here — this contract declares no `oauth2`
  // scheme at all (todo-api/openapi.yaml has no `security` block), so
  // mock/authz/gateway.ts enforces nothing and every request reaches here.
  http.get("/api/todos", ({ request }) => {
    const url = new URL(request.url);
    const limit = Math.min(Number(url.searchParams.get("limit") ?? 20), 100);
    const offset = Number(url.searchParams.get("offset") ?? 0);
    const page = todos.slice(offset, offset + limit);
    return HttpResponse.json({
      count: todos.length,
      next: offset + limit < todos.length ? `/todos?limit=${limit}&offset=${offset + limit}` : null,
      previous: offset > 0 ? `/todos?limit=${limit}&offset=${Math.max(0, offset - limit)}` : null,
      data: page,
    });
  }),

  http.post("/api/todos", async ({ request }) => {
    const body = (await request.json()) as { title?: string };
    if (!body?.title || body.title.trim().length === 0) {
      return HttpResponse.json(errorBody(400, "The title was blank"), { status: 400 });
    }
    const created: Todo = { id: String(nextId++), title: body.title, done: false };
    todos = [...todos, created];
    return HttpResponse.json(created, { status: 201 });
  }),

  http.delete("/api/todos/:todoId", ({ params }) => {
    const before = todos.length;
    todos = todos.filter((t) => t.id !== params.todoId);
    if (todos.length === before) {
      return HttpResponse.json(errorBody(404, "No todo with that id"), { status: 404 });
    }
    return new HttpResponse(null, { status: 204 });
  }),

  http.post("/api/todos/:todoId/toggle", ({ params }) => {
    const todo = todos.find((t) => t.id === params.todoId);
    if (!todo) {
      return HttpResponse.json(errorBody(404, "No todo with that id"), { status: 404 });
    }
    const updated: Todo = { ...todo, done: !todo.done };
    todos = todos.map((t) => (t.id === todo.id ? updated : t));
    return HttpResponse.json(updated);
  }),
];
