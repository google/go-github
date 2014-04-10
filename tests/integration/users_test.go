// Copyright 2014 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package tests

import (
	"fmt"
	"math/rand"
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

func TestUsers_Emails(t *testing.T) {
	if !checkAuth("TestUsers_Emails") {
		return
	}

	emails, _, err := client.Users.ListEmails()
	if err != nil {
		t.Fatalf("Users.ListEmails() returned error: %v", err)
	}

	// create random address not currently in user's emails
	var email string
	for {
		email = fmt.Sprintf("test-%d@example.com", rand.Int())
		for _, e := range emails {
			if e.Email != nil && *e.Email == email {
				continue
			}
		}
		break
	}

	// Add new address
	_, _, err = client.Users.AddEmails([]string{email})
	if err != nil {
		t.Fatalf("Users.AddEmails() returned error: %v", err)
	}

	// List emails again and verify new email is present
	emails, _, err = client.Users.ListEmails()
	if err != nil {
		t.Fatalf("Users.ListEmails() returned error: %v", err)
	}

	var found bool
	for _, e := range emails {
		if e.Email != nil && *e.Email == email {
			found = true
			break
		}
	}

	if !found {
		t.Fatalf("Users.ListEmails() does not contain new addres: %v", email)
	}

	// Remove new address
	_, err = client.Users.DeleteEmails([]string{email})
	if err != nil {
		t.Fatalf("Users.DeleteEmails() returned error: %v", err)
	}

	// List emails again and verify new email was removed
	emails, _, err = client.Users.ListEmails()
	if err != nil {
		t.Fatalf("Users.ListEmails() returned error: %v", err)
	}

	for _, e := range emails {
		if e.Email != nil && *e.Email == email {
			t.Fatalf("Users.ListEmails() still contains address %v after removing it", email)
		}
	}
}
