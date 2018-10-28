package main

import (
	"github.com/gorilla/mux"
	"blackbox-kaizen/app"
	"os"
	"fmt"
	"net/http"
)

func main() {

	router := mux.NewRouter()
	router.Use(app.JwtAuthentication) // attach JWT auth middleware

	port := os.Getenv("PORT") // Get port from .env file
	if port == "" {
		port = "8000"
	}
	fmt.Println(port)

	err := http.ListenAndServe(":" + port, router) // Launch the app, visit localhost:8000
	if err != nil {
		fmt.Print(err)
	}
}