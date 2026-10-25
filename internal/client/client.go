package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type Client struct {
	BaseURL    string
	APIKey     string
	HTTPClient *http.Client
}

func New(baseURL, apiKey string) *Client {
	return &Client{
		BaseURL:    baseURL,
		APIKey:     apiKey,
		HTTPClient: &http.Client{Timeout: 30 * time.Second},
	}
}

type Post struct {
	ID     int    `json:"id"`
	Title  string `json:"title"`
	Body   string `json:"body"`
	UserID int    `json:"userId"`
}

// CreatePost calls POST /posts
func (c *Client) CreatePost(ctx context.Context, title, body string) (*Post, error) {
	payload := map[string]interface{}{
		"title":  title,
		"body":   body,
		"userId": 1,
	}
	buf, _ := json.Marshal(payload)

	req, err := http.NewRequestWithContext(ctx, "POST", c.BaseURL+"/posts", bytes.NewReader(buf))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	if c.APIKey != "" {
		req.Header.Set("Authorization", "Bearer "+c.APIKey)
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		return nil, fmt.Errorf("api error: %s", resp.Status)
	}

	var p Post
	if err := json.NewDecoder(resp.Body).Decode(&p); err != nil {
		return nil, err
	}
	return &p, nil
}

// GetPost calls GET /posts/{id}
func (c *Client) GetPost(ctx context.Context, id int) (*Post, error) {
	req, _ := http.NewRequestWithContext(ctx, "GET", fmt.Sprintf("%s/posts/%d", c.BaseURL, id), nil)
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, nil // signal: gone
	}
	if resp.StatusCode >= 300 {
		return nil, fmt.Errorf("api error: %s", resp.Status)
	}

	var p Post
	json.NewDecoder(resp.Body).Decode(&p)
	return &p, nil
}

// DeletePost calls DELETE /posts/{id}
func (c *Client) DeletePost(ctx context.Context, id int) error {
	req, _ := http.NewRequestWithContext(ctx, "DELETE", fmt.Sprintf("%s/posts/%d", c.BaseURL, id), nil)
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 && resp.StatusCode != http.StatusNotFound {
		return fmt.Errorf("api error: %s", resp.Status)
	}
	return nil
}
