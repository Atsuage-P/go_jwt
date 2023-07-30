package main

import (
	"fmt"
	"go_oauth/registry"
	"net/http"
)

func main() {
	handler := registry.AuthRegistry()
	http.HandleFunc("/signup", handler.SignUpHandler)
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Println("Error starting the server:", err)
	}
}
