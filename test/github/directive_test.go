package github

import (
	"fmt"
	"testing"

	govgh "github.com/gitrules/gitrules/github"
	"github.com/gitrules/gitrules/lib/base"
	"github.com/gitrules/gitrules/lib/form"
	"github.com/gitrules/gitrules/lib/testutil"
	"github.com/gitrules/gitrules/proto/account"
	"github.com/gitrules/gitrules/runtime"
	"github.com/gitrules/gitrules/test"
	"github.com/google/go-github/v58/github"
	"github.com/migueleliasweb/go-github-mock/src/mock"
)

var (
	testDirectiveOrganizerGithubUser = "organizer"
	testDirectiveEditIssue           = []any{
		&github.Issue{},
		&github.Issue{},
	}
	testDirectivePostComments = []any{
		&github.IssueComment{},
		&github.IssueComment{},
	}
)

func TestDirective(t *testing.T) {
	base.LogVerbosely()

	// init governance
	ctx := testutil.NewCtx(t, runtime.TestWithCache)
	cty := test.NewTestCommunity(t, ctx, 2)

	testIssueAmount := 20.0
	testTransferAmount := 10.0
	testDirectiveGetIssues := []any{
		[]*github.Issue{
			{
				ID:     github.Int64(111),
				Number: github.Int(1),
				Title:  github.String("Issue directive"),
				URL:    github.String("https://test/issue/1"),
				Labels: []*github.Label{{Name: github.String(govgh.DirectiveLabel)}},
				Locked: github.Bool(false),
				State:  github.String("open"),
				Body: github.String(
					fmt.Sprintf("issue %v credits to @%v", testIssueAmount, cty.MemberUser(0)),
				),
				User: &github.User{Login: github.String(testDirectiveOrganizerGithubUser)},
			},
			{
				ID:     github.Int64(222),
				Number: github.Int(2),
				Title:  github.String("Transfer directive"),
				URL:    github.String("https://test/issue/2"),
				Labels: []*github.Label{{Name: github.String(govgh.DirectiveLabel)}},
				Locked: github.Bool(false),
				State:  github.String("open"),
				Body: github.String(
					fmt.Sprintf("transfer %v credits from @%v to @%v", testTransferAmount, cty.MemberUser(0), cty.MemberUser(1)),
				),
				User: &github.User{Login: github.String(testDirectiveOrganizerGithubUser)},
			},
		},
	}

	// init mock github
	mockedHTTPClient := mock.NewMockedHTTPClient(
		// fetch issues
		mock.WithRequestMatch(mock.GetReposIssuesByOwnerByRepo,
			testDirectiveGetIssues...),
		// issue + transfer directives execution
		mock.WithRequestMatch(mock.PostReposIssuesCommentsByOwnerByRepoByIssueNumber,
			testDirectivePostComments...),
		mock.WithRequestMatch(mock.PatchReposIssuesByOwnerByRepoByIssueNumber,
			testDirectiveEditIssue...),
	)
	ghRepo := govgh.Repo{Owner: "owner1", Name: "repo1"}
	ghClient := github.NewClient(mockedHTTPClient)

	// process directives
	chg := govgh.ProcessDirectiveIssues(ctx, ghRepo, ghClient, cty.Organizer(), []string{testDirectiveOrganizerGithubUser})
	if len(chg.Result) != 2 {
		t.Fatalf("expecting 2 directives")
	}
	fmt.Println(form.SprintJSON(chg.Result))

	c1 := account.Get(ctx, cty.Gov(), cty.MemberAccountID(0)).Balance(account.PluralAsset).Quantity
	if c1 != 10.0 {
		t.Errorf("expecting %v, got %v", 10.0, c1)
	}
	c2 := account.Get(ctx, cty.Gov(), cty.MemberAccountID(1)).Balance(account.PluralAsset).Quantity
	if c2 != 10.0 {
		t.Errorf("expecting %v, got %v", 10.0, c2)
	}

	// <-(chan int(nil))
}
