// Copyright 2015 The go-github AUTHORS. All rights reserved.
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

func TestOrganizationsService_ListHooks(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/orgs/o/hooks", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testFormValues(t, r, values{"page": "2"})
		fmt.Fprint(w, `[{"id":1}, {"id":2}]`)
	})

	opt := &ListOptions{Page: 2}

	ctx := context.Background()
	hooks, _, err := client.Organizations.ListHooks(ctx, "o", opt)
	if err != nil {
		t.Errorf("Organizations.ListHooks returned error: %v", err)
	}

	want := []*Hook{{ID: Int64(1)}, {ID: Int64(2)}}
	if !reflect.DeepEqual(hooks, want) {
		t.Errorf("Organizations.ListHooks returned %+v, want %+v", hooks, want)
	}
}

func TestOrganizationsService_ListHooks_invalidOrg(t *testing.T) {
	client, _, _, teardown := setup()
	defer teardown()

	ctx := context.Background()
	_, _, err := client.Organizations.ListHooks(ctx, "%", nil)
	testURLParseError(t, err)
}

func TestOrganizationsService_CreateHook(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	input := &Hook{CreatedAt: &referenceTime}

	mux.HandleFunc("/orgs/o/hooks", func(w http.ResponseWriter, r *http.Request) {
		v := new(createHookRequest)
		json.NewDecoder(r.Body).Decode(v)

		testMethod(t, r, "POST")
		want := &createHookRequest{Name: "web"}
		if !reflect.DeepEqual(v, want) {
			t.Errorf("Request body = %+v, want %+v", v, want)
		}

		fmt.Fprint(w, `{"id":1}`)
	})

	ctx := context.Background()
	hook, _, err := client.Organizations.CreateHook(ctx, "o", input)
	if err != nil {
		t.Errorf("Organizations.CreateHook returned error: %v", err)
	}

	want := &Hook{ID: Int64(1)}
	if !reflect.DeepEqual(hook, want) {
		t.Errorf("Organizations.CreateHook returned %+v, want %+v", hook, want)
	}
}

func TestOrganizationsService_GetHook(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/orgs/o/hooks/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{"id":1}`)
	})

	ctx := context.Background()
	hook, _, err := client.Organizations.GetHook(ctx, "o", 1)
	if err != nil {
		t.Errorf("Organizations.GetHook returned error: %v", err)
	}

	want := &Hook{ID: Int64(1)}
	if !reflect.DeepEqual(hook, want) {
		t.Errorf("Organizations.GetHook returned %+v, want %+v", hook, want)
	}
}

func TestOrganizationsService_GetHook_invalidOrg(t *testing.T) {
	client, _, _, teardown := setup()
	defer teardown()

	ctx := context.Background()
	_, _, err := client.Organizations.GetHook(ctx, "%", 1)
	testURLParseError(t, err)
}

func TestOrganizationsService_EditHook(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	input := &Hook{}

	mux.HandleFunc("/orgs/o/hooks/1", func(w http.ResponseWriter, r *http.Request) {
		v := new(Hook)
		json.NewDecoder(r.Body).Decode(v)

		testMethod(t, r, "PATCH")
		if !reflect.DeepEqual(v, input) {
			t.Errorf("Request body = %+v, want %+v", v, input)
		}

		fmt.Fprint(w, `{"id":1}`)
	})

	ctx := context.Background()
	hook, _, err := client.Organizations.EditHook(ctx, "o", 1, input)
	if err != nil {
		t.Errorf("Organizations.EditHook returned error: %v", err)
	}

	want := &Hook{ID: Int64(1)}
	if !reflect.DeepEqual(hook, want) {
		t.Errorf("Organizations.EditHook returned %+v, want %+v", hook, want)
	}
}

func TestOrganizationsService_EditHook_invalidOrg(t *testing.T) {
	client, _, _, teardown := setup()
	defer teardown()

	ctx := context.Background()
	_, _, err := client.Organizations.EditHook(ctx, "%", 1, nil)
	testURLParseError(t, err)
}

func TestOrganizationsService_PingHook(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/orgs/o/hooks/1/pings", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
	})

	ctx := context.Background()
	_, err := client.Organizations.PingHook(ctx, "o", 1)
	if err != nil {
		t.Errorf("Organizations.PingHook returned error: %v", err)
	}
}

func TestOrganizationsService_DeleteHook(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/orgs/o/hooks/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
	})

	ctx := context.Background()
	_, err := client.Organizations.DeleteHook(ctx, "o", 1)
	if err != nil {
		t.Errorf("Organizations.DeleteHook returned error: %v", err)
	}
}

func TestOrganizationsService_DeleteHook_invalidOrg(t *testing.T) {
	client, _, _, teardown := setup()
	defer teardown()

	ctx := context.Background()
	_, err := client.Organizations.DeleteHook(ctx, "%", 1)
	testURLParseError(t, err)
}
