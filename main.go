package main

import (
    "fmt"
    "net/http"
    "os"
    "sort"
    "strings"
    "time"
    "html/template"
)

type IndexData struct {
    Now time.Time
}
func logRequest(r *http.Request) {
    uri := r.RequestURI
    method := r.Method
    fmt.Println("Got request!", method, uri)
}
func main() {
    port := os.Getenv("PORT")
    if port == "" {
        port = "8080"
    }

    http.HandleFunc("/headers", func(w http.ResponseWriter, r *http.Request) {
        logRequest(r)
        fmt.Fprintf(w, "Hello! you've requested %s\n", r.URL.Path)

        headerNames := make([]string, 0, len(r.Header))
        for name := range r.Header {
            headerNames = append(headerNames, name)
        }

        sort.Strings(headerNames)

        for _, name := range headerNames {
            headers := r.Header[name]
            for _, h := range headers {
                fmt.Fprintf(w, " %v = %v\n", name, h)
            }
        }
    })


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


    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        logRequest(r)
        now := time.Now()
        data := IndexData{Now: now}

        tmpl := `<!DOCTYPE html>
        <html>
        <head>
        <title>Index</title>
        </head>
        <body>
        <h1>Hello from Go!</h1>
        <p>Current time: {{.Now}}</p>
        <p>This page is to test Cloudflare caching.</p>
        <p>Random number for caching test: {{randInt 1000000}}</p>
        </body>
        </html>`

        funcMap := template.FuncMap{
            "randInt": func(max int) int {
                return int(time.Now().UnixNano() % int64(max))
            },
        }

        t, err := template.New("index").Funcs(funcMap).Parse(tmpl)
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }

        w.Header().Set("Cache-Control", "public, max-age=30")
        err = t.Execute(w, data)
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
        }
    })
    bindAddr := fmt.Sprintf(":%s", port)
    fmt.Printf("\n ==> Server listening at %s 🚀\n", bindAddr)

    if err := http.ListenAndServe(bindAddr, nil); err != nil {
        panic(err)
    }
}

