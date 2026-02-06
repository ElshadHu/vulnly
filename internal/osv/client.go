package osv

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

const (
	QueryEndpoint = "https://api.osv.dev/v1/querybatch"
	GetEndpoint   = "https://api.osv.dev/v1/vulns"
)

var httpClient = &http.Client{
	Timeout: 30 * time.Second,
}

// MakeRequest sends a batched query to OSV
func MakeRequest(query BatchedQuery) (*BatchedResponse, error) {
	body, err := json.Marshal(query)
	if err != nil {
		return nil, fmt.Errorf("marshal query %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, QueryEndpoint, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("osv request %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("osv returned status %d", resp.StatusCode)
	}
	var result BatchedResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}
	return &result, nil
}

// Get fetches full vulnerability details by ID
func Get(id string) (*Vulnerability, error) {
	url := fmt.Sprintf("%s/%s", GetEndpoint, id)
	resp, err := httpClient.Get(url)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode == http.StatusNotFound {
		return nil, nil
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("osv vulns returned status %d", resp.StatusCode)
	}
	var vuln Vulnerability
	if err := json.NewDecoder(resp.Body).Decode(&vuln); err != nil {
		return nil, err
	}
	return &vuln, nil
}
