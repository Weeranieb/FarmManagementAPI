package service

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/weeranieb/boonmafarm-backend/src/internal/config"
)

type ClaudeClient interface {
	SendVisionRequest(imagesBase64 []string, mimeTypes []string, prompt string) (*ClaudeResponse, error)
}

type claudeClient struct {
	apiKey    string
	model     string
	maxTokens int
	httpClient *http.Client
}

func NewClaudeClient(conf *config.Config) ClaudeClient {
	return &claudeClient{
		apiKey:     conf.AI.AnthropicAPIKey,
		model:      conf.AI.Model,
		maxTokens:  conf.AI.MaxTokens,
		httpClient: &http.Client{},
	}
}

type ClaudeResponse struct {
	Content []struct {
		Type string `json:"type"`
		Text string `json:"text"`
	} `json:"content"`
}

func (c *claudeClient) SendVisionRequest(imagesBase64 []string, mimeTypes []string, prompt string) (*ClaudeResponse, error) {
	if c.apiKey == "" {
		return nil, fmt.Errorf("anthropic API key is not configured")
	}

	content := make([]map[string]any, 0, len(imagesBase64)+1)
	for i, img := range imagesBase64 {
		mt := "image/jpeg"
		if i < len(mimeTypes) && mimeTypes[i] != "" {
			mt = mimeTypes[i]
		}
		content = append(content, map[string]any{
			"type": "image",
			"source": map[string]any{
				"type":       "base64",
				"media_type": mt,
				"data":       img,
			},
		})
	}
	content = append(content, map[string]any{
		"type": "text",
		"text": prompt,
	})

	body := map[string]any{
		"model":      c.model,
		"max_tokens": c.maxTokens,
		"messages": []map[string]any{
			{
				"role":    "user",
				"content": content,
			},
		},
	}

	jsonBody, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequest("POST", "https://api.anthropic.com/v1/messages", bytes.NewReader(jsonBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", c.apiKey)
	req.Header.Set("anthropic-version", "2023-06-01")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("claude API error (status %d): %s", resp.StatusCode, string(respBody))
	}

	var claudeResp ClaudeResponse
	if err := json.Unmarshal(respBody, &claudeResp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &claudeResp, nil
}

func EncodeImageToBase64(data []byte) string {
	return base64.StdEncoding.EncodeToString(data)
}

func DetectMimeType(filename string) string {
	lower := strings.ToLower(filename)
	switch {
	case strings.HasSuffix(lower, ".png"):
		return "image/png"
	case strings.HasSuffix(lower, ".gif"):
		return "image/gif"
	case strings.HasSuffix(lower, ".webp"):
		return "image/webp"
	default:
		return "image/jpeg"
	}
}
