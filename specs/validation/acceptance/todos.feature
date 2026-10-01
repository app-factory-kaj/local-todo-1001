Feature: Todo list

  @story-1
  Rule: The todo list shows every todo and its done state

    Scenario: Viewing a list with open and done todos
      Given a todo named "Buy milk" exists and is not done
      And a todo named "Walk the dog" exists and is done
      When Jamie the user opens the todo list
      Then Jamie sees "Buy milk" marked as not done
      And Jamie sees "Walk the dog" marked as done

  @story-2
  Rule: A new todo can be added with a title

    Scenario: Adding a todo
      Given the todo list has 2 todos
      When Jamie the user adds a todo titled "Pay rent"
      Then the todo list has 3 todos
      And "Pay rent" appears in the list, not done

    @negative
    Scenario: A blank title is refused
      Given the todo list has 2 todos
      When Jamie the user tries to add a todo with a blank title
      Then the todo list still has 2 todos

  @story-3
  Rule: A todo's done state can be toggled

    Scenario: Marking a todo done
      Given a todo named "Buy milk" exists and is not done
      When Jamie the user toggles "Buy milk"
      Then "Buy milk" is marked as done

    Scenario: Reopening a done todo
      Given a todo named "Walk the dog" exists and is done
      When Jamie the user toggles "Walk the dog"
      Then "Walk the dog" is marked as not done

  @story-4
  Rule: A todo can be deleted

    Scenario: Deleting a todo
      Given a todo named "Buy milk" exists
      And the todo list has 3 todos
      When Jamie the user deletes "Buy milk"
      Then the todo list has 2 todos
      And "Buy milk" no longer appears in the list
