package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
)

type Payload struct {
	Body string `json:"body"`
}

func ownerAndRepo(url string) (string, string) {
	url = strings.TrimPrefix(strings.TrimPrefix(url, "https://"), "git@")
	paths := strings.FieldsFunc(url, func(r rune) bool { return r == '/' || r == ':' })
	return paths[1], strings.TrimSuffix(paths[2], ".git")
}

func main() {
	authToken := os.Getenv("personal_access_token")
	message := os.Getenv("body")
	owner, repo := ownerAndRepo(os.Getenv("repository_url"))
	issueNumber := os.Getenv("issue_number")
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/issues/%s/comments", owner, repo, issueNumber)

	data := Payload{
		message,
	}
	payloadBytes, err := json.Marshal(data)
	if err != nil {
		os.Exit(1)
	}
	body := bytes.NewReader(payloadBytes)

	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		os.Exit(1)
	}
	req.SetBasicAuth(authToken, "x-oauth-basic")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		os.Exit(1)
	}
	defer resp.Body.Close()

	os.Exit(0)
}
