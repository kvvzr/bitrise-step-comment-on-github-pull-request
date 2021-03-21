package main

import (
	"os"
	"strconv"
	"strings"

	"github.com/bitrise-io/go-utils/log"
	"github.com/bitrise-tools/go-steputils/stepconf"
	"github.com/kvvzr/bitrise-step-comment-on-github-pull-request/github"
)

type Config struct {
	AuthToken        stepconf.Secret `env:"personal_access_token,required"`
	Body             string          `env:"body,required"`
	RepositoryURL    string          `env:"repository_url,required"`
	IssueNumber      string          `env:"issue_number,required"`
	APIBaseURL       string          `env:"api_base_url,required"`
	UpdateCommentTag string          `env:"update_comment_tag"`
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

	github.Initialize(string(conf.AuthToken))

	issueNumber, err := strconv.Atoi(conf.IssueNumber)
	if err != nil {
		log.Errorf("Failed to convert IssueNumber %s to string: %w\n", conf.IssueNumber, err)
	}
	comment, err := github.CreateComment(owner, repo, issueNumber, conf.Body)
	if err != nil {
		log.Errorf("Github API call failed: %w\n", conf.IssueNumber, err)
	} else {
		log.Successf("Success: %v\n", comment)
	}

	os.Exit(0)
}
