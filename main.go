package main

import (
	"fmt"
	"go_jwt/env"
	"go_jwt/registry"
	"net/http"
	"time"
)

func main() {
	cnf := env.LoadEnv()

	handler := registry.AuthRegistry(cnf)
	http.HandleFunc("/signup", handler.SignUpHandler)
	http.HandleFunc("/login", handler.LoginHandler)
	http.HandleFunc("/hello", handler.HelloHandler)

	const ReadTime = 30
	const Port = "8080"
	s := http.Server{
		Addr:        ":" + Port,
		ReadTimeout: ReadTime * time.Second,
		Handler:     nil,
	}
	fmt.Println("Starting server at port:", Port)
	if err := s.ListenAndServe(); err != nil {
		fmt.Println("Error starting the server:", err)
	}
}
