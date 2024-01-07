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
	"time"

	"github.com/ebitengine/oto/v3"
	"github.com/hajimehoshi/go-mp3"
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

// audioPlayer plays audio from a byte slice.
func audioPlayer(audioData []byte) {
	// Convert the audio data into a reader
	audioDataReader := bytes.NewReader(audioData)

	// Decode the audio data
	decodedMp3, err := mp3.NewDecoder(audioDataReader)
	if err != nil {
		panic("mp3.NewDecoder failed: " + err.Error())
	}

	// Prepare an Oto context (this will use your default audio device) that will
	// play all our sounds. Its configuration can't be changed later.

	op := &oto.NewContextOptions{}

	// Usually 44100 or 48000. Other values might cause distortions in Oto
	op.SampleRate = 44100

	// Number of channels (aka locations) to play sounds from. Either 1 or 2.
	// 1 is mono sound, and 2 is stereo (most speakers are stereo).
	op.ChannelCount = 2

	// Format of the source. go-mp3's format is signed 16bit integers.
	op.Format = oto.FormatSignedInt16LE

	// Remember that you should **not** create more than one context
	otoCtx, readyChan, err := oto.NewContext(op)
	if err != nil {
		panic("oto.NewContext failed: " + err.Error())
	}
	// It might take a bit for the hardware audio devices to be ready, so we wait on the channel.
	<-readyChan

	// Create a new 'player' that will handle our sound. Paused by default.
	player := otoCtx.NewPlayer(decodedMp3)

	// Play starts playing the sound and returns without waiting for it (Play() is async).
	player.Play()

	// We can wait for the sound to finish playing using something like this
	for player.IsPlaying() {
		time.Sleep(time.Millisecond)
	}

	// Now that the sound finished playing, we can restart from the beginning (or go to any location in the sound) using seek
	// newPos, err := player.(io.Seeker).Seek(0, io.SeekStart)
	// if err != nil{
	//     panic("player.Seek failed: " + err.Error())
	// }
	// println("Player is now at position:", newPos)
	// player.Play()

	// If you don't want the player/sound anymore simply close
	err = player.Close()
	if err != nil {
		panic("player.Close failed: " + err.Error())
	}
}
