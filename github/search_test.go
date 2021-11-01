// Copyright 2013 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package github

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestSearchService_Repositories(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/search/repositories", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testFormValues(t, r, values{
			"q":        "blah",
			"sort":     "forks",
			"order":    "desc",
			"page":     "2",
			"per_page": "2",
		})

		fmt.Fprint(w, `{"total_count": 4, "incomplete_results": false, "items": [{"id":1},{"id":2}]}`)
	})

	opts := &SearchOptions{Sort: "forks", Order: "desc", ListOptions: ListOptions{Page: 2, PerPage: 2}}
	ctx := context.Background()
	result, _, err := client.Search.Repositories(ctx, "blah", opts)
	if err != nil {
		t.Errorf("Search.Repositories returned error: %v", err)
	}

	want := &RepositoriesSearchResult{
		Total:             Int(4),
		IncompleteResults: Bool(false),
		Repositories:      []*Repository{{ID: Int64(1)}, {ID: Int64(2)}},
	}
	if !cmp.Equal(result, want) {
		t.Errorf("Search.Repositories returned %+v, want %+v", result, want)
	}
}

func TestSearchService_Repositories_coverage(t *testing.T) {
	client, _, _, teardown := setup()
	defer teardown()

	ctx := context.Background()

	const methodName = "Repositories"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Search.Repositories(ctx, "\n", nil)
		return err
	})
}

func TestSearchService_Topics(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/search/topics", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testFormValues(t, r, values{
			"q":        "blah",
			"page":     "2",
			"per_page": "2",
		})

		fmt.Fprint(w, `{"total_count": 4, "incomplete_results": false, "items": [{"name":"blah"},{"name":"blahblah"}]}`)
	})

	opts := &SearchOptions{ListOptions: ListOptions{Page: 2, PerPage: 2}}
	ctx := context.Background()
	result, _, err := client.Search.Topics(ctx, "blah", opts)
	if err != nil {
		t.Errorf("Search.Topics returned error: %v", err)
	}

	want := &TopicsSearchResult{
		Total:             Int(4),
		IncompleteResults: Bool(false),
		Topics:            []*TopicResult{{Name: String("blah")}, {Name: String("blahblah")}},
	}
	if !cmp.Equal(result, want) {
		t.Errorf("Search.Topics returned %+v, want %+v", result, want)
	}
}

func TestSearchService_Topics_coverage(t *testing.T) {
	client, _, _, teardown := setup()
	defer teardown()

	ctx := context.Background()

	const methodName = "Topics"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Search.Topics(ctx, "\n", nil)
		return err
	})
}

func TestSearchService_Commits(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/search/commits", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testFormValues(t, r, values{
			"q":     "blah",
			"sort":  "author-date",
			"order": "desc",
		})

		fmt.Fprint(w, `{"total_count": 4, "incomplete_results": false, "items": [{"sha":"random_hash1"},{"sha":"random_hash2"}]}`)
	})

	opts := &SearchOptions{Sort: "author-date", Order: "desc"}
	ctx := context.Background()
	result, _, err := client.Search.Commits(ctx, "blah", opts)
	if err != nil {
		t.Errorf("Search.Commits returned error: %v", err)
	}

	want := &CommitsSearchResult{
		Total:             Int(4),
		IncompleteResults: Bool(false),
		Commits:           []*CommitResult{{SHA: String("random_hash1")}, {SHA: String("random_hash2")}},
	}
	if !cmp.Equal(result, want) {
		t.Errorf("Search.Commits returned %+v, want %+v", result, want)
	}
}

func TestSearchService_Commits_coverage(t *testing.T) {
	client, _, _, teardown := setup()
	defer teardown()

	ctx := context.Background()

	const methodName = "Commits"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Search.Commits(ctx, "\n", nil)
		return err
	})
}

func TestSearchService_Issues(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/search/issues", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testFormValues(t, r, values{
			"q":        "blah",
			"sort":     "forks",
			"order":    "desc",
			"page":     "2",
			"per_page": "2",
		})

		fmt.Fprint(w, `{"total_count": 4, "incomplete_results": true, "items": [{"number":1},{"number":2}]}`)
	})

	opts := &SearchOptions{Sort: "forks", Order: "desc", ListOptions: ListOptions{Page: 2, PerPage: 2}}
	ctx := context.Background()
	result, _, err := client.Search.Issues(ctx, "blah", opts)
	if err != nil {
		t.Errorf("Search.Issues returned error: %v", err)
	}

	want := &IssuesSearchResult{
		Total:             Int(4),
		IncompleteResults: Bool(true),
		Issues:            []*Issue{{Number: Int(1)}, {Number: Int(2)}},
	}
	if !cmp.Equal(result, want) {
		t.Errorf("Search.Issues returned %+v, want %+v", result, want)
	}
}

func TestSearchService_Issues_coverage(t *testing.T) {
	client, _, _, teardown := setup()
	defer teardown()

	ctx := context.Background()

	const methodName = "Issues"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Search.Issues(ctx, "\n", nil)
		return err
	})
}

func TestSearchService_Issues_withQualifiersNoOpts(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	const q = "gopher is:issue label:bug language:c++ pushed:>=2018-01-01 stars:>=200"

	var requestURI string
	mux.HandleFunc("/search/issues", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testFormValues(t, r, values{
			"q": q,
		})
		requestURI = r.RequestURI

		fmt.Fprint(w, `{"total_count": 4, "incomplete_results": true, "items": [{"number":1},{"number":2}]}`)
	})

	opts := &SearchOptions{}
	ctx := context.Background()
	result, _, err := client.Search.Issues(ctx, q, opts)
	if err != nil {
		t.Errorf("Search.Issues returned error: %v", err)
	}

	if want := "/api-v3/search/issues?q=gopher+is%3Aissue+label%3Abug+language%3Ac%2B%2B+pushed%3A%3E%3D2018-01-01+stars%3A%3E%3D200"; requestURI != want {
		t.Fatalf("URI encoding failed: got %v, want %v", requestURI, want)
	}

	want := &IssuesSearchResult{
		Total:             Int(4),
		IncompleteResults: Bool(true),
		Issues:            []*Issue{{Number: Int(1)}, {Number: Int(2)}},
	}
	if !cmp.Equal(result, want) {
		t.Errorf("Search.Issues returned %+v, want %+v", result, want)
	}
}

func TestSearchService_Issues_withQualifiersAndOpts(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	const q = "gopher is:issue label:bug language:c++ pushed:>=2018-01-01 stars:>=200"

	var requestURI string
	mux.HandleFunc("/search/issues", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testFormValues(t, r, values{
			"q":    q,
			"sort": "forks",
		})
		requestURI = r.RequestURI

		fmt.Fprint(w, `{"total_count": 4, "incomplete_results": true, "items": [{"number":1},{"number":2}]}`)
	})

	opts := &SearchOptions{Sort: "forks"}
	ctx := context.Background()
	result, _, err := client.Search.Issues(ctx, q, opts)
	if err != nil {
		t.Errorf("Search.Issues returned error: %v", err)
	}

	if want := "/api-v3/search/issues?q=gopher+is%3Aissue+label%3Abug+language%3Ac%2B%2B+pushed%3A%3E%3D2018-01-01+stars%3A%3E%3D200&sort=forks"; requestURI != want {
		t.Fatalf("URI encoding failed: got %v, want %v", requestURI, want)
	}

	want := &IssuesSearchResult{
		Total:             Int(4),
		IncompleteResults: Bool(true),
		Issues:            []*Issue{{Number: Int(1)}, {Number: Int(2)}},
	}
	if !cmp.Equal(result, want) {
		t.Errorf("Search.Issues returned %+v, want %+v", result, want)
	}
}

func TestSearchService_Users(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/search/users", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testFormValues(t, r, values{
			"q":        "blah",
			"sort":     "forks",
			"order":    "desc",
			"page":     "2",
			"per_page": "2",
		})

		fmt.Fprint(w, `{"total_count": 4, "incomplete_results": false, "items": [{"id":1},{"id":2}]}`)
	})

	opts := &SearchOptions{Sort: "forks", Order: "desc", ListOptions: ListOptions{Page: 2, PerPage: 2}}
	ctx := context.Background()
	result, _, err := client.Search.Users(ctx, "blah", opts)
	if err != nil {
		t.Errorf("Search.Issues returned error: %v", err)
	}

	want := &UsersSearchResult{
		Total:             Int(4),
		IncompleteResults: Bool(false),
		Users:             []*User{{ID: Int64(1)}, {ID: Int64(2)}},
	}
	if !cmp.Equal(result, want) {
		t.Errorf("Search.Users returned %+v, want %+v", result, want)
	}
}

func TestSearchService_Users_coverage(t *testing.T) {
	client, _, _, teardown := setup()
	defer teardown()

	ctx := context.Background()

	const methodName = "Users"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Search.Users(ctx, "\n", nil)
		return err
	})
}

func TestSearchService_Code(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/search/code", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testFormValues(t, r, values{
			"q":        "blah",
			"sort":     "forks",
			"order":    "desc",
			"page":     "2",
			"per_page": "2",
		})

		fmt.Fprint(w, `{"total_count": 4, "incomplete_results": false, "items": [{"name":"1"},{"name":"2"}]}`)
	})

	opts := &SearchOptions{Sort: "forks", Order: "desc", ListOptions: ListOptions{Page: 2, PerPage: 2}}
	ctx := context.Background()
	result, _, err := client.Search.Code(ctx, "blah", opts)
	if err != nil {
		t.Errorf("Search.Code returned error: %v", err)
	}

	want := &CodeSearchResult{
		Total:             Int(4),
		IncompleteResults: Bool(false),
		CodeResults:       []*CodeResult{{Name: String("1")}, {Name: String("2")}},
	}
	if !cmp.Equal(result, want) {
		t.Errorf("Search.Code returned %+v, want %+v", result, want)
	}
}

func TestSearchService_Code_coverage(t *testing.T) {
	client, _, _, teardown := setup()
	defer teardown()

	ctx := context.Background()

	const methodName = "Code"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Search.Code(ctx, "\n", nil)
		return err
	})
}

func TestSearchService_CodeTextMatch(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/search/code", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")

		textMatchResponse := `
		{
			"total_count": 1,
			"incomplete_results": false,
			"items": [
				{
					"name":"gopher1",
					"text_matches": [
						{
							"fragment": "I'm afraid my friend what you have found\nIs a gopher who lives to feed",
							"matches": [
								{
									"text": "gopher",
									"indices": [
										14,
										21
							  	]
								}
						  ]
					  }
				  ]
				}
			]
		}
    `

		fmt.Fprint(w, textMatchResponse)
	})

	opts := &SearchOptions{Sort: "forks", Order: "desc", ListOptions: ListOptions{Page: 2, PerPage: 2}, TextMatch: true}
	ctx := context.Background()
	result, _, err := client.Search.Code(ctx, "blah", opts)
	if err != nil {
		t.Errorf("Search.Code returned error: %v", err)
	}

	wantedCodeResult := &CodeResult{
		Name: String("gopher1"),
		TextMatches: []*TextMatch{{
			Fragment: String("I'm afraid my friend what you have found\nIs a gopher who lives to feed"),
			Matches:  []*Match{{Text: String("gopher"), Indices: []int{14, 21}}},
		},
		},
	}

	want := &CodeSearchResult{
		Total:             Int(1),
		IncompleteResults: Bool(false),
		CodeResults:       []*CodeResult{wantedCodeResult},
	}
	if !cmp.Equal(result, want) {
		t.Errorf("Search.Code returned %+v, want %+v", result, want)
	}
}

func TestSearchService_Labels(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/search/labels", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testFormValues(t, r, values{
			"repository_id": "1234",
			"q":             "blah",
			"sort":          "updated",
			"order":         "desc",
			"page":          "2",
			"per_page":      "2",
		})

		fmt.Fprint(w, `{"total_count": 4, "incomplete_results": false, "items": [{"id": 1234, "name":"bug", "description": "some text"},{"id": 4567, "name":"feature"}]}`)
	})

	opts := &SearchOptions{Sort: "updated", Order: "desc", ListOptions: ListOptions{Page: 2, PerPage: 2}}
	ctx := context.Background()
	result, _, err := client.Search.Labels(ctx, 1234, "blah", opts)
	if err != nil {
		t.Errorf("Search.Code returned error: %v", err)
	}

	want := &LabelsSearchResult{
		Total:             Int(4),
		IncompleteResults: Bool(false),
		Labels: []*LabelResult{
			{ID: Int64(1234), Name: String("bug"), Description: String("some text")},
			{ID: Int64(4567), Name: String("feature")},
		},
	}
	if !cmp.Equal(result, want) {
		t.Errorf("Search.Labels returned %+v, want %+v", result, want)
	}
}

func TestSearchService_Labels_coverage(t *testing.T) {
	client, _, _, teardown := setup()
	defer teardown()

	ctx := context.Background()

	const methodName = "Labels"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Search.Labels(ctx, -1234, "\n", nil)
		return err
	})
}

func TestMatch_Marshal(t *testing.T) {
	testJSONMarshal(t, &Match{}, "{}")

	u := &Match{
		Text:    String("txt"),
		Indices: []int{1},
	}

	want := `{
		"text": "txt",
		"indices": [1]
	}`

	testJSONMarshal(t, u, want)
}

func TestTextMatch_Marshal(t *testing.T) {
	testJSONMarshal(t, &TextMatch{}, "{}")

	u := &TextMatch{
		ObjectURL:  String("ourl"),
		ObjectType: String("otype"),
		Property:   String("prop"),
		Fragment:   String("fragment"),
		Matches: []*Match{
			{
				Text:    String("txt"),
				Indices: []int{1},
			},
		},
	}

	want := `{
		"object_url": "ourl",
		"object_type": "otype",
		"property": "prop",
		"fragment": "fragment",
		"matches": [{
			"text": "txt",
			"indices": [1]
		}]
	}`

	testJSONMarshal(t, u, want)
}

func TestTopicResult_Marshal(t *testing.T) {
	testJSONMarshal(t, &TopicResult{}, "{}")

	u := &TopicResult{
		Name:             String("name"),
		DisplayName:      String("displayName"),
		ShortDescription: String("shortDescription"),
		Description:      String("description"),
		CreatedBy:        String("createdBy"),
		UpdatedAt:        String("2021-10-26"),
		Featured:         Bool(false),
		Curated:          Bool(true),
		Score:            Float64(99.9),
	}

	want := `{
		"name": "name",
		"display_name": "displayName",
		"short_description": "shortDescription",
		"description": "description",
		"created_by": "createdBy",
		"updated_at": "2021-10-26",
		"featured": false,
		"curated": true,
		"score": 99.9
	}`

	testJSONMarshal(t, u, want)
}

func TestRepositoriesSearchResult_Marshal(t *testing.T) {
	testJSONMarshal(t, &RepositoriesSearchResult{}, "{}")

	u := &RepositoriesSearchResult{
		Total:             Int(0),
		IncompleteResults: Bool(true),
		Repositories:      []*Repository{{ID: Int64(1)}},
	}

	want := `{
		"total_count" : 0,
		"incomplete_results" : true,
		"items" : [{"id":1}]
	}`

	testJSONMarshal(t, u, want)
}

func TestCommitsSearchResult_Marshal(t *testing.T) {
	testJSONMarshal(t, &CommitsSearchResult{}, "{}")

	c := &CommitsSearchResult{
		Total:             Int(0),
		IncompleteResults: Bool(true),
		Commits: []*CommitResult{{
			SHA: String("s"),
		}},
	}

	want := `{
		"total_count" : 0,
		"incomplete_results" : true,
		"items" : [{"sha" : "s"}]
	}`

	testJSONMarshal(t, c, want)
}

func TestTopicsSearchResult_Marshal(t *testing.T) {
	testJSONMarshal(t, &TopicsSearchResult{}, "{}")

	u := &TopicsSearchResult{
		Total:             Int(2),
		IncompleteResults: Bool(false),
		Topics: []*TopicResult{
			{
				Name:             String("t1"),
				DisplayName:      String("tt"),
				ShortDescription: String("t desc"),
				Description:      String("desc"),
				CreatedBy:        String("mi"),
				CreatedAt:        &Timestamp{referenceTime},
				UpdatedAt:        String("2006-01-02T15:04:05Z"),
				Featured:         Bool(true),
				Curated:          Bool(true),
				Score:            Float64(123),
			},
		},
	}

	want := `{
		"total_count" : 2,
		"incomplete_results" : false,
		"items" : [
			{
				"name" : "t1",
				"display_name":"tt",
				"short_description":"t desc",
				"description":"desc",
				"created_by":"mi",
				"created_at":` + referenceTimeStr + `,
				"updated_at":"2006-01-02T15:04:05Z",
				"featured":true,
				"curated":true,
				"score":123
			}
		]
	}`

	testJSONMarshal(t, u, want)
}

func TestLabelResult_Marshal(t *testing.T) {
	testJSONMarshal(t, &LabelResult{}, "{}")

	u := &LabelResult{
		ID:          Int64(11),
		URL:         String("url"),
		Name:        String("label"),
		Color:       String("green"),
		Default:     Bool(true),
		Description: String("desc"),
		Score:       Float64(123),
	}

	want := `{
		"id":11,
		"url":"url",
		"name":"label",
		"color":"green",
		"default":true,
		"description":"desc",
		"score":123
	}`

	testJSONMarshal(t, u, want)
}

func TestSearchOptions_Marshal(t *testing.T) {
	testJSONMarshal(t, &SearchOptions{}, "{}")

	u := &SearchOptions{
		Sort:      "author-date",
		Order:     "asc",
		TextMatch: false,
		ListOptions: ListOptions{
			Page:    int(1),
			PerPage: int(10),
		},
	}

	want := `{	
		"sort": "author-date",
		"order": "asc",
		"page": 1,
		"perPage": 10
      }`

	testJSONMarshal(t, u, want)
}

func TestIssuesSearchResult_Marshal(t *testing.T) {
	testJSONMarshal(t, &IssuesSearchResult{}, "{}")

	u := &IssuesSearchResult{
		Total:             Int(48),
		IncompleteResults: Bool(false),
		Issues: []*Issue{
			{
				ID:                Int64(1),
				Number:            Int(1),
				State:             String("s"),
				Locked:            Bool(false),
				Title:             String("title"),
				Body:              String("body"),
				AuthorAssociation: String("aa"),
				User:              &User{ID: Int64(1)},
				Labels:            []*Label{{ID: Int64(1)}},
				Assignee:          &User{ID: Int64(1)},
				Comments:          Int(1),
				ClosedAt:          &referenceTime,
				CreatedAt:         &referenceTime,
				UpdatedAt:         &referenceTime,
				ClosedBy:          &User{ID: Int64(1)},
				URL:               String("url"),
				HTMLURL:           String("hurl"),
				CommentsURL:       String("curl"),
				EventsURL:         String("eurl"),
				LabelsURL:         String("lurl"),
				RepositoryURL:     String("rurl"),
				Milestone:         &Milestone{ID: Int64(1)},
				PullRequestLinks:  &PullRequestLinks{URL: String("url")},
				Repository:        &Repository{ID: Int64(1)},
				Reactions:         &Reactions{TotalCount: Int(1)},
				Assignees:         []*User{{ID: Int64(1)}},
				NodeID:            String("nid"),
				TextMatches:       []*TextMatch{{ObjectURL: String("ourl")}},
				ActiveLockReason:  String("alr"),
			},
		},
	}

	want := `{
		"total_count": 48,
		"incomplete_results": false,
		"items": [
			{
				"id": 1,
				"number": 1,
				"state": "s",
				"locked": false,
				"title": "title",
				"body": "body",
				"author_association": "aa",
				"user": {
					"id": 1
				},
				"labels": [
					{
						"id": 1
					}
				],
				"assignee": {
					"id": 1
				},
				"comments": 1,
				"closed_at": ` + referenceTimeStr + `,
				"created_at": ` + referenceTimeStr + `,
				"updated_at": ` + referenceTimeStr + `,
				"closed_by": {
					"id": 1
				},
				"url": "url",
				"html_url": "hurl",
				"comments_url": "curl",
				"events_url": "eurl",
				"labels_url": "lurl",
				"repository_url": "rurl",
				"milestone": {
					"id": 1
				},
				"pull_request": {
					"url": "url"
				},
				"repository": {
					"id": 1
				},
				"reactions": {
					"total_count": 1
				},
				"assignees": [
					{
						"id": 1
					}
				],
				"node_id": "nid",
				"text_matches": [
					{
						"object_url": "ourl"
					}
				],
				"active_lock_reason": "alr"
			}
		]
	}`

	testJSONMarshal(t, u, want)
}
