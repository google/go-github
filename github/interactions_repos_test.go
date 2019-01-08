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

func TestInteractionsService_GetRestrictionsForRepo(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/repos/o/r/interaction-limits", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testHeader(t, r, "Accept", mediaTypeInteractionRestrictionsPreview)
		fmt.Fprint(w, `{"origin":"repository"}`)
	})

	repoInteractions, _, err := client.Interactions.GetRestrictionsForRepo(context.Background(), "o", "r")
	if err != nil {
		t.Errorf("Interactions.GetRestrictionsForRepo returned error: %v", err)
	}

	want := &InteractionRestriction{Origin: String("repository")}
	if !reflect.DeepEqual(repoInteractions, want) {
		t.Errorf("Interactions.GetRestrictionsForRepo returned %+v, want %+v", repoInteractions, want)
	}
}

func TestInteractionsService_UpdateRestrictionsForRepo(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	input := &InteractionRestriction{Limit: String("existing_users")}

	mux.HandleFunc("/repos/o/r/interaction-limits", func(w http.ResponseWriter, r *http.Request) {
		v := new(InteractionRestriction)
		json.NewDecoder(r.Body).Decode(v)

		testMethod(t, r, "PUT")
		testHeader(t, r, "Accept", mediaTypeInteractionRestrictionsPreview)
		if !reflect.DeepEqual(v, input) {
			t.Errorf("Request body = %+v, want %+v", v, input)
		}
		fmt.Fprint(w, `{"origin":"repository"}`)
	})

	repoInteractions, _, err := client.Interactions.UpdateRestrictionsForRepo(context.Background(), "o", "r", input.GetLimit())
	if err != nil {
		t.Errorf("Interactions.UpdateRestrictionsForRepo returned error: %v", err)
	}

	want := &InteractionRestriction{Origin: String("repository")}
	if !reflect.DeepEqual(repoInteractions, want) {
		t.Errorf("Interactions.UpdateRestrictionsForRepo returned %+v, want %+v", repoInteractions, want)
	}
}

func TestInteractionsService_RemoveRestrictionsFromRepo(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/repos/o/r/interaction-limits", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
		testHeader(t, r, "Accept", mediaTypeInteractionRestrictionsPreview)
	})

	_, err := client.Interactions.RemoveRestrictionsFromRepo(context.Background(), "o", "r")
	if err != nil {
		t.Errorf("Interactions.RemoveRestrictionsFromRepo returned error: %v", err)
	}
}
