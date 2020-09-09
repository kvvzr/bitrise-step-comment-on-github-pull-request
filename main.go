package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/bitrise-io/go-utils/log"
	"github.com/bitrise-tools/go-steputils/stepconf"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

type Config struct {
	AuthToken     stepconf.Secret `env:"personal_access_token,required"`
	Body          string `env:"body,required"`
	RepositoryURL string `env:"repository_url,required"`
	IssueNumber   string `env:"issue_number,required"`
	APIBaseURL    string `env:"api_base_url,required"`
}

type Payload struct {
	Body string `json:"body"`
}

func ownerAndRepo(url string) (string, string) {
	url = strings.TrimPrefix(strings.TrimPrefix(url, "https://"), "git@")
	paths := strings.FieldsFunc(url, func(r rune) bool { return r == '/' || r == ':' })
	return paths[1], strings.TrimSuffix(paths[2], ".git")
}

func main() {
	var conf Config
	if err := stepconf.Parse(&conf); err != nil {
		log.Errorf("Error: %s\n", err)
		os.Exit(1)
	}
	stepconf.Print(conf)

	owner, repo := ownerAndRepo(conf.RepositoryURL)
	url := fmt.Sprintf("%s/repos/%s/%s/issues/%s/comments", conf.APIBaseURL, owner, repo, conf.IssueNumber)

	data := Payload{
		conf.Body,
	}
	payloadBytes, err := json.Marshal(data)
	if err != nil {
		log.Errorf("Error: %s\n", err)
		os.Exit(1)
	}
	body := bytes.NewReader(payloadBytes)

	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		log.Errorf("Error: %s\n", err)
		os.Exit(1)
	}
	req.Header.Set("Authorization", "token " + string(conf.AuthToken))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Errorf("Error: %s\n", err)
		os.Exit(1)
	}
	respBody, _ := ioutil.ReadAll(resp.Body)
	log.Successf("Success: %s\n", respBody)
	defer resp.Body.Close()

	os.Exit(0)
}
