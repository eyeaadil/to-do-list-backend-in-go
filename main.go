package main

import (
	"log"
	"net/http"
	"to-do-list-api/config"
	"to-do-list-api/routes"

	"github.com/gorilla/mux"
)


func main() {
	// Connect to MongoDB
	config.ConnectDB()

	router := mux.NewRouter()

	// Register routes
	routes.RegisterRoutes(router)

	// Start server
	log.Println("Server running on port 8080")
	http.ListenAndServe(":8080", router)
}
