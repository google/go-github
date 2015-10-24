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
	// 	"time"
)

func TestOAuthService_List_all(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/authorizations", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `[{"id":1}]`)
	})

	opt := &ListOptions{Page: 1, PerPage: 2}
	auths, _, err := client.OAuth.ListAuthorizations(opt)

	if err != nil {
		t.Errorf("OAuth.ListAuthorizations returned error: %v", err)
	}

	want := []Authorization{{ID: Int(1)}}

	if !reflect.DeepEqual(auths, want) {
		t.Errorf("OAuth.ListAuthorizations returned %+v, want %+v", *auths[0].ID, *want[0].ID)
	}
}

func TestOAuthService_Get(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/authorizations/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{"id":1}`)
	})

	auth, _, err := client.OAuth.Get(1)

	if err != nil {
		t.Errorf("OAuth.Get returned error: %v", err)
	}

	want := &Authorization{ID: Int(1)}
	if !reflect.DeepEqual(auth, want) {
		t.Errorf("OAuth.Get returned\nauth %+v\nwant %+v", auth, want)
	}
}

func TestOAuthService_Create(t *testing.T) {
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

	auth, _, err := client.OAuth.Create(input)
	if err != nil {
		t.Errorf("OAuth.Create returned error: %v", err)
	}

	want := &Authorization{ID: Int(1)}
	if !reflect.DeepEqual(auth, want) {
		t.Errorf("Authorization.Create returned %+v, want %+v", auth, want)
	}
}

func TestOAuthService_GetOrCreate(t *testing.T) {
	setup()
	defer teardown()

	input := &AuthorizationRequest{
		Note: String("test"),
	}

	client_id := "abcdefghijklmnopqrstabcdefghijklmnopqrst"
	url := "/authorizations/clients/" + client_id
	mux.HandleFunc(url, func(w http.ResponseWriter, r *http.Request) {
		v := new(AuthorizationRequest)
		json.NewDecoder(r.Body).Decode(v)

		testMethod(t, r, "PUT")
		if !reflect.DeepEqual(v, input) {
			t.Errorf("Request body = %+v, want %+v", v, input)
		}

		fmt.Fprint(w, `{"ID":1}`)
	})

	auth, _, err := client.OAuth.GetOrCreate(client_id, input)
	if err != nil {
		t.Errorf("OAuth.GetOrCreate returned error: %v", err)
	}

	want := &Authorization{ID: Int(1)}
	if !reflect.DeepEqual(auth, want) {
		t.Errorf("Authorization.GetOrCreate returned %+v, want %+v", auth, want)
	}
}

func TestOAuthService_GetOrCreateFingerprint(t *testing.T) {
	setup()
	defer teardown()

	input := &AuthorizationRequest{
		Note: String("test"),
	}

	client_id := "abcdefghijklmnopqrstabcdefghijklmnopqrst"
	fingerprint := "jklmnop12345678"

	url := "/authorizations/clients/" + client_id + "/" + fingerprint
	mux.HandleFunc(url, func(w http.ResponseWriter, r *http.Request) {
		v := new(AuthorizationRequest)
		json.NewDecoder(r.Body).Decode(v)

		testMethod(t, r, "PUT")
		if !reflect.DeepEqual(v, input) {
			t.Errorf("Request body = %+v, want %+v", v, input)
		}

		fmt.Fprint(w, `{"ID":1}`)
	})

	auth, _, err := client.OAuth.GetOrCreateFingerprint(client_id, fingerprint, input)
	if err != nil {
		t.Errorf("OAuth.GetOrCreateFingerprint returned error: %v", err)
	}

	want := &Authorization{ID: Int(1)}
	if !reflect.DeepEqual(auth, want) {
		t.Errorf("Authorization.GetOrCreateFingerprint returned %+v, want %+v", auth, want)
	}
}

func TestOAuthService_Update(t *testing.T) {
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

	auth, _, err := client.OAuth.Update(1, input)
	if err != nil {
		t.Errorf("OAuth.Update returned error: %v", err)
	}

	want := &Authorization{ID: Int(1)}
	if !reflect.DeepEqual(auth, want) {
		t.Errorf("Authorization.Update returned %+v, want %+v", auth, want)
	}
}

func TestOAuthService_Delete(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/authorizations/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
		fmt.Fprint(w, `{"id":1}`)
	})

	auth, _, err := client.OAuth.Delete(1)

	if err != nil {
		t.Errorf("OAuth.Delete returned error: %v", err)
	}

	want := &Authorization{ID: Int(1)}
	if !reflect.DeepEqual(auth, want) {
		t.Errorf("OAuth.Delete returned\nauth %+v\nwant %+v", auth, want)
	}
}

func TestOAuthService_Check(t *testing.T) {
	setup()
	defer teardown()

	client_id := "abcdefghijklmnopqrstabcdefghijklmnopqrst"
	token := "jklmnop12345678"

	url := "/applications/" + client_id + "/tokens/" + token
	mux.HandleFunc(url, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{"id":1}`)
	})

	auth, _, err := client.OAuth.Check(client_id, token)

	if err != nil {
		t.Errorf("OAuth.Check returned error: %v", err)
	}

	want := &Authorization{ID: Int(1)}
	if !reflect.DeepEqual(auth, want) {
		t.Errorf("OAuth.Check returned\nauth %+v\nwant %+v", auth, want)
	}
}

func TestOAuthService_Reset(t *testing.T) {
	setup()
	defer teardown()

	client_id := "abcdefghijklmnopqrstabcdefghijklmnopqrst"
	token := "jklmnop12345678"

	url := "/applications/" + client_id + "/tokens/" + token
	mux.HandleFunc(url, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		fmt.Fprint(w, `{"ID":1}`)
	})

	auth, _, err := client.OAuth.Reset(client_id, token)

	if err != nil {
		t.Errorf("OAuth.Reset returned error: %v", err)
	}

	want := &Authorization{ID: Int(1)}
	if !reflect.DeepEqual(auth, want) {
		t.Errorf("OAuth.Reset returned\nauth %+v\nwant %+v", auth, want)
	}
}
func TestOAuthService_RevokeAll(t *testing.T) {
	setup()
	defer teardown()

	client_id := "abcdefghijklmnopqrstabcdefghijklmnopqrst"
	url := "/applications/" + client_id + "/tokens"
	mux.HandleFunc(url, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
		fmt.Fprint(w, `{"ID":1}`)
	})

	auth, _, err := client.OAuth.RevokeAll(client_id)

	if err != nil {
		t.Errorf("OAuth.RevokeAll returned error: %v", err)
	}

	want := &Authorization{ID: Int(1)}
	if !reflect.DeepEqual(auth, want) {
		t.Errorf("OAuth.RevokeAll returned\nauth %+v\nwant %+v", auth, want)
	}
}

func TestOAuthService_Revoke(t *testing.T) {
	setup()
	defer teardown()

	client_id := "abcdefghijklmnopqrstabcdefghijklmnopqrst"
	token := "jklmnop12345678"

	url := "/applications/" + client_id + "/tokens/" + token
	mux.HandleFunc(url, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
		fmt.Fprint(w, `{"ID":1}`)
	})

	auth, _, err := client.OAuth.Revoke(client_id, token)

	if err != nil {
		t.Errorf("OAuth.Revoke returned error: %v", err)
	}

	want := &Authorization{ID: Int(1)}
	if !reflect.DeepEqual(auth, want) {
		t.Errorf("OAuth.Revoke returned\nauth %+v\nwant %+v", auth, want)
	}
}
