// Copyright 2018 The go-github AUTHORS. All rights reserved.
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

func TestReactionService_GetInteractions(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/repos/o/r/interaction-limits", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testHeader(t, r, "Accept", mediaTypeRepositoryInteractionsPreview)
		fmt.Fprint(w, `{"origin":"repository"}`)
	})

	repoInteractions, _, err := client.Repositories.GetInteractions(context.Background(), "o", "r")
	if err != nil {
		t.Errorf("Repositories.GetInteractions returned error: %v", err)
	}

	want := &RepositoryInteraction{Origin: String("repository")}
	if !reflect.DeepEqual(repoInteractions, want) {
		t.Errorf("Repositories.GetInteractions returned %+v, want %+v", repoInteractions, want)
	}
}
