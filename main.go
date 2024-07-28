package main

import (
	"fmt"
	"go_jwt/env"
	"go_jwt/registry"
	"net/http"
)

func main() {
	cnf := env.LoadEnv()

	handler := registry.AuthRegistry(cnf)
	http.HandleFunc("/signup", handler.SignUpHandler)
	http.HandleFunc("/login", handler.LoginHandler)
	http.HandleFunc("/hello", handler.HelloHandler)
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Println("Error starting the server:", err)
	}
}
