package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

// Gemini native API types
type geminiRequest struct {
	Contents          []gemContent     `json:"contents"`
	GenerationConfig  generationConfig `json:"generationConfig"`
	SystemInstruction *gemContent      `json:"systemInstruction,omitempty"`
}

type gemContent struct {
	Parts []gemPart `json:"parts"`
	Role  string    `json:"role,omitempty"`
}

type gemPart struct {
	Text string `json:"text"`
}

type generationConfig struct {
	ResponseMimeType string  `json:"responseMimeType"`
	Temperature      float64 `json:"temperature"`
	MaxOutputTokens  int     `json:"maxOutputTokens"`
}

type geminiResponse struct {
	Candidates []struct {
		Content struct {
			Parts []struct {
				Text string `json:"text"`
			} `json:"parts"`
		} `json:"content"`
	} `json:"candidates"`
	Error *struct {
		Message string `json:"message"`
	} `json:"error"`
}

// OpenRouter / OpenAI-compatible types
type openaiRequest struct {
	Model       string          `json:"model"`
	Messages    []openaiMessage `json:"messages"`
	Temperature float64         `json:"temperature"`
	MaxTokens   int             `json:"max_tokens"`
	ResponseFormat *responseFormat `json:"response_format,omitempty"`
}

type openaiMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type responseFormat struct {
	Type string `json:"type"`
}

type openaiResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
	Error *struct {
		Message string `json:"message"`
	} `json:"error"`
}

func callLLM(apiKey, model, input, instruction string, truncated bool) (string, error) {
	prompt := buildPrompt(input, instruction, truncated)

	// Detect backend from env
	orKey := getOpenRouterKey()

	if apiKey != "" {
		// Gemini native
		return callGemini(apiKey, model, prompt)
	} else if orKey != "" {
		// OpenRouter
		orModel := model
		if orModel == "gemini-2.0-flash" {
			orModel = "google/gemini-2.0-flash-001"
		}
		return callOpenRouter(orKey, orModel, prompt)
	}

	return "", fmt.Errorf("no API key found")
}

func getOpenRouterKey() string {
	return os.Getenv("OPENROUTER_API_KEY")
}

func callGemini(apiKey, model, prompt string) (string, error) {
	reqBody := geminiRequest{
		Contents: []gemContent{
			{Parts: []gemPart{{Text: prompt}}, Role: "user"},
		},
		SystemInstruction: &gemContent{
			Parts: []gemPart{{Text: systemPrompt}},
		},
		GenerationConfig: generationConfig{
			ResponseMimeType: "application/json",
			Temperature:      0.1,
			MaxOutputTokens:  8192,
		},
	}

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("marshaling request: %w", err)
	}

	url := fmt.Sprintf(
		"https://generativelanguage.googleapis.com/v1beta/models/%s:generateContent?key=%s",
		model, apiKey,
	)

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Post(url, "application/json", bytes.NewReader(jsonBody))
	if err != nil {
		return "", fmt.Errorf("API request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("reading response: %w", err)
	}

	if resp.StatusCode != 200 {
		return "", fmt.Errorf("API error (HTTP %d): %s", resp.StatusCode, string(body))
	}

	var gemResp geminiResponse
	if err := json.Unmarshal(body, &gemResp); err != nil {
		return "", fmt.Errorf("parsing response: %w", err)
	}

	if gemResp.Error != nil {
		return "", fmt.Errorf("API error: %s", gemResp.Error.Message)
	}

	if len(gemResp.Candidates) == 0 || len(gemResp.Candidates[0].Content.Parts) == 0 {
		return "", fmt.Errorf("empty response")
	}

	return gemResp.Candidates[0].Content.Parts[0].Text, nil
}

func callOpenRouter(apiKey, model, prompt string) (string, error) {
	reqBody := openaiRequest{
		Model: model,
		Messages: []openaiMessage{
			{Role: "system", Content: systemPrompt},
			{Role: "user", Content: prompt},
		},
		Temperature: 0.1,
		MaxTokens:   8192,
		ResponseFormat: &responseFormat{Type: "json_object"},
	}

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("marshaling request: %w", err)
	}

	req, err := http.NewRequest("POST", "https://openrouter.ai/api/v1/chat/completions", bytes.NewReader(jsonBody))
	if err != nil {
		return "", fmt.Errorf("creating request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("API request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("reading response: %w", err)
	}

	if resp.StatusCode != 200 {
		return "", fmt.Errorf("API error (HTTP %d): %s", resp.StatusCode, string(body))
	}

	var orResp openaiResponse
	if err := json.Unmarshal(body, &orResp); err != nil {
		return "", fmt.Errorf("parsing response: %w", err)
	}

	if orResp.Error != nil {
		return "", fmt.Errorf("API error: %s", orResp.Error.Message)
	}

	if len(orResp.Choices) == 0 {
		return "", fmt.Errorf("empty response")
	}

	result := orResp.Choices[0].Message.Content

	// Validate JSON
	if !json.Valid([]byte(result)) {
		return "", fmt.Errorf("LLM returned invalid JSON")
	}

	return result, nil
}
