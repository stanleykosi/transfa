/**
 * @description
 * Main entry point for the Scheduler microservice.
 *
 * This file initializes and starts the service. Unlike other services, its primary
 * role is to run scheduled background tasks (cron jobs) rather than serving a
 * HTTP API. The included health check is for monitoring purposes.
 * In later steps, this will be expanded to include the cron job scheduler and job definitions.
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
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, `{"status": "ok"}`)
	})

	serviceName := "Scheduler Service"
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Default port if not specified, mainly for health checks.
	}

	log.Printf("%s is starting...", serviceName)
	// In a real implementation, the cron scheduler would be started here.
	// For now, we just start the health check server.

	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("Failed to start server for %s: %v", serviceName, err)
	}
}
