package main

import (
	"fmt"
	"log"
	"net/http"
	"user_db/config"
	"user_db/models"

	_ "user_db/docs" // This is required for Swagger documentation generation

	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"
)

// @title User Microservice API
// @version 1.0
// @description A simple microservice for managing user data.
// @host localhost:8080
// @BasePath /

// @contact.name Petr Dittrich
// @contact.url https://github.com/bodynek
// @contact.email bodynek@gmail.com
func main() {
	db, err := config.ConnectDB()
	if err != nil {
		log.Fatalf("Could not connect to the database: %v", err)
	}

	r := mux.NewRouter()

	// Routes
	r.HandleFunc("/save", models.SaveUser(db)).Methods("POST")
	r.HandleFunc("/{id}", models.GetUser(db)).Methods("GET")

	// Serve Swagger UI
	r.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	port := config.GetServerPort()
	log.Printf("Server is listening on port %s...\n", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), r))
}
