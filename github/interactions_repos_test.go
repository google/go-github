// Copyright 2018 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package github

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"testing"
)

func TestInteractionsService_GetInteractions(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/repos/o/r/interaction-limits", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testHeader(t, r, "Accept", mediaTypeRepositoryInteractionsPreview)
		fmt.Fprint(w, `{"origin":"repository"}`)
	})

	repoInteractions, _, err := client.Interactions.GetInteractions(context.Background(), "o", "r")
	if err != nil {
		t.Errorf("Interactions.GetInteractions returned error: %v", err)
	}

	want := &Interaction{Origin: String("repository")}
	if !reflect.DeepEqual(repoInteractions, want) {
		t.Errorf("Interactions.GetInteractions returned %+v, want %+v", repoInteractions, want)
	}
}

func TestInteractionsService_UpdateInteractions(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	input := &Interaction{Limit: String("existing_users")}

	mux.HandleFunc("/repos/o/r/interaction-limits", func(w http.ResponseWriter, r *http.Request) {
		v := new(Interaction)
		json.NewDecoder(r.Body).Decode(v)

		testMethod(t, r, "PUT")
		testHeader(t, r, "Accept", mediaTypeRepositoryInteractionsPreview)
		if !reflect.DeepEqual(v, input) {
			t.Errorf("Request body = %+v, want %+v", v, input)
		}
		fmt.Fprint(w, `{"origin":"repository"}`)
	})

	repoInteractions, _, err := client.Interactions.UpdateInteractions(context.Background(), "o", "r", input)
	if err != nil {
		t.Errorf("Interactions.UpdateInteractions returned error: %v", err)
	}

	want := &Interaction{Origin: String("repository")}
	if !reflect.DeepEqual(repoInteractions, want) {
		t.Errorf("Interactions.UpdateInteractions returned %+v, want %+v", repoInteractions, want)
	}
}
