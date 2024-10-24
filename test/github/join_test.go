package github

import (
	"fmt"
	"testing"

	govgh "github.com/gitrules/gitrules/github"
	"github.com/gitrules/gitrules/lib/base"
	"github.com/gitrules/gitrules/lib/git"
	"github.com/gitrules/gitrules/lib/testutil"
	"github.com/gitrules/gitrules/proto/id"
	"github.com/gitrules/gitrules/runtime"
	"github.com/gitrules/gitrules/test"
	"github.com/google/go-github/v58/github"
	"github.com/migueleliasweb/go-github-mock/src/mock"
)

var (
	testProcessJoinRequestsOrganizerGithubUser = "organizer"
	testProcessJoinRequestsApplicantGithubUser = "applicant"
	testProcessJoinRequestsGetComments         = []any{
		[]*github.IssueComment{
			{
				User: &github.User{Login: github.String(testProcessJoinRequestsOrganizerGithubUser)},
				Body: github.String("Approve."),
			},
		},
	}
	testProcessJoinRequestsCreateComments = []any{
		&github.IssueComment{
			User: &github.User{Login: github.String(testProcessJoinRequestsOrganizerGithubUser)},
			Body: github.String("Approve."),
		},
	}
	testProcessJoinRequestsEditIssue = []any{
		&github.Issue{},
	}
)

func TestProcessJoinRequests(t *testing.T) {
	base.LogVerbosely()

	// init governance
	ctx := testutil.NewCtx(t, runtime.TestWithCache)
	cty := test.NewTestCommunity(t, ctx, 2)

	// init join applicant's identity
	applicantID := id.NewTestID(ctx, t, git.MainBranch, true)
	id.Init(ctx, applicantID.OwnerAddress())

	testProcessJoinRequestsGetIssues := []any{
		[]*github.Issue{
			{ // issue without governance
				ID:     github.Int64(111),
				Number: github.Int(1),
				Title:  github.String("Issue 1"),
				URL:    github.String("https://test/issue/1"),
				Labels: nil,
				Locked: github.Bool(false),
				State:  github.String("open"),
				Body: github.String(
					fmt.Sprintf("### Your public repo\n\n%v\n\n### Your public branch\n\n%v\n\n### Your email (optional)\n\n%v",
						applicantID.Public.Dir(), git.MainBranch, "test@test"),
				),
				User:     &github.User{Login: github.String(testProcessJoinRequestsApplicantGithubUser)},
				Comments: github.Int(1),
			},
		},
	}

	// init mock github
	mockedHTTPClient := mock.NewMockedHTTPClient(
		mock.WithRequestMatch(mock.GetReposIssuesByOwnerByRepo, testProcessJoinRequestsGetIssues...),
		mock.WithRequestMatch(mock.GetReposIssuesCommentsByOwnerByRepoByIssueNumber, testProcessJoinRequestsGetComments...),
		mock.WithRequestMatch(mock.PostReposIssuesCommentsByOwnerByRepoByIssueNumber, testProcessJoinRequestsCreateComments...),
		mock.WithRequestMatch(mock.PatchReposIssuesByOwnerByRepoByIssueNumber, testProcessJoinRequestsEditIssue...),
	)
	ghRepo := govgh.Repo{Owner: "owner1", Name: "repo1"}
	ghClient := github.NewClient(mockedHTTPClient)

	// process join requests
	chg := govgh.ProcessJoinRequestIssues(
		ctx,
		ghRepo,
		ghClient,
		cty.Organizer(),
		[]string{testProcessJoinRequestsOrganizerGithubUser},
		true,
	)
	if len(chg.Result.Joined) != 1 {
		t.Fatalf("expecting 1 join")
	}
	if chg.Result.Joined[0] != testProcessJoinRequestsApplicantGithubUser {
		t.Errorf("expecting %v, got %v", testProcessJoinRequestsApplicantGithubUser, chg.Result.Joined[0])
	}

	// <-(chan int(nil))
}
