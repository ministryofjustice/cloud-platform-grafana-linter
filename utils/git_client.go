package utils

import (
	"context"
	"os"

	"github.com/google/go-github/v62/github"
	"golang.org/x/oauth2"
)

var (
	token = os.Getenv("GITHUB_TOKEN")
)

func GitHubClient() (*github.Client, context.Context) {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)

	return client, ctx
}
