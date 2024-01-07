package main

import (
	"io/ioutil"
	"os"
	"strings"
	"testing"
)

func TestRunChat(t *testing.T) {
	// Create a temporary file to simulate stdin
	tempInputFile, err := ioutil.TempFile("", "test_stdin")
	if err != nil {
		t.Fatalf("Failed to create temp file for stdin: %v", err)
	}
	defer os.Remove(tempInputFile.Name()) // Clean up

	// Write test input to temp file
	testInput := "Hello, world!"
	if _, err := tempInputFile.WriteString(testInput); err != nil {
		t.Fatalf("Failed to write to temp stdin file: %v", err)
	}

	// Reset file offset to the beginning
	if _, err := tempInputFile.Seek(0, 0); err != nil {
		t.Fatalf("Failed to seek temp stdin file: %v", err)
	}

	// Backup the real stdin and defer restoration
	oldStdin := os.Stdin
	defer func() { os.Stdin = oldStdin }()
	os.Stdin = tempInputFile

	// Capture the stdout
	oldStdout := os.Stdout
	tempOutputFile, err := ioutil.TempFile("", "test_stdout")
	if err != nil {
		t.Fatalf("Failed to create temp file for stdout: %v", err)
	}
	defer func() {
		os.Stdout = oldStdout
		os.Remove(tempOutputFile.Name()) // Clean up
	}()
	os.Stdout = tempOutputFile

	// Call the function under test
	main()

	// Read the output
	if _, err := tempOutputFile.Seek(0, 0); err != nil {
		t.Fatalf("Failed to seek temp stdout file: %v", err)
	}
	output, err := ioutil.ReadAll(tempOutputFile)
	if err != nil {
		t.Fatalf("Failed to read from temp stdout file: %v", err)
	}

	// Verify the output
	expectedOutputPart := "Hello, world!" // Adjust based on expected output
	if !strings.Contains(string(output), expectedOutputPart) {
		t.Errorf("Expected output to contain %q, got %q", expectedOutputPart, output)
	}
}
