// Copyright 2022 The go-github AUTHORS. All rights reserved.
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

func TestSecretScanningService_ListAlertsForEnterprise(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/enterprises/e/secret-scanning/alerts", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testFormValues(t, r, values{"state": "open", "secret_type": "mailchimp_api_key"})

		fmt.Fprint(w, `[{
			"number": 1,
			"created_at": "1996-06-20T00:00:00Z",
			"url": "https://api.github.com/repos/o/r/secret-scanning/alerts/1",
			"html_url": "https://github.com/o/r/security/secret-scanning/1",
			"locations_url": "https://api.github.com/repos/o/r/secret-scanning/alerts/1/locations",
			"state": "open",
			"resolution": null,
			"resolved_at": null,
			"resolved_by": null,
			"secret_type": "mailchimp_api_key",
			"secret": "XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX-us2",
			"repository": {
				"id": 1,
				"name": "n",
				"url": "url"
			}
		}]`)
	})

	ctx := context.Background()
	opts := &SecretScanningAlertListOptions{State: "open", SecretType: "mailchimp_api_key"}

	alerts, _, err := client.SecretScanning.ListAlertsForEnterprise(ctx, "e", opts)
	if err != nil {
		t.Errorf("SecretScanning.ListAlertsForEnterprise returned error: %v", err)
	}

	date := Timestamp{time.Date(1996, time.June, 20, 00, 00, 00, 0, time.UTC)}
	want := []*SecretScanningAlert{
		{
			Number:       Ptr(1),
			CreatedAt:    &date,
			URL:          Ptr("https://api.github.com/repos/o/r/secret-scanning/alerts/1"),
			HTMLURL:      Ptr("https://github.com/o/r/security/secret-scanning/1"),
			LocationsURL: Ptr("https://api.github.com/repos/o/r/secret-scanning/alerts/1/locations"),
			State:        Ptr("open"),
			Resolution:   nil,
			ResolvedAt:   nil,
			ResolvedBy:   nil,
			SecretType:   Ptr("mailchimp_api_key"),
			Secret:       Ptr("XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX-us2"),
			Repository: &Repository{
				ID:   Ptr(int64(1)),
				URL:  Ptr("url"),
				Name: Ptr("n"),
			},
		},
	}

	if !cmp.Equal(alerts, want) {
		t.Errorf("SecretScanning.ListAlertsForEnterprise returned %+v, want %+v", alerts, want)
	}

	const methodName = "ListAlertsForEnterprise"

	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.SecretScanning.ListAlertsForEnterprise(ctx, "\n", opts)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		_, resp, err := client.SecretScanning.ListAlertsForEnterprise(ctx, "e", opts)
		return resp, err
	})
}

func TestSecretScanningService_ListAlertsForOrg(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/orgs/o/secret-scanning/alerts", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testFormValues(t, r, values{"state": "open", "secret_type": "mailchimp_api_key"})

		fmt.Fprint(w, `[{
			"number": 1,
			"created_at": "1996-06-20T00:00:00Z",
			"url": "https://api.github.com/repos/o/r/secret-scanning/alerts/1",
			"html_url": "https://github.com/o/r/security/secret-scanning/1",
			"locations_url": "https://api.github.com/repos/o/r/secret-scanning/alerts/1/locations",
			"state": "open",
			"resolution": null,
			"resolved_at": null,
			"resolved_by": null,
			"secret_type": "mailchimp_api_key",
			"secret": "XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX-us2"
		}]`)
	})

	ctx := context.Background()
	opts := &SecretScanningAlertListOptions{State: "open", SecretType: "mailchimp_api_key"}

	alerts, _, err := client.SecretScanning.ListAlertsForOrg(ctx, "o", opts)
	if err != nil {
		t.Errorf("SecretScanning.ListAlertsForOrg returned error: %v", err)
	}

	date := Timestamp{time.Date(1996, time.June, 20, 00, 00, 00, 0, time.UTC)}
	want := []*SecretScanningAlert{
		{
			Number:       Ptr(1),
			CreatedAt:    &date,
			URL:          Ptr("https://api.github.com/repos/o/r/secret-scanning/alerts/1"),
			HTMLURL:      Ptr("https://github.com/o/r/security/secret-scanning/1"),
			LocationsURL: Ptr("https://api.github.com/repos/o/r/secret-scanning/alerts/1/locations"),
			State:        Ptr("open"),
			Resolution:   nil,
			ResolvedAt:   nil,
			ResolvedBy:   nil,
			SecretType:   Ptr("mailchimp_api_key"),
			Secret:       Ptr("XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX-us2"),
		},
	}

	if !cmp.Equal(alerts, want) {
		t.Errorf("SecretScanning.ListAlertsForOrg returned %+v, want %+v", alerts, want)
	}

	const methodName = "ListAlertsForOrg"

	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.SecretScanning.ListAlertsForOrg(ctx, "\n", opts)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		_, resp, err := client.SecretScanning.ListAlertsForOrg(ctx, "o", opts)
		return resp, err
	})
}

func TestSecretScanningService_ListAlertsForOrgListOptions(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/orgs/o/secret-scanning/alerts", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testFormValues(t, r, values{"state": "open", "secret_type": "mailchimp_api_key", "per_page": "1", "page": "1"})

		fmt.Fprint(w, `[{
			"number": 1,
			"created_at": "1996-06-20T00:00:00Z",
			"url": "https://api.github.com/repos/o/r/secret-scanning/alerts/1",
			"html_url": "https://github.com/o/r/security/secret-scanning/1",
			"locations_url": "https://api.github.com/repos/o/r/secret-scanning/alerts/1/locations",
			"state": "open",
			"resolution": null,
			"resolved_at": null,
			"resolved_by": null,
			"secret_type": "mailchimp_api_key",
			"secret": "XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX-us2"
		}]`)
	})

	ctx := context.Background()

	// Testing pagination by index
	opts := &SecretScanningAlertListOptions{State: "open", SecretType: "mailchimp_api_key", ListOptions: ListOptions{Page: 1, PerPage: 1}}

	alerts, _, err := client.SecretScanning.ListAlertsForOrg(ctx, "o", opts)
	if err != nil {
		t.Errorf("SecretScanning.ListAlertsForOrg returned error: %v", err)
	}

	date := Timestamp{time.Date(1996, time.June, 20, 00, 00, 00, 0, time.UTC)}
	want := []*SecretScanningAlert{
		{
			Number:       Ptr(1),
			CreatedAt:    &date,
			URL:          Ptr("https://api.github.com/repos/o/r/secret-scanning/alerts/1"),
			HTMLURL:      Ptr("https://github.com/o/r/security/secret-scanning/1"),
			LocationsURL: Ptr("https://api.github.com/repos/o/r/secret-scanning/alerts/1/locations"),
			State:        Ptr("open"),
			Resolution:   nil,
			ResolvedAt:   nil,
			ResolvedBy:   nil,
			SecretType:   Ptr("mailchimp_api_key"),
			Secret:       Ptr("XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX-us2"),
		},
	}

	if !cmp.Equal(alerts, want) {
		t.Errorf("SecretScanning.ListAlertsForOrg returned %+v, want %+v", alerts, want)
	}

	const methodName = "ListAlertsForOrg"

	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.SecretScanning.ListAlertsForOrg(ctx, "\n", opts)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		_, resp, err := client.SecretScanning.ListAlertsForOrg(ctx, "o", opts)
		return resp, err
	})
}

func TestSecretScanningService_ListAlertsForRepo(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/repos/o/r/secret-scanning/alerts", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testFormValues(t, r, values{"state": "open", "secret_type": "mailchimp_api_key"})

		fmt.Fprint(w, `[{
			"number": 1,
			"created_at": "1996-06-20T00:00:00Z",
			"url": "https://api.github.com/repos/o/r/secret-scanning/alerts/1",
			"html_url": "https://github.com/o/r/security/secret-scanning/1",
			"locations_url": "https://api.github.com/repos/o/r/secret-scanning/alerts/1/locations",
			"state": "open",
			"resolution": null,
			"resolved_at": null,
			"resolved_by": null,
			"secret_type": "mailchimp_api_key",
			"secret": "XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX-us2"
		}]`)
	})

	ctx := context.Background()
	opts := &SecretScanningAlertListOptions{State: "open", SecretType: "mailchimp_api_key"}

	alerts, _, err := client.SecretScanning.ListAlertsForRepo(ctx, "o", "r", opts)
	if err != nil {
		t.Errorf("SecretScanning.ListAlertsForRepo returned error: %v", err)
	}

	date := Timestamp{time.Date(1996, time.June, 20, 00, 00, 00, 0, time.UTC)}
	want := []*SecretScanningAlert{
		{
			Number:       Ptr(1),
			CreatedAt:    &date,
			URL:          Ptr("https://api.github.com/repos/o/r/secret-scanning/alerts/1"),
			HTMLURL:      Ptr("https://github.com/o/r/security/secret-scanning/1"),
			LocationsURL: Ptr("https://api.github.com/repos/o/r/secret-scanning/alerts/1/locations"),
			State:        Ptr("open"),
			Resolution:   nil,
			ResolvedAt:   nil,
			ResolvedBy:   nil,
			SecretType:   Ptr("mailchimp_api_key"),
			Secret:       Ptr("XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX-us2"),
		},
	}

	if !cmp.Equal(alerts, want) {
		t.Errorf("SecretScanning.ListAlertsForRepo returned %+v, want %+v", alerts, want)
	}

	const methodName = "ListAlertsForRepo"

	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.SecretScanning.ListAlertsForRepo(ctx, "\n", "\n", opts)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		_, resp, err := client.SecretScanning.ListAlertsForRepo(ctx, "o", "r", opts)
		return resp, err
	})
}

func TestSecretScanningService_GetAlert(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/repos/o/r/secret-scanning/alerts/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")

		fmt.Fprint(w, `{
			"number": 1,
			"created_at": "1996-06-20T00:00:00Z",
			"url": "https://api.github.com/repos/o/r/secret-scanning/alerts/1",
			"html_url": "https://github.com/o/r/security/secret-scanning/1",
			"locations_url": "https://api.github.com/repos/o/r/secret-scanning/alerts/1/locations",
			"state": "open",
			"resolution": null,
			"resolved_at": null,
			"resolved_by": null,
			"secret_type": "mailchimp_api_key",
			"secret": "XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX-us2"
		}`)
	})

	ctx := context.Background()

	alert, _, err := client.SecretScanning.GetAlert(ctx, "o", "r", 1)
	if err != nil {
		t.Errorf("SecretScanning.GetAlert returned error: %v", err)
	}

	date := Timestamp{time.Date(1996, time.June, 20, 00, 00, 00, 0, time.UTC)}
	want := &SecretScanningAlert{
		Number:       Ptr(1),
		CreatedAt:    &date,
		URL:          Ptr("https://api.github.com/repos/o/r/secret-scanning/alerts/1"),
		HTMLURL:      Ptr("https://github.com/o/r/security/secret-scanning/1"),
		LocationsURL: Ptr("https://api.github.com/repos/o/r/secret-scanning/alerts/1/locations"),
		State:        Ptr("open"),
		Resolution:   nil,
		ResolvedAt:   nil,
		ResolvedBy:   nil,
		SecretType:   Ptr("mailchimp_api_key"),
		Secret:       Ptr("XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX-us2"),
	}

	if !cmp.Equal(alert, want) {
		t.Errorf("SecretScanning.GetAlert returned %+v, want %+v", alert, want)
	}

	const methodName = "GetAlert"

	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.SecretScanning.GetAlert(ctx, "\n", "\n", 0)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		_, resp, err := client.SecretScanning.GetAlert(ctx, "o", "r", 1)
		return resp, err
	})
}

func TestSecretScanningService_UpdateAlert(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/repos/o/r/secret-scanning/alerts/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PATCH")

		v := new(SecretScanningAlertUpdateOptions)
		assertNilError(t, json.NewDecoder(r.Body).Decode(v))

		want := &SecretScanningAlertUpdateOptions{State: "resolved", Resolution: Ptr("used_in_tests")}

		if !cmp.Equal(v, want) {
			t.Errorf("Request body = %+v, want %+v", v, want)
		}

		fmt.Fprint(w, `{
			"number": 1,
			"created_at": "1996-06-20T00:00:00Z",
			"url": "https://api.github.com/repos/o/r/secret-scanning/alerts/1",
			"html_url": "https://github.com/o/r/security/secret-scanning/1",
			"locations_url": "https://api.github.com/repos/o/r/secret-scanning/alerts/1/locations",
			"state": "resolved",
			"resolution": "used_in_tests",
			"resolution_comment": "resolution comment",
			"resolved_at": "1996-06-20T00:00:00Z",
			"resolved_by": null,
			"secret_type": "mailchimp_api_key",
			"secret": "XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX-us2"
		}`)
	})

	ctx := context.Background()
	opts := &SecretScanningAlertUpdateOptions{State: "resolved", Resolution: Ptr("used_in_tests")}

	alert, _, err := client.SecretScanning.UpdateAlert(ctx, "o", "r", 1, opts)
	if err != nil {
		t.Errorf("SecretScanning.UpdateAlert returned error: %v", err)
	}

	date := Timestamp{time.Date(1996, time.June, 20, 00, 00, 00, 0, time.UTC)}
	want := &SecretScanningAlert{
		Number:            Ptr(1),
		CreatedAt:         &date,
		URL:               Ptr("https://api.github.com/repos/o/r/secret-scanning/alerts/1"),
		HTMLURL:           Ptr("https://github.com/o/r/security/secret-scanning/1"),
		LocationsURL:      Ptr("https://api.github.com/repos/o/r/secret-scanning/alerts/1/locations"),
		State:             Ptr("resolved"),
		Resolution:        Ptr("used_in_tests"),
		ResolutionComment: Ptr("resolution comment"),
		ResolvedAt:        &date,
		ResolvedBy:        nil,
		SecretType:        Ptr("mailchimp_api_key"),
		Secret:            Ptr("XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX-us2"),
	}

	if !cmp.Equal(alert, want) {
		t.Errorf("SecretScanning.UpdateAlert returned %+v, want %+v", alert, want)
	}

	const methodName = "UpdateAlert"

	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.SecretScanning.UpdateAlert(ctx, "\n", "\n", 1, opts)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		_, resp, err := client.SecretScanning.UpdateAlert(ctx, "o", "r", 1, opts)
		return resp, err
	})
}

func TestSecretScanningService_ListLocationsForAlert(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/repos/o/r/secret-scanning/alerts/1/locations", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testFormValues(t, r, values{"page": "1", "per_page": "100"})

		fmt.Fprint(w, `[{
			"type": "commit",
			"details": {
			  "path": "/example/secrets.txt",
			  "start_line": 1,
			  "end_line": 1,
			  "start_column": 1,
			  "end_column": 64,
			  "blob_sha": "af5626b4a114abcb82d63db7c8082c3c4756e51b",
			  "blob_url": "https://api.github.com/repos/o/r/git/blobs/af5626b4a114abcb82d63db7c8082c3c4756e51b",
			  "commit_sha": "f14d7debf9775f957cf4f1e8176da0786431f72b",
			  "commit_url": "https://api.github.com/repos/o/r/git/commits/f14d7debf9775f957cf4f1e8176da0786431f72b"
			}
		}]`)
	})

	ctx := context.Background()
	opts := &ListOptions{Page: 1, PerPage: 100}

	locations, _, err := client.SecretScanning.ListLocationsForAlert(ctx, "o", "r", 1, opts)
	if err != nil {
		t.Errorf("SecretScanning.ListLocationsForAlert returned error: %v", err)
	}

	want := []*SecretScanningAlertLocation{
		{
			Type: Ptr("commit"),
			Details: &SecretScanningAlertLocationDetails{
				Path:        Ptr("/example/secrets.txt"),
				Startline:   Ptr(1),
				EndLine:     Ptr(1),
				StartColumn: Ptr(1),
				EndColumn:   Ptr(64),
				BlobSHA:     Ptr("af5626b4a114abcb82d63db7c8082c3c4756e51b"),
				BlobURL:     Ptr("https://api.github.com/repos/o/r/git/blobs/af5626b4a114abcb82d63db7c8082c3c4756e51b"),
				CommitSHA:   Ptr("f14d7debf9775f957cf4f1e8176da0786431f72b"),
				CommitURL:   Ptr("https://api.github.com/repos/o/r/git/commits/f14d7debf9775f957cf4f1e8176da0786431f72b"),
			},
		},
	}

	if !cmp.Equal(locations, want) {
		t.Errorf("SecretScanning.ListLocationsForAlert returned %+v, want %+v", locations, want)
	}

	const methodName = "ListLocationsForAlert"

	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.SecretScanning.ListLocationsForAlert(ctx, "\n", "\n", 1, opts)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		_, resp, err := client.SecretScanning.ListLocationsForAlert(ctx, "o", "r", 1, opts)
		return resp, err
	})
}

func TestSecretScanningAlert_Marshal(t *testing.T) {
	t.Parallel()
	testJSONMarshal(t, &SecretScanningAlert{}, `{}`)

	u := &SecretScanningAlert{
		Number:       Ptr(1),
		CreatedAt:    &Timestamp{referenceTime},
		URL:          Ptr("https://api.github.com/teams/2/discussions/3/comments"),
		HTMLURL:      Ptr("https://api.github.com/teams/2/discussions/3/comments"),
		LocationsURL: Ptr("https://api.github.com/teams/2/discussions/3/comments"),
		State:        Ptr("test_state"),
		Resolution:   Ptr("test_resolution"),
		ResolvedAt:   &Timestamp{referenceTime},
		ResolvedBy: &User{
			Login:     Ptr("test"),
			ID:        Ptr(int64(10)),
			NodeID:    Ptr("A123"),
			AvatarURL: Ptr("https://api.github.com/teams/2/discussions/3/comments"),
		},
		SecretType: Ptr("test"),
		Secret:     Ptr("test"),
	}

	want := `{
		"number": 1,
		"created_at": ` + referenceTimeStr + `,
		"url": "https://api.github.com/teams/2/discussions/3/comments",
		"html_url": "https://api.github.com/teams/2/discussions/3/comments",
		"locations_url": "https://api.github.com/teams/2/discussions/3/comments",
		"state": "test_state",
		"resolution": "test_resolution",
		"resolved_at": ` + referenceTimeStr + `,
		"resolved_by": {
			"login": "test",
			"id": 10,
			"node_id": "A123",
			"avatar_url": "https://api.github.com/teams/2/discussions/3/comments"
		},
		"secret_type": "test",
		"secret": "test"
	}`

	testJSONMarshal(t, u, want)
}

func TestSecretScanningAlertLocation_Marshal(t *testing.T) {
	t.Parallel()
	testJSONMarshal(t, &SecretScanningAlertLocation{}, `{}`)

	u := &SecretScanningAlertLocation{
		Type: Ptr("test"),
		Details: &SecretScanningAlertLocationDetails{
			Path:        Ptr("test_path"),
			Startline:   Ptr(10),
			EndLine:     Ptr(20),
			StartColumn: Ptr(30),
			EndColumn:   Ptr(40),
			BlobSHA:     Ptr("test_sha"),
			BlobURL:     Ptr("https://api.github.com/repos/o/r/git/commits/f14d7debf9775f957cf4f1e8176da0786431f72b"),
			CommitSHA:   Ptr("test_sha"),
			CommitURL:   Ptr("https://api.github.com/repos/o/r/git/commits/f14d7debf9775f957cf4f1e8176da0786431f72b"),
		},
	}

	want := `{
		"type": "test",
		"details": {
			"path": "test_path",
			"start_line": 10,
			"end_line": 20,
			"start_column": 30,
			"end_column": 40,
			"blob_sha": "test_sha",
			"blob_url": "https://api.github.com/repos/o/r/git/commits/f14d7debf9775f957cf4f1e8176da0786431f72b",
			"commit_sha": "test_sha",
			"commit_url": "https://api.github.com/repos/o/r/git/commits/f14d7debf9775f957cf4f1e8176da0786431f72b"
		}
	}`

	testJSONMarshal(t, u, want)
}

func TestSecretScanningAlertLocationDetails_Marshal(t *testing.T) {
	t.Parallel()
	testJSONMarshal(t, &SecretScanningAlertLocationDetails{}, `{}`)

	u := &SecretScanningAlertLocationDetails{
		Path:        Ptr("test_path"),
		Startline:   Ptr(10),
		EndLine:     Ptr(20),
		StartColumn: Ptr(30),
		EndColumn:   Ptr(40),
		BlobSHA:     Ptr("test_sha"),
		BlobURL:     Ptr("https://api.github.com/repos/o/r/git/commits/f14d7debf9775f957cf4f1e8176da0786431f72b"),
		CommitSHA:   Ptr("test_sha"),
		CommitURL:   Ptr("https://api.github.com/repos/o/r/git/commits/f14d7debf9775f957cf4f1e8176da0786431f72b"),
	}

	want := `{
		"path": "test_path",
		"start_line": 10,
		"end_line": 20,
		"start_column": 30,
		"end_column": 40,
		"blob_sha": "test_sha",
		"blob_url": "https://api.github.com/repos/o/r/git/commits/f14d7debf9775f957cf4f1e8176da0786431f72b",
		"commit_sha": "test_sha",
		"commit_url": "https://api.github.com/repos/o/r/git/commits/f14d7debf9775f957cf4f1e8176da0786431f72b"
	}`

	testJSONMarshal(t, u, want)
}

func TestSecretScanningAlertUpdateOptions_Marshal(t *testing.T) {
	t.Parallel()
	testJSONMarshal(t, &SecretScanningAlertUpdateOptions{}, `{}`)

	u := &SecretScanningAlertUpdateOptions{
		State:      "open",
		Resolution: Ptr("false_positive"),
	}

	want := `{
		"state": "open",
		"resolution": "false_positive"
	}`

	testJSONMarshal(t, u, want)
}
