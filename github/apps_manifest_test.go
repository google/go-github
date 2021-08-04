// Copyright 2019 The go-github AUTHORS. All rights reserved.
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

const (
	manifestJSON = `{
	"id": 1,
  "client_id": "a" ,
  "client_secret": "b",
  "webhook_secret": "c",
  "pem": "key"
}
`
)

func TestGetConfig(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/app-manifests/code/conversions", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		fmt.Fprint(w, manifestJSON)
	})

	ctx := context.Background()
	cfg, _, err := client.Apps.CompleteAppManifest(ctx, "code")
	if err != nil {
		t.Errorf("AppManifest.GetConfig returned error: %v", err)
	}

	want := &AppConfig{
		ID:            Int64(1),
		ClientID:      String("a"),
		ClientSecret:  String("b"),
		WebhookSecret: String("c"),
		PEM:           String("key"),
	}

	if !cmp.Equal(cfg, want) {
		t.Errorf("GetConfig returned %+v, want %+v", cfg, want)
	}

	const methodName = "CompleteAppManifest"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Apps.CompleteAppManifest(ctx, "\n")
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Apps.CompleteAppManifest(ctx, "code")
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestAppConfig_Marshal(t *testing.T) {
	testJSONMarshal(t, &AppConfig{}, "{}")

	u := &AppConfig{
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
		Name:          String("n"),
		Description:   String("d"),
		ExternalURL:   String("eu"),
		HTMLURL:       String("hu"),
		CreatedAt:     &Timestamp{referenceTime},
		UpdatedAt:     &Timestamp{referenceTime},
		ClientID:      String("ci"),
		ClientSecret:  String("cs"),
		WebhookSecret: String("ws"),
		PEM:           String("pem"),
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
		"client_id": "ci",
		"client_secret": "cs",
		"webhook_secret": "ws",
		"pem": "pem"
	}`

	testJSONMarshal(t, u, want)
}
