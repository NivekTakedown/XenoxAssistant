package text

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

// RequestPayload defines the structure of the API request
type RequestPayload struct {
	Contents []Content `json:"contents"`
}

// Content defines the content part of the request
type Content struct {
	Parts []Part `json:"parts"`
}

// Part defines the text part of the content
type Part struct {
	Text string `json:"text"`
}

// ResponsePayload defines the structure of the API response
type ResponsePayload struct {
	GeneratedText string `json:"generated_text"`
}

// HandleText sends a request to the Generative Language API and returns the generated content
func HandleText(inputText string) (string, error) {
	apiKey := os.Getenv("GENERATIVE_LANGUAGE_API_KEY")
	if apiKey == "" {
		return "", fmt.Errorf("API key not set in GENERATIVE_LANGUAGE_API_KEY environment variable")
	}

	url := fmt.Sprintf("https://generativelanguage.googleapis.com/v1beta/models/gemini-1.5-flash:generateContent?key=%s", apiKey)

	payload := RequestPayload{
		Contents: []Content{
			{
				Parts: []Part{
					{Text: inputText},
				},
			},
		},
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(bodyBytes))
	}

	var responsePayload ResponsePayload
	err = json.NewDecoder(resp.Body).Decode(&responsePayload)
	if err != nil {
		return "", err
	}

	return responsePayload.GeneratedText, nil
}
