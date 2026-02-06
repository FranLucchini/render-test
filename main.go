package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
)

// getEnv retrieves an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	// Load .env file
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: Error loading .env file, using system environment variables")
	} else {
		log.Println(".env file loaded successfully")
	}
	value := os.Getenv(key)

	if value == "" {
		return defaultValue
	}
	// Print key and value
	return value
}

// handleRoot handles both GET and POST requests to the root path
func handleRoot(w http.ResponseWriter, r *http.Request) {
	verifyToken := os.Getenv("VERIFY_TOKEN")
	switch r.Method {
	case http.MethodGet:
		// Parse query parameters
		mode := r.URL.Query().Get("hub.mode")
		challenge := r.URL.Query().Get("hub.challenge")
		token := r.URL.Query().Get("hub.verify_token")

		// Verify webhook
		if mode == "subscribe" && token == verifyToken {
			log.Println("WEBHOOK VERIFIED")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(challenge))
		} else {
			log.Println("Webhook verification failed")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("holo"))
		}

	case http.MethodPost:
		// Log timestamp
		timestamp := time.Now().Format("2006-01-02 15:04:05")
		fmt.Printf("\n\nWebhook received %s\n\n", timestamp)

		// Read and parse body
		body, err := io.ReadAll(r.Body)
		if err != nil {
			log.Printf("Error reading body: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		defer r.Body.Close()

		// Pretty print JSON body
		var bodyJSON interface{}
		if err := json.Unmarshal(body, &bodyJSON); err != nil {
			log.Printf("Error parsing JSON: %v", err)
			// Still log the raw body if JSON parsing fails
			fmt.Println(string(body))
		} else {
			prettyJSON, _ := json.MarshalIndent(bodyJSON, "", "  ")
			fmt.Println(string(prettyJSON))
		}

		w.WriteHeader(http.StatusOK)

	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func main() {
	// Get port from environment or use default
	port := getEnv("PORT", "3000")

	// Register handler
	http.HandleFunc("/", handleRoot)

	// Start server
	fmt.Printf("\nListening on port %s\n\n", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}
