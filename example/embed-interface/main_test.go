// Copyright 2021 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/google/go-github/v34/github"
)

type fakeOrgSvc struct {
	github.OrganizationsServiceInterface

	orgs []*github.Organization
}

func (f *fakeOrgSvc) List(ctx context.Context, org string, opts *github.ListOptions) ([]*github.Organization, *github.Response, error) {
	if org != "octocat" {
		return nil, nil, errors.New("unexpected org")
	}

	return f.orgs, nil, nil
}

func TestFetchOrganizations(t *testing.T) {
	want := []*github.Organization{
		{Name: github.String("octocat")},
	}

	orgService := &fakeOrgSvc{orgs: want}
	got, err := fetchOrganizations(orgService, "octocat")
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("got = %#v, want = %#v", got, want)
	}
}
