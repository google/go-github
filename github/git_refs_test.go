// Copyright 2013 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package github

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestGitService_GetRef_singleRef(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/repos/o/r/git/ref/heads/b", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `
		  {
		    "ref": "refs/heads/b",
		    "url": "https://api.github.com/repos/o/r/git/refs/heads/b",
		    "object": {
		      "type": "commit",
		      "sha": "aa218f56b14c9653891f9e74264a383fa43fefbd",
		      "url": "https://api.github.com/repos/o/r/git/commits/aa218f56b14c9653891f9e74264a383fa43fefbd"
		    }
		  }`)
	})

	ctx := context.Background()
	ref, _, err := client.Git.GetRef(ctx, "o", "r", "refs/heads/b")
	if err != nil {
		t.Fatalf("Git.GetRef returned error: %v", err)
	}

	want := &Reference{
		Ref: String("refs/heads/b"),
		URL: String("https://api.github.com/repos/o/r/git/refs/heads/b"),
		Object: &GitObject{
			Type: String("commit"),
			SHA:  String("aa218f56b14c9653891f9e74264a383fa43fefbd"),
			URL:  String("https://api.github.com/repos/o/r/git/commits/aa218f56b14c9653891f9e74264a383fa43fefbd"),
		},
	}
	if !cmp.Equal(ref, want) {
		t.Errorf("Git.GetRef returned %+v, want %+v", ref, want)
	}

	// without 'refs/' prefix
	if _, _, err := client.Git.GetRef(ctx, "o", "r", "heads/b"); err != nil {
		t.Errorf("Git.GetRef returned error: %v", err)
	}

	const methodName = "GetRef"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Git.GetRef(ctx, "\n", "\n", "\n")
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Git.GetRef(ctx, "o", "r", "refs/heads/b")
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestGitService_GetRef_noRefs(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/repos/o/r/git/refs/heads/b", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		w.WriteHeader(http.StatusNotFound)
	})

	ctx := context.Background()
	ref, resp, err := client.Git.GetRef(ctx, "o", "r", "refs/heads/b")
	if err == nil {
		t.Errorf("Expected HTTP 404 response")
	}
	if got, want := resp.Response.StatusCode, http.StatusNotFound; got != want {
		t.Errorf("Git.GetRef returned status %d, want %d", got, want)
	}
	if ref != nil {
		t.Errorf("Git.GetRef return %+v, want nil", ref)
	}

	const methodName = "GetRef"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Git.GetRef(ctx, "o", "r", "refs/heads/b")
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Git.GetRef(ctx, "o", "r", "refs/heads/b")
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestGitService_ListMatchingRefs_singleRef(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/repos/o/r/git/matching-refs/heads/b", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `
		  [
		    {
		      "ref": "refs/heads/b",
		      "url": "https://api.github.com/repos/o/r/git/refs/heads/b",
		      "object": {
		        "type": "commit",
		        "sha": "aa218f56b14c9653891f9e74264a383fa43fefbd",
		        "url": "https://api.github.com/repos/o/r/git/commits/aa218f56b14c9653891f9e74264a383fa43fefbd"
		      }
		    }
		  ]`)
	})

	opts := &ReferenceListOptions{Ref: "refs/heads/b"}
	ctx := context.Background()
	refs, _, err := client.Git.ListMatchingRefs(ctx, "o", "r", opts)
	if err != nil {
		t.Fatalf("Git.ListMatchingRefs returned error: %v", err)
	}

	ref := refs[0]
	want := &Reference{
		Ref: String("refs/heads/b"),
		URL: String("https://api.github.com/repos/o/r/git/refs/heads/b"),
		Object: &GitObject{
			Type: String("commit"),
			SHA:  String("aa218f56b14c9653891f9e74264a383fa43fefbd"),
			URL:  String("https://api.github.com/repos/o/r/git/commits/aa218f56b14c9653891f9e74264a383fa43fefbd"),
		},
	}
	if !cmp.Equal(ref, want) {
		t.Errorf("Git.ListMatchingRefs returned %+v, want %+v", ref, want)
	}

	// without 'refs/' prefix
	opts = &ReferenceListOptions{Ref: "heads/b"}
	if _, _, err := client.Git.ListMatchingRefs(ctx, "o", "r", opts); err != nil {
		t.Errorf("Git.ListMatchingRefs returned error: %v", err)
	}

	const methodName = "ListMatchingRefs"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Git.ListMatchingRefs(ctx, "\n", "\n", opts)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Git.ListMatchingRefs(ctx, "o", "r", opts)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestGitService_ListMatchingRefs_multipleRefs(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/repos/o/r/git/matching-refs/heads/b", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `
		  [
		    {
			    "ref": "refs/heads/booger",
			    "url": "https://api.github.com/repos/o/r/git/refs/heads/booger",
			    "object": {
			      "type": "commit",
			      "sha": "aa218f56b14c9653891f9e74264a383fa43fefbd",
			      "url": "https://api.github.com/repos/o/r/git/commits/aa218f56b14c9653891f9e74264a383fa43fefbd"
			    }
		  	},
		    {
		      "ref": "refs/heads/bandsaw",
		      "url": "https://api.github.com/repos/o/r/git/refs/heads/bandsaw",
		      "object": {
		        "type": "commit",
		        "sha": "612077ae6dffb4d2fbd8ce0cccaa58893b07b5ac",
		        "url": "https://api.github.com/repos/o/r/git/commits/612077ae6dffb4d2fbd8ce0cccaa58893b07b5ac"
		      }
		    }
		  ]
		`)
	})

	opts := &ReferenceListOptions{Ref: "refs/heads/b"}
	ctx := context.Background()
	refs, _, err := client.Git.ListMatchingRefs(ctx, "o", "r", opts)
	if err != nil {
		t.Errorf("Git.ListMatchingRefs returned error: %v", err)
	}

	want := &Reference{
		Ref: String("refs/heads/booger"),
		URL: String("https://api.github.com/repos/o/r/git/refs/heads/booger"),
		Object: &GitObject{
			Type: String("commit"),
			SHA:  String("aa218f56b14c9653891f9e74264a383fa43fefbd"),
			URL:  String("https://api.github.com/repos/o/r/git/commits/aa218f56b14c9653891f9e74264a383fa43fefbd"),
		},
	}
	if !cmp.Equal(refs[0], want) {
		t.Errorf("Git.ListMatchingRefs returned %+v, want %+v", refs[0], want)
	}

	const methodName = "ListMatchingRefs"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Git.ListMatchingRefs(ctx, "\n", "\n", opts)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Git.ListMatchingRefs(ctx, "o", "r", opts)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestGitService_ListMatchingRefs_noRefs(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/repos/o/r/git/matching-refs/heads/b", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, "[]")
	})

	opts := &ReferenceListOptions{Ref: "refs/heads/b"}
	ctx := context.Background()
	refs, _, err := client.Git.ListMatchingRefs(ctx, "o", "r", opts)
	if err != nil {
		t.Errorf("Git.ListMatchingRefs returned error: %v", err)
	}

	if len(refs) != 0 {
		t.Errorf("Git.ListMatchingRefs returned %+v, want an empty slice", refs)
	}

	const methodName = "ListMatchingRefs"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Git.ListMatchingRefs(ctx, "\n", "\n", opts)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Git.ListMatchingRefs(ctx, "o", "r", opts)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestGitService_ListMatchingRefs_allRefs(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/repos/o/r/git/matching-refs/", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `
		  [
		    {
		      "ref": "refs/heads/branchA",
		      "url": "https://api.github.com/repos/o/r/git/refs/heads/branchA",
		      "object": {
		        "type": "commit",
		        "sha": "aa218f56b14c9653891f9e74264a383fa43fefbd",
		        "url": "https://api.github.com/repos/o/r/git/commits/aa218f56b14c9653891f9e74264a383fa43fefbd"
		      }
		    },
		    {
		      "ref": "refs/heads/branchB",
		      "url": "https://api.github.com/repos/o/r/git/refs/heads/branchB",
		      "object": {
		        "type": "commit",
		        "sha": "aa218f56b14c9653891f9e74264a383fa43fefbd",
		        "url": "https://api.github.com/repos/o/r/git/commits/aa218f56b14c9653891f9e74264a383fa43fefbd"
		      }
		    }
		  ]`)
	})

	ctx := context.Background()
	refs, _, err := client.Git.ListMatchingRefs(ctx, "o", "r", nil)
	if err != nil {
		t.Errorf("Git.ListMatchingRefs returned error: %v", err)
	}

	want := []*Reference{
		{
			Ref: String("refs/heads/branchA"),
			URL: String("https://api.github.com/repos/o/r/git/refs/heads/branchA"),
			Object: &GitObject{
				Type: String("commit"),
				SHA:  String("aa218f56b14c9653891f9e74264a383fa43fefbd"),
				URL:  String("https://api.github.com/repos/o/r/git/commits/aa218f56b14c9653891f9e74264a383fa43fefbd"),
			},
		},
		{
			Ref: String("refs/heads/branchB"),
			URL: String("https://api.github.com/repos/o/r/git/refs/heads/branchB"),
			Object: &GitObject{
				Type: String("commit"),
				SHA:  String("aa218f56b14c9653891f9e74264a383fa43fefbd"),
				URL:  String("https://api.github.com/repos/o/r/git/commits/aa218f56b14c9653891f9e74264a383fa43fefbd"),
			},
		},
	}
	if !cmp.Equal(refs, want) {
		t.Errorf("Git.ListMatchingRefs returned %+v, want %+v", refs, want)
	}

	const methodName = "ListMatchingRefs"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Git.ListMatchingRefs(ctx, "\n", "\n", nil)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Git.ListMatchingRefs(ctx, "o", "r", nil)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestGitService_ListMatchingRefs_options(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/repos/o/r/git/matching-refs/t", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testFormValues(t, r, values{"page": "2"})
		fmt.Fprint(w, `[{"ref": "r"}]`)
	})

	opts := &ReferenceListOptions{Ref: "t", ListOptions: ListOptions{Page: 2}}
	ctx := context.Background()
	refs, _, err := client.Git.ListMatchingRefs(ctx, "o", "r", opts)
	if err != nil {
		t.Errorf("Git.ListMatchingRefs returned error: %v", err)
	}

	want := []*Reference{{Ref: String("r")}}
	if !cmp.Equal(refs, want) {
		t.Errorf("Git.ListMatchingRefs returned %+v, want %+v", refs, want)
	}

	const methodName = "ListMatchingRefs"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Git.ListMatchingRefs(ctx, "\n", "\n", opts)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Git.ListMatchingRefs(ctx, "o", "r", opts)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestGitService_CreateRef(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	args := &createRefRequest{
		Ref: String("refs/heads/b"),
		SHA: String("aa218f56b14c9653891f9e74264a383fa43fefbd"),
	}

	mux.HandleFunc("/repos/o/r/git/refs", func(w http.ResponseWriter, r *http.Request) {
		v := new(createRefRequest)
		json.NewDecoder(r.Body).Decode(v)

		testMethod(t, r, "POST")
		if !cmp.Equal(v, args) {
			t.Errorf("Request body = %+v, want %+v", v, args)
		}
		fmt.Fprint(w, `
		  {
		    "ref": "refs/heads/b",
		    "url": "https://api.github.com/repos/o/r/git/refs/heads/b",
		    "object": {
		      "type": "commit",
		      "sha": "aa218f56b14c9653891f9e74264a383fa43fefbd",
		      "url": "https://api.github.com/repos/o/r/git/commits/aa218f56b14c9653891f9e74264a383fa43fefbd"
		    }
		  }`)
	})

	ctx := context.Background()
	ref, _, err := client.Git.CreateRef(ctx, "o", "r", &Reference{
		Ref: String("refs/heads/b"),
		Object: &GitObject{
			SHA: String("aa218f56b14c9653891f9e74264a383fa43fefbd"),
		},
	})
	if err != nil {
		t.Errorf("Git.CreateRef returned error: %v", err)
	}

	want := &Reference{
		Ref: String("refs/heads/b"),
		URL: String("https://api.github.com/repos/o/r/git/refs/heads/b"),
		Object: &GitObject{
			Type: String("commit"),
			SHA:  String("aa218f56b14c9653891f9e74264a383fa43fefbd"),
			URL:  String("https://api.github.com/repos/o/r/git/commits/aa218f56b14c9653891f9e74264a383fa43fefbd"),
		},
	}
	if !cmp.Equal(ref, want) {
		t.Errorf("Git.CreateRef returned %+v, want %+v", ref, want)
	}

	// without 'refs/' prefix
	_, _, err = client.Git.CreateRef(ctx, "o", "r", &Reference{
		Ref: String("heads/b"),
		Object: &GitObject{
			SHA: String("aa218f56b14c9653891f9e74264a383fa43fefbd"),
		},
	})
	if err != nil {
		t.Errorf("Git.CreateRef returned error: %v", err)
	}

	const methodName = "CreateRef"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Git.CreateRef(ctx, "\n", "\n", &Reference{
			Ref: String("refs/heads/b"),
			Object: &GitObject{
				SHA: String("aa218f56b14c9653891f9e74264a383fa43fefbd"),
			},
		})
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Git.CreateRef(ctx, "o", "r", &Reference{
			Ref: String("refs/heads/b"),
			Object: &GitObject{
				SHA: String("aa218f56b14c9653891f9e74264a383fa43fefbd"),
			},
		})
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestGitService_UpdateRef(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	args := &updateRefRequest{
		SHA:   String("aa218f56b14c9653891f9e74264a383fa43fefbd"),
		Force: Bool(true),
	}

	mux.HandleFunc("/repos/o/r/git/refs/heads/b", func(w http.ResponseWriter, r *http.Request) {
		v := new(updateRefRequest)
		json.NewDecoder(r.Body).Decode(v)

		testMethod(t, r, "PATCH")
		if !cmp.Equal(v, args) {
			t.Errorf("Request body = %+v, want %+v", v, args)
		}
		fmt.Fprint(w, `
		  {
		    "ref": "refs/heads/b",
		    "url": "https://api.github.com/repos/o/r/git/refs/heads/b",
		    "object": {
		      "type": "commit",
		      "sha": "aa218f56b14c9653891f9e74264a383fa43fefbd",
		      "url": "https://api.github.com/repos/o/r/git/commits/aa218f56b14c9653891f9e74264a383fa43fefbd"
		    }
		  }`)
	})

	ctx := context.Background()
	ref, _, err := client.Git.UpdateRef(ctx, "o", "r", &Reference{
		Ref:    String("refs/heads/b"),
		Object: &GitObject{SHA: String("aa218f56b14c9653891f9e74264a383fa43fefbd")},
	}, true)
	if err != nil {
		t.Errorf("Git.UpdateRef returned error: %v", err)
	}

	want := &Reference{
		Ref: String("refs/heads/b"),
		URL: String("https://api.github.com/repos/o/r/git/refs/heads/b"),
		Object: &GitObject{
			Type: String("commit"),
			SHA:  String("aa218f56b14c9653891f9e74264a383fa43fefbd"),
			URL:  String("https://api.github.com/repos/o/r/git/commits/aa218f56b14c9653891f9e74264a383fa43fefbd"),
		},
	}
	if !cmp.Equal(ref, want) {
		t.Errorf("Git.UpdateRef returned %+v, want %+v", ref, want)
	}

	// without 'refs/' prefix
	_, _, err = client.Git.UpdateRef(ctx, "o", "r", &Reference{
		Ref:    String("heads/b"),
		Object: &GitObject{SHA: String("aa218f56b14c9653891f9e74264a383fa43fefbd")},
	}, true)
	if err != nil {
		t.Errorf("Git.UpdateRef returned error: %v", err)
	}

	const methodName = "UpdateRef"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Git.UpdateRef(ctx, "\n", "\n", &Reference{
			Ref:    String("refs/heads/b"),
			Object: &GitObject{SHA: String("aa218f56b14c9653891f9e74264a383fa43fefbd")},
		}, true)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Git.UpdateRef(ctx, "o", "r", &Reference{
			Ref:    String("refs/heads/b"),
			Object: &GitObject{SHA: String("aa218f56b14c9653891f9e74264a383fa43fefbd")},
		}, true)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestGitService_DeleteRef(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/repos/o/r/git/refs/heads/b", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
	})

	ctx := context.Background()
	_, err := client.Git.DeleteRef(ctx, "o", "r", "refs/heads/b")
	if err != nil {
		t.Errorf("Git.DeleteRef returned error: %v", err)
	}

	// without 'refs/' prefix
	if _, err := client.Git.DeleteRef(ctx, "o", "r", "heads/b"); err != nil {
		t.Errorf("Git.DeleteRef returned error: %v", err)
	}

	const methodName = "DeleteRef"
	testBadOptions(t, methodName, func() (err error) {
		_, err = client.Git.DeleteRef(ctx, "\n", "\n", "\n")
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		return client.Git.DeleteRef(ctx, "o", "r", "refs/heads/b")
	})
}

func TestGitService_GetRef_pathEscape(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/repos/o/r/git/ref/heads/b", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		if strings.Contains(r.URL.RawPath, "%2F") {
			t.Errorf("RawPath still contains escaped / as %%2F: %v", r.URL.RawPath)
		}
		fmt.Fprint(w, `
		  {
		    "ref": "refs/heads/b",
		    "url": "https://api.github.com/repos/o/r/git/refs/heads/b",
		    "object": {
		      "type": "commit",
		      "sha": "aa218f56b14c9653891f9e74264a383fa43fefbd",
		      "url": "https://api.github.com/repos/o/r/git/commits/aa218f56b14c9653891f9e74264a383fa43fefbd"
		    }
		  }`)
	})

	ctx := context.Background()
	_, _, err := client.Git.GetRef(ctx, "o", "r", "refs/heads/b")
	if err != nil {
		t.Fatalf("Git.GetRef returned error: %v", err)
	}

	const methodName = "GetRef"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Git.GetRef(ctx, "\n", "\n", "\n")
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Git.GetRef(ctx, "o", "r", "refs/heads/b")
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestGitService_UpdateRef_pathEscape(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	args := &updateRefRequest{
		SHA:   String("aa218f56b14c9653891f9e74264a383fa43fefbd"),
		Force: Bool(true),
	}

	mux.HandleFunc("/repos/o/r/git/refs/heads/b#1", func(w http.ResponseWriter, r *http.Request) {
		v := new(updateRefRequest)
		json.NewDecoder(r.Body).Decode(v)

		testMethod(t, r, "PATCH")
		if !cmp.Equal(v, args) {
			t.Errorf("Request body = %+v, want %+v", v, args)
		}
		fmt.Fprint(w, `
		  {
		    "ref": "refs/heads/b#1",
		    "url": "https://api.github.com/repos/o/r/git/refs/heads/b%231",
		    "object": {
		      "type": "commit",
		      "sha": "aa218f56b14c9653891f9e74264a383fa43fefbd",
		      "url": "https://api.github.com/repos/o/r/git/commits/aa218f56b14c9653891f9e74264a383fa43fefbd"
		    }
		  }`)
	})

	ctx := context.Background()
	ref, _, err := client.Git.UpdateRef(ctx, "o", "r", &Reference{
		Ref:    String("refs/heads/b#1"),
		Object: &GitObject{SHA: String("aa218f56b14c9653891f9e74264a383fa43fefbd")},
	}, true)
	if err != nil {
		t.Errorf("Git.UpdateRef returned error: %v", err)
	}

	want := &Reference{
		Ref: String("refs/heads/b#1"),
		URL: String("https://api.github.com/repos/o/r/git/refs/heads/b%231"),
		Object: &GitObject{
			Type: String("commit"),
			SHA:  String("aa218f56b14c9653891f9e74264a383fa43fefbd"),
			URL:  String("https://api.github.com/repos/o/r/git/commits/aa218f56b14c9653891f9e74264a383fa43fefbd"),
		},
	}
	if !cmp.Equal(ref, want) {
		t.Errorf("Git.UpdateRef returned %+v, want %+v", ref, want)
	}
}

func TestReference_Marshal(t *testing.T) {
	testJSONMarshal(t, &Reference{}, "{}")

	u := &Reference{
		Ref: String("ref"),
		URL: String("url"),
		Object: &GitObject{
			Type: String("type"),
			SHA:  String("sha"),
			URL:  String("url"),
		},
		NodeID: String("nid"),
	}

	want := `{
		"ref": "ref",
		"url": "url",
		"object": {
			"type": "type",
			"sha": "sha",
			"url": "url"
		},
		"node_id": "nid"
	}`

	testJSONMarshal(t, u, want)
}

func TestGitObject_Marshal(t *testing.T) {
	testJSONMarshal(t, &GitObject{}, "{}")

	u := &GitObject{
		Type: String("type"),
		SHA:  String("sha"),
		URL:  String("url"),
	}

	want := `{
		"type": "type",
		"sha": "sha",
		"url": "url"
	}`

	testJSONMarshal(t, u, want)
}

func TestCreateRefRequest_Marshal(t *testing.T) {
	testJSONMarshal(t, &createRefRequest{}, "{}")

	u := &createRefRequest{
		Ref: String("ref"),
		SHA: String("sha"),
	}

	want := `{
		"ref": "ref",
		"sha": "sha"
	}`

	testJSONMarshal(t, u, want)
}

func TestUpdateRefRequest_Marshal(t *testing.T) {
	testJSONMarshal(t, &updateRefRequest{}, "{}")

	u := &updateRefRequest{
		SHA:   String("sha"),
		Force: Bool(true),
	}

	want := `{
		"sha": "sha",
		"force": true
	}`

	testJSONMarshal(t, u, want)
}
