// Copyright 2013 The go-github AUTHORS. All rights reserved.
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

	want := []Authorization{{ID: Int(1)}, {ID: Int(2)}}

	mux.HandleFunc("/authorizations", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")

		j, _ := json.Marshal(want)
		fmt.Fprint(w, string(j))
	})

	output, _, err := client.Authorizations.List()

	if err != nil {
		t.Errorf("Authorizations.List returned error: %v", err)
	}

	if !reflect.DeepEqual(output, want) {
		t.Errorf("Authorizations.List returned %+v, want %+v", output, want)
	}
}

func TestAuthorizationsService_Get(t *testing.T) {

	setup()
	defer teardown()

	id := 1
	want := &Authorization{ID: Int(id)}
	path := fmt.Sprintf("/authorizations/%v", id)

	mux.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")

		j, _ := json.Marshal(want)
		fmt.Fprint(w, string(j))
	})

	output, _, err := client.Authorizations.Get(id)

	if err != nil {
		t.Errorf("Authorizations.Get returned error: %v", err)
	}

	if !reflect.DeepEqual(output, want) {
		t.Errorf("Authorizations.Get returned %+v, want %+v", output, want)
	}
}

func TestAuthorizationsService_Create(t *testing.T) {

	setup()
	defer teardown()

	input := &AuthorizationRequest{Note: String("12345")}
	want := &Authorization{ID: Int(1), Note: input.Note}

	mux.HandleFunc("/authorizations", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		v := new(AuthorizationRequest)
		json.NewDecoder(r.Body).Decode(v)

		if !reflect.DeepEqual(v, input) {
			t.Errorf("Request body = %+v, want %+v", v, input)
		}

		j, _ := json.Marshal(want)
		fmt.Fprint(w, string(j))
	})

	output, _, err := client.Authorizations.Create(input)

	if err != nil {
		t.Errorf("Authorizations.Create returned error: %v", err)
	}

	if !reflect.DeepEqual(output, want) {
		t.Errorf("Authorizations.Create returned %+v, want %+v", output, want)
	}
}

func TestAuthorizationsService_GetOrCreateForApp(t *testing.T) {

	setup()
	defer teardown()

	input := &AuthorizationRequest{ClientID: String("abcde"), ClientSecret: String("clientSecret"), Note: String("12345")}
	expectedOnServer := &AuthorizationRequest{ClientSecret: input.ClientSecret, Note: input.Note}
	want := &Authorization{ID: Int(1), Note: input.Note}
	path := fmt.Sprintf("/authorizations/clients/%v", *input.ClientID)

	mux.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PUT")
		v := new(AuthorizationRequest)
		json.NewDecoder(r.Body).Decode(v)

		if !reflect.DeepEqual(v, expectedOnServer) {
			t.Errorf("Request body = %+v, want %+v", v, input)
		}

		j, _ := json.Marshal(want)

		fmt.Fprint(w, string(j))
	})

	output, _, err := client.Authorizations.GetOrCreateForApp(input)

	if err != nil {
		t.Errorf("Authorizations.GetOrCreateForApp returned error: %v", err)
	}

	if !reflect.DeepEqual(output, want) {
		t.Errorf("Authorizations.GetOrCreateForApp returned %+v, want %+v", output, want)
	}
}

func TestAuthorizationsService_Edit(t *testing.T) {

	setup()
	defer teardown()

	id := 1
	path := fmt.Sprintf("/authorizations/%v", id)
	input := &AuthorizationUpdate{ID: Int(id), Note: String("12345")}

	// The "ID" field does not get serialized into the body, it goes in the path, so we omit it here
	expectedOnServer := &AuthorizationUpdate{Note: input.Note}
	want := &Authorization{ID: input.ID, Note: input.Note}

	mux.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PATCH")

		v := new(AuthorizationUpdate)
		json.NewDecoder(r.Body).Decode(v)

		if !reflect.DeepEqual(v, expectedOnServer) {
			t.Errorf("Request body = %+v, want %+v", v, input)
		}

		j, _ := json.Marshal(want)

		fmt.Fprint(w, string(j))
	})

	output, _, err := client.Authorizations.Edit(input)

	if err != nil {
		t.Errorf("Authorizations.Edit returned error: %v", err)
	}

	if !reflect.DeepEqual(output, want) {
		t.Errorf("Authorizations.Edit returned %+v, want %+v", output, want)
	}
}

func TestAuthorizationsService_Delete(t *testing.T) {

	setup()
	defer teardown()

	id := 1
	path := fmt.Sprintf("/authorizations/%v", id)

	mux.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
		testBodyIsEmpty(t, r)
		w.WriteHeader(http.StatusNoContent)
	})

	resp, err := client.Authorizations.Delete(id)

	if err != nil {
		t.Errorf("Authorizations.Delete returned error: %v", err)
	}

	if resp.StatusCode != http.StatusNoContent {
		t.Errorf("Authorizations.Delete should have returned status code %v, but returned %v", http.StatusNoContent, resp.StatusCode)
	}
}

func TestAuthorizationsService_CheckAppAuthorization(t *testing.T) {

	setup()
	defer teardown()

	clientID := "abcde"
	token := "12345"
	path := "/applications/" + clientID + "/tokens/" + token

	app := &App{ClientID: &clientID}

	want := &Authorization{ID: Int(1), App: app}

	mux.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testBodyIsEmpty(t, r)

		j, _ := json.Marshal(want)
		fmt.Fprint(w, string(j))
	})

	output, _, err := client.Authorizations.CheckAppAuthorization(clientID, token)

	if err != nil {
		t.Errorf("Authorizations.CheckAppAuthorization returned error: %v", err)
	}

	if !reflect.DeepEqual(output, want) {
		t.Errorf("Authorizations.CheckAppAuthorization returned %+v, want %+v", output, want)
	}
}

func TestAuthorizationsService_ResetAppAuthorization(t *testing.T) {

	setup()
	defer teardown()

	clientID := "abcde"
	token := "12345"
	path := "/applications/" + clientID + "/tokens/" + token

	app := &App{ClientID: &clientID}

	want := &Authorization{ID: Int(1), App: app}

	mux.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		testBodyIsEmpty(t, r)

		j, _ := json.Marshal(want)
		fmt.Fprint(w, string(j))
	})

	output, _, err := client.Authorizations.ResetAppAuthorization(clientID, token)

	if err != nil {
		t.Errorf("Authorizations.ResetAppAuthorization returned error: %v", err)
	}

	if !reflect.DeepEqual(output, want) {
		t.Errorf("Authorizations.ResetAppAuthorization returned %+v, want %+v", output, want)
	}
}

func TestAuthorizationsService_RevokeAppAuthorization(t *testing.T) {

	setup()
	defer teardown()

	clientID := "abcde"
	token := "12345"
	path := "/applications/" + clientID + "/tokens/" + token

	mux.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
		testBodyIsEmpty(t, r)
		w.WriteHeader(http.StatusNoContent)
	})

	resp, err := client.Authorizations.RevokeAppAuthorization(clientID, token)

	if err != nil {
		t.Errorf("Authorizations.RevokeAppAuthorization returned error: %v", err)
	}

	if resp.StatusCode != http.StatusNoContent {
		t.Errorf("Authorizations.RevokeAppAuthorization should have returned status code %v, but returned %v", http.StatusNoContent, resp.StatusCode)
	}
}
