// Copyright 2014 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build integration
// +build integration

package integration

import (
	"context"
	"io"
	"net/http"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-github/v63/github"
)

func TestRepositories_CRUD(t *testing.T) {
	if !checkAuth("TestRepositories_CRUD") {
		return
	}

	// create a random repository
	repo, err := createRandomTestRepository("", true)
	if err != nil {
		t.Fatalf("createRandomTestRepository returned error: %v", err)
	}

	// update the repository description
	repo.Description = github.String("description")
	repo.DefaultBranch = nil // FIXME: this shouldn't be necessary
	_, _, err = client.Repositories.Edit(context.Background(), *repo.Owner.Login, *repo.Name, repo)
	if err != nil {
		t.Fatalf("Repositories.Edit() returned error: %v", err)
	}

	// delete the repository
	_, err = client.Repositories.Delete(context.Background(), *repo.Owner.Login, *repo.Name)
	if err != nil {
		t.Fatalf("Repositories.Delete() returned error: %v", err)
	}

	// verify that the repository was deleted
	_, resp, err := client.Repositories.Get(context.Background(), *repo.Owner.Login, *repo.Name)
	if err == nil {
		t.Fatalf("Test repository still exists after deleting it.")
	}
	if err != nil && resp.StatusCode != http.StatusNotFound {
		t.Fatalf("Repositories.Get() returned error: %v", err)
	}
}

func TestRepositories_BranchesTags(t *testing.T) {
	// branches
	branches, _, err := client.Repositories.ListBranches(context.Background(), "git", "git", nil)
	if err != nil {
		t.Fatalf("Repositories.ListBranches() returned error: %v", err)
	}

	if len(branches) == 0 {
		t.Fatalf("Repositories.ListBranches('git', 'git') returned no branches")
	}

	_, _, err = client.Repositories.GetBranch(context.Background(), "git", "git", *branches[0].Name, 0)
	if err != nil {
		t.Fatalf("Repositories.GetBranch() returned error: %v", err)
	}

	// tags
	tags, _, err := client.Repositories.ListTags(context.Background(), "git", "git", nil)
	if err != nil {
		t.Fatalf("Repositories.ListTags() returned error: %v", err)
	}

	if len(tags) == 0 {
		t.Fatalf("Repositories.ListTags('git', 'git') returned no tags")
	}
}

func TestRepositories_EditBranches(t *testing.T) {
	if !checkAuth("TestRepositories_EditBranches") {
		return
	}

	// create a random repository
	repo, err := createRandomTestRepository("", true)
	if err != nil {
		t.Fatalf("createRandomTestRepository returned error: %v", err)
	}

	branch, _, err := client.Repositories.GetBranch(context.Background(), *repo.Owner.Login, *repo.Name, "master", 0)
	if err != nil {
		t.Fatalf("Repositories.GetBranch() returned error: %v", err)
	}

	if *branch.Protected {
		t.Fatalf("Branch %v of repo %v is already protected", "master", *repo.Name)
	}

	protectionRequest := &github.ProtectionRequest{
		RequiredStatusChecks: &github.RequiredStatusChecks{
			Strict:   true,
			Contexts: &[]string{"continuous-integration"},
		},
		RequiredPullRequestReviews: &github.PullRequestReviewsEnforcementRequest{
			DismissStaleReviews: true,
		},
		EnforceAdmins: true,
		// TODO: Only organization repositories can have users and team restrictions.
		//       In order to be able to test these Restrictions, need to add support
		//       for creating temporary organization repositories.
		Restrictions:     nil,
		BlockCreations:   github.Bool(false),
		LockBranch:       github.Bool(false),
		AllowForkSyncing: github.Bool(false),
	}

	protection, _, err := client.Repositories.UpdateBranchProtection(context.Background(), *repo.Owner.Login, *repo.Name, "master", protectionRequest)
	if err != nil {
		t.Fatalf("Repositories.UpdateBranchProtection() returned error: %v", err)
	}

	want := &github.Protection{
		RequiredStatusChecks: &github.RequiredStatusChecks{
			Strict:   true,
			Contexts: &[]string{"continuous-integration"},
		},
		RequiredPullRequestReviews: &github.PullRequestReviewsEnforcement{
			DismissStaleReviews:          true,
			RequiredApprovingReviewCount: 0,
		},
		EnforceAdmins: &github.AdminEnforcement{
			URL:     github.String("https://api.github.com/repos/" + *repo.Owner.Login + "/" + *repo.Name + "/branches/master/protection/enforce_admins"),
			Enabled: true,
		},
		Restrictions: nil,
		BlockCreations: &github.BlockCreations{
			Enabled: github.Bool(false),
		},
		LockBranch: &github.LockBranch{
			Enabled: github.Bool(false),
		},
		AllowForkSyncing: &github.AllowForkSyncing{
			Enabled: github.Bool(false),
		},
	}
	if !cmp.Equal(protection, want) {
		t.Errorf("Repositories.UpdateBranchProtection() returned %+v, want %+v", protection, want)
	}

	_, err = client.Repositories.Delete(context.Background(), *repo.Owner.Login, *repo.Name)
	if err != nil {
		t.Fatalf("Repositories.Delete() returned error: %v", err)
	}
}

func TestRepositories_ListByAuthenticatedUser(t *testing.T) {
	if !checkAuth("TestRepositories_ListByAuthenticatedUser") {
		return
	}

	_, _, err := client.Repositories.ListByAuthenticatedUser(context.Background(), nil)
	if err != nil {
		t.Fatalf("Repositories.ListByAuthenticatedUser() returned error: %v", err)
	}
}

func TestRepositories_ListByUser(t *testing.T) {
	_, _, err := client.Repositories.ListByUser(context.Background(), "google", nil)
	if err != nil {
		t.Fatalf("Repositories.ListByUser('google') returned error: %v", err)
	}

	opt := github.RepositoryListByUserOptions{Sort: "created"}
	repos, _, err := client.Repositories.ListByUser(context.Background(), "google", &opt)
	if err != nil {
		t.Fatalf("Repositories.List('google') with Sort opt returned error: %v", err)
	}
	for i, repo := range repos {
		if i > 0 && (*repos[i-1].CreatedAt).Time.Before((*repo.CreatedAt).Time) {
			t.Fatalf("Repositories.ListByUser('google') with default descending Sort returned incorrect order")
		}
	}
}

func TestRepositories_DownloadReleaseAsset(t *testing.T) {
	if !checkAuth("TestRepositories_DownloadReleaseAsset") {
		return
	}

	rc, _, err := client.Repositories.DownloadReleaseAsset(context.Background(), "andersjanmyr", "goose", 484892, http.DefaultClient)
	if err != nil {
		t.Fatalf("Repositories.DownloadReleaseAsset(andersjanmyr, goose, 484892, true) returned error: %v", err)
	}
	defer func() { _ = rc.Close() }()
	_, err = io.Copy(io.Discard, rc)
	if err != nil {
		t.Fatalf("Repositories.DownloadReleaseAsset(andersjanmyr, goose, 484892, true) returned error: %v", err)
	}
}

func TestRepositories_Autolinks(t *testing.T) {
	if !checkAuth("TestRepositories_Autolinks") {
		return
	}

	// create a random repository
	repo, err := createRandomTestRepository("", true)
	if err != nil {
		t.Fatalf("createRandomTestRepository returned error: %v", err)
	}

	opts := &github.AutolinkOptions{
		KeyPrefix:      github.String("TICKET-"),
		URLTemplate:    github.String("https://example.com/TICKET?query=<num>"),
		IsAlphanumeric: github.Bool(false),
	}

	actionlink, _, err := client.Repositories.AddAutolink(context.Background(), *repo.Owner.Login, *repo.Name, opts)
	if err != nil {
		t.Fatalf("Repositories.AddAutolink() returned error: %v", err)
	}

	if !cmp.Equal(actionlink.KeyPrefix, opts.KeyPrefix) ||
		!cmp.Equal(actionlink.URLTemplate, opts.URLTemplate) ||
		!cmp.Equal(actionlink.IsAlphanumeric, opts.IsAlphanumeric) {
		t.Errorf("Repositories.AddAutolink() returned %+v, want %+v", actionlink, opts)
	}

	_, err = client.Repositories.Delete(context.Background(), *repo.Owner.Login, *repo.Name)
	if err != nil {
		t.Fatalf("Repositories.Delete() returned error: %v", err)
	}
}
