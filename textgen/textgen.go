package main

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os"

	openai "github.com/sashabaranov/go-openai"
)

// OpenAIClient defines the interface for an OpenAI client.
type OpenAIClient interface {
	CreateChatCompletion(ctx context.Context, req openai.ChatCompletionRequest) (*openai.ChatCompletionResponse, error)
}

// runChat handles the chat logic, can be tested with different io.Reader and io.Writer
func runChat(input io.Reader, output io.Writer, client OpenAIClient) error {
	scanner := bufio.NewScanner(input)
	for scanner.Scan() {
		userInput := scanner.Text()

		resp, err := client.CreateChatCompletion(
			context.Background(),
			openai.ChatCompletionRequest{
				Model: openai.GPT3Dot5Turbo,
				Messages: []openai.ChatCompletionMessage{
					{
						Role:    openai.ChatMessageRoleUser,
						Content: userInput,
					},
				},
			},
		)

		if err != nil {
			return fmt.Errorf("API request error: %v", err)
		}

		aiResponse := resp.Choices[0].Message.Content
		if _, err := fmt.Fprintln(output, aiResponse); err != nil {
			return fmt.Errorf("output write error: %v", err)
		}
	}
	return nil
}

func main() {
	apiKey := os.Getenv("OPENAI_API_KEY")
	client := openai.NewClient(apiKey) // Ensure this returns a type that satisfies OpenAIClient
	if err := runChat(os.Stdin, os.Stdout, client); err != nil {
		fmt.Fprintf(os.Stderr, "runChat error: %v\n", err)
	}
}

// MockOpenAIClient is a mock implementation of OpenAIClient for testing.
type MockOpenAIClient struct{}

// CreateChatCompletion simulates the CreateChatCompletion method of an OpenAIClient.
func (m *MockOpenAIClient) CreateChatCompletion(ctx context.Context, req openai.ChatCompletionRequest) (*openai.ChatCompletionResponse, error) {
	// Return a mock response
	return &openai.ChatCompletionResponse{
		Choices: []openai.Choice{
			{
				Message: openai.Message{Content: "mock response"},
			},
		},
	}, nil
}
