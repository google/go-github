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

func TestAdminUsers_Create(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/admin/users", func(w http.ResponseWriter, r *http.Request) {
		v := new(createUserRequest)
		json.NewDecoder(r.Body).Decode(v)

		testMethod(t, r, "POST")
		want := &createUserRequest{Login: String("github"), Email: String("email@domain.com")}
		if !reflect.DeepEqual(v, want) {
			t.Errorf("Request body = %+v, want %+v", v, want)
		}

		fmt.Fprint(w, `{"login":"github","id":1}`)
	})

	org, _, err := client.Admin.CreateUser(context.Background(), "github", "email@domain.com")
	if err != nil {
		t.Errorf("Admin.CreateUser returned error: %v", err)
	}

	want := &User{ID: Int64(1), Login: String("github")}
	if !reflect.DeepEqual(org, want) {
		t.Errorf("Admin.CreateUser returned %+v, want %+v", org, want)
	}
}

func TestAdminUsers_Delete(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/admin/users/github", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
	})

	_, err := client.Admin.DeleteUser(context.Background(), "github")
	if err != nil {
		t.Errorf("Admin.DeleteUser returned error: %v", err)
	}
}
