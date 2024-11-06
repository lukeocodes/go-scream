package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"math/rand/v2"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
)

// AuthResponse represents the authentication response from Bluesky
type AuthResponse struct {
	AccessJwt string `json:"accessJwt"`
	Did       string `json:"did"`
}

// ErrorResponse represents the error response structure from Bluesky
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}

const (
	authURL = "https://bsky.social/xrpc/com.atproto.server.createSession"
	postURL = "https://bsky.social/xrpc/com.atproto.repo.createRecord"
)

func main() {
	// Only load .env file in development environment
	if os.Getenv("ENVIRONMENT") != "production" {
		err := godotenv.Load()
		if err != nil {
			log.Fatal("Error loading .env file")
		}
	}

	// Load environment variables
	username := os.Getenv("BLUESKY_USERNAME")
	if username == "" {
		log.Fatal("BLUESKY_USERNAME environment variable not set")
	}

	password := os.Getenv("BLUESKY_PASSWORD")
	if password == "" {
		log.Fatal("BLUESKY_PASSWORD environment variable not set")
	}

	// Authenticate and obtain access token
	authResponse, err := authenticate(username, password)
	if err != nil {
		log.Fatalf("Authentication failed: %v", err)
	}

	scream := getScream()

	// Post message using access token
	err = postMessage(authResponse.AccessJwt, authResponse.Did, scream)
	if err != nil {
		log.Fatalf("Failed to post message: %v", err)
	}

	fmt.Println("Message posted successfully!")
}

func authenticate(identifier string, password string) (*AuthResponse, error) {
	authBody := map[string]string{
		"identifier": identifier,
		"password":   password,
	}
	bodyBytes, err := json.Marshal(authBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal auth request body: %w", err)
	}

	req, err := http.NewRequest("POST", authURL, bytes.NewBuffer(bodyBytes))
	if err != nil {
		return nil, fmt.Errorf("failed to create auth request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("auth request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		var authResponse AuthResponse
		if err := json.NewDecoder(resp.Body).Decode(&authResponse); err != nil {
			return nil, fmt.Errorf("failed to decode auth response: %w", err)
		}

		fmt.Println("Authentication successful!")
		return &authResponse, nil
	}

	var errResponse ErrorResponse
	if err := json.NewDecoder(resp.Body).Decode(&errResponse); err != nil {
		return nil, fmt.Errorf("failed to decode error response: %w", err)
	}
	return nil, fmt.Errorf("auth error (%d): %s - %s", resp.StatusCode, errResponse.Error, errResponse.Message)
}

func postMessage(accessToken, did, message string) error {
	postBody := map[string]interface{}{
		"repo":       did,
		"collection": "app.bsky.feed.post",
		"record": map[string]interface{}{
			"$type":     "app.bsky.feed.post",
			"text":      message,
			"createdAt": time.Now().UTC().Format(time.RFC3339),
		},
	}
	bodyBytes, err := json.Marshal(postBody)
	if err != nil {
		return fmt.Errorf("failed to marshal post request body: %w", err)
	}

	req, err := http.NewRequest("POST", postURL, bytes.NewBuffer(bodyBytes))
	if err != nil {
		return fmt.Errorf("failed to create post request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("post request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		fmt.Println("Post successful!")
		return nil
	}

	var errResponse ErrorResponse
	if err := json.NewDecoder(resp.Body).Decode(&errResponse); err != nil {
		return fmt.Errorf("failed to decode error response: %w", err)
	}
	return fmt.Errorf("post error (%d): %s - %s", resp.StatusCode, errResponse.Error, errResponse.Message)
}

func getScream() string {
	numAs := 1 + rand.IntN(20)  // Random number of A's (1-20)
	numHs := 1 + rand.IntN(100) // Random number of H's (1-100)

	scream := "A"
	for i := 1; i < numAs; i++ {
		scream += "A"
	}
	for i := 0; i < numHs; i++ {
		scream += "H"
	}

	fmt.Println("Scream:", scream)

	return scream
}
