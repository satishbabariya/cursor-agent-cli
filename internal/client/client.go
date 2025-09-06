package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const (
	BaseURL = "https://api.cursor.com/v0"
)

// Client represents the Cursor API client
type Client struct {
	BaseURL    string
	HTTPClient *http.Client
	APIKey     string
}

// NewClient creates a new Cursor API client
func NewClient(apiKey string) *Client {
	return &Client{
		BaseURL:    BaseURL,
		HTTPClient: &http.Client{Timeout: 30 * time.Second},
		APIKey:     apiKey,
	}
}

// Agent represents a background agent
type Agent struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Status    string    `json:"status"`
	Source    Source    `json:"source"`
	Target    Target    `json:"target"`
	Summary   string    `json:"summary"`
	CreatedAt time.Time `json:"createdAt"`
}

// Source represents the source repository information
type Source struct {
	Repository string `json:"repository"`
	Ref        string `json:"ref"`
}

// Target represents the target branch and PR information
type Target struct {
	BranchName   string `json:"branchName"`
	URL          string `json:"url"`
	PrURL        string `json:"prUrl"`
	AutoCreatePr bool   `json:"autoCreatePr"`
}

// ListAgentsResponse represents the response from the list agents API
type ListAgentsResponse struct {
	Agents     []Agent `json:"agents"`
	NextCursor string  `json:"nextCursor,omitempty"`
}

// Message represents a conversation message
type Message struct {
	ID   string `json:"id"`
	Type string `json:"type"`
	Text string `json:"text"`
}

// ConversationResponse represents the response from the conversation API
type ConversationResponse struct {
	ID       string    `json:"id"`
	Messages []Message `json:"messages"`
}

// APIKeyInfo represents the API key information
type APIKeyInfo struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"createdAt"`
	UserEmail string    `json:"userEmail"`
}

// FollowupRequest represents a follow-up request
type FollowupRequest struct {
	Prompt Prompt `json:"prompt"`
}

// Prompt represents a prompt with text and optional images
type Prompt struct {
	Text   string  `json:"text"`
	Images []Image `json:"images,omitempty"`
}

// Image represents an image with data and dimensions
type Image struct {
	Data      string    `json:"data"`
	Dimension Dimension `json:"dimension"`
}

// Dimension represents image dimensions
type Dimension struct {
	Width  int `json:"width"`
	Height int `json:"height"`
}

// FollowupResponse represents the response from adding a follow-up
type FollowupResponse struct {
	ID string `json:"id"`
}

// makeRequest makes an HTTP request to the API
func (c *Client) makeRequest(method, endpoint string, body interface{}) (*http.Response, error) {
	var reqBody io.Reader
	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("error marshaling request body: %w", err)
		}
		reqBody = bytes.NewBuffer(jsonBody)
	}

	req, err := http.NewRequest(method, c.BaseURL+endpoint, reqBody)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+c.APIKey)
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making request: %w", err)
	}

	return resp, nil
}

// ListAgents retrieves a paginated list of all background agents
func (c *Client) ListAgents(limit int, cursor string) (*ListAgentsResponse, error) {
	endpoint := "/agents"

	// Add query parameters
	if limit > 0 || cursor != "" {
		endpoint += "?"
		if limit > 0 {
			endpoint += fmt.Sprintf("limit=%d", limit)
		}
		if cursor != "" {
			if limit > 0 {
				endpoint += "&"
			}
			endpoint += fmt.Sprintf("cursor=%s", cursor)
		}
	}

	resp, err := c.makeRequest("GET", endpoint, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	var result ListAgentsResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	return &result, nil
}

// GetAgentStatus gets the current status and results of a specific background agent
func (c *Client) GetAgentStatus(agentID string) (*Agent, error) {
	endpoint := fmt.Sprintf("/agents/%s", agentID)

	resp, err := c.makeRequest("GET", endpoint, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	var result Agent
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	return &result, nil
}

// GetAgentConversation retrieves the conversation history of a background agent
func (c *Client) GetAgentConversation(agentID string) (*ConversationResponse, error) {
	endpoint := fmt.Sprintf("/agents/%s/conversation", agentID)

	resp, err := c.makeRequest("GET", endpoint, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	var result ConversationResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	return &result, nil
}

// AddFollowup sends an additional instruction to a running background agent
func (c *Client) AddFollowup(agentID string, prompt string) (*FollowupResponse, error) {
	endpoint := fmt.Sprintf("/agents/%s/followup", agentID)

	request := FollowupRequest{
		Prompt: Prompt{
			Text: prompt,
		},
	}

	resp, err := c.makeRequest("POST", endpoint, request)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	var result FollowupResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	return &result, nil
}

// GetAPIKeyInfo retrieves information about the current API key
func (c *Client) GetAPIKeyInfo() (*APIKeyInfo, error) {
	endpoint := "/me"

	resp, err := c.makeRequest("GET", endpoint, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	var result APIKeyInfo
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	return &result, nil
}
