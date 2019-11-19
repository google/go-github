// Copyright 2019 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package github

import (
	"context"
	"fmt"
	"net/http"
	"reflect"
	"testing"
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
		testHeader(t, r, "Accept", mediaTypeAppManifestPreview)
		fmt.Fprint(w, manifestJSON)
	})

	cfg, _, err := client.Apps.CompleteAppManifest(context.Background(), "code")
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

	if !reflect.DeepEqual(cfg, want) {
		t.Errorf("GetConfig returned %+v, want %+v", cfg, want)
	}
}
