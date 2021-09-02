// Copyright 2021 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package github

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestOrganizationsService_ListPackages(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/orgs/test/packages", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `[{
			"id": 197,
			"name": "hello_docker",
			"package_type": "container",
			"owner": {
			  "login": "github",
			  "id": 9919,
			  "node_id": "MDEyOk9yZ2FuaXphdGlvbjk5MTk=",
			  "avatar_url": "https://avatars.githubusercontent.com/u/9919?v=4",
			  "gravatar_id": "",
			  "url": "https://api.github.com/users/github",
			  "html_url": "https://github.com/github",
			  "followers_url": "https://api.github.com/users/github/followers",
			  "following_url": "https://api.github.com/users/github/following{/other_user}",
			  "gists_url": "https://api.github.com/users/github/gists{/gist_id}",
			  "starred_url": "https://api.github.com/users/github/starred{/owner}{/repo}",
			  "subscriptions_url": "https://api.github.com/users/github/subscriptions",
			  "organizations_url": "https://api.github.com/users/github/orgs",
			  "repos_url": "https://api.github.com/users/github/repos",
			  "events_url": "https://api.github.com/users/github/events{/privacy}",
			  "received_events_url": "https://api.github.com/users/github/received_events",
			  "type": "Organization",
			  "site_admin": false
			},
			"version_count": 1,
			"visibility": "private",
			"url": "https://api.github.com/orgs/github/packages/container/hello_docker",
			"created_at": `+referenceTimeStr+`,
			"updated_at": `+referenceTimeStr+`,
			"html_url": "https://github.com/orgs/github/packages/container/package/hello_docker"
		  }
		  ]`)
	})

	ctx := context.Background()
	packages, _, err := client.Organizations.ListPackages(ctx, "test", nil)
	if err != nil {
		t.Errorf("Organizations.ListPackages returned error: %v", err)
	}

	want := []*Package{{
		ID:           Int64(197),
		Name:         String("hello_docker"),
		PackageType:  String("container"),
		VersionCount: Int64(1),
		Visibility:   String("private"),
		URL:          String("https://api.github.com/orgs/github/packages/container/hello_docker"),
		HTMLURL:      String("https://github.com/orgs/github/packages/container/package/hello_docker"),
		CreatedAt:    &Timestamp{referenceTime},
		UpdatedAt:    &Timestamp{referenceTime},
		Owner: &User{
			Login:             String("github"),
			ID:                Int64(9919),
			NodeID:            String("MDEyOk9yZ2FuaXphdGlvbjk5MTk="),
			AvatarURL:         String("https://avatars.githubusercontent.com/u/9919?v=4"),
			GravatarID:        String(""),
			URL:               String("https://api.github.com/users/github"),
			HTMLURL:           String("https://github.com/github"),
			FollowersURL:      String("https://api.github.com/users/github/followers"),
			FollowingURL:      String("https://api.github.com/users/github/following{/other_user}"),
			GistsURL:          String("https://api.github.com/users/github/gists{/gist_id}"),
			StarredURL:        String("https://api.github.com/users/github/starred{/owner}{/repo}"),
			SubscriptionsURL:  String("https://api.github.com/users/github/subscriptions"),
			OrganizationsURL:  String("https://api.github.com/users/github/orgs"),
			ReposURL:          String("https://api.github.com/users/github/repos"),
			EventsURL:         String("https://api.github.com/users/github/events{/privacy}"),
			ReceivedEventsURL: String("https://api.github.com/users/github/received_events"),
			Type:              String("Organization"),
			SiteAdmin:         Bool(false),
		},
	}}
	if !cmp.Equal(packages, want) {
		t.Errorf("Organizations.ListPackages returned %+v, want %+v", packages, want)
	}

	const methodName = "ListPackages"
	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Organizations.ListPackages(ctx, "test", nil)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestOrganizationsService_GetPackage(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/orgs/test/packages/container/hello_docker", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{
			"id": 197,
			"name": "hello_docker",
			"package_type": "container",
			"version_count": 1,
			"visibility": "private",
			"url": "https://api.github.com/orgs/github/packages/container/hello_docker",
			"created_at": `+referenceTimeStr+`,
			"updated_at": `+referenceTimeStr+`,
			"html_url": "https://github.com/orgs/github/packages/container/package/hello_docker"
		  }`)
	})

	ctx := context.Background()
	packages, _, err := client.Organizations.GetPackage(ctx, "test", "container", "hello_docker")
	if err != nil {
		t.Errorf("Organizations.GetPackage returned error: %v", err)
	}

	want := &Package{
		ID:           Int64(197),
		Name:         String("hello_docker"),
		PackageType:  String("container"),
		VersionCount: Int64(1),
		Visibility:   String("private"),
		URL:          String("https://api.github.com/orgs/github/packages/container/hello_docker"),
		HTMLURL:      String("https://github.com/orgs/github/packages/container/package/hello_docker"),
		CreatedAt:    &Timestamp{referenceTime},
		UpdatedAt:    &Timestamp{referenceTime},
	}
	if !cmp.Equal(packages, want) {
		t.Errorf("Organizations.GetPackage returned %+v, want %+v", packages, want)
	}

	const methodName = "GetPackage"
	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Organizations.GetPackage(ctx, "", "", "")
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestOrganizationsService_DeletePackage(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/orgs/test/packages/container/hello_docker", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
	})

	ctx := context.Background()
	_, err := client.Organizations.DeletePackage(ctx, "test", "container", "hello_docker")
	if err != nil {
		t.Errorf("Organizations.DeletePackage returned error: %v", err)
	}

	const methodName = "DeletePackage"
	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Organizations.GetPackage(ctx, "", "", "")
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestOrganizationsService_RestorePackage(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/orgs/test/packages/container/hello_docker/restore", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
	})

	ctx := context.Background()
	_, err := client.Organizations.RestorePackage(ctx, "test", "container", "hello_docker")
	if err != nil {
		t.Errorf("Organizations.RestorePackage returned error: %v", err)
	}

	const methodName = "RestorePackage"
	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		return client.Organizations.RestorePackage(ctx, "", "container", "hello_docker")
	})
}

func TestOrganizationsService_ListPackagesVersions(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/orgs/test/packages/container/hello_docker/versions", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testFormValues(t, r, values{"per_page": "10", "page": "2", "state": "deleted"})
		fmt.Fprint(w, `[
			{
			  "id": 45763,
			  "name": "sha256:08a44bab0bddaddd8837a8b381aebc2e4b933768b981685a9e088360af0d3dd9",
			  "url": "https://api.github.com/users/octocat/packages/container/hello_docker/versions/45763",
			  "package_html_url": "https://github.com/users/octocat/packages/container/package/hello_docker",
			  "created_at": `+referenceTimeStr+`,
			  "updated_at": `+referenceTimeStr+`,
			  "html_url": "https://github.com/users/octocat/packages/container/hello_docker/45763",
			  "metadata": {
				"package_type": "container",
				"container": {
				  "tags": [
					"latest"
				  ]
				}
			  }
			}]`)
	})

	ctx := context.Background()
	packages, _, err := client.Organizations.PackageGetAllVersions(ctx, "test", "container", "hello_docker", &PackageListOptions{PerPage: 10, Page: 2, State: "deleted"})
	if err != nil {
		t.Errorf("Organizations.PackageGetAllVersions returned error: %v", err)
	}

	want := []*PackageVersion{{
		ID:             Int64(45763),
		Name:           String("sha256:08a44bab0bddaddd8837a8b381aebc2e4b933768b981685a9e088360af0d3dd9"),
		URL:            String("https://api.github.com/users/octocat/packages/container/hello_docker/versions/45763"),
		PackageHTMLURL: String("https://github.com/users/octocat/packages/container/package/hello_docker"),
		CreatedAt:      &Timestamp{referenceTime},
		UpdatedAt:      &Timestamp{referenceTime},
		HTMLURL:        String("https://github.com/users/octocat/packages/container/hello_docker/45763"),
		Metadata: &PackageMetadata{
			PackageType: String("container"),
			Container: &PackageContainerMetadata{
				Tags: []*string{String("latest")},
			},
		},
	}}
	if !cmp.Equal(packages, want) {
		t.Errorf("Organizations.PackageGetAllVersions returned %+v, want %+v", packages, want)
	}

	const methodName = "PackageGetAllVersions"
	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Organizations.PackageGetAllVersions(ctx, "", "", "", nil)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestOrganizationsService_PackageGetVersion(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/orgs/test/packages/container/hello_docker/versions/45763", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `
			{
			  "id": 45763,
			  "name": "sha256:08a44bab0bddaddd8837a8b381aebc2e4b933768b981685a9e088360af0d3dd9",
			  "url": "https://api.github.com/users/octocat/packages/container/hello_docker/versions/45763",
			  "package_html_url": "https://github.com/users/octocat/packages/container/package/hello_docker",
			  "created_at": `+referenceTimeStr+`,
			  "updated_at": `+referenceTimeStr+`,
			  "html_url": "https://github.com/users/octocat/packages/container/hello_docker/45763",
			  "metadata": {
				"package_type": "container",
				"container": {
				  "tags": [
					"latest"
				  ]
				}
			  }
			}`)
	})

	ctx := context.Background()
	packages, _, err := client.Organizations.PackageGetVersion(ctx, "test", "container", "hello_docker", 45763)
	if err != nil {
		t.Errorf("Organizations.PackageGetVersion returned error: %v", err)
	}

	want := &PackageVersion{
		ID:             Int64(45763),
		Name:           String("sha256:08a44bab0bddaddd8837a8b381aebc2e4b933768b981685a9e088360af0d3dd9"),
		URL:            String("https://api.github.com/users/octocat/packages/container/hello_docker/versions/45763"),
		PackageHTMLURL: String("https://github.com/users/octocat/packages/container/package/hello_docker"),
		CreatedAt:      &Timestamp{referenceTime},
		UpdatedAt:      &Timestamp{referenceTime},
		HTMLURL:        String("https://github.com/users/octocat/packages/container/hello_docker/45763"),
		Metadata: &PackageMetadata{
			PackageType: String("container"),
			Container: &PackageContainerMetadata{
				Tags: []*string{String("latest")},
			},
		},
	}
	if !cmp.Equal(packages, want) {
		t.Errorf("Organizations.PackageGetVersion returned %+v, want %+v", packages, want)
	}

	const methodName = "PackageGetVersion"
	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Organizations.PackageGetVersion(ctx, "", "", "", 45763)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestOrganizationsService_PackageDeleteVersion(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/orgs/test/packages/container/hello_docker/versions/45763", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
	})

	ctx := context.Background()
	_, err := client.Organizations.PackageDeleteVersion(ctx, "test", "container", "hello_docker", 45763)
	if err != nil {
		t.Errorf("Organizations.PackageDeleteVersion returned error: %v", err)
	}

	const methodName = "PackageDeleteVersion"
	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		return client.Organizations.PackageDeleteVersion(ctx, "", "", "", 45763)

	})
}

func TestOrganizationsService_PackageRestoreVersion(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/orgs/test/packages/container/hello_docker/versions/45763/restore", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
	})

	ctx := context.Background()
	_, err := client.Organizations.PackageRestoreVersion(ctx, "test", "container", "hello_docker", 45763)
	if err != nil {
		t.Errorf("Organizations.PackageRestoreVersion returned error: %v", err)
	}

	const methodName = "PackageRestoreVersion"
	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		return client.Organizations.PackageRestoreVersion(ctx, "", "", "", 45763)

	})
}
