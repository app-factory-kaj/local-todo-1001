screen TodoList "See, add, complete and remove todos"
  navbar "Todo"
  heading "My Todos"
  row
    input "What needs doing?"
    button "Add" primary
  table "Done | Title | "
    row "checkbox | Buy milk | button:Delete"
    row "checkbox | Walk the dog | button:Delete"

flow "Manage todos"
  description "A visitor lists, adds, completes and deletes todos on the one shared list"
  TodoList
