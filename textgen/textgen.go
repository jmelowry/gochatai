package main

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os"

	openai "github.com/sashabaranov/go-openai"
)

// runChat handles the chat logic, can be tested with different io.Reader and io.Writer
func runChat(input io.Reader, output io.Writer, apiKey string) error {
	client := openai.NewClient(apiKey)

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
	if err := runChat(os.Stdin, os.Stdout, apiKey); err != nil {
		fmt.Fprintf(os.Stderr, "runChat error: %v\n", err)
	}
}
