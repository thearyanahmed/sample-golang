package main

import (
    "database/sql"
    "net/http"
    "strings"
    "fmt"
    "log"
    "math/rand"
    "os"
    _ "time"

    _ "github.com/go-sql-driver/mysql"
)

func logRequest(r *http.Request) {
    uri := r.RequestURI
    method := r.Method
    fmt.Println("Got request!", method, uri)
}

type User struct {
    ID   int
    Name string
    Age  int
}


type DBConfig struct {
	User     string
	Password string
	Host     string
	Port     int
	DBName   string
}

func main() {

    // Open a connection to the MySQL database
// Load the database configuration from environment variables
	config := DBConfig{
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		Host:     os.Getenv("DB_HOST"),
		Port:     25060, // Replace with your port
		DBName:   os.Getenv("DB_NAME"),
	}

	// Build the MySQL connection string using the struct
	connStr := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?ssl-mode=REQUIRED", config.User, config.Password, config.Host, config.Port, config.DBName)
    fmt.Println(connStr)
    db, err := sql.Open("mysql", connStr)
    if err != nil {
        log.Fatal("Error opening database connection: ", err)
    }
    defer db.Close()

    // Check if the connection is successful by pinging the database
    err = db.Ping()
    if err != nil {
        log.Fatal("Error pinging database: ", err)
    }


    http.HandleFunc("/setup", func(w http.ResponseWriter, r *http.Request) {
        logRequest(r)

        // Create the users table (if not already created)
        createTableSQL := `CREATE TABLE IF NOT EXISTS users (
            id INT AUTO_INCREMENT PRIMARY KEY,
            name VARCHAR(100),
            age INT
        );`
        _, err = db.Exec(createTableSQL)
        if err != nil {
            fmt.Fprintf(w, "something went wrong %v", err.Error())
        }

        fmt.Fprintf(w, "done")
    })

    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        logRequest(r)
        fmt.Fprintf(w, "Hello! you've requested %s\n", r.URL.Path)
    })


    http.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {
        logRequest(r)

        // Retrieve and display the inserted values
        rows, err := db.Query("SELECT id, name, age FROM users")
        if err != nil {
            log.Fatal("Error retrieving data: ", err)
        }
        defer rows.Close()

        var users []User
        for rows.Next() {
            var user User
            if err := rows.Scan(&user.ID, &user.Name, &user.Age); err != nil {
                log.Fatal("Error scanning row: ", err)
            }
            users = append(users, user)
        }

        // Check for any row iteration errors
        if err := rows.Err(); err != nil {
            log.Fatal("Error during rows iteration: ", err)
        }

        // Print the retrieved users
        fmt.Println("Retrieved Users:")
        for _, user := range users {
            fmt.Printf("ID: %d, Name: %s, Age: %d\n", user.ID, user.Name, user.Age)
        }

        fmt.Fprint(w, "\nusers get query")
    })


    http.HandleFunc("/seed/users", func(w http.ResponseWriter, r *http.Request) {
        logRequest(r)

        // Insert 10 random values into the users table
        insertSQL := "INSERT INTO users (name, age) VALUES (?, ?)"
        for i := 0; i < 10; i++ {
            name := fmt.Sprintf("User%d", rand.Intn(1000))
            age := rand.Intn(60) + 18 // Random age between 18 and 77
            _, err := db.Exec(insertSQL, name, age)
            if err != nil {
                log.Fatal("Error inserting data: ", err)
            }
        }

        fmt.Fprint(w, "requestID.String()")
    })

    port := os.Getenv("PORT")
    if port == "" {
        port = "8001"
    }

    for _, encodedRoute := range strings.Split(os.Getenv("ROUTES"), ",") {
        if encodedRoute == "" {
            continue
        }
        pathAndBody := strings.SplitN(encodedRoute, "=", 2)
        path, body := pathAndBody[0], pathAndBody[1]
        http.HandleFunc("/"+path, func(w http.ResponseWriter, r *http.Request) {
            fmt.Fprint(w, body)
        })
    }

    bindAddr := fmt.Sprintf(":%s", port)
    fmt.Println()
    fmt.Printf("==> Server listening at %s 🚀\n", bindAddr)

    if err := http.ListenAndServe(bindAddr, nil); err != nil {
        panic(err)
    }
}
