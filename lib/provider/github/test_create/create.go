package main

import (
	"context"
	"fmt"
	"os"

	"github.com/gitrules/gitrules/lib/must"
	"github.com/google/go-github/v66/github"
	"golang.org/x/oauth2"
)

func main() {
	ctx := context.Background()

	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: os.Getenv("GITHUB_ACCESS_TOKEN")},
	)
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)

	repo := &github.Repository{
		Name:    github.String("xxx2"),
		Private: github.Bool(true),
	}
	repo, _, err := client.Repositories.Create(ctx, "", repo)
	fmt.Printf("err=%v | type=%T", err, err)
	must.NoError(ctx, err)

	fmt.Printf("HTMLURL: %v\n", repo.GetHTMLURL())
	fmt.Printf("CloneURL: %v\n", repo.GetCloneURL())
	fmt.Printf("MirrorURL: %v\n", repo.GetMirrorURL())
	fmt.Printf("GitURL: %v\n", repo.GetGitURL())
	fmt.Printf("SSHURL: %v\n", repo.GetSSHURL())
}

// output:
//
// HTMLURL: https://github.com/petar/xxx2
// CloneURL: https://github.com/petar/xxx2.git
// MirrorURL:
// GitURL: git://github.com/petar/xxx2.git
// SSHURL: git@github.com:petar/xxx2.git
