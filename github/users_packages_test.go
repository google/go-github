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

func TestUsersService_ListPackages(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/user/packages", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testFormValues(t, r, values{"package_type": "container", "visibility": "private"})
		fmt.Fprint(w, `[{
			"id": 197,
			"name": "hello_docker",
			"package_type": "container",
			"version_count": 1,
			"visibility": "private",
			"url": "https://api.github.com/orgs/github/packages/container/hello_docker",
			"created_at": `+referenceTimeStr+`,
			"updated_at": `+referenceTimeStr+`,
			"html_url": "https://github.com/orgs/github/packages/container/package/hello_docker"
		  }]`)
	})

	ctx := context.Background()
	packages, _, err := client.Users.ListPackages(ctx, "", &PackageListOptions{PackageType: String("container"), Visibility: String("private")})
	if err != nil {
		t.Errorf("Users.ListPackages returned error: %v", err)
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
	}}
	if !cmp.Equal(packages, want) {
		t.Errorf("Users.ListPackages returned %+v, want %+v", packages, want)
	}

	const methodName = "ListPackages"
	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Users.ListPackages(ctx, "", nil)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestUsersService_GetPackage(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/users/test/packages/container/hello_docker", func(w http.ResponseWriter, r *http.Request) {
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
	packages, _, err := client.Users.GetPackage(ctx, "test", "container", "hello_docker")
	if err != nil {
		t.Errorf("Users.GetPackage returned error: %v", err)
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
		t.Errorf("Users.GetPackage returned %+v, want %+v", packages, want)
	}

	const methodName = "GetPackage"
	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Users.GetPackage(ctx, "", "", "")
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestUsersService_DeletePackage(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/user/packages/container/hello_docker", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
	})

	ctx := context.Background()
	_, err := client.Users.DeletePackage(ctx, "", "container", "hello_docker")
	if err != nil {
		t.Errorf("Users.DeletePackage returned error: %v", err)
	}

	const methodName = "DeletePackage"
	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Users.GetPackage(ctx, "", "", "")
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestUsersService_RestorePackage(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/user/packages/container/hello_docker/restore", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
	})

	ctx := context.Background()
	_, err := client.Users.RestorePackage(ctx, "", "container", "hello_docker")
	if err != nil {
		t.Errorf("Users.RestorePackage returned error: %v", err)
	}

	const methodName = "RestorePackage"
	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		return client.Users.RestorePackage(ctx, "", "container", "hello_docker")
	})
}

func TestUsersService_ListPackagesVersions(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/user/packages/container/hello_docker/versions", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
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
	packages, _, err := client.Users.PackageGetAllVersions(ctx, "", "container", "hello_docker")
	if err != nil {
		t.Errorf("Users.PackageGetAllVersions returned error: %v", err)
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
				Tags: []string{"latest"},
			},
		},
	}}
	if !cmp.Equal(packages, want) {
		t.Errorf("Users.PackageGetAllVersions returned %+v, want %+v", packages, want)
	}

	const methodName = "PackageGetAllVersions"
	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Users.PackageGetAllVersions(ctx, "", "", "")
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestUsersService_PackageGetVersion(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/user/packages/container/hello_docker/versions/45763", func(w http.ResponseWriter, r *http.Request) {
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
	packages, _, err := client.Users.PackageGetVersion(ctx, "", "container", "hello_docker", 45763)
	if err != nil {
		t.Errorf("Users.PackageGetVersion returned error: %v", err)
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
				Tags: []string{"latest"},
			},
		},
	}
	if !cmp.Equal(packages, want) {
		t.Errorf("Users.PackageGetVersion returned %+v, want %+v", packages, want)
	}

	const methodName = "PackageGetVersion"
	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Users.PackageGetVersion(ctx, "", "", "", 45763)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestUsersService_PackageDeleteVersion(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/user/packages/container/hello_docker/versions/45763", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
	})

	ctx := context.Background()
	_, err := client.Users.PackageDeleteVersion(ctx, "", "container", "hello_docker", 45763)
	if err != nil {
		t.Errorf("Users.PackageDeleteVersion returned error: %v", err)
	}

	const methodName = "PackageDeleteVersion"
	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		return client.Users.PackageDeleteVersion(ctx, "", "", "", 45763)

	})
}

func TestUsersService_PackageRestoreVersion(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/user/packages/container/hello_docker/versions/45763/restore", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
	})

	ctx := context.Background()
	_, err := client.Users.PackageRestoreVersion(ctx, "", "container", "hello_docker", 45763)
	if err != nil {
		t.Errorf("Users.PackageRestoreVersion returned error: %v", err)
	}

	const methodName = "PackageRestoreVersion"
	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		return client.Users.PackageRestoreVersion(ctx, "", "", "", 45763)

	})
}
