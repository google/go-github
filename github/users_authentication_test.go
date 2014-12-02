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

func TestUsersService_Unsuspend(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/users/willnorris/suspended", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
	})

	_, err := client.Users.Unsuspend("willnorris")
	if err != nil {
		t.Errorf("Users.Unsuspend returned error: %v", err)
	}
}
