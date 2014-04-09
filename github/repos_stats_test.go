// Copyright 2014 The go-github AUTHORS. All rights reserved.
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

func TestRepositoriesService_Participation(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/repos/o/r/stats/participation", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")

		fmt.Fprint(w, `
{
  "all": [
    11,21,15,2,8,1,8,23,17,21,11,10,33,
    91,38,34,22,23,32,3,43,87,71,18,13,5,
    13,16,66,27,12,45,110,117,13,8,18,9,19,
    26,39,12,20,31,46,91,45,10,24,9,29,7
  ],
  "owner": [
    3,2,3,0,2,0,5,14,7,9,1,5,0,
    48,19,2,0,1,10,2,23,40,35,8,8,2,
    10,6,30,0,2,9,53,104,3,3,10,4,7,
    11,21,4,4,22,26,63,11,2,14,1,10,3
  ]
}
`)
	})

	participation, _, err := client.Repositories.ListParticipation("o", "r")
	if err != nil {
		t.Errorf("RepositoriesService.ListParticipation returned error: %v", err)
	}

	want := &RepositoryParticipation{
		All: []int{
			11, 21, 15, 2, 8, 1, 8, 23, 17, 21, 11, 10, 33,
			91, 38, 34, 22, 23, 32, 3, 43, 87, 71, 18, 13, 5,
			13, 16, 66, 27, 12, 45, 110, 117, 13, 8, 18, 9, 19,
			26, 39, 12, 20, 31, 46, 91, 45, 10, 24, 9, 29, 7,
		},
		Owner: []int{
			3, 2, 3, 0, 2, 0, 5, 14, 7, 9, 1, 5, 0,
			48, 19, 2, 0, 1, 10, 2, 23, 40, 35, 8, 8, 2,
			10, 6, 30, 0, 2, 9, 53, 104, 3, 3, 10, 4, 7,
			11, 21, 4, 4, 22, 26, 63, 11, 2, 14, 1, 10, 3,
		},
	}

	if !reflect.DeepEqual(participation, want) {
		t.Errorf("RepositoriesService.ListParticipation returned %+v, want %+v", participation, want)
	}
}
