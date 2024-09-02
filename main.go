package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	_ "github.com/lib/pq"
)

func main() {
	// Connection string (replace with your own)
	connStr := "postgres://username:password@localhost/dbname?sslmode=disable"

	// Open a connection to the database
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}
	defer db.Close()

	// Define the handler for the `/` endpoint
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Ping the database to check the connection
		err := db.Ping()
		if err != nil {
			http.Error(w, "Failed to ping the database", http.StatusInternalServerError)
			return
		}

		// Write "pong" to the response if the ping was successful
		fmt.Fprintln(w, "pong")
	})

	// Start the HTTP server
	port := ":8080"
	fmt.Printf("Server is running on port %s\n", port)
	err = http.ListenAndServe(port, nil)
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
