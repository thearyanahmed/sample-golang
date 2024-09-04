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
	connStr := "postgresql://db:AVNS_6iU_A7Lltq6S-GdVwE0@app-446cac4d-3248-4af4-b226-879f62f65df2-do-user-16220553-0.d.db.ondigitalocean.com:25060/db?sslmode=require"

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

		// Query the PostgreSQL version
		var version string
		err = db.QueryRow("SELECT version()").Scan(&version)
		if err != nil {
			http.Error(w, "Failed to retrieve PostgreSQL version", http.StatusInternalServerError)
			return
		}

		// Write the response with "pong" and the PostgreSQL version
		fmt.Fprintf(w, "pong\nPostgreSQL version: %s\n", version)
	})

	// Start the HTTP server
	port := ":8080"
	fmt.Printf("Server is running on port %s\n", port)
	err = http.ListenAndServe(port, nil)
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
