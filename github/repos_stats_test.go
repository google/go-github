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
)

func TestRepositoriesService_ListContributorsStats(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/repos/o/r/stats/contributors", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")

		fmt.Fprint(w, `
[
  {
    "author": {
      "id": 1
    },
    "total": 135,
    "weeks": [
      {
        "w": 1367712000,
        "a": 6898,
        "d": 77,
        "c": 10
      }
    ]
  }
]
`)
	})

	stats, _, err := client.Repositories.ListContributorsStats("o", "r")
	if err != nil {
		t.Errorf("RepositoriesService.ListContributorsStats returned error: %v", err)
	}

	want := &[]ContributorStats{
		ContributorStats{
			Author: &Contributor{
				ID: Int(1),
			},
			Total: Int(135),
			Weeks: []WeeklyHash{
				WeeklyHash{
					Week:      Int(1367712000),
					Additions: Int(6898),
					Deletions: Int(77),
					Commits:   Int(10),
				},
			},
		},
	}

	if !reflect.DeepEqual(stats, want) {
		t.Errorf("RepositoriesService.ListContributorsStats returned %+v, want %+v", stats, want)
	}
}

func TestRepositoriesService_ListCommitActivity(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/repos/o/r/stats/commit_activity", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")

		fmt.Fprint(w, `
[
  {
    "days": [0, 3, 26, 20, 39, 1, 0],
    "total": 89,
    "week": 1336280400
  }
]
`)
	})

	activity, _, err := client.Repositories.ListCommitActivity("o", "r")
	if err != nil {
		t.Errorf("RepositoriesService.ListCommitActivity returned error: %v", err)
	}

	want := &[]WeeklyCommitActivity{
		WeeklyCommitActivity{
			Days:  []int{0, 3, 26, 20, 39, 1, 0},
			Total: Int(89),
			Week:  Int(1336280400),
		},
	}

	if !reflect.DeepEqual(activity, want) {
		t.Errorf("RepositoriesService.ListCommitActivity returned %+v, want %+v", activity, want)
	}
}
