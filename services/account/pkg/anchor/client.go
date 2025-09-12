/**
 * @description
 * This file provides an API client for interacting with the Anchor BaaS.
 * It encapsulates the details of making HTTP requests to Anchor's endpoints,
 * handling authentication, and parsing responses.
 *
 * @dependencies
 * - Go standard library packages for handling HTTP, JSON, and contexts.
 * - Internal domain models for structuring API responses.
 */
package anchor

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// Client is a client for interacting with the Anchor API.
type Client struct {
	BaseURL    string
	APIKey     string
	HTTPClient *http.Client
}

// Defines the JSON:API structure for creating a deposit account request.
type createAccountRequest struct {
	Data struct {
		Type       string `json:"type"`
		Attributes struct {
			ProductName string `json:"productName"`
		} `json:"attributes"`
		Relationships struct {
			Customer struct {
				Data struct {
					ID   string `json:"id"`
					Type string `json:"type"`
				} `json:"data"`
			} `json:"customer"`
		} `json:"relationships"`
	} `json:"data"`
}

// Defines the structure for the successful response from Anchor.
type createAccountResponse struct {
	Data struct {
		ID string `json:"id"`
	} `json:"data"`
}

// NewClient creates a new Anchor API client.
func NewClient(baseURL, apiKey string) *Client {
	return &Client{
		BaseURL: baseURL,
		APIKey:  apiKey,
		HTTPClient: &http.Client{
			Timeout: time.Minute,
		},
	}
}

// CreateDepositAccount creates a new DepositAccount in Anchor for a given customer.
func (c *Client) CreateDepositAccount(ctx context.Context, anchorCustomerID, customerType, productName string) (string, error) {
	// Construct the request payload based on Anchor's JSON:API spec.
	reqPayload := createAccountRequest{}
	reqPayload.Data.Type = "DepositAccount"
	reqPayload.Data.Attributes.ProductName = productName
	reqPayload.Data.Relationships.Customer.Data.ID = anchorCustomerID
	reqPayload.Data.Relationships.Customer.Data.Type = customerType

	reqBody, err := json.Marshal(reqPayload)
	if err != nil {
		return "", fmt.Errorf("failed to marshal create account request: %w", err)
	}

	// Create the HTTP request.
	url := fmt.Sprintf("%s/api/v1/accounts", c.BaseURL)
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(reqBody))
	if err != nil {
		return "", fmt.Errorf("failed to create new http request: %w", err)
	}

	// Set required headers.
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-anchor-key", c.APIKey)

	// Execute the request.
	res, err := c.HTTPClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to execute request to anchor: %w", err)
	}
	defer res.Body.Close()

	// Check for a successful status code.
	if res.StatusCode != http.StatusCreated && res.StatusCode != http.StatusOK { // Anchor might return 200 or 201
		bodyBytes, _ := io.ReadAll(res.Body)
		return "", fmt.Errorf("anchor API returned non-success status: %s, body: %s", res.Status, string(bodyBytes))
	}

	// Decode the successful response.
	var successRes createAccountResponse
	if err := json.NewDecoder(res.Body).Decode(&successRes); err != nil {
		return "", fmt.Errorf("failed to decode successful response from anchor: %w", err)
	}

	if successRes.Data.ID == "" {
		return "", fmt.Errorf("anchor response did not contain an account ID")
	}

	return successRes.Data.ID, nil
}