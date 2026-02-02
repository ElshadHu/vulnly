package upload

// Package upload provides functionality to upload scan results to the API

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// Client handles HTTP requests to Vulnly API
type Client struct {
	baseURL    string
	token      string
	httpClient *http.Client
}

// NewClient creates a new API client
func NewClient(baseURL, token string) *Client {
	return &Client{
		baseURL: baseURL,
		token:   token,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// IngestRequest is the payload for POST /api/ingest
type IngestRequest struct {
	Project         string          `json:"project"`
	Commit          string          `json:"commit,omitempty"`
	Branch          string          `json:"branch,omitempty"`
	Ecosystem       string          `json:"ecosystem"`
	Dependencies    []Dependency    `json:"dependencies"`
	Vulnerabilities []Vulnerability `json:"vulnerabilities"`
	Summary         VulnSummary     `json:"summary"`
}
type Dependency struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}
type Vulnerability struct {
	ID           string `json:"id"`
	Package      string `json:"package"`
	Version      string `json:"version"`
	Severity     string `json:"severity"`
	FixedVersion string `json:"fixed_version"`
	Description  string `json:"description"`
}
type VulnSummary struct {
	Total    int `json:"total"`
	Critical int `json:"critical"`
	High     int `json:"high"`
	Medium   int `json:"medium"`
	Low      int `json:"low"`
}

// IngestResponse is returned from POST /api/ingest
type IngestResponse struct {
	ScanID    string    `json:"scan_id"`
	CreatedAt time.Time `json:"created_at"`
}

// Ingest uploads scan results to the backend
func (c *Client) Ingest(req *IngestRequest) (*IngestResponse, error) {
	body, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}
	httpReq, err := http.NewRequest(
		http.MethodPost,
		c.baseURL+"/api/ingest",
		bytes.NewReader(body),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+c.token)
	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("unexpected status: %d", resp.StatusCode)
	}
	var result IngestResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}
	return &result, nil
}
