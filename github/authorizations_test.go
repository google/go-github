// Copyright 2015 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package github

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"testing"
)

func TestAuthorizationsService_List(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/authorizations", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testFormValues(t, r, values{"page": "1", "per_page": "2"})
		fmt.Fprint(w, `[{"id":1}]`)
	})

	opt := &ListOptions{Page: 1, PerPage: 2}
	got, _, err := client.Authorizations.List(opt)
	if err != nil {
		t.Errorf("Authorizations.List returned error: %v", err)
	}

	want := []*Authorization{{ID: Int(1)}}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Authorizations.List returned %+v, want %+v", *got[0].ID, *want[0].ID)
	}
}

func TestAuthorizationsService_Get(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/authorizations/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{"id":1}`)
	})

	got, _, err := client.Authorizations.Get(1)
	if err != nil {
		t.Errorf("Authorizations.Get returned error: %v", err)
	}

	want := &Authorization{ID: Int(1)}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Authorizations.Get returned auth %+v, want %+v", got, want)
	}
}

func TestAuthorizationsService_Create(t *testing.T) {
	setup()
	defer teardown()

	input := &AuthorizationRequest{
		Note: String("test"),
	}

	mux.HandleFunc("/authorizations", func(w http.ResponseWriter, r *http.Request) {
		v := new(AuthorizationRequest)
		json.NewDecoder(r.Body).Decode(v)

		testMethod(t, r, "POST")
		if !reflect.DeepEqual(v, input) {
			t.Errorf("Request body = %+v, want %+v", v, input)
		}

		fmt.Fprint(w, `{"ID":1}`)
	})

	got, _, err := client.Authorizations.Create(input)
	if err != nil {
		t.Errorf("Authorizations.Create returned error: %v", err)
	}

	want := &Authorization{ID: Int(1)}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Authorization.Create returned %+v, want %+v", got, want)
	}
}

func TestAuthorizationsService_GetOrCreateForApp(t *testing.T) {
	setup()
	defer teardown()

	input := &AuthorizationRequest{
		Note: String("test"),
	}

	mux.HandleFunc("/authorizations/clients/id", func(w http.ResponseWriter, r *http.Request) {
		v := new(AuthorizationRequest)
		json.NewDecoder(r.Body).Decode(v)

		testMethod(t, r, "PUT")
		if !reflect.DeepEqual(v, input) {
			t.Errorf("Request body = %+v, want %+v", v, input)
		}

		fmt.Fprint(w, `{"ID":1}`)
	})

	got, _, err := client.Authorizations.GetOrCreateForApp("id", input)
	if err != nil {
		t.Errorf("Authorizations.GetOrCreateForApp returned error: %v", err)
	}

	want := &Authorization{ID: Int(1)}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Authorization.GetOrCreateForApp returned %+v, want %+v", got, want)
	}
}

func TestAuthorizationsService_GetOrCreateForApp_Fingerprint(t *testing.T) {
	setup()
	defer teardown()

	input := &AuthorizationRequest{
		Note:        String("test"),
		Fingerprint: String("fp"),
	}

	mux.HandleFunc("/authorizations/clients/id/fp", func(w http.ResponseWriter, r *http.Request) {
		v := new(AuthorizationRequest)
		json.NewDecoder(r.Body).Decode(v)

		testMethod(t, r, "PUT")
		if !reflect.DeepEqual(v, input) {
			t.Errorf("Request body = %+v, want %+v", v, input)
		}

		fmt.Fprint(w, `{"ID":1}`)
	})

	got, _, err := client.Authorizations.GetOrCreateForApp("id", input)
	if err != nil {
		t.Errorf("Authorizations.GetOrCreateForApp returned error: %v", err)
	}

	want := &Authorization{ID: Int(1)}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Authorization.GetOrCreateForApp returned %+v, want %+v", got, want)
	}
}

func TestAuthorizationsService_Edit(t *testing.T) {
	setup()
	defer teardown()

	input := &AuthorizationUpdateRequest{
		Note: String("test"),
	}

	mux.HandleFunc("/authorizations/1", func(w http.ResponseWriter, r *http.Request) {
		v := new(AuthorizationUpdateRequest)
		json.NewDecoder(r.Body).Decode(v)

		testMethod(t, r, "PATCH")
		if !reflect.DeepEqual(v, input) {
			t.Errorf("Request body = %+v, want %+v", v, input)
		}

		fmt.Fprint(w, `{"ID":1}`)
	})

	got, _, err := client.Authorizations.Edit(1, input)
	if err != nil {
		t.Errorf("Authorizations.Edit returned error: %v", err)
	}

	want := &Authorization{ID: Int(1)}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Authorization.Update returned %+v, want %+v", got, want)
	}
}

func TestAuthorizationsService_Delete(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/authorizations/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
		w.WriteHeader(http.StatusNoContent)
	})

	_, err := client.Authorizations.Delete(1)
	if err != nil {
		t.Errorf("Authorizations.Delete returned error: %v", err)
	}
}

func TestAuthorizationsService_Check(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/applications/id/tokens/t", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{"id":1}`)
	})

	got, _, err := client.Authorizations.Check("id", "t")
	if err != nil {
		t.Errorf("Authorizations.Check returned error: %v", err)
	}

	want := &Authorization{ID: Int(1)}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Authorizations.Check returned auth %+v, want %+v", got, want)
	}
}

func TestAuthorizationsService_Reset(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/applications/id/tokens/t", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		fmt.Fprint(w, `{"ID":1}`)
	})

	got, _, err := client.Authorizations.Reset("id", "t")
	if err != nil {
		t.Errorf("Authorizations.Reset returned error: %v", err)
	}

	want := &Authorization{ID: Int(1)}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Authorizations.Reset returned auth %+v, want %+v", got, want)
	}
}

func TestAuthorizationsService_Revoke(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/applications/id/tokens/t", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
		w.WriteHeader(http.StatusNoContent)
	})

	_, err := client.Authorizations.Revoke("id", "t")
	if err != nil {
		t.Errorf("Authorizations.Revoke returned error: %v", err)
	}
}

func TestListGrants(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/applications/grants", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testHeader(t, r, "Accept", mediaTypeOAuthGrantAuthorizationsPreview)
		fmt.Fprint(w, `[{"id": 1}]`)
	})

	got, _, err := client.Authorizations.ListGrants()
	if err != nil {
		t.Errorf("OAuthAuthorizations.ListGrants returned error: %v", err)
	}

	want := []*Grant{{ID: Int(1)}}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("OAuthAuthorizations.ListGrants = %+v, want %+v", got, want)
	}
}

func TestGetGrant(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/applications/grants/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testHeader(t, r, "Accept", mediaTypeOAuthGrantAuthorizationsPreview)
		fmt.Fprint(w, `{"id": 1}`)
	})

	got, _, err := client.Authorizations.GetGrant(1)
	if err != nil {
		t.Errorf("OAuthAuthorizations.GetGrant returned error: %v", err)
	}

	want := &Grant{ID: Int(1)}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("OAuthAuthorizations.GetGrant = %+v, want %+v", got, want)
	}
}

func TestDeleteGrant(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/applications/grants/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
		testHeader(t, r, "Accept", mediaTypeOAuthGrantAuthorizationsPreview)
	})

	_, err := client.Authorizations.DeleteGrant(1)
	if err != nil {
		t.Errorf("OAuthAuthorizations.DeleteGrant returned error: %v", err)
	}
}

func TestAuthorizationsService_CreateImpersonation(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/admin/users/u/authorizations", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		fmt.Fprint(w, `{"id":1}`)
	})

	req := &AuthorizationRequest{Scopes: []Scope{ScopePublicRepo}}
	got, _, err := client.Authorizations.CreateImpersonation("u", req)
	if err != nil {
		t.Errorf("Authorizations.CreateImpersonation returned error: %+v", err)
	}

	want := &Authorization{ID: Int(1)}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Authorizations.CreateImpersonation returned %+v, want %+v", *got.ID, *want.ID)
	}
}

func TestAuthorizationsService_DeleteImpersonation(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/admin/users/u/authorizations", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
	})

	_, err := client.Authorizations.DeleteImpersonation("u")
	if err != nil {
		t.Errorf("Authorizations.DeleteImpersonation returned error: %+v", err)
	}
}
