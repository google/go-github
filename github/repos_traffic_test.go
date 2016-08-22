// Copyright 2013 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package github

import (
	"fmt"
	"net/http"
	"reflect"
	"testing"
	"time"
)

func TestRepositoriesService_ListReferrers(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/repos/1/2/traffic/popular/referrers", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testHeader(t, r, "Accept", mediaTypeTrafficPreview)
		fmt.Fprintf(w, `[{
			"referrer": "Google",
			"count": 4,
			"uniques": 3
 		}]`)
	})
	referrers, _, err := client.Repositories.ListReferrers("1", "2")
	if err != nil {
		t.Errorf("Repositories.ListPaths returned error: %+v", err)
	}

	want := []*Referrer{{
		Referrer: String("Google"),
		Count:    Int(4),
		Uniques:  Int(3),
	}}
	if !reflect.DeepEqual(referrers, want) {
		t.Errorf("Repositories.ListReferrers returned %+v, want %+v", referrers, want)
	}

}

func TestRepositoriesService_ListPaths(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/repos/1/2/traffic/popular/paths", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testHeader(t, r, "Accept", mediaTypeTrafficPreview)
		fmt.Fprintf(w, `[{
			"path": "/github/hubot",
			"title": "github/hubot: A customizable life embetterment robot.",
			"count": 3542,
			"uniques": 2225
 		}]`)
	})
	paths, _, err := client.Repositories.ListPaths("1", "2")
	if err != nil {
		t.Errorf("Repositories.ListPaths returned error: %+v", err)
	}

	want := []*Path{{
		Path:    String("/github/hubot"),
		Title:   String("github/hubot: A customizable life embetterment robot."),
		Count:   Int(3542),
		Uniques: Int(2225),
	}}
	if !reflect.DeepEqual(paths, want) {
		t.Errorf("Repositories.ListPaths returned %+v, want %+v", paths, want)
	}

}

func TestRepositoriesService_ListViews(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/repos/1/2/traffic/views", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testHeader(t, r, "Accept", mediaTypeTrafficPreview)
		fmt.Fprintf(w, `{"count": 7,
			"uniques": 6,
			"views": [{
				"timestamp": 1464710400000,
				"count": 7,
				"uniques": 6
		}]}`)
	})

	views, _, err := client.Repositories.ListViews("1", "2", nil)
	if err != nil {
		t.Errorf("Repositories.ListPaths returned error: %+v", err)
	}

	want := &Views{
		Views: &[]Datapoint{{
			Timestamp: &Time{time.Unix(1464710400, 0)},
			Count:     Int(7),
			Uniques:   Int(6),
		}},
		Count:   Int(7),
		Uniques: Int(6),
	}

	if !reflect.DeepEqual(views, want) {
		t.Errorf("Repositories.ListViews returned %+v, want %+v", views, want)
	}

}

func TestRepositoriesService_ListClones(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/repos/1/2/traffic/clones", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testHeader(t, r, "Accept", mediaTypeTrafficPreview)
		fmt.Fprintf(w, `{"count": 7,
			"uniques": 6,
			"clones": [{
				"timestamp": 1464710400000,
				"count": 7,
				"uniques": 6
		}]}`)
	})

	clones, _, err := client.Repositories.ListClones("1", "2", nil)
	if err != nil {
		t.Errorf("Repositories.ListPaths returned error: %+v", err)
	}

	want := &Clones{
		Clones: &[]Datapoint{{
			Timestamp: &Time{time.Unix(1464710400, 0)},
			Count:     Int(7),
			Uniques:   Int(6),
		}},
		Count:   Int(7),
		Uniques: Int(6),
	}

	if !reflect.DeepEqual(clones, want) {
		t.Errorf("Repositories.ListViews returned %+v, want %+v", clones, want)
	}

}
