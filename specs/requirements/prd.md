# local-todo-1001 — PRD

## Problem Statement

People who just want to jot down and track a short list of tasks are often forced into todo apps that demand an account, persistent storage, and a login flow before they can type a single item. They need a lightweight way to list, add, complete, and remove tasks right now, without any setup or commitment.

## Solution

A minimal todo application: a small Go REST API that keeps todos (title, done) in memory, and a single-page React web app that calls it to list, create, toggle, and delete todos. There is no sign-in and no database — the whole system is intentionally as small as possible.

## Actors

- **User** — anyone who opens the web app. There are no accounts or roles; every visitor has the same capabilities against the same todo list.

## User Stories

1. As a User, I want to see the list of all todos, so that I know what's outstanding and what's done.
2. As a User, I want to add a new todo with a title, so that I can capture a new task.
3. As a User, I want to toggle a todo's done state, so that I can mark it complete or reopen it.
4. As a User, I want to delete a todo, so that I can remove tasks I no longer need.

## Product Decisions

- **Authentication**: none. The app has no sign-in; it is a single, unauthenticated experience for every visitor (per the brief).
- **Storage**: in-memory only, inside the Go API process; no database is used (per the brief). *assumed*: a server restart permanently clears all todos, and this is expected behavior rather than a defect.
- **Scope of the list**: a single shared todo list visible to and editable by every visitor — there are no per-user or per-browser lists. *assumed*
- **Title validation**: the create-todo action rejects a blank/empty title, both in the API and in the web app's form. *assumed*
- **External services**: none. The product needs no third-party integrations.
- **Agents**: none. No story in this product involves summarizing, drafting, or interpreting free-form content — todos are stored and displayed verbatim.

## Out of Scope

- User accounts, sign-in, or per-user todo lists.
- Persistent storage (database, file, or otherwise) — data does not survive a server restart.
- Editing an existing todo's title.
- Due dates, priorities, categories, tags, search, or filtering.
- Multi-list or multi-board support.
- Offline support or real-time sync across multiple open browser tabs/clients.

## Open Questions

None at this time.