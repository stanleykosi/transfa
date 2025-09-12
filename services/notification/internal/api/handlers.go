/**
 * @description
 * This file contains the HTTP handlers for the Notification service. Handlers are responsible
 * for parsing incoming requests, calling the appropriate application service method,
 * and writing the HTTP response.
 *
 * @dependencies
 * - "io": For reading the request body.
 * - "log": For logging errors.
 * - "net/http": For standard HTTP handling.
 * - "transfa/services/notification/internal/app": Imports the application service layer.
 */
package api

import (
	"io"
	"log"
	"net/http"

	"transfa/services/notification/internal/app"
)

// NotificationHandler holds dependencies for the notification-related HTTP handlers.
type NotificationHandler struct {
	service *app.Service
}

// NewNotificationHandler creates a new handler with the given application service.
func NewNotificationHandler(service *app.Service) *NotificationHandler {
	return &NotificationHandler{
		service: service,
	}
}

// AnchorWebhookHandler handles incoming POST requests from Anchor.
func (h *NotificationHandler) AnchorWebhookHandler(w http.ResponseWriter, r *http.Request) {
	// 1. Read the entire request body. It's needed for signature verification.
	payload, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error reading webhook body: %v", err)
		http.Error(w, "Cannot read request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// 2. Get the signature from the header.
	signature := r.Header.Get("x-anchor-signature")
	if signature == "" {
		log.Println("Missing x-anchor-signature header")
		http.Error(w, "Missing signature header", http.StatusBadRequest)
		return
	}

	// 3. Pass the payload and signature to the application service for processing.
	if err := h.service.ProcessAnchorWebhook(r.Context(), payload, signature); err != nil {
		// Log the specific error, but return a generic error to the client.
		// The error could be due to a bad signature or an internal processing failure.
		log.Printf("Error processing Anchor webhook: %v", err)
		http.Error(w, "Webhook processing failed", http.StatusInternalServerError)
		return
	}

	// 4. Respond with 200 OK to acknowledge receipt of the webhook.
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"received"}`))
}