package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

const (
	apiURL      = "https://api.prod.testnet.anoma.net/api/v1/explorer"
	claimURL    = "https://api.prod.testnet.anoma.net/api/v1/explorer/claim"
	bearerToken = "YOUR_TOKEN"
)

type ExplorerResponse struct {
	Data []struct {
		ID string `json:"id"`
	} `json:"data"`
}

func main() {
	ids, err := fetchIDs()
	if err != nil {
		log.Fatalf("Error fetching IDs: %v", err)
	}

	fmt.Printf("Found %d IDs\n", len(ids))
	for _, id := range ids {
		resp, err := claimID(id)
		if err != nil {
			log.Printf("Error claiming ID %s: %v", id, err)
			continue
		}
		fmt.Printf("ID: %s\nResponse: %s\n\n", id, resp)
	}
}

func fetchIDs() ([]string, error) {
	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		return nil, err
	}
	setHeaders(req)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var explorer ExplorerResponse
	body, _ := io.ReadAll(resp.Body)
	if err := json.Unmarshal(body, &explorer); err != nil {
		return nil, err
	}

	var ids []string
	for _, item := range explorer.Data {
		ids = append(ids, item.ID)
	}
	return ids, nil
}

func claimID(id string) (string, error) {
	payload := map[string]string{"id": id}
	bodyBytes, _ := json.Marshal(payload)

	req, err := http.NewRequest("POST", claimURL, bytes.NewBuffer(bodyBytes))
	if err != nil {
		return "", err
	}
	setHeaders(req)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	responseText, _ := io.ReadAll(resp.Body)
	return string(responseText), nil
}

func setHeaders(req *http.Request) {
	req.Header.Set("Authorization", "Bearer "+bearerToken)
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9,en-IN;q=0.8")
	req.Header.Set("Origin", "https://testnet.anoma.net")
	req.Header.Set("Referer", "https://testnet.anoma.net/")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64)")
}
