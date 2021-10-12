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
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/app", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{"id":1}`)
	})

	ctx := context.Background()
	app, _, err := client.Apps.Get(ctx, "")
	if err != nil {
		t.Errorf("Apps.Get returned error: %v", err)
	}

	want := &App{ID: Int64(1)}
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
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/apps/a", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{"html_url":"https://github.com/apps/a"}`)
	})

	ctx := context.Background()
	app, _, err := client.Apps.Get(ctx, "a")
	if err != nil {
		t.Errorf("Apps.Get returned error: %v", err)
	}

	want := &App{HTMLURL: String("https://github.com/apps/a")}
	if !cmp.Equal(app, want) {
		t.Errorf("Apps.Get returned %+v, want %+v", *app.HTMLURL, *want.HTMLURL)
	}
}

func TestAppsService_ListInstallations(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

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
                                       "organization_hooks": "write",
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
		ID:                  Int64(1),
		AppID:               Int64(1),
		TargetID:            Int64(1),
		TargetType:          String("Organization"),
		SingleFileName:      String("config.yml"),
		RepositorySelection: String("selected"),
		Permissions: &InstallationPermissions{
			Actions:                       String("read"),
			Administration:                String("read"),
			Checks:                        String("read"),
			Contents:                      String("read"),
			ContentReferences:             String("read"),
			Deployments:                   String("read"),
			Environments:                  String("read"),
			Issues:                        String("write"),
			Metadata:                      String("read"),
			Members:                       String("read"),
			OrganizationAdministration:    String("write"),
			OrganizationHooks:             String("write"),
			OrganizationPlan:              String("read"),
			OrganizationPreReceiveHooks:   String("write"),
			OrganizationProjects:          String("read"),
			OrganizationSecrets:           String("read"),
			OrganizationSelfHostedRunners: String("read"),
			OrganizationUserBlocking:      String("write"),
			Packages:                      String("read"),
			Pages:                         String("read"),
			PullRequests:                  String("write"),
			RepositoryHooks:               String("write"),
			RepositoryProjects:            String("read"),
			RepositoryPreReceiveHooks:     String("read"),
			Secrets:                       String("read"),
			SecretScanningAlerts:          String("read"),
			SecurityEvents:                String("read"),
			SingleFile:                    String("write"),
			Statuses:                      String("write"),
			TeamDiscussions:               String("read"),
			VulnerabilityAlerts:           String("read"),
			Workflows:                     String("write")},
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
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/app/installations/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{"id":1, "app_id":1, "target_id":1, "target_type": "Organization"}`)
	})

	ctx := context.Background()
	installation, _, err := client.Apps.GetInstallation(ctx, 1)
	if err != nil {
		t.Errorf("Apps.GetInstallation returned error: %v", err)
	}

	want := &Installation{ID: Int64(1), AppID: Int64(1), TargetID: Int64(1), TargetType: String("Organization")}
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
	client, mux, _, teardown := setup()
	defer teardown()

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

	want := []*Installation{{ID: Int64(1), AppID: Int64(1), TargetID: Int64(1), TargetType: String("Organization")}}
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
	client, mux, _, teardown := setup()
	defer teardown()

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
	client, mux, _, teardown := setup()
	defer teardown()

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
	client, mux, _, teardown := setup()
	defer teardown()

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
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/app/installations/1/access_tokens", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		fmt.Fprint(w, `{"token":"t"}`)
	})

	ctx := context.Background()
	token, _, err := client.Apps.CreateInstallationToken(ctx, 1, nil)
	if err != nil {
		t.Errorf("Apps.CreateInstallationToken returned error: %v", err)
	}

	want := &InstallationToken{Token: String("t")}
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
	client, mux, _, teardown := setup()
	defer teardown()

	installationTokenOptions := &InstallationTokenOptions{
		RepositoryIDs: []int64{1234},
		Repositories:  []string{"foo"},
		Permissions: &InstallationPermissions{
			Contents: String("write"),
			Issues:   String("read"),
		},
	}

	mux.HandleFunc("/app/installations/1/access_tokens", func(w http.ResponseWriter, r *http.Request) {
		v := new(InstallationTokenOptions)
		json.NewDecoder(r.Body).Decode(v)

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

	want := &InstallationToken{Token: String("t")}
	if !cmp.Equal(token, want) {
		t.Errorf("Apps.CreateInstallationToken returned %+v, want %+v", token, want)
	}
}

func TestAppsService_CreateAttachement(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/content_references/11/attachments", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		testHeader(t, r, "Accept", mediaTypeContentAttachmentsPreview)

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"id":1,"title":"title1","body":"body1"}`))
	})

	ctx := context.Background()
	got, _, err := client.Apps.CreateAttachment(ctx, 11, "title1", "body1")
	if err != nil {
		t.Errorf("CreateAttachment returned error: %v", err)
	}

	want := &Attachment{ID: Int64(1), Title: String("title1"), Body: String("body1")}
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
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/orgs/o/installation", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{"id":1, "app_id":1, "target_id":1, "target_type": "Organization"}`)
	})

	ctx := context.Background()
	installation, _, err := client.Apps.FindOrganizationInstallation(ctx, "o")
	if err != nil {
		t.Errorf("Apps.FindOrganizationInstallation returned error: %v", err)
	}

	want := &Installation{ID: Int64(1), AppID: Int64(1), TargetID: Int64(1), TargetType: String("Organization")}
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
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/repos/o/r/installation", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{"id":1, "app_id":1, "target_id":1, "target_type": "Organization"}`)
	})

	ctx := context.Background()
	installation, _, err := client.Apps.FindRepositoryInstallation(ctx, "o", "r")
	if err != nil {
		t.Errorf("Apps.FindRepositoryInstallation returned error: %v", err)
	}

	want := &Installation{ID: Int64(1), AppID: Int64(1), TargetID: Int64(1), TargetType: String("Organization")}
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
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/repositories/1/installation", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{"id":1, "app_id":1, "target_id":1, "target_type": "Organization"}`)
	})

	ctx := context.Background()
	installation, _, err := client.Apps.FindRepositoryInstallationByID(ctx, 1)
	if err != nil {
		t.Errorf("Apps.FindRepositoryInstallationByID returned error: %v", err)
	}

	want := &Installation{ID: Int64(1), AppID: Int64(1), TargetID: Int64(1), TargetType: String("Organization")}
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
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/users/u/installation", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{"id":1, "app_id":1, "target_id":1, "target_type": "User"}`)
	})

	ctx := context.Background()
	installation, _, err := client.Apps.FindUserInstallation(ctx, "u")
	if err != nil {
		t.Errorf("Apps.FindUserInstallation returned error: %v", err)
	}

	want := &Installation{ID: Int64(1), AppID: Int64(1), TargetID: Int64(1), TargetType: String("User")}
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
	testJSONMarshal(t, &ContentReference{}, "{}")

	u := &ContentReference{
		ID:        Int64(1),
		NodeID:    String("nid"),
		Reference: String("r"),
	}

	want := `{
		"id": 1,
		"node_id": "nid",
		"reference": "r"
	}`

	testJSONMarshal(t, u, want)
}

func TestAttachment_Marshal(t *testing.T) {
	testJSONMarshal(t, &Attachment{}, "{}")

	u := &Attachment{
		ID:    Int64(1),
		Title: String("t"),
		Body:  String("b"),
	}

	want := `{
		"id": 1,
		"title": "t",
		"body": "b"
	}`

	testJSONMarshal(t, u, want)
}

func TestInstallationPermissions_Marshal(t *testing.T) {
	testJSONMarshal(t, &InstallationPermissions{}, "{}")

	u := &InstallationPermissions{
		Actions:                       String("a"),
		Administration:                String("ad"),
		Checks:                        String("c"),
		Contents:                      String("co"),
		ContentReferences:             String("cr"),
		Deployments:                   String("d"),
		Environments:                  String("e"),
		Issues:                        String("i"),
		Metadata:                      String("md"),
		Members:                       String("m"),
		OrganizationAdministration:    String("oa"),
		OrganizationHooks:             String("oh"),
		OrganizationPlan:              String("op"),
		OrganizationPreReceiveHooks:   String("opr"),
		OrganizationProjects:          String("op"),
		OrganizationSecrets:           String("os"),
		OrganizationSelfHostedRunners: String("osh"),
		OrganizationUserBlocking:      String("oub"),
		Packages:                      String("pkg"),
		Pages:                         String("pg"),
		PullRequests:                  String("pr"),
		RepositoryHooks:               String("rh"),
		RepositoryProjects:            String("rp"),
		RepositoryPreReceiveHooks:     String("rprh"),
		Secrets:                       String("s"),
		SecretScanningAlerts:          String("ssa"),
		SecurityEvents:                String("se"),
		SingleFile:                    String("sf"),
		Statuses:                      String("s"),
		TeamDiscussions:               String("td"),
		VulnerabilityAlerts:           String("va"),
		Workflows:                     String("w"),
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
	testJSONMarshal(t, &Installation{}, "{}")

	u := &Installation{
		ID:       Int64(1),
		NodeID:   String("nid"),
		AppID:    Int64(1),
		AppSlug:  String("as"),
		TargetID: Int64(1),
		Account: &User{
			Login:           String("l"),
			ID:              Int64(1),
			URL:             String("u"),
			AvatarURL:       String("a"),
			GravatarID:      String("g"),
			Name:            String("n"),
			Company:         String("c"),
			Blog:            String("b"),
			Location:        String("l"),
			Email:           String("e"),
			Hireable:        Bool(true),
			Bio:             String("b"),
			TwitterUsername: String("t"),
			PublicRepos:     Int(1),
			Followers:       Int(1),
			Following:       Int(1),
			CreatedAt:       &Timestamp{referenceTime},
			SuspendedAt:     &Timestamp{referenceTime},
		},
		AccessTokensURL:     String("atu"),
		RepositoriesURL:     String("ru"),
		HTMLURL:             String("hu"),
		TargetType:          String("tt"),
		SingleFileName:      String("sfn"),
		RepositorySelection: String("rs"),
		Events:              []string{"e"},
		SingleFilePaths:     []string{"s"},
		Permissions: &InstallationPermissions{
			Actions:                       String("a"),
			Administration:                String("ad"),
			Checks:                        String("c"),
			Contents:                      String("co"),
			ContentReferences:             String("cr"),
			Deployments:                   String("d"),
			Environments:                  String("e"),
			Issues:                        String("i"),
			Metadata:                      String("md"),
			Members:                       String("m"),
			OrganizationAdministration:    String("oa"),
			OrganizationHooks:             String("oh"),
			OrganizationPlan:              String("op"),
			OrganizationPreReceiveHooks:   String("opr"),
			OrganizationProjects:          String("op"),
			OrganizationSecrets:           String("os"),
			OrganizationSelfHostedRunners: String("osh"),
			OrganizationUserBlocking:      String("oub"),
			Packages:                      String("pkg"),
			Pages:                         String("pg"),
			PullRequests:                  String("pr"),
			RepositoryHooks:               String("rh"),
			RepositoryProjects:            String("rp"),
			RepositoryPreReceiveHooks:     String("rprh"),
			Secrets:                       String("s"),
			SecretScanningAlerts:          String("ssa"),
			SecurityEvents:                String("se"),
			SingleFile:                    String("sf"),
			Statuses:                      String("s"),
			TeamDiscussions:               String("td"),
			VulnerabilityAlerts:           String("va"),
			Workflows:                     String("w"),
		},
		CreatedAt:              &Timestamp{referenceTime},
		UpdatedAt:              &Timestamp{referenceTime},
		HasMultipleSingleFiles: Bool(false),
		SuspendedBy: &User{
			Login:           String("l"),
			ID:              Int64(1),
			URL:             String("u"),
			AvatarURL:       String("a"),
			GravatarID:      String("g"),
			Name:            String("n"),
			Company:         String("c"),
			Blog:            String("b"),
			Location:        String("l"),
			Email:           String("e"),
			Hireable:        Bool(true),
			Bio:             String("b"),
			TwitterUsername: String("t"),
			PublicRepos:     Int(1),
			Followers:       Int(1),
			Following:       Int(1),
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
	testJSONMarshal(t, &InstallationTokenOptions{}, "{}")

	u := &InstallationTokenOptions{
		RepositoryIDs: []int64{1},
		Permissions: &InstallationPermissions{
			Actions:                       String("a"),
			Administration:                String("ad"),
			Checks:                        String("c"),
			Contents:                      String("co"),
			ContentReferences:             String("cr"),
			Deployments:                   String("d"),
			Environments:                  String("e"),
			Issues:                        String("i"),
			Metadata:                      String("md"),
			Members:                       String("m"),
			OrganizationAdministration:    String("oa"),
			OrganizationHooks:             String("oh"),
			OrganizationPlan:              String("op"),
			OrganizationPreReceiveHooks:   String("opr"),
			OrganizationProjects:          String("op"),
			OrganizationSecrets:           String("os"),
			OrganizationSelfHostedRunners: String("osh"),
			OrganizationUserBlocking:      String("oub"),
			Packages:                      String("pkg"),
			Pages:                         String("pg"),
			PullRequests:                  String("pr"),
			RepositoryHooks:               String("rh"),
			RepositoryProjects:            String("rp"),
			RepositoryPreReceiveHooks:     String("rprh"),
			Secrets:                       String("s"),
			SecretScanningAlerts:          String("ssa"),
			SecurityEvents:                String("se"),
			SingleFile:                    String("sf"),
			Statuses:                      String("s"),
			TeamDiscussions:               String("td"),
			VulnerabilityAlerts:           String("va"),
			Workflows:                     String("w"),
		},
	}

	want := `{
		"repository_ids": [1],
		"permissions": {
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
	testJSONMarshal(t, &InstallationToken{}, "{}")

	u := &InstallationToken{
		Token:     String("t"),
		ExpiresAt: &referenceTime,
		Permissions: &InstallationPermissions{
			Actions:                       String("a"),
			Administration:                String("ad"),
			Checks:                        String("c"),
			Contents:                      String("co"),
			ContentReferences:             String("cr"),
			Deployments:                   String("d"),
			Environments:                  String("e"),
			Issues:                        String("i"),
			Metadata:                      String("md"),
			Members:                       String("m"),
			OrganizationAdministration:    String("oa"),
			OrganizationHooks:             String("oh"),
			OrganizationPlan:              String("op"),
			OrganizationPreReceiveHooks:   String("opr"),
			OrganizationProjects:          String("op"),
			OrganizationSecrets:           String("os"),
			OrganizationSelfHostedRunners: String("osh"),
			OrganizationUserBlocking:      String("oub"),
			Packages:                      String("pkg"),
			Pages:                         String("pg"),
			PullRequests:                  String("pr"),
			RepositoryHooks:               String("rh"),
			RepositoryProjects:            String("rp"),
			RepositoryPreReceiveHooks:     String("rprh"),
			Secrets:                       String("s"),
			SecretScanningAlerts:          String("ssa"),
			SecurityEvents:                String("se"),
			SingleFile:                    String("sf"),
			Statuses:                      String("s"),
			TeamDiscussions:               String("td"),
			VulnerabilityAlerts:           String("va"),
			Workflows:                     String("w"),
		},
		Repositories: []*Repository{
			{
				ID:   Int64(1),
				URL:  String("u"),
				Name: String("n"),
			},
		},
	}

	want := `{
		"token": "t",
		"expires_at": ` + referenceTimeStr + `,
		"permissions": {
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
	testJSONMarshal(t, &App{}, "{}")

	u := &App{
		ID:     Int64(1),
		Slug:   String("s"),
		NodeID: String("nid"),
		Owner: &User{
			Login:           String("l"),
			ID:              Int64(1),
			URL:             String("u"),
			AvatarURL:       String("a"),
			GravatarID:      String("g"),
			Name:            String("n"),
			Company:         String("c"),
			Blog:            String("b"),
			Location:        String("l"),
			Email:           String("e"),
			Hireable:        Bool(true),
			Bio:             String("b"),
			TwitterUsername: String("t"),
			PublicRepos:     Int(1),
			Followers:       Int(1),
			Following:       Int(1),
			CreatedAt:       &Timestamp{referenceTime},
			SuspendedAt:     &Timestamp{referenceTime},
		},
		Name:        String("n"),
		Description: String("d"),
		ExternalURL: String("eu"),
		HTMLURL:     String("hu"),
		CreatedAt:   &Timestamp{referenceTime},
		UpdatedAt:   &Timestamp{referenceTime},
		Permissions: &InstallationPermissions{
			Actions:                       String("a"),
			Administration:                String("ad"),
			Checks:                        String("c"),
			Contents:                      String("co"),
			ContentReferences:             String("cr"),
			Deployments:                   String("d"),
			Environments:                  String("e"),
			Issues:                        String("i"),
			Metadata:                      String("md"),
			Members:                       String("m"),
			OrganizationAdministration:    String("oa"),
			OrganizationHooks:             String("oh"),
			OrganizationPlan:              String("op"),
			OrganizationPreReceiveHooks:   String("opr"),
			OrganizationProjects:          String("op"),
			OrganizationSecrets:           String("os"),
			OrganizationSelfHostedRunners: String("osh"),
			OrganizationUserBlocking:      String("oub"),
			Packages:                      String("pkg"),
			Pages:                         String("pg"),
			PullRequests:                  String("pr"),
			RepositoryHooks:               String("rh"),
			RepositoryProjects:            String("rp"),
			RepositoryPreReceiveHooks:     String("rprh"),
			Secrets:                       String("s"),
			SecretScanningAlerts:          String("ssa"),
			SecurityEvents:                String("se"),
			SingleFile:                    String("sf"),
			Statuses:                      String("s"),
			TeamDiscussions:               String("td"),
			VulnerabilityAlerts:           String("va"),
			Workflows:                     String("w"),
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
