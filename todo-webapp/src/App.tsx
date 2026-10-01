import { Route, Routes } from "react-router-dom";
import AppLayout from "./layouts/AppLayout";
import TodoListPage from "./pages/TodoListPage";

// One screen (wireframes.dsl: TodoList), so no appRoutes.tsx indirection — a
// single route under the app shell, no additional routing.
export default function App() {
  return (
    <Routes>
      <Route element={<AppLayout />}>
        <Route path="/" element={<TodoListPage />} />
      </Route>
    </Routes>
  );
}
