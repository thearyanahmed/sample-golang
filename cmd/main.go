package main

import (
    "bytes"
    "encoding/json"
    "fmt"
    "net/http"
    "os"
    "github.com/joho/godotenv"
    "io"
    "strings"
)

func main ( ) {
    if len(os.Args) < 2 {
        fmt.Println("Usage: go run main.go [command] [flags]")
        os.Exit(1)
    }

    switch os.Args[1] {
    case "update-metadata":
        updateCustomMetadata(os.Args[2], os.Args[3])
    case "test-email-obs":
        testEmailObs(os.Args[2])
    case "test-cache":
        fmt.Println("todo!")
    default:
        fmt.Println("Unknown command:", os.Args[1])
        os.Exit(1)
    } 
}


func testEmailObs(currentState string) {
    if currentState == "" {
        fmt.Println("Error: --current-state flag is required.")
        os.Exit(1)
    }

    err := godotenv.Load()
    if err != nil {
        fmt.Println("Error loading .env file")
        os.Exit(1)
    }
    // Replace with your test domain
    testDomain := os.Getenv("TEST_DOMAIN")
    if testDomain == "" {
        fmt.Println("test domain not set")
        os.Exit(1)
    }

    resp, err := http.Get(testDomain)
    if err != nil {
        fmt.Println("Error making HTTP request:", err)
        os.Exit(1)
    }
    defer resp.Body.Close()

    bodyBytes, err := io.ReadAll(resp.Body)
    if err != nil {
        fmt.Println("Error reading response body:", err)
        os.Exit(1)
    }
    bodyString := string(bodyBytes)

    emailObsEnabled := strings.Contains(bodyString, "__cf_email__") // Or your specific check

    if currentState == "disabled" && emailObsEnabled {
        fmt.Println("Assertion failed: Email obfuscation should be disabled, but it is enabled.")
        os.Exit(1)
    } else if currentState == "enabled" && !emailObsEnabled {
        fmt.Println("Assertion failed: Email obfuscation should be enabled, but it is disabled.")
        os.Exit(1)
    }

    fmt.Println("Email obfuscation test successful.")
}

func updateCustomMetadata(emailObs, cache string) {
    err := godotenv.Load()
    if err != nil {
        fmt.Println("Error loading .env file")
    }

    zoneID := os.Getenv("CF_ZONE_ID")
    hostnameID := os.Getenv("CF_HOSTNAME_ID")
    authEmail := os.Getenv("CF_AUTH_EMAIL")
    authKey := os.Getenv("CF_AUTH_KEY")

    if zoneID == "" || hostnameID == "" || authEmail == "" || authKey == "" {
        fmt.Println("Error: Missing required environment variables.")
        fmt.Println("Please set CF_ZONE_ID, CF_HOSTNAME_ID, CF_AUTH_EMAIL, CF_AUTH_KEY")
        os.Exit(1)
    }

    apiEndpoint := fmt.Sprintf("https://api.cloudflare.com/client/v4/zones/%s/custom_hostnames/%s", zoneID, hostnameID)

    requestBody := map[string]any{
        "ssl": map[string]string{
            "method": "http",
            "type":   "dv",
        },
        "custom_metadata": map[string]string{
            "cf_cache":          cache,
            "email_obfuscation": emailObs,
            "some_other_data":   "helloworld1",
        },
    }

    requestBodyBytes, err := json.Marshal(requestBody)
    if err != nil {
        fmt.Println("Error marshaling JSON:", err)
        os.Exit(1)
    }

    req, err := http.NewRequest("PATCH", apiEndpoint, bytes.NewBuffer(requestBodyBytes))
    if err != nil {
        fmt.Println("Error creating request:", err)
        os.Exit(1)
    }

    req.Header.Set("X-Auth-Email", authEmail)
    req.Header.Set("X-Auth-Key", authKey)
    req.Header.Set("Content-Type", "application/json")

    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        fmt.Println("Error sending request:", err)
        os.Exit(1)
    }
    defer resp.Body.Close()

    if resp.StatusCode >= http.StatusBadRequest {
        fmt.Println("Request failed with status code:", resp.StatusCode)

        buf := new(bytes.Buffer)
        buf.ReadFrom(resp.Body)
        newStr := buf.String()
        fmt.Println("Response Body: ", newStr)

        os.Exit(1)
    }

    fmt.Println("Request successful")
}
