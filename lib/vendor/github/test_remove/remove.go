package main

import (
	"context"
	"fmt"
	"os"

	"github.com/google/go-github/v58/github"
	"golang.org/x/oauth2"
)

func main() {
	ctx := context.Background()

	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: os.Getenv("GITHUB_ACCESS_TOKEN")},
	)
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)

	_, err := client.Repositories.Delete(ctx, "petar", "xxx1")
	ghErr, ok := err.(*github.ErrorResponse)
	if ghErr != nil && ok && ghErr.Response.StatusCode == 404 {
		fmt.Println("repo not found")
		return
	}
	fmt.Printf("err: %v | type: %T\n", err, err)
}
