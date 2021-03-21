package github

import (
	"context"

	"github.com/google/go-github/v33/github"
	"golang.org/x/oauth2"
)

var githubClient *github.Client
var ctx = context.Background()

func Initialize(accessToken string) {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: accessToken},
	)
	tc := oauth2.NewClient(ctx, ts)

	githubClient = github.NewClient(tc)
}

func CreateComment(owner, repo string, issueNumber int, body string) (*github.IssueComment, error) {
	issueComment := github.IssueComment{Body: &body}
	comment, _, err := githubClient.Issues.CreateComment(ctx, owner, repo, issueNumber, &issueComment)

	return comment, err
}
