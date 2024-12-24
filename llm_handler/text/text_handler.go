package text

import (
	"context"
	"fmt"
	"os"

	"google.golang.org/genai"
)

// HandleText sends a request to the Gemini API and returns the generated content
func HandleText(inputText string) (string, error) {
	apiKey := os.Getenv("GENERATIVE_LANGUAGE_API_KEY")
	if apiKey == "" {
		return "", fmt.Errorf("API key not set in GENERATIVE_LANGUAGE_API_KEY environment variable")
	}

	ctx := context.Background()
	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey:  apiKey,
		Backend: genai.BackendGoogleAI,
	})
	if err != nil {
		return "", fmt.Errorf("failed to create Gemini client: %v", err)
	}
	defer client.Close()

	parts := []*genai.Part{
		{Text: inputText},
	}

	result, err := client.GenerateContent(ctx, "gemini-1.0-pro", []*genai.Content{{Parts: parts}}, nil)
	if err != nil {
		return "", fmt.Errorf("Gemini API request failed: %v", err)
	}

	if len(result.Candidates) == 0 || len(result.Candidates[0].Content.Parts) == 0 {
		return "", fmt.Errorf("no response generated")
	}

	return result.Candidates[0].Content.Parts[0].Text, nil
}
