// Copyright 2013 Google. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file or at
// https://developers.google.com/open-source/licenses/bsd

package github

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"testing"
)

func TestTreesService_List_authenticatedUser(t *testing.T) {
	setup()
	defer teardown()

	url_ := fmt.Sprintf("/repos/%v/%v/git/trees/%v", "user", "repo", "coffebabecoffebabecoffebabe")

	mux.HandleFunc(url_, func(w http.ResponseWriter, r *http.Request) {
		if m := "GET"; m != r.Method {
			t.Errorf("Request method = %v, want %v", r.Method, m)
		}
		fmt.Fprint(w, `{
			  "sha": "9fb037999f264ba9a7fc6274d15fa3ae2ab98312",
			  "url": "https://api.github.com/repos/octocat/Hello-World/trees/9fb037999f264ba9a7fc6274d15fa3ae2ab98312",
			  "tree": [
			    {
			      "path": "file.rb",
			      "mode": "100644",
			      "type": "blob",
			      "size": 30,
			      "sha": "44b4fc6d56897b048c772eb4087f854f46256132",
			      "url": "https://api.github.com/repos/octocat/Hello-World/git/blobs/44b4fc6d56897b048c772eb4087f854f46256132"
			    },
			    {
			      "path": "subdir",
			      "mode": "040000",
			      "type": "tree",
			      "sha": "f484d249c660418515fb01c2b9662073663c242e",
			      "url": "https://api.github.com/repos/octocat/Hello-World/git/blobs/f484d249c660418515fb01c2b9662073663c242e"
			    }
			  ]
			}`)
	})

	trees, err := client.Trees.List("user", "repo", "coffebabecoffebabecoffebabe", nil)
	if err != nil {
		t.Errorf("Trees.List returned error: %v", err)
	}

	want := Tree{
		SHA: `9fb037999f264ba9a7fc6274d15fa3ae2ab98312`,
		Trees: []GitTree{
			GitTree{
				Path: "file.rb",
				Mode: "100644",
				Type: "blob",
				Size: 30,
				SHA:  "44b4fc6d56897b048c772eb4087f854f46256132",
			},
			GitTree{
				Path: "subdir",
				Mode: "040000",
				Type: "tree",
				SHA:  "f484d249c660418515fb01c2b9662073663c242e",
			},
		},
	}
	if !reflect.DeepEqual(*trees, want) {
		t.Errorf("Tree.List returned %+v, want %+v", *trees, want)
	}
}

func TestTreesService_Create_authenticatedUser(t *testing.T) {
	setup()
	defer teardown()

	url_ := fmt.Sprintf("/repos/%v/%v/git/trees/%v", "user", "repo", "coffebabecoffebabecoffebabe")

	input := &CreateTree{
		BaseTree: "9fb037999f264ba9a7fc6274d15fa3ae2ab98312",
		Tree: []GitTree{
			GitTree{
				Path: "file.rb",
				Mode: "100644",
				Type: "blob",
				SHA:  "44b4fc6d56897b048c772eb4087f854f46256132",
			},
		},
	}

	mux.HandleFunc(url_, func(w http.ResponseWriter, r *http.Request) {
		v := new(CreateTree)
		json.NewDecoder(r.Body).Decode(v)

		if m := "POST"; m != r.Method {
			t.Errorf("Request method = %v, want %v", r.Method, m)
		}

		if !reflect.DeepEqual(v, input) {
			t.Errorf("Request body = %+v, want %+v", v, input)
		}

		fmt.Fprint(w, `{
		  "sha": "cd8274d15fa3ae2ab983129fb037999f264ba9a7",
		  "url": "https://api.github.com/repo/octocat/Hello-World/trees/cd8274d15fa3ae2ab983129fb037999f264ba9a7",
		  "tree": [
		    {
		      "path": "file.rb",
		      "mode": "100644",
		      "type": "blob",
		      "size": 132,
		      "sha": "7c258a9869f33c1e1e1f74fbb32f07c86cb5a75b",
		      "url": "https://api.github.com/octocat/Hello-World/git/blobs/7c258a9869f33c1e1e1f74fbb32f07c86cb5a75b"
		    }
		  ]
		}`)
	})

	tree, err := client.Trees.Create("user", "repo", "coffebabecoffebabecoffebabe", input)
	if err != nil {
		t.Errorf("Trees.Create returned error: %v", err)
	}

	want := Tree{
		SHA: "cd8274d15fa3ae2ab983129fb037999f264ba9a7",
		Trees: []GitTree{
			{
				Path: "file.rb",
				Mode: "100644",
				Type: "blob",
				Size: 132,
				SHA:  "7c258a9869f33c1e1e1f74fbb32f07c86cb5a75b",
			},
		},
	}
	if !reflect.DeepEqual(*tree, want) {
		t.Errorf("Tree.Create returned %+v, want %+v", *tree, want)
	}
}
