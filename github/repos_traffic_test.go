// Copyright 2016 The go-github AUTHORS. All rights reserved.
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

func TestRepositoriesService_ListTrafficReferrers(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/repos/o/r/traffic/popular/referrers", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testHeader(t, r, "Accept", mediaTypeTrafficPreview)
		fmt.Fprintf(w, `[{
			"referrer": "Google",
			"count": 4,
			"uniques": 3
 		}]`)
	})
	referrers, _, err := client.Repositories.ListTrafficReferrers("o", "r")
	if err != nil {
		t.Errorf("Repositories.ListPaths returned error: %+v", err)
	}

	want := []*TrafficReferrer{{
		Referrer: String("Google"),
		Count:    Int(4),
		Uniques:  Int(3),
	}}
	if !reflect.DeepEqual(referrers, want) {
		t.Errorf("Repositories.ListReferrers returned %+v, want %+v", referrers, want)
	}

}

func TestRepositoriesService_ListTrafficPaths(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/repos/o/r/traffic/popular/paths", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testHeader(t, r, "Accept", mediaTypeTrafficPreview)
		fmt.Fprintf(w, `[{
			"path": "/github/hubot",
			"title": "github/hubot: A customizable life embetterment robot.",
			"count": 3542,
			"uniques": 2225
 		}]`)
	})
	paths, _, err := client.Repositories.ListTrafficPaths("o", "r")
	if err != nil {
		t.Errorf("Repositories.ListPaths returned error: %+v", err)
	}

	want := []*TrafficPath{{
		Path:    String("/github/hubot"),
		Title:   String("github/hubot: A customizable life embetterment robot."),
		Count:   Int(3542),
		Uniques: Int(2225),
	}}
	if !reflect.DeepEqual(paths, want) {
		t.Errorf("Repositories.ListPaths returned %+v, want %+v", paths, want)
	}

}

func TestRepositoriesService_ListTrafficViews(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/repos/o/r/traffic/views", func(w http.ResponseWriter, r *http.Request) {
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

	views, _, err := client.Repositories.ListTrafficViews("o", "r", nil)
	if err != nil {
		t.Errorf("Repositories.ListPaths returned error: %+v", err)
	}

	want := &TrafficViews{
		Views: []*TrafficData{{
			Timestamp: &TimestampMS{time.Unix(1464710400, 0)},
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

func TestRepositoriesService_ListTrafficClones(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/repos/o/r/traffic/clones", func(w http.ResponseWriter, r *http.Request) {
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

	clones, _, err := client.Repositories.ListTrafficClones("o", "r", nil)
	if err != nil {
		t.Errorf("Repositories.ListPaths returned error: %+v", err)
	}

	want := &TrafficClones{
		Clones: []*TrafficData{{
			Timestamp: &TimestampMS{time.Unix(1464710400, 0)},
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
