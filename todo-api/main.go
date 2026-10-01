// Command todo-api is a minimal in-memory todo list REST service.
// See specs/design/components/todo-api/openapi.yaml for the contract.
package main

import (
	"log"
	"net/http"
	"os"

	"todo-api/internal/server"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "9090"
	}

	handler := server.New()

	log.Printf("todo-api listening on :%s", port)
	if err := http.ListenAndServe(":"+port, handler); err != nil {
		log.Fatal(err)
	}
}
