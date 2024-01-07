package main

import (
	"bytes"
	"os"
	"strings"
	"testing"
)

func TestRunChatFunction(t *testing.T) {
	// Mock the environment variable for the API key
	os.Setenv("OPENAI_API_KEY", "test_api_key")

	// Create a reader with the input you want to test
	input := "Hello, world!"
	reader := strings.NewReader(input)

	// Capture the output in a buffer
	var output bytes.Buffer

	// Run the chat function with the mocked input and output
	err := runChat(reader, &output, os.Getenv("OPENAI_API_KEY"))
	if err != nil {
		t.Fatalf("runChat returned an error: %v", err)
	}

	// Check the output
	expectedOutputPart := "Hello, world!" // Adjust this based on the expected output from the API
	actualOutput := output.String()
	if !strings.Contains(actualOutput, expectedOutputPart) {
		t.Errorf("Expected output to contain %q, got %q", expectedOutputPart, actualOutput)
	}
}
