package knowledge

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

// EmbeddingProvider generates embeddings for code/text
type EmbeddingProvider interface {
	Embed(text string) ([]float32, error)
}

// OpenAIEmbeddings implements EmbeddingProvider using OpenAI's text-embedding-ada-002 model
type OpenAIEmbeddings struct {
	apiKey     string
	model      string
	httpClient *http.Client
}

// NewOpenAIEmbeddings creates a new OpenAI embeddings provider
// API key is read from OPENAI_API_KEY environment variable
func NewOpenAIEmbeddings() (*OpenAIEmbeddings, error) {
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("OPENAI_API_KEY environment variable not set")
	}

	return &OpenAIEmbeddings{
		apiKey:     apiKey,
		model:      "text-embedding-ada-002",
		httpClient: &http.Client{},
	}, nil
}

// openAIEmbeddingRequest represents the API request format
type openAIEmbeddingRequest struct {
	Input string `json:"input"`
	Model string `json:"model"`
}

// openAIEmbeddingResponse represents the API response format
type openAIEmbeddingResponse struct {
	Data []struct {
		Embedding []float32 `json:"embedding"`
	} `json:"data"`
	Usage struct {
		TotalTokens int `json:"total_tokens"`
	} `json:"usage"`
}

// Embed generates an embedding vector for the given text
func (o *OpenAIEmbeddings) Embed(text string) ([]float32, error) {
	// Create request body
	reqBody := openAIEmbeddingRequest{
		Input: text,
		Model: o.model,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	// Create HTTP request
	req, err := http.NewRequest("POST", "https://api.openai.com/v1/embeddings", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+o.apiKey)

	// Send request
	resp, err := o.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	// Read response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	// Check status code
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("OpenAI API error (status %d): %s", resp.StatusCode, string(body))
	}

	// Parse response
	var apiResp openAIEmbeddingResponse
	if err := json.Unmarshal(body, &apiResp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	if len(apiResp.Data) == 0 {
		return nil, fmt.Errorf("no embeddings returned from API")
	}

	return apiResp.Data[0].Embedding, nil
}

// cosineSimilarity computes the cosine similarity between two vectors
// Returns a value between -1 (opposite) and 1 (identical)
func cosineSimilarity(a, b []float32) float32 {
	if len(a) != len(b) {
		return 0
	}

	var dotProduct, normA, normB float32
	for i := 0; i < len(a); i++ {
		dotProduct += a[i] * b[i]
		normA += a[i] * a[i]
		normB += b[i] * b[i]
	}

	// Avoid division by zero
	if normA == 0 || normB == 0 {
		return 0
	}

	// Compute cosine similarity
	similarity := dotProduct / (sqrt32(normA) * sqrt32(normB))
	return similarity
}

// sqrt32 computes square root of a float32
func sqrt32(x float32) float32 {
	// Newton's method for square root
	if x == 0 {
		return 0
	}

	z := x
	for i := 0; i < 10; i++ {
		z = z - (z*z-x)/(2*z)
	}
	return z
}
