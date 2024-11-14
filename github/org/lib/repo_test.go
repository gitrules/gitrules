package lib

import (
	"testing"

	"github.com/gitrules/gitrules/github/common"
)

func TestParseGithubRepoHTTPSURL(t *testing.T) {
	repo, err := common.ParseGithubRepoHTTPSURL("https://github.com/abc/xyz.git")
	if err != nil {
		t.Error(err)
	}
	if repo.Owner != "abc" {
		t.Errorf("expecting %v, got %v", "abc", repo.Owner)
	}
	if repo.Name != "xyz" {
		t.Errorf("expecting %v, got %v", "xyz", repo.Name)
	}
}

func TestParseGithubRepoSSHURL(t *testing.T) {
	repo, err := common.ParseGithubRepoSSHURL("git@github.com:abc/x.y.z.git")
	if err != nil {
		t.Error(err)
	}
	if repo.Owner != "abc" {
		t.Errorf("expecting %v, got %v", "abc", repo.Owner)
	}
	if repo.Name != "x.y.z" {
		t.Errorf("expecting %v, got %v", "x.y.z", repo.Name)
	}
}
