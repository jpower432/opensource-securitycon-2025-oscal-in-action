package server

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

// Assisted by: Cursor Agent

// LokiClient represents a client for querying Loki logs
type LokiClient struct {
	baseURL    string
	httpClient *http.Client
	username   string
	password   string
}

// LokiQueryResponse represents the response from Loki query API
type LokiQueryResponse struct {
	Status string `json:"status"`
	Data   struct {
		ResultType string `json:"resultType"`
		Result     []struct {
			Stream map[string]string `json:"stream"`
			Values [][]string        `json:"values"`
		} `json:"result"`
	} `json:"data"`
}

// LokiLogEntry represents a single log entry from Loki
type LokiLogEntry struct {
	Timestamp string            `json:"timestamp"`
	Message   string            `json:"message"`
	Labels    map[string]string `json:"labels"`
}

// NewLokiClient creates a new Loki client
func NewLokiClient(baseURL string) *LokiClient {
	return &LokiClient{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// NewLokiClientWithAuth creates a new Loki client with authentication
func NewLokiClientWithAuth(baseURL, instanceID, apiKey string) *LokiClient {
	return &LokiClient{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		username: instanceID, // Use instance ID as username
		password: apiKey,     // Use API key as password
	}
}

// QueryLogs queries Loki for logs matching the given policy ID
func (lc *LokiClient) QueryLogs(ctx context.Context, policyID string, limit int) ([]LokiLogEntry, error) {
	// Construct the query URL
	queryURL := fmt.Sprintf("%s/loki/api/v1/query_range", lc.baseURL)

	// Create query parameters
	params := url.Values{}
	params.Set("query", fmt.Sprintf(`{service_name=~".+"} | policy_id="%s"`, policyID))
	params.Set("limit", strconv.Itoa(limit))

	fullURL := fmt.Sprintf("%s?%s", queryURL, params.Encode())

	// Make the HTTP request
	req, err := http.NewRequestWithContext(ctx, "GET", fullURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Add authentication headers (basic auth with instance ID and API key)
	if lc.username != "" && lc.password != "" {
		req.SetBasicAuth(lc.username, lc.password)
	}

	resp, err := lc.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Loki API returned status %d", resp.StatusCode)
	}

	// Parse the response
	var queryResp LokiQueryResponse
	if err := json.NewDecoder(resp.Body).Decode(&queryResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	// Convert to LokiLogEntry format
	var entries []LokiLogEntry
	for _, result := range queryResp.Data.Result {
		for _, value := range result.Values {
			if len(value) >= 2 {
				entry := LokiLogEntry{
					Timestamp: value[0],
					Message:   value[1],
					Labels:    result.Stream,
				}
				entries = append(entries, entry)
			}
		}
	}

	return entries, nil
}
