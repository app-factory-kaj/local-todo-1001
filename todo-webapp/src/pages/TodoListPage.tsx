import { useCallback, useEffect, useState, type JSX } from "react";
import {
  Alert,
  Box,
  Button,
  Checkbox,
  CircularProgress,
  ListingTable,
  PageContent,
  PageTitle,
  Stack,
  TextField,
  Typography,
} from "@wso2/oxygen-ui";
import { todoApi } from "../api";
import type { components } from "../generated/todo-api";

type Todo = components["schemas"]["Todo"];

// wireframes.dsl TodoList: a navbar, a heading, an add row (input + primary
// button), and a table of Done | Title | (actions) with a checkbox-toggle and
// a Delete button per row. The whole shared list lives in todo-api; this page
// holds no state of its own beyond what it fetches.
export default function TodoListPage(): JSX.Element {
  const [todos, setTodos] = useState<Todo[]>([]);
  const [loading, setLoading] = useState(true);
  const [loadError, setLoadError] = useState<string | null>(null);

  const [newTitle, setNewTitle] = useState("");
  const [addAttempted, setAddAttempted] = useState(false);
  const [addError, setAddError] = useState<string | null>(null);
  const [adding, setAdding] = useState(false);

  // One todo id mid-toggle/mid-delete, so its own row shows a spinner rather
  // than the whole table blocking on one request.
  const [pendingId, setPendingId] = useState<string | null>(null);

  const loadTodos = useCallback(async () => {
    setLoading(true);
    setLoadError(null);
    const { data, error } = await todoApi.GET("/todos", { params: { query: { limit: 100 } } });
    if (error) {
      // listTodos' contract declares no error response — any non-2xx here is
      // unexpected (network failure, 5xx), so there is no typed .message to read.
      setLoadError("Couldn't load the todo list. Please try again.");
    } else {
      setTodos(data.data);
    }
    setLoading(false);
  }, []);

  useEffect(() => {
    void loadTodos();
  }, [loadTodos]);

  const titleIsBlank = newTitle.trim().length === 0;

  const handleAdd = useCallback(async () => {
    setAddAttempted(true);
    if (titleIsBlank) return; // refused client-side, same rule the server enforces
    setAdding(true);
    setAddError(null);
    const { data, error } = await todoApi.POST("/todos", { body: { title: newTitle.trim() } });
    if (error) {
      setAddError(error.message);
    } else {
      setTodos((current) => [...current, data]);
      setNewTitle("");
      setAddAttempted(false);
    }
    setAdding(false);
  }, [newTitle, titleIsBlank]);

  const handleToggle = useCallback(async (todo: Todo) => {
    setPendingId(todo.id);
    const { data, error } = await todoApi.POST("/todos/{todoId}/toggle", {
      params: { path: { todoId: todo.id } },
    });
    if (!error) {
      setTodos((current) => current.map((t) => (t.id === todo.id ? data : t)));
    }
    setPendingId(null);
  }, []);

  const handleDelete = useCallback(async (todo: Todo) => {
    setPendingId(todo.id);
    const { error } = await todoApi.DELETE("/todos/{todoId}", {
      params: { path: { todoId: todo.id } },
    });
    if (!error) {
      setTodos((current) => current.filter((t) => t.id !== todo.id));
    }
    setPendingId(null);
  }, []);

  return (
    <PageContent>
      <PageTitle>
        <PageTitle.Header>My Todos</PageTitle.Header>
      </PageTitle>

      <Stack direction="row" spacing={2} alignItems="flex-start" sx={{ mb: 3 }}>
        <TextField
          label="What needs doing?"
          value={newTitle}
          onChange={(e) => setNewTitle(e.target.value)}
          error={addAttempted && titleIsBlank}
          helperText={addAttempted && titleIsBlank ? "Title can't be blank" : " "}
          fullWidth
        />
        <Button variant="contained" onClick={() => void handleAdd()} disabled={titleIsBlank || adding}>
          Add
        </Button>
      </Stack>

      {addError && (
        <Alert severity="error" sx={{ mb: 2 }}>
          {addError}
        </Alert>
      )}

      {loadError && (
        <Alert severity="error" sx={{ mb: 2 }}>
          {loadError}
        </Alert>
      )}

      {loading ? (
        <Box sx={{ display: "flex", justifyContent: "center", p: 4 }}>
          <CircularProgress />
        </Box>
      ) : (
        <ListingTable.Container>
          <ListingTable>
            <ListingTable.Head>
              <ListingTable.Row>
                <ListingTable.Cell>Done</ListingTable.Cell>
                <ListingTable.Cell>Title</ListingTable.Cell>
                <ListingTable.Cell />
              </ListingTable.Row>
            </ListingTable.Head>
            <ListingTable.Body>
              {todos.map((todo) => (
                <ListingTable.Row key={todo.id}>
                  <ListingTable.Cell>
                    <Checkbox
                      checked={todo.done}
                      onChange={() => void handleToggle(todo)}
                      disabled={pendingId === todo.id}
                      inputProps={{ "aria-label": `Mark "${todo.title}" as ${todo.done ? "not done" : "done"}` }}
                    />
                  </ListingTable.Cell>
                  <ListingTable.Cell>
                    <Typography
                      sx={{ textDecoration: todo.done ? "line-through" : "none" }}
                      color={todo.done ? "text.secondary" : "text.primary"}
                    >
                      {todo.title}
                    </Typography>
                  </ListingTable.Cell>
                  <ListingTable.Cell>
                    <Button
                      variant="outlined"
                      color="error"
                      size="small"
                      onClick={() => void handleDelete(todo)}
                      disabled={pendingId === todo.id}
                    >
                      Delete
                    </Button>
                  </ListingTable.Cell>
                </ListingTable.Row>
              ))}
            </ListingTable.Body>
          </ListingTable>
          {todos.length === 0 && (
            <ListingTable.EmptyState title="No todos yet" description="Add your first todo above." />
          )}
        </ListingTable.Container>
      )}
    </PageContent>
  );
}
