/**
 * @description
 * This file provides a dedicated HTTP client for interacting with the Anchor BaaS API.
 * It encapsulates all the logic for making authenticated requests to Anchor, handling
 * request body construction, and parsing responses.
 *
 * Key features:
 * - Centralized client for all Anchor API interactions.
 * - Handles authentication by adding the API key to headers.
 * - Provides methods for creating customers and triggering verification.
 * - Uses standard library packages for HTTP communication and JSON handling.
 *
 * @dependencies
 * - "bytes", "context", "encoding/json", "fmt", "io", "log", "net/http", "time"
 * - "transfa/services/customer/internal/domain": For event data structures.
 */
package anchor

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"transfa/services/customer/internal/domain"
)

// Client is a client for interacting with the Anchor API.
type Client struct {
	baseURL    string
	apiKey     string
	httpClient *http.Client
}

// NewClient creates a new Anchor API client.
func NewClient(baseURL, apiKey string) *Client {
	return &Client{
		baseURL: baseURL,
		apiKey:  apiKey,
		httpClient: &http.Client{
			Timeout: 15 * time.Second,
		},
	}
}

// CreateIndividualCustomer sends a request to Anchor to create a new individual customer.
func (c *Client) CreateIndividualCustomer(ctx context.Context, event domain.UserCreatedEvent) (string, error) {
	// Minimal payload, assuming email and phone are not available from the initial event.
	// In a real-world scenario, you might need to fetch these from your DB or Clerk.
	payload := map[string]interface{}{
		"data": map[string]interface{}{
			"type": "IndividualCustomer",
			"attributes": map[string]interface{}{
				"fullName": map[string]string{
					"firstName": event.KYCDetails.FullName, // Assuming full name is just in one field
					"lastName":  "...",                   // Placeholder, Anchor requires this.
				},
				"address": map[string]string{
					"addressLine_1": "...", // Placeholder
					"city":          "...",
					"state":         "Lagos",
					"country":       "NG",
				},
				"email":       fmt.Sprintf("%s@transfa.com", event.UserID.String()), // Placeholder email
				"phoneNumber": "08000000000",                                     // Placeholder phone
			},
		},
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("failed to marshal anchor create customer payload: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", c.baseURL+"/api/v1/customers", bytes.NewBuffer(body))
	if err != nil {
		return "", fmt.Errorf("failed to create anchor request: %w", err)
	}
	c.setHeaders(req)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to call anchor create customer api: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		respBody, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("anchor api returned non-2xx status: %d - %s", resp.StatusCode, string(respBody))
	}

	var anchorResp struct {
		Data struct {
			ID string `json:"id"`
		} `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&anchorResp); err != nil {
		return "", fmt.Errorf("failed to decode anchor create customer response: %w", err)
	}

	if anchorResp.Data.ID == "" {
		return "", fmt.Errorf("anchor customer id not found in response")
	}

	return anchorResp.Data.ID, nil
}

// TriggerIndividualVerification sends a request to Anchor to start the KYC verification process.
func (c *Client) TriggerIndividualVerification(ctx context.Context, anchorCustomerID string, kycDetails *domain.KYCDetails) error {
	payload := map[string]interface{}{
		"data": map[string]interface{}{
			"type": "Verification",
			"attributes": map[string]interface{}{
				"level": "TIER_2", // Using TIER_2 as it corresponds to BVN level verification
				"level2": map[string]string{
					"bvn":         kycDetails.BVN,
					"dateOfBirth": kycDetails.DateOfBirth,
					"gender":      "Other", // Placeholder as gender isn't in spec
				},
			},
		},
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal anchor verification payload: %w", err)
	}

	url := fmt.Sprintf("%s/api/v1/customers/%s/verification/individual", c.baseURL, anchorCustomerID)
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(body))
	if err != nil {
		return fmt.Errorf("failed to create anchor verification request: %w", err)
	}
	c.setHeaders(req)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to call anchor verification api: %w", err)
	}
	defer resp.Body.Close()

	// Anchor API returns an empty JSON object `{}` on success with 200 OK.
	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("anchor verification api returned non-200 status: %d - %s", resp.StatusCode, string(respBody))
	}

	log.Printf("Successfully triggered KYC verification for Anchor customer %s", anchorCustomerID)
	return nil
}

// setHeaders adds the necessary authentication and content-type headers to an HTTP request.
func (c *Client) setHeaders(req *http.Request) {
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("x-anchor-key", c.apiKey)
}

// NOTE: CreateBusinessCustomer and TriggerBusinessVerification would be implemented here
// in a similar fashion, but are omitted as the provided UserCreatedEvent does not
// contain enough information to satisfy the Anchor API's requirements for business onboarding.