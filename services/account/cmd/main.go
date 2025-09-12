/**
 * @description
 * Main entry point for the Account microservice.
 *
 * This file initializes and starts the HTTP server for the service.
 * For this boilerplate setup, it only includes a basic health check endpoint.
 * In later steps, this will be expanded to include configuration loading,
 * database connections, message broker setup, and route initialization.
 *
 * @dependencies
 * - Standard "fmt", "log", "net/http", "os" packages.
 */
package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {
	// A simple health check handler to verify the service is running.
	// This is essential for container orchestration systems like Kubernetes.
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		// Simple JSON response indicating the service is healthy.
		fmt.Fprintln(w, `{"status": "ok"}`)
	})

	serviceName := "Account Service"
	// Read the port from the environment variable for configurability in different environments.
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Default to port 8080 if not specified.
	}

	log.Printf("%s is starting on port %s...", serviceName, port)

	// Start the HTTP server and listen for incoming requests.
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("Failed to start server for %s: %v", serviceName, err)
	}
}
