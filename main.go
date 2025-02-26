package main

import (
	"fmt"
	"runtime"
        "sync"
         "time"
	"net/http"
	"os"
	"strings"
)

func logRequest(r *http.Request) {
	uri := r.RequestURI
	method := r.Method
	fmt.Println("Got request!", method, uri)
}
func main() {


	port := os.Getenv("PORT")
	if port == "" {
		port = "8001"
	}
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		logRequest(r)
		fmt.Fprintf(w, "Hello! you've requested %s\n", r.URL.Path)
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

	bindAddr := fmt.Sprintf(":%s", port)
	fmt.Println()
	fmt.Printf("==> Server listening at %s 🚀\n", bindAddr)

	        graduallyIncreaseLoad(1*time.Minute, 5*time.Second, 10)

	if err := http.ListenAndServe(bindAddr, nil); err != nil {
		panic(err)
	}
}


// graduallyIncreaseLoad increases CPU load by spawning more goroutines over time.
func graduallyIncreaseLoad(duration time.Duration, increment time.Duration, goroutinesPerIncrement int) {
        startTime := time.Now()
        var wg sync.WaitGroup

        for time.Since(startTime) < duration {
                for i := 0; i < goroutinesPerIncrement; i++ {
                        wg.Add(1)
                        go func() {
                                defer wg.Done()
                                for {
                                        // Simple CPU-intensive loop.
                                        _ = 1 + 1
                                        // Allow other goroutines to run.
                                        runtime.Gosched()
                                }
                        }()
                }
                time.Sleep(increment)
        }

        // Wait for all goroutines to finish (though they won't, in this example).
        // This is mainly to prevent the program from exiting before the load is applied.
        time.Sleep(duration) // Keep running for the desired duration
        //wg.Wait() //Uncomment to wait, but the go routines will never finish.
}
