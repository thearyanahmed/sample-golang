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
	connStr := "postgresql://db:AVNS_uslm4DS6tbBSr9xX7Nr@app-e9154229-f057-4817-8098-efe18e6e84b1-do-user-16220553-0.f.db.ondigitalocean.com:25060/db?sslmode=require"

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
