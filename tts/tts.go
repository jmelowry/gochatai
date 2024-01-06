package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

// Define command-line arguments
var (
	silent bool
	noSave bool
	voice  string
	apiKey string
	url    string = "https://api.openai.com/v1/audio/speech" // Default API endpoint

)

func init() {
	flag.BoolVar(&silent, "silent", false, "Run in silent mode with no terminal output except errors.")
	flag.BoolVar(&noSave, "no-save", false, "Do not save the output file.")
	flag.StringVar(&voice, "voice", "alloy", "Specify the voice to use.") // Set default to 'alloy'
	apiKey = os.Getenv("OPENAI_API_KEY")
}

func main() {
	flag.Parse()

	// Validate API key
	if apiKey == "" {
		fmt.Println("OPENAI_API_KEY environment variable not set.")
		os.Exit(1)
	}

	// Validate voice
	validVoices := map[string]bool{
		"nova": true, "shimmer": true, "echo": true, "onyx": true, "fable": true, "alloy": true,
	}
	if !validVoices[voice] {
		fmt.Println("Invalid voice option. Valid options are 'nova', 'shimmer', 'echo', 'onyx', 'fable', 'alloy'.")
		os.Exit(1)
	}

	// Read input from stdin
	reader := bufio.NewReader(os.Stdin)
	input, err := ioutil.ReadAll(reader)
	if err != nil {
		fmt.Println("Error reading input:", err)
		os.Exit(1)
	}

	// Making the API request and getting the binary response
	responseBytes, err := makeAPIRequest(string(input))
	if err != nil {
		fmt.Println("Error making API request:", err)
		os.Exit(1)
	}

	// Save the binary data to an audio file
	if !noSave {
		outputFileName := "output.mp3" // or any other format you expect
		err := ioutil.WriteFile(outputFileName, responseBytes, 0644)
		if err != nil {
			fmt.Println("Error writing audio file:", err)
			os.Exit(1)
		}
		fmt.Println("Audio file saved as:", outputFileName)
	}
}

// OpenAI API request
func makeAPIRequest(inputText string) ([]byte, error) {
	// Prepare the request body
	requestBody, err := json.Marshal(map[string]string{
		"model": "tts-1",   // You can choose between tts-1 and tts-1-hd
		"voice": voice,     // The voice to use, e.g., "alloy", "echo", etc.
		"input": inputText, // The text to be converted to speech
	})
	if err != nil {
		return nil, err
	}

	// Create the HTTP request
	url := "https://api.openai.com/v1/audio/speech" // Updated API endpoint
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, err
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	// Send the request
	client := &http.Client{}

	// Send the request and get the response
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Read the response as binary data
	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return responseBody, nil
}
