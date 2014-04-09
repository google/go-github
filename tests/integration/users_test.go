// Copyright 2014 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package tests

import (
	"fmt"
	"testing"
)

func TestUsers_List(t *testing.T) {
	u, _, err := client.Users.ListAll(nil)
	if err != nil {
		t.Fatalf("Users.ListAll returned error: %v", err)
	}

	if len(u) == 0 {
		t.Errorf("Users.ListAll returned no users")
	}

	// mojombo is user #1
	if want := "mojombo"; want != *u[0].Login {
		t.Errorf("user[0].Login was %q, wanted %q", *u[0].Login, want)
	}
}

func TestUsers_Get(t *testing.T) {
	u, _, err := client.Users.Get("octocat")
	if err != nil {
		t.Fatalf("Users.Get('octocat') returned error: %v", err)
	}

	if want := "octocat"; want != *u.Login {
		t.Errorf("user.Login was %q, wanted %q", *u.Login, want)
	}
	if want := "The Octocat"; want != *u.Name {
		t.Errorf("user.Name was %q, wanted %q", *u.Name, want)
	}

	if checkAuth("TestUsers_Get") {
		u, _, err := client.Users.Get("")
		if err != nil {
			t.Fatalf("Users.Get('') returned error: %v", err)
		}

		if *u.Login == "" {
			t.Errorf("wanted non-empty values for user.Login")
		}
	}
}
