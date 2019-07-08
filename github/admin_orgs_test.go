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

func TestAdminOrgs_Create(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	input := &Organization{
		Login: String("github"),
	}

	mux.HandleFunc("/admin/organizations", func(w http.ResponseWriter, r *http.Request) {
		v := new(createOrgRequest)
		json.NewDecoder(r.Body).Decode(v)

		testMethod(t, r, "POST")
		want := &createOrgRequest{Login: String("github"), Admin: String("ghAdmin")}
		if !reflect.DeepEqual(v, want) {
			t.Errorf("Request body = %+v, want %+v", v, want)
		}

		fmt.Fprint(w, `{"login":"github","id":1}`)
	})

	org, _, err := client.Admin.CreateOrg(context.Background(), input, "ghAdmin")
	if err != nil {
		t.Errorf("Admin.CreateOrg returned error: %v", err)
	}

	want := &Organization{ID: Int64(1), Login: String("github")}
	if !reflect.DeepEqual(org, want) {
		t.Errorf("Admin.CreateOrg returned %+v, want %+v", org, want)
	}
}
