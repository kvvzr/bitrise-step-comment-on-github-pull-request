package github

import (
	"context"
	"errors"
	"strings"

	"github.com/google/go-github/v33/github"
	"golang.org/x/oauth2"
)

var client *github.Client
var ctx = context.Background()

func Initialize(accessToken string) {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: accessToken},
	)
	tc := oauth2.NewClient(ctx, ts)

	client = github.NewClient(tc)
}

func CreateComment(owner, repo string, issueNumber int, body string) (*github.IssueComment, error) {
	commentToCreate := github.IssueComment{Body: &body}
	comment, _, err := client.Issues.CreateComment(ctx, owner, repo, issueNumber, &commentToCreate)

	return comment, err
}

func GetComments(owner, repo string, issueNumber int) ([]*github.IssueComment, error) {
	comments, _, err := client.Issues.ListComments(ctx, owner, repo, issueNumber, &github.IssueListCommentsOptions{})

	return comments, err
}

func GetFirstCommentWithTag(owner, repo string, issueNumber int, tag string) (*github.IssueComment, error) {
	comments, err := GetComments(owner, repo, issueNumber)
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

func UpdateComment(owner, repo string, commentID int64, body string) (*github.IssueComment, error) {
	commentToUpdate := github.IssueComment{Body: &body}
	comment, _, err := client.Issues.EditComment(ctx, owner, repo, int64(commentID), &commentToUpdate)

	return comment, err
}
