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
	client, mux, _, teardown := setup()
	defer teardown()

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
			"secret": "XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX-us2"
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
			Number:       Int(1),
			CreatedAt:    &date,
			URL:          String("https://api.github.com/repos/o/r/secret-scanning/alerts/1"),
			HTMLURL:      String("https://github.com/o/r/security/secret-scanning/1"),
			LocationsURL: String("https://api.github.com/repos/o/r/secret-scanning/alerts/1/locations"),
			State:        String("open"),
			Resolution:   nil,
			ResolvedAt:   nil,
			ResolvedBy:   nil,
			SecretType:   String("mailchimp_api_key"),
			Secret:       String("XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX-us2"),
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
	client, mux, _, teardown := setup()
	defer teardown()

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
			Number:       Int(1),
			CreatedAt:    &date,
			URL:          String("https://api.github.com/repos/o/r/secret-scanning/alerts/1"),
			HTMLURL:      String("https://github.com/o/r/security/secret-scanning/1"),
			LocationsURL: String("https://api.github.com/repos/o/r/secret-scanning/alerts/1/locations"),
			State:        String("open"),
			Resolution:   nil,
			ResolvedAt:   nil,
			ResolvedBy:   nil,
			SecretType:   String("mailchimp_api_key"),
			Secret:       String("XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX-us2"),
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
	client, mux, _, teardown := setup()
	defer teardown()

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
			Number:       Int(1),
			CreatedAt:    &date,
			URL:          String("https://api.github.com/repos/o/r/secret-scanning/alerts/1"),
			HTMLURL:      String("https://github.com/o/r/security/secret-scanning/1"),
			LocationsURL: String("https://api.github.com/repos/o/r/secret-scanning/alerts/1/locations"),
			State:        String("open"),
			Resolution:   nil,
			ResolvedAt:   nil,
			ResolvedBy:   nil,
			SecretType:   String("mailchimp_api_key"),
			Secret:       String("XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX-us2"),
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
	client, mux, _, teardown := setup()
	defer teardown()

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
			Number:       Int(1),
			CreatedAt:    &date,
			URL:          String("https://api.github.com/repos/o/r/secret-scanning/alerts/1"),
			HTMLURL:      String("https://github.com/o/r/security/secret-scanning/1"),
			LocationsURL: String("https://api.github.com/repos/o/r/secret-scanning/alerts/1/locations"),
			State:        String("open"),
			Resolution:   nil,
			ResolvedAt:   nil,
			ResolvedBy:   nil,
			SecretType:   String("mailchimp_api_key"),
			Secret:       String("XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX-us2"),
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
	client, mux, _, teardown := setup()
	defer teardown()

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
		Number:       Int(1),
		CreatedAt:    &date,
		URL:          String("https://api.github.com/repos/o/r/secret-scanning/alerts/1"),
		HTMLURL:      String("https://github.com/o/r/security/secret-scanning/1"),
		LocationsURL: String("https://api.github.com/repos/o/r/secret-scanning/alerts/1/locations"),
		State:        String("open"),
		Resolution:   nil,
		ResolvedAt:   nil,
		ResolvedBy:   nil,
		SecretType:   String("mailchimp_api_key"),
		Secret:       String("XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX-us2"),
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
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/repos/o/r/secret-scanning/alerts/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PATCH")

		v := new(SecretScanningAlertUpdateOptions)
		json.NewDecoder(r.Body).Decode(v)

		want := &SecretScanningAlertUpdateOptions{State: String("resolved"), Resolution: String("used_in_tests")}

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
			"resolved_at": "1996-06-20T00:00:00Z",
			"resolved_by": null,
			"secret_type": "mailchimp_api_key",
			"secret": "XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX-us2"
		}`)
	})

	ctx := context.Background()
	opts := &SecretScanningAlertUpdateOptions{State: String("resolved"), Resolution: String("used_in_tests")}

	alert, _, err := client.SecretScanning.UpdateAlert(ctx, "o", "r", 1, opts)
	if err != nil {
		t.Errorf("SecretScanning.UpdateAlert returned error: %v", err)
	}

	date := Timestamp{time.Date(1996, time.June, 20, 00, 00, 00, 0, time.UTC)}
	want := &SecretScanningAlert{
		Number:       Int(1),
		CreatedAt:    &date,
		URL:          String("https://api.github.com/repos/o/r/secret-scanning/alerts/1"),
		HTMLURL:      String("https://github.com/o/r/security/secret-scanning/1"),
		LocationsURL: String("https://api.github.com/repos/o/r/secret-scanning/alerts/1/locations"),
		State:        String("resolved"),
		Resolution:   String("used_in_tests"),
		ResolvedAt:   &date,
		ResolvedBy:   nil,
		SecretType:   String("mailchimp_api_key"),
		Secret:       String("XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX-us2"),
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
	client, mux, _, teardown := setup()
	defer teardown()

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
			Type: String("commit"),
			Details: &SecretScanningAlertLocationDetails{
				Path:        String("/example/secrets.txt"),
				Startline:   Int(1),
				EndLine:     Int(1),
				StartColumn: Int(1),
				EndColumn:   Int(64),
				BlobSHA:     String("af5626b4a114abcb82d63db7c8082c3c4756e51b"),
				BlobURL:     String("https://api.github.com/repos/o/r/git/blobs/af5626b4a114abcb82d63db7c8082c3c4756e51b"),
				CommitSHA:   String("f14d7debf9775f957cf4f1e8176da0786431f72b"),
				CommitURL:   String("https://api.github.com/repos/o/r/git/commits/f14d7debf9775f957cf4f1e8176da0786431f72b"),
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
	testJSONMarshal(t, &SecretScanningAlert{}, `{}`)

	u := &SecretScanningAlert{
		Number:       Int(1),
		CreatedAt:    &Timestamp{referenceTime},
		URL:          String("https://api.github.com/teams/2/discussions/3/comments"),
		HTMLURL:      String("https://api.github.com/teams/2/discussions/3/comments"),
		LocationsURL: String("https://api.github.com/teams/2/discussions/3/comments"),
		State:        String("test_state"),
		Resolution:   String("test_resolution"),
		ResolvedAt:   &Timestamp{referenceTime},
		ResolvedBy: &User{
			Login:     String("test"),
			ID:        Int64(10),
			NodeID:    String("A123"),
			AvatarURL: String("https://api.github.com/teams/2/discussions/3/comments"),
		},
		SecretType: String("test"),
		Secret:     String("test"),
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
	testJSONMarshal(t, &SecretScanningAlertLocation{}, `{}`)

	u := &SecretScanningAlertLocation{
		Type: String("test"),
		Details: &SecretScanningAlertLocationDetails{
			Path:        String("test_path"),
			Startline:   Int(10),
			EndLine:     Int(20),
			StartColumn: Int(30),
			EndColumn:   Int(40),
			BlobSHA:     String("test_sha"),
			BlobURL:     String("https://api.github.com/repos/o/r/git/commits/f14d7debf9775f957cf4f1e8176da0786431f72b"),
			CommitSHA:   String("test_sha"),
			CommitURL:   String("https://api.github.com/repos/o/r/git/commits/f14d7debf9775f957cf4f1e8176da0786431f72b"),
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
	testJSONMarshal(t, &SecretScanningAlertLocationDetails{}, `{}`)

	u := &SecretScanningAlertLocationDetails{
		Path:        String("test_path"),
		Startline:   Int(10),
		EndLine:     Int(20),
		StartColumn: Int(30),
		EndColumn:   Int(40),
		BlobSHA:     String("test_sha"),
		BlobURL:     String("https://api.github.com/repos/o/r/git/commits/f14d7debf9775f957cf4f1e8176da0786431f72b"),
		CommitSHA:   String("test_sha"),
		CommitURL:   String("https://api.github.com/repos/o/r/git/commits/f14d7debf9775f957cf4f1e8176da0786431f72b"),
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
