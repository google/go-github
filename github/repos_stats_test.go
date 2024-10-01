// Copyright 2014 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package github

import (
	"context"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
)

func TestRepositoriesService_ListContributorsStats(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/repos/o/r/stats/contributors", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")

		fmt.Fprint(w, `
[
  {
    "author": {
			"id": 1,
			"node_id": "nodeid-1"
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

	ctx := context.Background()
	stats, _, err := client.Repositories.ListContributorsStats(ctx, "o", "r")
	if err != nil {
		t.Errorf("RepositoriesService.ListContributorsStats returned error: %v", err)
	}

	want := []*ContributorStats{
		{
			Author: &Contributor{
				ID:     Int64(1),
				NodeID: String("nodeid-1"),
			},
			Total: Int(135),
			Weeks: []*WeeklyStats{
				{
					Week:      &Timestamp{time.Date(2013, time.May, 05, 00, 00, 00, 0, time.UTC).Local()},
					Additions: Int(6898),
					Deletions: Int(77),
					Commits:   Int(10),
				},
			},
		},
	}

	if !cmp.Equal(stats, want) {
		t.Errorf("RepositoriesService.ListContributorsStats returned %+v, want %+v", stats, want)
	}

	const methodName = "ListContributorsStats"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Repositories.ListContributorsStats(ctx, "\n", "\n")
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Repositories.ListContributorsStats(ctx, "o", "r")
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestRepositoriesService_ListCommitActivity(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

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

	ctx := context.Background()
	activity, _, err := client.Repositories.ListCommitActivity(ctx, "o", "r")
	if err != nil {
		t.Errorf("RepositoriesService.ListCommitActivity returned error: %v", err)
	}

	want := []*WeeklyCommitActivity{
		{
			Days:  []int{0, 3, 26, 20, 39, 1, 0},
			Total: Int(89),
			Week:  &Timestamp{time.Date(2012, time.May, 06, 05, 00, 00, 0, time.UTC).Local()},
		},
	}

	if !cmp.Equal(activity, want) {
		t.Errorf("RepositoriesService.ListCommitActivity returned %+v, want %+v", activity, want)
	}

	const methodName = "ListCommitActivity"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Repositories.ListCommitActivity(ctx, "\n", "\n")
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Repositories.ListCommitActivity(ctx, "o", "r")
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestRepositoriesService_ListCodeFrequency(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/repos/o/r/stats/code_frequency", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")

		fmt.Fprint(w, `[[1302998400, 1124, -435]]`)
	})

	ctx := context.Background()
	code, _, err := client.Repositories.ListCodeFrequency(ctx, "o", "r")
	if err != nil {
		t.Errorf("RepositoriesService.ListCodeFrequency returned error: %v", err)
	}

	want := []*WeeklyStats{{
		Week:      &Timestamp{time.Date(2011, time.April, 17, 00, 00, 00, 0, time.UTC).Local()},
		Additions: Int(1124),
		Deletions: Int(-435),
	}}

	if !cmp.Equal(code, want) {
		t.Errorf("RepositoriesService.ListCodeFrequency returned %+v, want %+v", code, want)
	}

	const methodName = "ListCodeFrequency"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Repositories.ListCodeFrequency(ctx, "\n", "\n")
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Repositories.ListCodeFrequency(ctx, "o", "r")
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestRepositoriesService_Participation(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

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

	ctx := context.Background()
	participation, _, err := client.Repositories.ListParticipation(ctx, "o", "r")
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

	if !cmp.Equal(participation, want) {
		t.Errorf("RepositoriesService.ListParticipation returned %+v, want %+v", participation, want)
	}

	const methodName = "ListParticipation"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Repositories.ListParticipation(ctx, "\n", "\n")
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Repositories.ListParticipation(ctx, "o", "r")
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestRepositoriesService_ListPunchCard(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/repos/o/r/stats/punch_card", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")

		fmt.Fprint(w, `[
		  [0, 0, 5],
		  [0, 1, 43],
		  [0, 2, 21]
		]`)
	})

	ctx := context.Background()
	card, _, err := client.Repositories.ListPunchCard(ctx, "o", "r")
	if err != nil {
		t.Errorf("RepositoriesService.ListPunchCard returned error: %v", err)
	}

	want := []*PunchCard{
		{Day: Int(0), Hour: Int(0), Commits: Int(5)},
		{Day: Int(0), Hour: Int(1), Commits: Int(43)},
		{Day: Int(0), Hour: Int(2), Commits: Int(21)},
	}

	if !cmp.Equal(card, want) {
		t.Errorf("RepositoriesService.ListPunchCard returned %+v, want %+v", card, want)
	}

	const methodName = "ListPunchCard"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Repositories.ListPunchCard(ctx, "\n", "\n")
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Repositories.ListPunchCard(ctx, "o", "r")
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestRepositoriesService_AcceptedError(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/repos/o/r/stats/contributors", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		// This response indicates the fork will happen asynchronously.
		w.WriteHeader(http.StatusAccepted)
		fmt.Fprint(w, `{"id":1}`)
	})

	ctx := context.Background()
	stats, _, err := client.Repositories.ListContributorsStats(ctx, "o", "r")
	if err == nil {
		t.Errorf("RepositoriesService.AcceptedError should have returned an error")
	}

	if _, ok := err.(*AcceptedError); !ok {
		t.Errorf("RepositoriesService.AcceptedError returned an AcceptedError: %v", err)
	}

	if stats != nil {
		t.Errorf("RepositoriesService.AcceptedError expected stats to be nil: %v", stats)
	}

	const methodName = "ListContributorsStats"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Repositories.ListContributorsStats(ctx, "o", "r")
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Repositories.ListContributorsStats(ctx, "o", "r")
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestRepositoryParticipation_Marshal(t *testing.T) {
	t.Parallel()
	testJSONMarshal(t, &RepositoryParticipation{}, "{}")

	u := &RepositoryParticipation{
		All:   []int{1},
		Owner: []int{1},
	}

	want := `{
		"all": [1],
		"owner": [1]
	}`

	testJSONMarshal(t, u, want)
}

func TestWeeklyCommitActivity_Marshal(t *testing.T) {
	t.Parallel()
	testJSONMarshal(t, &WeeklyCommitActivity{}, "{}")

	u := &WeeklyCommitActivity{
		Days:  []int{1},
		Total: Int(1),
		Week:  &Timestamp{referenceTime},
	}

	want := `{
		"days": [
			1
		],
		"total": 1,
		"week": ` + referenceTimeStr + `
	}`

	testJSONMarshal(t, u, want)
}

func TestWeeklyStats_Marshal(t *testing.T) {
	t.Parallel()
	testJSONMarshal(t, &WeeklyStats{}, "{}")

	u := &WeeklyStats{
		Week:      &Timestamp{referenceTime},
		Additions: Int(1),
		Deletions: Int(1),
		Commits:   Int(1),
	}

	want := `{
		"w": ` + referenceTimeStr + `,
		"a": 1,
		"d": 1,
		"c": 1
	}`

	testJSONMarshal(t, u, want)
}

func TestContributorStats_Marshal(t *testing.T) {
	t.Parallel()
	testJSONMarshal(t, &ContributorStats{}, "{}")

	u := &ContributorStats{
		Author: &Contributor{ID: Int64(1)},
		Total:  Int(1),
		Weeks: []*WeeklyStats{
			{
				Week:      &Timestamp{referenceTime},
				Additions: Int(1),
				Deletions: Int(1),
				Commits:   Int(1),
			},
		},
	}

	want := `{
		"author": {
			"id": 1
		},
		"total": 1,
		"weeks": [
			{
				"w": ` + referenceTimeStr + `,
				"a": 1,
				"d": 1,
				"c": 1
			}
		]
	}`

	testJSONMarshal(t, u, want)
}
