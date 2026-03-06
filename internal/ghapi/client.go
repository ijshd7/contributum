package ghapi

import (
	"context"
	"os"

	"github.com/google/go-github/v60/github"
	"github.com/joho/godotenv"
	"golang.org/x/oauth2"
)

type Client struct {
	gh *github.Client
}

func NewClient(token string) *Client {
	if token == "" {
		return &Client{gh: github.NewClient(nil)}
	}
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
	tc := oauth2.NewClient(context.Background(), ts)
	return &Client{gh: github.NewClient(tc)}
}

func TokenFromEnv() string {
	_ = godotenv.Load()
	return os.Getenv("GITHUB_TOKEN")
}
