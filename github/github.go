package github

import (
	"context"
	"errors"
	"net/http"
	"os"
	"strings"

	"github.com/bitrise-io/go-utils/log"
	gogithub "github.com/google/go-github/v33/github"
	"golang.org/x/oauth2"
)

type GithubClient struct {
	Context context.Context
	*gogithub.Client
}

func NewClient(accessToken string) *GithubClient {
	ctx := context.Background()
	tc := NewAuthTokenClient(accessToken)

	return &GithubClient{
		ctx,
		gogithub.NewClient(tc),
	}
}

func NewEnterpriseClient(baseURL string, accessToken string) *GithubClient {
	ctx := context.Background()
	tc := NewAuthTokenClient(accessToken)
	enterpriseClient, err := gogithub.NewEnterpriseClient(baseURL, baseURL, tc)
	if err != nil {
		log.Errorf("Error: %s\n", err)
		os.Exit(1)
	}

	return &GithubClient{
		ctx,
		enterpriseClient,
	}
}

func NewAuthTokenClient(accessToken string) *http.Client {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: accessToken},
	)
	return oauth2.NewClient(ctx, ts)
}

func (c *GithubClient) CreateComment(owner, repo string, issueNumber int, body string) (*gogithub.IssueComment, error) {
	commentToCreate := gogithub.IssueComment{Body: &body}
	comment, _, err := c.Issues.CreateComment(c.Context, owner, repo, issueNumber, &commentToCreate)

	return comment, err
}

func (c *GithubClient) GetComments(owner, repo string, issueNumber int) ([]*gogithub.IssueComment, error) {
	comments, _, err := c.Issues.ListComments(c.Context, owner, repo, issueNumber, &gogithub.IssueListCommentsOptions{})

	return comments, err
}

func (c *GithubClient) GetFirstCommentWithTag(owner, repo string, issueNumber int, tag string) (*gogithub.IssueComment, error) {
	comments, err := c.GetComments(owner, repo, issueNumber)
	if err != nil {
		return nil, err
	}

	for _, comment := range comments {
		if strings.Contains(*comment.Body, tag) {
			return comment, nil
		}
	}

	return nil, errors.New("No comment containing tag found")
}

func (c *GithubClient) UpdateComment(owner, repo string, commentID int64, body string) (*gogithub.IssueComment, error) {
	commentToUpdate := gogithub.IssueComment{Body: &body}
	comment, _, err := c.Issues.EditComment(c.Context, owner, repo, int64(commentID), &commentToUpdate)

	return comment, err
}
