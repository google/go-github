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

func TestGistsService_ListComments(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/gists/1/comments", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `[{"id": "1"}]`)
	})

	comments, _, err := client.Gists.ListComments("1")

	if err != nil {
		t.Errorf("Gists.Comments returned error: %v", err)
	}

	want := []GistComment{GistComment{ID: "1"}}
	if !reflect.DeepEqual(comments, want) {
		t.Errorf("Gists.ListComments returned %+v, want %+v", comments, want)
	}
}

func TestGistsService_GetComment(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/gists/1/comments/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{"id": "1"}`)
	})

	comment, _, err := client.Gists.GetComment("1", "1")

	if err != nil {
		t.Errorf("Gists.GetComment returned error: %v", err)
	}

	want := &GistComment{ID: "1"}
	if !reflect.DeepEqual(comment, want) {
		t.Errorf("Gists.GetComment returned %+v, want %+v", comment, want)
	}
}

func TestGistsService_CreateComment(t *testing.T) {
	setup()
	defer teardown()

	input := &GistComment{
		ID:   "1",
		Body: "This is the comment body.",
	}

	mux.HandleFunc("/gists/1/comments", func(w http.ResponseWriter, r *http.Request) {
		v := new(GistComment)
		json.NewDecoder(r.Body).Decode(v)

		testMethod(t, r, "POST")
		if !reflect.DeepEqual(v, input) {
			t.Errorf("Request body = %+v, want %+v", v, input)
		}

		fmt.Fprint(w,
			`
			{
				"id": "1",
				"body": "This is the comment body.",
				"url": "",
				"user": {
					"login": "octocat",
				        "id": 1,
    					"avatar_url": "https://github.com/images/error/octocat_happy.gif",
    					"gravatar_id": "somehexcode",
    					"url": "https://api.github.com/users/octocat"
				}
			}`)
	})

	comment, _, err := client.Gists.CreateComment("1", input)
	if err != nil {
		t.Errorf("Gists.CreateComment returned error: %v", err)
	}

	user := User{Login: "octocat", ID: 1, AvatarURL: "https://github.com/images/error/octocat_happy.gif",
		GravatarID: "somehexcode",
		URL:        "https://api.github.com/users/octocat"}

	want := &GistComment{
		ID:   "1",
		Body: "This is the comment body.",
		URL:  "",
		User: &user,
	}
	if !reflect.DeepEqual(comment, want) {
		t.Errorf("Gists.CreateComment returned %+v, want %+v", comment, want)
	}
}

func TestGistsService_EditComment(t *testing.T) {
	setup()
	defer teardown()

	input := &GistComment{
		ID:   "1",
		Body: "New comment.",
	}

	mux.HandleFunc("/gists/1/comments/1", func(w http.ResponseWriter, r *http.Request) {
		v := new(GistComment)
		json.NewDecoder(r.Body).Decode(v)

		testMethod(t, r, "PATCH")
		if !reflect.DeepEqual(v, input) {
			t.Errorf("Request body = %+v, want %+v", v, input)
		}

		fmt.Fprint(w,
			`
			{
				"id": "1",
				"body": "New comment.",
				"url": "",
				"user": {
					"login": "octocat",
				        "id": 1,
    					"avatar_url": "https://github.com/images/error/octocat_happy.gif",
    					"gravatar_id": "somehexcode",
    					"url": "https://api.github.com/users/octocat"
				}
			}`)
	})

	comment, _, err := client.Gists.EditComment("1", "1", input)
	if err != nil {
		t.Errorf("Gists.EditComment returned error: %v", err)
	}

	user := User{Login: "octocat", ID: 1, AvatarURL: "https://github.com/images/error/octocat_happy.gif",
		GravatarID: "somehexcode",
		URL:        "https://api.github.com/users/octocat"}

	want := &GistComment{
		ID:   "1",
		Body: "New comment.",
		URL:  "",
		User: &user,
	}
	if !reflect.DeepEqual(comment, want) {
		t.Errorf("Gists.EditComment returned %+v, want %+v", comment, want)
	}
}

func TestGistsService_DeleteComment(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/gists/1/comments/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
	})

	_, err := client.Gists.DeleteComment("1", "1")
	if err != nil {
		t.Errorf("Gists.Delete returned error: %v", err)
	}
}
