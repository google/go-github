// Copyright 2019 The go-github AUTHORS. All rights reserved.
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

func TestInteractionsService_GetRestrictionsForOrgs(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/orgs/o/interaction-limits", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testHeader(t, r, "Accept", mediaTypeInteractionRestrictionsPreview)
		fmt.Fprint(w, `{"origin":"organization"}`)
	})

	organizationInteractions, _, err := client.Interactions.GetRestrictionsForOrg(context.Background(), "o")
	if err != nil {
		t.Errorf("Interactions.GetRestrictionsForOrg returned error: %v", err)
	}

	want := &InteractionRestriction{Origin: String("organization")}
	if !reflect.DeepEqual(organizationInteractions, want) {
		t.Errorf("Interactions.GetRestrictionsForOrg returned %+v, want %+v", organizationInteractions, want)
	}
}

func TestInteractionsService_UpdateRestrictionsForOrg(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	input := &InteractionRestriction{Limit: String("existing_users")}

	mux.HandleFunc("/orgs/o/interaction-limits", func(w http.ResponseWriter, r *http.Request) {
		v := new(InteractionRestriction)
		json.NewDecoder(r.Body).Decode(v)

		testMethod(t, r, "PUT")
		testHeader(t, r, "Accept", mediaTypeInteractionRestrictionsPreview)
		if !reflect.DeepEqual(v, input) {
			t.Errorf("Request body = %+v, want %+v", v, input)
		}
		fmt.Fprint(w, `{"origin":"organization"}`)
	})

	organizationInteractions, _, err := client.Interactions.UpdateRestrictionsForOrg(context.Background(), "o", input.GetLimit())
	if err != nil {
		t.Errorf("Interactions.UpdateRestrictionsForOrg returned error: %v", err)
	}

	want := &InteractionRestriction{Origin: String("organization")}
	if !reflect.DeepEqual(organizationInteractions, want) {
		t.Errorf("Interactions.UpdateRestrictionsForOrg returned %+v, want %+v", organizationInteractions, want)
	}
}

func TestInteractionsService_RemoveRestrictionsFromOrg(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/orgs/o/interaction-limits", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
		testHeader(t, r, "Accept", mediaTypeInteractionRestrictionsPreview)
	})

	_, err := client.Interactions.RemoveRestrictionsFromOrg(context.Background(), "o")
	if err != nil {
		t.Errorf("Interactions.RemoveRestrictionsFromOrg returned error: %v", err)
	}
}
