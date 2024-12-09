// Copyright 2016 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package github

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
)

func TestAppsService_Get_authenticatedApp(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/app", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{"id":1}`)
	})

	ctx := context.Background()
	app, _, err := client.Apps.Get(ctx, "")
	if err != nil {
		t.Errorf("Apps.Get returned error: %v", err)
	}

	want := &App{ID: Ptr(int64(1))}
	if !cmp.Equal(app, want) {
		t.Errorf("Apps.Get returned %+v, want %+v", app, want)
	}

	const methodName = "Get"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Apps.Get(ctx, "\n")
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Apps.Get(ctx, "")
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestAppsService_Get_specifiedApp(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/apps/a", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{"html_url":"https://github.com/apps/a"}`)
	})

	ctx := context.Background()
	app, _, err := client.Apps.Get(ctx, "a")
	if err != nil {
		t.Errorf("Apps.Get returned error: %v", err)
	}

	want := &App{HTMLURL: Ptr("https://github.com/apps/a")}
	if !cmp.Equal(app, want) {
		t.Errorf("Apps.Get returned %+v, want %+v", *app.HTMLURL, *want.HTMLURL)
	}
}

func TestAppsService_ListInstallationRequests(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/app/installation-requests", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testFormValues(t, r, values{
			"page":     "1",
			"per_page": "2",
		})
		fmt.Fprint(w, `[{
			"id": 1,
			"account": { "id": 2 },
			"requester": { "id": 3 },
			"created_at": "2018-01-01T00:00:00Z"
		}]`,
		)
	})

	opt := &ListOptions{Page: 1, PerPage: 2}
	ctx := context.Background()
	installationRequests, _, err := client.Apps.ListInstallationRequests(ctx, opt)
	if err != nil {
		t.Errorf("Apps.ListInstallations returned error: %v", err)
	}

	date := Timestamp{Time: time.Date(2018, time.January, 1, 0, 0, 0, 0, time.UTC)}
	want := []*InstallationRequest{{
		ID:        Ptr(int64(1)),
		Account:   &User{ID: Ptr(int64(2))},
		Requester: &User{ID: Ptr(int64(3))},
		CreatedAt: &date,
	}}
	if !cmp.Equal(installationRequests, want) {
		t.Errorf("Apps.ListInstallationRequests returned %+v, want %+v", installationRequests, want)
	}

	const methodName = "ListInstallationRequests"
	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Apps.ListInstallationRequests(ctx, opt)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestAppsService_ListInstallations(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/app/installations", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testFormValues(t, r, values{
			"page":     "1",
			"per_page": "2",
		})
		fmt.Fprint(w, `[{
                                   "id":1,
                                   "app_id":1,
                                   "target_id":1,
                                   "target_type": "Organization",
                                   "permissions": {
                                       "actions": "read",
                                       "administration": "read",
                                       "checks": "read",
                                       "contents": "read",
                                       "content_references": "read",
                                       "deployments": "read",
                                       "environments": "read",
                                       "issues": "write",
                                       "metadata": "read",
                                       "members": "read",
                                       "organization_administration": "write",
                                       "organization_custom_roles": "write",
                                       "organization_hooks": "write",
                                       "organization_packages": "write",
                                       "organization_personal_access_tokens": "read",
                                       "organization_personal_access_token_requests": "read",
                                       "organization_plan": "read",
                                       "organization_pre_receive_hooks": "write",
                                       "organization_projects": "read",
                                       "organization_secrets": "read",
                                       "organization_self_hosted_runners": "read",
                                       "organization_user_blocking": "write",
                                       "packages": "read",
                                       "pages": "read",
                                       "pull_requests": "write",
                                       "repository_hooks": "write",
                                       "repository_projects": "read",
                                       "repository_pre_receive_hooks": "read",
                                       "secrets": "read",
                                       "secret_scanning_alerts": "read",
                                       "security_events": "read",
                                       "single_file": "write",
                                       "statuses": "write",
                                       "team_discussions": "read",
                                       "vulnerability_alerts": "read",
                                       "workflows": "write"
                                   },
                                  "events": [
                                      "push",
                                      "pull_request"
                                  ],
                                 "single_file_name": "config.yml",
                                 "repository_selection": "selected",
                                 "created_at": "2018-01-01T00:00:00Z",
                                 "updated_at": "2018-01-01T00:00:00Z"}]`,
		)
	})

	opt := &ListOptions{Page: 1, PerPage: 2}
	ctx := context.Background()
	installations, _, err := client.Apps.ListInstallations(ctx, opt)
	if err != nil {
		t.Errorf("Apps.ListInstallations returned error: %v", err)
	}

	date := Timestamp{Time: time.Date(2018, time.January, 1, 0, 0, 0, 0, time.UTC)}
	want := []*Installation{{
		ID:                  Ptr(int64(1)),
		AppID:               Ptr(int64(1)),
		TargetID:            Ptr(int64(1)),
		TargetType:          Ptr("Organization"),
		SingleFileName:      Ptr("config.yml"),
		RepositorySelection: Ptr("selected"),
		Permissions: &InstallationPermissions{
			Actions:                                 Ptr("read"),
			Administration:                          Ptr("read"),
			Checks:                                  Ptr("read"),
			Contents:                                Ptr("read"),
			ContentReferences:                       Ptr("read"),
			Deployments:                             Ptr("read"),
			Environments:                            Ptr("read"),
			Issues:                                  Ptr("write"),
			Metadata:                                Ptr("read"),
			Members:                                 Ptr("read"),
			OrganizationAdministration:              Ptr("write"),
			OrganizationCustomRoles:                 Ptr("write"),
			OrganizationHooks:                       Ptr("write"),
			OrganizationPackages:                    Ptr("write"),
			OrganizationPersonalAccessTokens:        Ptr("read"),
			OrganizationPersonalAccessTokenRequests: Ptr("read"),
			OrganizationPlan:                        Ptr("read"),
			OrganizationPreReceiveHooks:             Ptr("write"),
			OrganizationProjects:                    Ptr("read"),
			OrganizationSecrets:                     Ptr("read"),
			OrganizationSelfHostedRunners:           Ptr("read"),
			OrganizationUserBlocking:                Ptr("write"),
			Packages:                                Ptr("read"),
			Pages:                                   Ptr("read"),
			PullRequests:                            Ptr("write"),
			RepositoryHooks:                         Ptr("write"),
			RepositoryProjects:                      Ptr("read"),
			RepositoryPreReceiveHooks:               Ptr("read"),
			Secrets:                                 Ptr("read"),
			SecretScanningAlerts:                    Ptr("read"),
			SecurityEvents:                          Ptr("read"),
			SingleFile:                              Ptr("write"),
			Statuses:                                Ptr("write"),
			TeamDiscussions:                         Ptr("read"),
			VulnerabilityAlerts:                     Ptr("read"),
			Workflows:                               Ptr("write")},
		Events:    []string{"push", "pull_request"},
		CreatedAt: &date,
		UpdatedAt: &date,
	}}
	if !cmp.Equal(installations, want) {
		t.Errorf("Apps.ListInstallations returned %+v, want %+v", installations, want)
	}

	const methodName = "ListInstallations"
	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Apps.ListInstallations(ctx, opt)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestAppsService_GetInstallation(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/app/installations/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{"id":1, "app_id":1, "target_id":1, "target_type": "Organization"}`)
	})

	ctx := context.Background()
	installation, _, err := client.Apps.GetInstallation(ctx, 1)
	if err != nil {
		t.Errorf("Apps.GetInstallation returned error: %v", err)
	}

	want := &Installation{ID: Ptr(int64(1)), AppID: Ptr(int64(1)), TargetID: Ptr(int64(1)), TargetType: Ptr("Organization")}
	if !cmp.Equal(installation, want) {
		t.Errorf("Apps.GetInstallation returned %+v, want %+v", installation, want)
	}

	const methodName = "GetInstallation"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Apps.GetInstallation(ctx, -1)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Apps.GetInstallation(ctx, 1)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestAppsService_ListUserInstallations(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/user/installations", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testFormValues(t, r, values{
			"page":     "1",
			"per_page": "2",
		})
		fmt.Fprint(w, `{"installations":[{"id":1, "app_id":1, "target_id":1, "target_type": "Organization"}]}`)
	})

	opt := &ListOptions{Page: 1, PerPage: 2}
	ctx := context.Background()
	installations, _, err := client.Apps.ListUserInstallations(ctx, opt)
	if err != nil {
		t.Errorf("Apps.ListUserInstallations returned error: %v", err)
	}

	want := []*Installation{{ID: Ptr(int64(1)), AppID: Ptr(int64(1)), TargetID: Ptr(int64(1)), TargetType: Ptr("Organization")}}
	if !cmp.Equal(installations, want) {
		t.Errorf("Apps.ListUserInstallations returned %+v, want %+v", installations, want)
	}

	const methodName = "ListUserInstallations"
	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Apps.ListUserInstallations(ctx, opt)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestAppsService_SuspendInstallation(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/app/installations/1/suspended", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PUT")

		w.WriteHeader(http.StatusNoContent)
	})

	ctx := context.Background()
	if _, err := client.Apps.SuspendInstallation(ctx, 1); err != nil {
		t.Errorf("Apps.SuspendInstallation returned error: %v", err)
	}

	const methodName = "SuspendInstallation"
	testBadOptions(t, methodName, func() (err error) {
		_, err = client.Apps.SuspendInstallation(ctx, -1)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		return client.Apps.SuspendInstallation(ctx, 1)
	})
}

func TestAppsService_UnsuspendInstallation(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/app/installations/1/suspended", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")

		w.WriteHeader(http.StatusNoContent)
	})

	ctx := context.Background()
	if _, err := client.Apps.UnsuspendInstallation(ctx, 1); err != nil {
		t.Errorf("Apps.UnsuspendInstallation returned error: %v", err)
	}

	const methodName = "UnsuspendInstallation"
	testBadOptions(t, methodName, func() (err error) {
		_, err = client.Apps.UnsuspendInstallation(ctx, -1)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		return client.Apps.UnsuspendInstallation(ctx, 1)
	})
}

func TestAppsService_DeleteInstallation(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/app/installations/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
		w.WriteHeader(http.StatusNoContent)
	})

	ctx := context.Background()
	_, err := client.Apps.DeleteInstallation(ctx, 1)
	if err != nil {
		t.Errorf("Apps.DeleteInstallation returned error: %v", err)
	}

	const methodName = "DeleteInstallation"
	testBadOptions(t, methodName, func() (err error) {
		_, err = client.Apps.DeleteInstallation(ctx, -1)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		return client.Apps.DeleteInstallation(ctx, 1)
	})
}

func TestAppsService_CreateInstallationToken(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/app/installations/1/access_tokens", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		fmt.Fprint(w, `{"token":"t"}`)
	})

	ctx := context.Background()
	token, _, err := client.Apps.CreateInstallationToken(ctx, 1, nil)
	if err != nil {
		t.Errorf("Apps.CreateInstallationToken returned error: %v", err)
	}

	want := &InstallationToken{Token: Ptr("t")}
	if !cmp.Equal(token, want) {
		t.Errorf("Apps.CreateInstallationToken returned %+v, want %+v", token, want)
	}

	const methodName = "CreateInstallationToken"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Apps.CreateInstallationToken(ctx, -1, nil)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Apps.CreateInstallationToken(ctx, 1, nil)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestAppsService_CreateInstallationTokenWithOptions(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	installationTokenOptions := &InstallationTokenOptions{
		RepositoryIDs: []int64{1234},
		Repositories:  []string{"foo"},
		Permissions: &InstallationPermissions{
			Contents: Ptr("write"),
			Issues:   Ptr("read"),
		},
	}

	mux.HandleFunc("/app/installations/1/access_tokens", func(w http.ResponseWriter, r *http.Request) {
		v := new(InstallationTokenOptions)
		assertNilError(t, json.NewDecoder(r.Body).Decode(v))

		if !cmp.Equal(v, installationTokenOptions) {
			t.Errorf("request sent %+v, want %+v", v, installationTokenOptions)
		}

		testMethod(t, r, "POST")
		fmt.Fprint(w, `{"token":"t"}`)
	})

	ctx := context.Background()
	token, _, err := client.Apps.CreateInstallationToken(ctx, 1, installationTokenOptions)
	if err != nil {
		t.Errorf("Apps.CreateInstallationToken returned error: %v", err)
	}

	want := &InstallationToken{Token: Ptr("t")}
	if !cmp.Equal(token, want) {
		t.Errorf("Apps.CreateInstallationToken returned %+v, want %+v", token, want)
	}
}

func TestAppsService_CreateInstallationTokenListReposWithOptions(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	installationTokenListRepoOptions := &InstallationTokenListRepoOptions{
		Repositories: []string{"foo"},
		Permissions: &InstallationPermissions{
			Contents: Ptr("write"),
			Issues:   Ptr("read"),
		},
	}

	mux.HandleFunc("/app/installations/1/access_tokens", func(w http.ResponseWriter, r *http.Request) {
		v := new(InstallationTokenListRepoOptions)
		assertNilError(t, json.NewDecoder(r.Body).Decode(v))

		if !cmp.Equal(v, installationTokenListRepoOptions) {
			t.Errorf("request sent %+v, want %+v", v, installationTokenListRepoOptions)
		}

		testMethod(t, r, "POST")
		fmt.Fprint(w, `{"token":"t"}`)
	})

	ctx := context.Background()
	token, _, err := client.Apps.CreateInstallationTokenListRepos(ctx, 1, installationTokenListRepoOptions)
	if err != nil {
		t.Errorf("Apps.CreateInstallationTokenListRepos returned error: %v", err)
	}

	want := &InstallationToken{Token: Ptr("t")}
	if !cmp.Equal(token, want) {
		t.Errorf("Apps.CreateInstallationTokenListRepos returned %+v, want %+v", token, want)
	}
}

func TestAppsService_CreateInstallationTokenListReposWithNoOptions(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/app/installations/1/access_tokens", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		fmt.Fprint(w, `{"token":"t"}`)
	})

	ctx := context.Background()
	token, _, err := client.Apps.CreateInstallationTokenListRepos(ctx, 1, nil)
	if err != nil {
		t.Errorf("Apps.CreateInstallationTokenListRepos returned error: %v", err)
	}

	want := &InstallationToken{Token: Ptr("t")}
	if !cmp.Equal(token, want) {
		t.Errorf("Apps.CreateInstallationTokenListRepos returned %+v, want %+v", token, want)
	}

	const methodName = "CreateInstallationTokenListRepos"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Apps.CreateInstallationTokenListRepos(ctx, -1, nil)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Apps.CreateInstallationTokenListRepos(ctx, 1, nil)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestAppsService_CreateAttachment(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/content_references/11/attachments", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		testHeader(t, r, "Accept", mediaTypeContentAttachmentsPreview)

		w.WriteHeader(http.StatusOK)
		assertWrite(t, w, []byte(`{"id":1,"title":"title1","body":"body1"}`))
	})

	ctx := context.Background()
	got, _, err := client.Apps.CreateAttachment(ctx, 11, "title1", "body1")
	if err != nil {
		t.Errorf("CreateAttachment returned error: %v", err)
	}

	want := &Attachment{ID: Ptr(int64(1)), Title: Ptr("title1"), Body: Ptr("body1")}
	if !cmp.Equal(got, want) {
		t.Errorf("CreateAttachment = %+v, want %+v", got, want)
	}

	const methodName = "CreateAttachment"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Apps.CreateAttachment(ctx, -11, "\n", "\n")
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Apps.CreateAttachment(ctx, 11, "title1", "body1")
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestAppsService_FindOrganizationInstallation(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/orgs/o/installation", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{"id":1, "app_id":1, "target_id":1, "target_type": "Organization"}`)
	})

	ctx := context.Background()
	installation, _, err := client.Apps.FindOrganizationInstallation(ctx, "o")
	if err != nil {
		t.Errorf("Apps.FindOrganizationInstallation returned error: %v", err)
	}

	want := &Installation{ID: Ptr(int64(1)), AppID: Ptr(int64(1)), TargetID: Ptr(int64(1)), TargetType: Ptr("Organization")}
	if !cmp.Equal(installation, want) {
		t.Errorf("Apps.FindOrganizationInstallation returned %+v, want %+v", installation, want)
	}

	const methodName = "FindOrganizationInstallation"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Apps.FindOrganizationInstallation(ctx, "\n")
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Apps.FindOrganizationInstallation(ctx, "o")
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestAppsService_FindRepositoryInstallation(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/repos/o/r/installation", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{"id":1, "app_id":1, "target_id":1, "target_type": "Organization"}`)
	})

	ctx := context.Background()
	installation, _, err := client.Apps.FindRepositoryInstallation(ctx, "o", "r")
	if err != nil {
		t.Errorf("Apps.FindRepositoryInstallation returned error: %v", err)
	}

	want := &Installation{ID: Ptr(int64(1)), AppID: Ptr(int64(1)), TargetID: Ptr(int64(1)), TargetType: Ptr("Organization")}
	if !cmp.Equal(installation, want) {
		t.Errorf("Apps.FindRepositoryInstallation returned %+v, want %+v", installation, want)
	}

	const methodName = "FindRepositoryInstallation"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Apps.FindRepositoryInstallation(ctx, "\n", "\n")
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Apps.FindRepositoryInstallation(ctx, "o", "r")
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestAppsService_FindRepositoryInstallationByID(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/repositories/1/installation", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{"id":1, "app_id":1, "target_id":1, "target_type": "Organization"}`)
	})

	ctx := context.Background()
	installation, _, err := client.Apps.FindRepositoryInstallationByID(ctx, 1)
	if err != nil {
		t.Errorf("Apps.FindRepositoryInstallationByID returned error: %v", err)
	}

	want := &Installation{ID: Ptr(int64(1)), AppID: Ptr(int64(1)), TargetID: Ptr(int64(1)), TargetType: Ptr("Organization")}
	if !cmp.Equal(installation, want) {
		t.Errorf("Apps.FindRepositoryInstallationByID returned %+v, want %+v", installation, want)
	}

	const methodName = "FindRepositoryInstallationByID"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Apps.FindRepositoryInstallationByID(ctx, -1)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Apps.FindRepositoryInstallationByID(ctx, 1)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestAppsService_FindUserInstallation(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/users/u/installation", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{"id":1, "app_id":1, "target_id":1, "target_type": "User"}`)
	})

	ctx := context.Background()
	installation, _, err := client.Apps.FindUserInstallation(ctx, "u")
	if err != nil {
		t.Errorf("Apps.FindUserInstallation returned error: %v", err)
	}

	want := &Installation{ID: Ptr(int64(1)), AppID: Ptr(int64(1)), TargetID: Ptr(int64(1)), TargetType: Ptr("User")}
	if !cmp.Equal(installation, want) {
		t.Errorf("Apps.FindUserInstallation returned %+v, want %+v", installation, want)
	}

	const methodName = "FindUserInstallation"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Apps.FindUserInstallation(ctx, "\n")
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Apps.FindUserInstallation(ctx, "u")
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestContentReference_Marshal(t *testing.T) {
	t.Parallel()
	testJSONMarshal(t, &ContentReference{}, "{}")

	u := &ContentReference{
		ID:        Ptr(int64(1)),
		NodeID:    Ptr("nid"),
		Reference: Ptr("r"),
	}

	want := `{
		"id": 1,
		"node_id": "nid",
		"reference": "r"
	}`

	testJSONMarshal(t, u, want)
}

func TestAttachment_Marshal(t *testing.T) {
	t.Parallel()
	testJSONMarshal(t, &Attachment{}, "{}")

	u := &Attachment{
		ID:    Ptr(int64(1)),
		Title: Ptr("t"),
		Body:  Ptr("b"),
	}

	want := `{
		"id": 1,
		"title": "t",
		"body": "b"
	}`

	testJSONMarshal(t, u, want)
}

func TestInstallationPermissions_Marshal(t *testing.T) {
	t.Parallel()
	testJSONMarshal(t, &InstallationPermissions{}, "{}")

	u := &InstallationPermissions{
		Actions:                       Ptr("a"),
		Administration:                Ptr("ad"),
		Checks:                        Ptr("c"),
		Contents:                      Ptr("co"),
		ContentReferences:             Ptr("cr"),
		Deployments:                   Ptr("d"),
		Environments:                  Ptr("e"),
		Issues:                        Ptr("i"),
		Metadata:                      Ptr("md"),
		Members:                       Ptr("m"),
		OrganizationAdministration:    Ptr("oa"),
		OrganizationCustomOrgRoles:    Ptr("ocr"),
		OrganizationHooks:             Ptr("oh"),
		OrganizationPlan:              Ptr("op"),
		OrganizationPreReceiveHooks:   Ptr("opr"),
		OrganizationProjects:          Ptr("op"),
		OrganizationSecrets:           Ptr("os"),
		OrganizationSelfHostedRunners: Ptr("osh"),
		OrganizationUserBlocking:      Ptr("oub"),
		Packages:                      Ptr("pkg"),
		Pages:                         Ptr("pg"),
		PullRequests:                  Ptr("pr"),
		RepositoryHooks:               Ptr("rh"),
		RepositoryProjects:            Ptr("rp"),
		RepositoryPreReceiveHooks:     Ptr("rprh"),
		Secrets:                       Ptr("s"),
		SecretScanningAlerts:          Ptr("ssa"),
		SecurityEvents:                Ptr("se"),
		SingleFile:                    Ptr("sf"),
		Statuses:                      Ptr("s"),
		TeamDiscussions:               Ptr("td"),
		VulnerabilityAlerts:           Ptr("va"),
		Workflows:                     Ptr("w"),
	}

	want := `{
		"actions": "a",
		"administration": "ad",
		"checks": "c",
		"contents": "co",
		"content_references": "cr",
		"deployments": "d",
		"environments": "e",
		"issues": "i",
		"metadata": "md",
		"members": "m",
		"organization_administration": "oa",
		"organization_custom_org_roles": "ocr",
		"organization_hooks": "oh",
		"organization_plan": "op",
		"organization_pre_receive_hooks": "opr",
		"organization_projects": "op",
		"organization_secrets": "os",
		"organization_self_hosted_runners": "osh",
		"organization_user_blocking": "oub",
		"packages": "pkg",
		"pages": "pg",
		"pull_requests": "pr",
		"repository_hooks": "rh",
		"repository_projects": "rp",
		"repository_pre_receive_hooks": "rprh",
		"secrets": "s",
		"secret_scanning_alerts": "ssa",
		"security_events": "se",
		"single_file": "sf",
		"statuses": "s",
		"team_discussions": "td",
		"vulnerability_alerts":"va",
		"workflows": "w"
	}`

	testJSONMarshal(t, u, want)
}

func TestInstallation_Marshal(t *testing.T) {
	t.Parallel()
	testJSONMarshal(t, &Installation{}, "{}")

	u := &Installation{
		ID:       Ptr(int64(1)),
		NodeID:   Ptr("nid"),
		AppID:    Ptr(int64(1)),
		AppSlug:  Ptr("as"),
		TargetID: Ptr(int64(1)),
		Account: &User{
			Login:           Ptr("l"),
			ID:              Ptr(int64(1)),
			URL:             Ptr("u"),
			AvatarURL:       Ptr("a"),
			GravatarID:      Ptr("g"),
			Name:            Ptr("n"),
			Company:         Ptr("c"),
			Blog:            Ptr("b"),
			Location:        Ptr("l"),
			Email:           Ptr("e"),
			Hireable:        Ptr(true),
			Bio:             Ptr("b"),
			TwitterUsername: Ptr("t"),
			PublicRepos:     Ptr(1),
			Followers:       Ptr(1),
			Following:       Ptr(1),
			CreatedAt:       &Timestamp{referenceTime},
			SuspendedAt:     &Timestamp{referenceTime},
		},
		AccessTokensURL:     Ptr("atu"),
		RepositoriesURL:     Ptr("ru"),
		HTMLURL:             Ptr("hu"),
		TargetType:          Ptr("tt"),
		SingleFileName:      Ptr("sfn"),
		RepositorySelection: Ptr("rs"),
		Events:              []string{"e"},
		SingleFilePaths:     []string{"s"},
		Permissions: &InstallationPermissions{
			Actions:                       Ptr("a"),
			ActionsVariables:              Ptr("ac"),
			Administration:                Ptr("ad"),
			Checks:                        Ptr("c"),
			Contents:                      Ptr("co"),
			ContentReferences:             Ptr("cr"),
			Deployments:                   Ptr("d"),
			Environments:                  Ptr("e"),
			Issues:                        Ptr("i"),
			Metadata:                      Ptr("md"),
			Members:                       Ptr("m"),
			OrganizationAdministration:    Ptr("oa"),
			OrganizationCustomOrgRoles:    Ptr("ocr"),
			OrganizationHooks:             Ptr("oh"),
			OrganizationPlan:              Ptr("op"),
			OrganizationPreReceiveHooks:   Ptr("opr"),
			OrganizationProjects:          Ptr("op"),
			OrganizationSecrets:           Ptr("os"),
			OrganizationSelfHostedRunners: Ptr("osh"),
			OrganizationUserBlocking:      Ptr("oub"),
			Packages:                      Ptr("pkg"),
			Pages:                         Ptr("pg"),
			PullRequests:                  Ptr("pr"),
			RepositoryHooks:               Ptr("rh"),
			RepositoryProjects:            Ptr("rp"),
			RepositoryPreReceiveHooks:     Ptr("rprh"),
			Secrets:                       Ptr("s"),
			SecretScanningAlerts:          Ptr("ssa"),
			SecurityEvents:                Ptr("se"),
			SingleFile:                    Ptr("sf"),
			Statuses:                      Ptr("s"),
			TeamDiscussions:               Ptr("td"),
			VulnerabilityAlerts:           Ptr("va"),
			Workflows:                     Ptr("w"),
		},
		CreatedAt:              &Timestamp{referenceTime},
		UpdatedAt:              &Timestamp{referenceTime},
		HasMultipleSingleFiles: Ptr(false),
		SuspendedBy: &User{
			Login:           Ptr("l"),
			ID:              Ptr(int64(1)),
			URL:             Ptr("u"),
			AvatarURL:       Ptr("a"),
			GravatarID:      Ptr("g"),
			Name:            Ptr("n"),
			Company:         Ptr("c"),
			Blog:            Ptr("b"),
			Location:        Ptr("l"),
			Email:           Ptr("e"),
			Hireable:        Ptr(true),
			Bio:             Ptr("b"),
			TwitterUsername: Ptr("t"),
			PublicRepos:     Ptr(1),
			Followers:       Ptr(1),
			Following:       Ptr(1),
			CreatedAt:       &Timestamp{referenceTime},
			SuspendedAt:     &Timestamp{referenceTime},
		},
		SuspendedAt: &Timestamp{referenceTime},
	}

	want := `{
		"id": 1,
		"node_id": "nid",
		"app_id": 1,
		"app_slug": "as",
		"target_id": 1,
		"account": {
			"login": "l",
			"id": 1,
			"avatar_url": "a",
			"gravatar_id": "g",
			"name": "n",
			"company": "c",
			"blog": "b",
			"location": "l",
			"email": "e",
			"hireable": true,
			"bio": "b",
			"twitter_username": "t",
			"public_repos": 1,
			"followers": 1,
			"following": 1,
			"created_at": ` + referenceTimeStr + `,
			"suspended_at": ` + referenceTimeStr + `,
			"url": "u"
		},
		"access_tokens_url": "atu",
		"repositories_url": "ru",
		"html_url": "hu",
		"target_type": "tt",
		"single_file_name": "sfn",
		"repository_selection": "rs",
		"events": [
			"e"
		],
		"single_file_paths": [
			"s"
		],
		"permissions": {
			"actions": "a",
			"actions_variables": "ac",
			"administration": "ad",
			"checks": "c",
			"contents": "co",
			"content_references": "cr",
			"deployments": "d",
			"environments": "e",
			"issues": "i",
			"metadata": "md",
			"members": "m",
			"organization_administration": "oa",
			"organization_custom_org_roles": "ocr",
			"organization_hooks": "oh",
			"organization_plan": "op",
			"organization_pre_receive_hooks": "opr",
			"organization_projects": "op",
			"organization_secrets": "os",
			"organization_self_hosted_runners": "osh",
			"organization_user_blocking": "oub",
			"packages": "pkg",
			"pages": "pg",
			"pull_requests": "pr",
			"repository_hooks": "rh",
			"repository_projects": "rp",
			"repository_pre_receive_hooks": "rprh",
			"secrets": "s",
			"secret_scanning_alerts": "ssa",
			"security_events": "se",
			"single_file": "sf",
			"statuses": "s",
			"team_discussions": "td",
			"vulnerability_alerts": "va",
			"workflows": "w"
		},
		"created_at": ` + referenceTimeStr + `,
		"updated_at": ` + referenceTimeStr + `,
		"has_multiple_single_files": false,
		"suspended_by": {
			"login": "l",
			"id": 1,
			"avatar_url": "a",
			"gravatar_id": "g",
			"name": "n",
			"company": "c",
			"blog": "b",
			"location": "l",
			"email": "e",
			"hireable": true,
			"bio": "b",
			"twitter_username": "t",
			"public_repos": 1,
			"followers": 1,
			"following": 1,
			"created_at": ` + referenceTimeStr + `,
			"suspended_at": ` + referenceTimeStr + `,
			"url": "u"
		},
		"suspended_at": ` + referenceTimeStr + `
	}`

	testJSONMarshal(t, u, want)
}

func TestInstallationTokenOptions_Marshal(t *testing.T) {
	t.Parallel()
	testJSONMarshal(t, &InstallationTokenOptions{}, "{}")

	u := &InstallationTokenOptions{
		RepositoryIDs: []int64{1},
		Permissions: &InstallationPermissions{
			Actions:                       Ptr("a"),
			ActionsVariables:              Ptr("ac"),
			Administration:                Ptr("ad"),
			Checks:                        Ptr("c"),
			Contents:                      Ptr("co"),
			ContentReferences:             Ptr("cr"),
			Deployments:                   Ptr("d"),
			Environments:                  Ptr("e"),
			Issues:                        Ptr("i"),
			Metadata:                      Ptr("md"),
			Members:                       Ptr("m"),
			OrganizationAdministration:    Ptr("oa"),
			OrganizationCustomOrgRoles:    Ptr("ocr"),
			OrganizationHooks:             Ptr("oh"),
			OrganizationPlan:              Ptr("op"),
			OrganizationPreReceiveHooks:   Ptr("opr"),
			OrganizationProjects:          Ptr("op"),
			OrganizationSecrets:           Ptr("os"),
			OrganizationSelfHostedRunners: Ptr("osh"),
			OrganizationUserBlocking:      Ptr("oub"),
			Packages:                      Ptr("pkg"),
			Pages:                         Ptr("pg"),
			PullRequests:                  Ptr("pr"),
			RepositoryHooks:               Ptr("rh"),
			RepositoryProjects:            Ptr("rp"),
			RepositoryPreReceiveHooks:     Ptr("rprh"),
			Secrets:                       Ptr("s"),
			SecretScanningAlerts:          Ptr("ssa"),
			SecurityEvents:                Ptr("se"),
			SingleFile:                    Ptr("sf"),
			Statuses:                      Ptr("s"),
			TeamDiscussions:               Ptr("td"),
			VulnerabilityAlerts:           Ptr("va"),
			Workflows:                     Ptr("w"),
		},
	}

	want := `{
		"repository_ids": [1],
		"permissions": {
			"actions": "a",
			"actions_variables": "ac",
			"administration": "ad",
			"checks": "c",
			"contents": "co",
			"content_references": "cr",
			"deployments": "d",
			"environments": "e",
			"issues": "i",
			"metadata": "md",
			"members": "m",
			"organization_administration": "oa",
			"organization_custom_org_roles": "ocr",
			"organization_hooks": "oh",
			"organization_plan": "op",
			"organization_pre_receive_hooks": "opr",
			"organization_projects": "op",
			"organization_secrets": "os",
			"organization_self_hosted_runners": "osh",
			"organization_user_blocking": "oub",
			"packages": "pkg",
			"pages": "pg",
			"pull_requests": "pr",
			"repository_hooks": "rh",
			"repository_projects": "rp",
			"repository_pre_receive_hooks": "rprh",
			"secrets": "s",
			"secret_scanning_alerts": "ssa",
			"security_events": "se",
			"single_file": "sf",
			"statuses": "s",
			"team_discussions": "td",
			"vulnerability_alerts": "va",
			"workflows": "w"
		}
	}`

	testJSONMarshal(t, u, want)
}

func TestInstallationToken_Marshal(t *testing.T) {
	t.Parallel()
	testJSONMarshal(t, &InstallationToken{}, "{}")

	u := &InstallationToken{
		Token:     Ptr("t"),
		ExpiresAt: &Timestamp{referenceTime},
		Permissions: &InstallationPermissions{
			Actions:                       Ptr("a"),
			ActionsVariables:              Ptr("ac"),
			Administration:                Ptr("ad"),
			Checks:                        Ptr("c"),
			Contents:                      Ptr("co"),
			ContentReferences:             Ptr("cr"),
			Deployments:                   Ptr("d"),
			Environments:                  Ptr("e"),
			Issues:                        Ptr("i"),
			Metadata:                      Ptr("md"),
			Members:                       Ptr("m"),
			OrganizationAdministration:    Ptr("oa"),
			OrganizationCustomOrgRoles:    Ptr("ocr"),
			OrganizationHooks:             Ptr("oh"),
			OrganizationPlan:              Ptr("op"),
			OrganizationPreReceiveHooks:   Ptr("opr"),
			OrganizationProjects:          Ptr("op"),
			OrganizationSecrets:           Ptr("os"),
			OrganizationSelfHostedRunners: Ptr("osh"),
			OrganizationUserBlocking:      Ptr("oub"),
			Packages:                      Ptr("pkg"),
			Pages:                         Ptr("pg"),
			PullRequests:                  Ptr("pr"),
			RepositoryHooks:               Ptr("rh"),
			RepositoryProjects:            Ptr("rp"),
			RepositoryPreReceiveHooks:     Ptr("rprh"),
			Secrets:                       Ptr("s"),
			SecretScanningAlerts:          Ptr("ssa"),
			SecurityEvents:                Ptr("se"),
			SingleFile:                    Ptr("sf"),
			Statuses:                      Ptr("s"),
			TeamDiscussions:               Ptr("td"),
			VulnerabilityAlerts:           Ptr("va"),
			Workflows:                     Ptr("w"),
		},
		Repositories: []*Repository{
			{
				ID:   Ptr(int64(1)),
				URL:  Ptr("u"),
				Name: Ptr("n"),
			},
		},
	}

	want := `{
		"token": "t",
		"expires_at": ` + referenceTimeStr + `,
		"permissions": {
			"actions": "a",
			"actions_variables": "ac",
			"administration": "ad",
			"checks": "c",
			"contents": "co",
			"content_references": "cr",
			"deployments": "d",
			"environments": "e",
			"issues": "i",
			"metadata": "md",
			"members": "m",
			"organization_administration": "oa",
			"organization_custom_org_roles": "ocr",
			"organization_hooks": "oh",
			"organization_plan": "op",
			"organization_pre_receive_hooks": "opr",
			"organization_projects": "op",
			"organization_secrets": "os",
			"organization_self_hosted_runners": "osh",
			"organization_user_blocking": "oub",
			"packages": "pkg",
			"pages": "pg",
			"pull_requests": "pr",
			"repository_hooks": "rh",
			"repository_projects": "rp",
			"repository_pre_receive_hooks": "rprh",
			"secrets": "s",
			"secret_scanning_alerts": "ssa",
			"security_events": "se",
			"single_file": "sf",
			"statuses": "s",
			"team_discussions": "td",
			"vulnerability_alerts": "va",
			"workflows": "w"
		},
		"repositories": [
			{
				"id": 1,
				"url": "u",
				"name": "n"
			}
		]
	}`

	testJSONMarshal(t, u, want)
}

func TestApp_Marshal(t *testing.T) {
	t.Parallel()
	testJSONMarshal(t, &App{}, "{}")

	u := &App{
		ID:     Ptr(int64(1)),
		Slug:   Ptr("s"),
		NodeID: Ptr("nid"),
		Owner: &User{
			Login:           Ptr("l"),
			ID:              Ptr(int64(1)),
			URL:             Ptr("u"),
			AvatarURL:       Ptr("a"),
			GravatarID:      Ptr("g"),
			Name:            Ptr("n"),
			Company:         Ptr("c"),
			Blog:            Ptr("b"),
			Location:        Ptr("l"),
			Email:           Ptr("e"),
			Hireable:        Ptr(true),
			Bio:             Ptr("b"),
			TwitterUsername: Ptr("t"),
			PublicRepos:     Ptr(1),
			Followers:       Ptr(1),
			Following:       Ptr(1),
			CreatedAt:       &Timestamp{referenceTime},
			SuspendedAt:     &Timestamp{referenceTime},
		},
		Name:        Ptr("n"),
		Description: Ptr("d"),
		ExternalURL: Ptr("eu"),
		HTMLURL:     Ptr("hu"),
		CreatedAt:   &Timestamp{referenceTime},
		UpdatedAt:   &Timestamp{referenceTime},
		Permissions: &InstallationPermissions{
			Actions:                       Ptr("a"),
			ActionsVariables:              Ptr("ac"),
			Administration:                Ptr("ad"),
			Checks:                        Ptr("c"),
			Contents:                      Ptr("co"),
			ContentReferences:             Ptr("cr"),
			Deployments:                   Ptr("d"),
			Environments:                  Ptr("e"),
			Issues:                        Ptr("i"),
			Metadata:                      Ptr("md"),
			Members:                       Ptr("m"),
			OrganizationAdministration:    Ptr("oa"),
			OrganizationCustomOrgRoles:    Ptr("ocr"),
			OrganizationHooks:             Ptr("oh"),
			OrganizationPlan:              Ptr("op"),
			OrganizationPreReceiveHooks:   Ptr("opr"),
			OrganizationProjects:          Ptr("op"),
			OrganizationSecrets:           Ptr("os"),
			OrganizationSelfHostedRunners: Ptr("osh"),
			OrganizationUserBlocking:      Ptr("oub"),
			Packages:                      Ptr("pkg"),
			Pages:                         Ptr("pg"),
			PullRequests:                  Ptr("pr"),
			RepositoryHooks:               Ptr("rh"),
			RepositoryProjects:            Ptr("rp"),
			RepositoryPreReceiveHooks:     Ptr("rprh"),
			Secrets:                       Ptr("s"),
			SecretScanningAlerts:          Ptr("ssa"),
			SecurityEvents:                Ptr("se"),
			SingleFile:                    Ptr("sf"),
			Statuses:                      Ptr("s"),
			TeamDiscussions:               Ptr("td"),
			VulnerabilityAlerts:           Ptr("va"),
			Workflows:                     Ptr("w"),
		},
		Events: []string{"s"},
	}

	want := `{
		"id": 1,
		"slug": "s",
		"node_id": "nid",
		"owner": {
			"login": "l",
			"id": 1,
			"avatar_url": "a",
			"gravatar_id": "g",
			"name": "n",
			"company": "c",
			"blog": "b",
			"location": "l",
			"email": "e",
			"hireable": true,
			"bio": "b",
			"twitter_username": "t",
			"public_repos": 1,
			"followers": 1,
			"following": 1,
			"created_at": ` + referenceTimeStr + `,
			"suspended_at": ` + referenceTimeStr + `,
			"url": "u"
		},
		"name": "n",
		"description": "d",
		"external_url": "eu",
		"html_url": "hu",
		"created_at": ` + referenceTimeStr + `,
		"updated_at": ` + referenceTimeStr + `,
		"permissions": {
			"actions": "a",
			"actions_variables": "ac",
			"administration": "ad",
			"checks": "c",
			"contents": "co",
			"content_references": "cr",
			"deployments": "d",
			"environments": "e",
			"issues": "i",
			"metadata": "md",
			"members": "m",
			"organization_administration": "oa",
			"organization_custom_org_roles": "ocr",
			"organization_hooks": "oh",
			"organization_plan": "op",
			"organization_pre_receive_hooks": "opr",
			"organization_projects": "op",
			"organization_secrets": "os",
			"organization_self_hosted_runners": "osh",
			"organization_user_blocking": "oub",
			"packages": "pkg",
			"pages": "pg",
			"pull_requests": "pr",
			"repository_hooks": "rh",
			"repository_projects": "rp",
			"repository_pre_receive_hooks": "rprh",
			"secrets": "s",
			"secret_scanning_alerts": "ssa",
			"security_events": "se",
			"single_file": "sf",
			"statuses": "s",
			"team_discussions": "td",
			"vulnerability_alerts": "va",
			"workflows": "w"
		},
		"events": ["s"]
	}`

	testJSONMarshal(t, u, want)
}
