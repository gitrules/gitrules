//go:build integration
// +build integration

package github_test

import (
	"context"
	"testing"

	govgh "github.com/gitrules/gitrules/github/lib"
	"github.com/google/go-github/v66/github"
)

func TestCreateLabel(t *testing.T) {
	ctx := context.Background()
	testLabel := "xyz:test-label"

	client.Issues.DeleteLabel(ctx, TestRepo.Owner, TestRepo.Name, testLabel)

	label := &github.Label{Name: github.String(testLabel)}

	_, _, err := client.Issues.CreateLabel(ctx, TestRepo.Owner, TestRepo.Name, label)
	if err != nil {
		t.Fatal(err)
	}

	_, _, err = client.Issues.CreateLabel(ctx, TestRepo.Owner, TestRepo.Name, label)
	if err == nil {
		t.Fatalf("error is expected")
	}

	if !govgh.IsLabelAlreadyExists(err) {
		t.Errorf("not expecting %v", err)
	}
}
