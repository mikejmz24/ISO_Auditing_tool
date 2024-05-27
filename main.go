package main

import (
	"fmt"
	"github.com/a-h/templ"
	"iso_auditing_tool/main/views/components"
	"net/http"
)

func main() {
	component := components.Greeting()
	http.Handle("/", templ.Handler(component))
	fmt.Printf("Listening on :3000")
	err := http.ListenAndServe(":3000", nil)
	if err != nil {
		fmt.Printf("Error starting server: %v\n", err)
	}
}
