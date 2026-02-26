// Copyright 2013 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package github

import (
	"fmt"
	"net/http"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestSearchService_Repositories(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

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
	ctx := t.Context()
	result, _, err := client.Search.Repositories(ctx, "blah", opts)
	if err != nil {
		t.Errorf("Search.Repositories returned error: %v", err)
	}

	want := &RepositoriesSearchResult{
		Total:             Ptr(4),
		IncompleteResults: Ptr(false),
		Repositories:      []*Repository{{ID: Ptr(int64(1))}, {ID: Ptr(int64(2))}},
	}
	if !cmp.Equal(result, want) {
		t.Errorf("Search.Repositories returned %+v, want %+v", result, want)
	}
	const methodName = "Repositories"
	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Search.Repositories(ctx, "blah", opts)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestSearchService_Repositories_coverage(t *testing.T) {
	t.Parallel()
	client, _, _ := setup(t)

	ctx := t.Context()

	const methodName = "Repositories"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Search.Repositories(ctx, "\n", nil)
		return err
	})
}

func TestSearchService_RepositoriesTextMatch(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/search/repositories", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		textMatchResponse := `
			{
				"total_count": 1,
				"incomplete_results": false,
				"items": [
					{
						"name":"gopher1"
					}
				]
			}
		`
		list := strings.Split(r.Header.Get("Accept"), ",")
		aMap := make(map[string]struct{})
		for _, s := range list {
			aMap[strings.TrimSpace(s)] = struct{}{}
		}
		if _, ok := aMap["application/vnd.github.v3.text-match+json"]; ok {
			textMatchResponse = `
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
		}

		fmt.Fprint(w, textMatchResponse)
	})

	opts := &SearchOptions{Sort: "forks", Order: "desc", ListOptions: ListOptions{Page: 2, PerPage: 2}, TextMatch: true}
	ctx := t.Context()
	result, _, err := client.Search.Repositories(ctx, "blah", opts)
	if err != nil {
		t.Errorf("Search.Code returned error: %v", err)
	}

	wantedRepoResult := &Repository{
		Name: Ptr("gopher1"),
		TextMatches: []*TextMatch{
			{
				Fragment: Ptr("I'm afraid my friend what you have found\nIs a gopher who lives to feed"),
				Matches:  []*Match{{Text: Ptr("gopher"), Indices: []int{14, 21}}},
			},
		},
	}

	want := &RepositoriesSearchResult{
		Total:             Ptr(1),
		IncompleteResults: Ptr(false),
		Repositories:      []*Repository{wantedRepoResult},
	}
	if !cmp.Equal(result, want) {
		t.Errorf("Search.Repo returned %+v, want %+v", result, want)
	}
}

func TestSearchService_Topics(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

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
	ctx := t.Context()
	result, _, err := client.Search.Topics(ctx, "blah", opts)
	if err != nil {
		t.Errorf("Search.Topics returned error: %v", err)
	}

	want := &TopicsSearchResult{
		Total:             Ptr(4),
		IncompleteResults: Ptr(false),
		Topics:            []*TopicResult{{Name: Ptr("blah")}, {Name: Ptr("blahblah")}},
	}
	if !cmp.Equal(result, want) {
		t.Errorf("Search.Topics returned %+v, want %+v", result, want)
	}
	const methodName = "Topics"
	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Search.Topics(ctx, "blah", opts)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestSearchService_Topics_coverage(t *testing.T) {
	t.Parallel()
	client, _, _ := setup(t)

	ctx := t.Context()

	const methodName = "Topics"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Search.Topics(ctx, "\n", nil)
		return err
	})
}

func TestSearchService_Commits(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

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
	ctx := t.Context()
	result, _, err := client.Search.Commits(ctx, "blah", opts)
	if err != nil {
		t.Errorf("Search.Commits returned error: %v", err)
	}

	want := &CommitsSearchResult{
		Total:             Ptr(4),
		IncompleteResults: Ptr(false),
		Commits:           []*CommitResult{{SHA: Ptr("random_hash1")}, {SHA: Ptr("random_hash2")}},
	}
	if !cmp.Equal(result, want) {
		t.Errorf("Search.Commits returned %+v, want %+v", result, want)
	}
	const methodName = "Commits"
	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Search.Commits(ctx, "blah", opts)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestSearchService_Commits_coverage(t *testing.T) {
	t.Parallel()
	client, _, _ := setup(t)

	ctx := t.Context()

	const methodName = "Commits"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Search.Commits(ctx, "\n", nil)
		return err
	})
}

func TestSearchService_Issues(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

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
	ctx := t.Context()
	result, _, err := client.Search.Issues(ctx, "blah", opts)
	if err != nil {
		t.Errorf("Search.Issues returned error: %v", err)
	}

	want := &IssuesSearchResult{
		Total:             Ptr(4),
		IncompleteResults: Ptr(true),
		Issues:            []*Issue{{Number: Ptr(1)}, {Number: Ptr(2)}},
	}
	if !cmp.Equal(result, want) {
		t.Errorf("Search.Issues returned %+v, want %+v", result, want)
	}
	const methodName = "Issues"
	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Search.Issues(ctx, "blah", opts)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestSearchService_Issues_advancedSearch(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/search/issues", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testFormValues(t, r, values{
			"q":               "blah",
			"sort":            "forks",
			"order":           "desc",
			"page":            "2",
			"per_page":        "2",
			"advanced_search": "true",
		})

		fmt.Fprint(w, `{"total_count": 4, "incomplete_results": true, "items": [{"number":1},{"number":2}]}`)
	})

	opts := &SearchOptions{Sort: "forks", Order: "desc", ListOptions: ListOptions{Page: 2, PerPage: 2}, AdvancedSearch: Ptr(true)}
	ctx := t.Context()
	result, _, err := client.Search.Issues(ctx, "blah", opts)
	if err != nil {
		t.Errorf("Search.Issues_advancedSearch returned error: %v", err)
	}

	want := &IssuesSearchResult{
		Total:             Ptr(4),
		IncompleteResults: Ptr(true),
		Issues:            []*Issue{{Number: Ptr(1)}, {Number: Ptr(2)}},
	}
	if !cmp.Equal(result, want) {
		t.Errorf("Search.Issues_advancedSearch returned %+v, want %+v", result, want)
	}
}

func TestSearchService_Issues_coverage(t *testing.T) {
	t.Parallel()
	client, _, _ := setup(t)

	ctx := t.Context()

	const methodName = "Issues"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Search.Issues(ctx, "\n", nil)
		return err
	})
}

func TestSearchService_Issues_withQualifiersNoOpts(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

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
	ctx := t.Context()
	result, _, err := client.Search.Issues(ctx, q, opts)
	if err != nil {
		t.Errorf("Search.Issues returned error: %v", err)
	}

	if want := "/api-v3/search/issues?q=gopher+is%3Aissue+label%3Abug+language%3Ac%2B%2B+pushed%3A%3E%3D2018-01-01+stars%3A%3E%3D200"; requestURI != want {
		t.Fatalf("URI encoding failed: got %v, want %v", requestURI, want)
	}

	want := &IssuesSearchResult{
		Total:             Ptr(4),
		IncompleteResults: Ptr(true),
		Issues:            []*Issue{{Number: Ptr(1)}, {Number: Ptr(2)}},
	}
	if !cmp.Equal(result, want) {
		t.Errorf("Search.Issues returned %+v, want %+v", result, want)
	}
}

func TestSearchService_Issues_withQualifiersAndOpts(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

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
	ctx := t.Context()
	result, _, err := client.Search.Issues(ctx, q, opts)
	if err != nil {
		t.Errorf("Search.Issues returned error: %v", err)
	}

	if want := "/api-v3/search/issues?q=gopher+is%3Aissue+label%3Abug+language%3Ac%2B%2B+pushed%3A%3E%3D2018-01-01+stars%3A%3E%3D200&sort=forks"; requestURI != want {
		t.Fatalf("URI encoding failed: got %v, want %v", requestURI, want)
	}

	want := &IssuesSearchResult{
		Total:             Ptr(4),
		IncompleteResults: Ptr(true),
		Issues:            []*Issue{{Number: Ptr(1)}, {Number: Ptr(2)}},
	}
	if !cmp.Equal(result, want) {
		t.Errorf("Search.Issues returned %+v, want %+v", result, want)
	}
}

func TestSearchService_Users(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

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
	ctx := t.Context()
	result, _, err := client.Search.Users(ctx, "blah", opts)
	if err != nil {
		t.Errorf("Search.Issues returned error: %v", err)
	}

	want := &UsersSearchResult{
		Total:             Ptr(4),
		IncompleteResults: Ptr(false),
		Users:             []*User{{ID: Ptr(int64(1))}, {ID: Ptr(int64(2))}},
	}
	if !cmp.Equal(result, want) {
		t.Errorf("Search.Users returned %+v, want %+v", result, want)
	}
	const methodName = "Users"
	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Search.Users(ctx, "blah", opts)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestSearchService_Users_coverage(t *testing.T) {
	t.Parallel()
	client, _, _ := setup(t)

	ctx := t.Context()

	const methodName = "Users"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Search.Users(ctx, "\n", nil)
		return err
	})
}

func TestSearchService_Code(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

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
	ctx := t.Context()
	result, _, err := client.Search.Code(ctx, "blah", opts)
	if err != nil {
		t.Errorf("Search.Code returned error: %v", err)
	}

	want := &CodeSearchResult{
		Total:             Ptr(4),
		IncompleteResults: Ptr(false),
		CodeResults:       []*CodeResult{{Name: Ptr("1")}, {Name: Ptr("2")}},
	}
	if !cmp.Equal(result, want) {
		t.Errorf("Search.Code returned %+v, want %+v", result, want)
	}
	const methodName = "Code"
	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Search.Code(ctx, "blah", opts)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestSearchService_Code_coverage(t *testing.T) {
	t.Parallel()
	client, _, _ := setup(t)

	ctx := t.Context()

	const methodName = "Code"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Search.Code(ctx, "\n", nil)
		return err
	})
}

func TestSearchService_CodeTextMatch(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

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
	ctx := t.Context()
	result, _, err := client.Search.Code(ctx, "blah", opts)
	if err != nil {
		t.Errorf("Search.Code returned error: %v", err)
	}

	wantedCodeResult := &CodeResult{
		Name: Ptr("gopher1"),
		TextMatches: []*TextMatch{
			{
				Fragment: Ptr("I'm afraid my friend what you have found\nIs a gopher who lives to feed"),
				Matches:  []*Match{{Text: Ptr("gopher"), Indices: []int{14, 21}}},
			},
		},
	}

	want := &CodeSearchResult{
		Total:             Ptr(1),
		IncompleteResults: Ptr(false),
		CodeResults:       []*CodeResult{wantedCodeResult},
	}
	if !cmp.Equal(result, want) {
		t.Errorf("Search.Code returned %+v, want %+v", result, want)
	}
}

func TestSearchService_Labels(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

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
	ctx := t.Context()
	result, _, err := client.Search.Labels(ctx, 1234, "blah", opts)
	if err != nil {
		t.Errorf("Search.Code returned error: %v", err)
	}

	want := &LabelsSearchResult{
		Total:             Ptr(4),
		IncompleteResults: Ptr(false),
		Labels: []*LabelResult{
			{ID: Ptr(int64(1234)), Name: Ptr("bug"), Description: Ptr("some text")},
			{ID: Ptr(int64(4567)), Name: Ptr("feature")},
		},
	}
	if !cmp.Equal(result, want) {
		t.Errorf("Search.Labels returned %+v, want %+v", result, want)
	}
	const methodName = "Labels"
	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Search.Labels(ctx, 1234, "blah", opts)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestSearchService_Labels_coverage(t *testing.T) {
	t.Parallel()
	client, _, _ := setup(t)

	ctx := t.Context()

	const methodName = "Labels"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Search.Labels(ctx, -1234, "\n", nil)
		return err
	})
}

func TestMatch_Marshal(t *testing.T) {
	t.Parallel()
	testJSONMarshal(t, &Match{}, "{}")

	u := &Match{
		Text:    Ptr("txt"),
		Indices: []int{1},
	}

	want := `{
		"text": "txt",
		"indices": [1]
	}`

	testJSONMarshal(t, u, want)
}

func TestTextMatch_Marshal(t *testing.T) {
	t.Parallel()
	testJSONMarshal(t, &TextMatch{}, "{}")

	u := &TextMatch{
		ObjectURL:  Ptr("ourl"),
		ObjectType: Ptr("otype"),
		Property:   Ptr("prop"),
		Fragment:   Ptr("fragment"),
		Matches: []*Match{
			{
				Text:    Ptr("txt"),
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
	t.Parallel()
	testJSONMarshal(t, &TopicResult{}, "{}")

	u := &TopicResult{
		Name:             Ptr("name"),
		DisplayName:      Ptr("displayName"),
		ShortDescription: Ptr("shortDescription"),
		Description:      Ptr("description"),
		CreatedBy:        Ptr("createdBy"),
		UpdatedAt:        Ptr("2021-10-26"),
		Featured:         Ptr(false),
		Curated:          Ptr(true),
		Score:            Ptr(99.9),
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
	t.Parallel()
	testJSONMarshal(t, &RepositoriesSearchResult{}, "{}")

	u := &RepositoriesSearchResult{
		Total:             Ptr(0),
		IncompleteResults: Ptr(true),
		Repositories:      []*Repository{{ID: Ptr(int64(1))}},
	}

	want := `{
		"total_count" : 0,
		"incomplete_results" : true,
		"items" : [{"id":1}]
	}`

	testJSONMarshal(t, u, want)
}

func TestCommitsSearchResult_Marshal(t *testing.T) {
	t.Parallel()
	testJSONMarshal(t, &CommitsSearchResult{}, "{}")

	c := &CommitsSearchResult{
		Total:             Ptr(0),
		IncompleteResults: Ptr(true),
		Commits: []*CommitResult{{
			SHA: Ptr("s"),
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
	t.Parallel()
	testJSONMarshal(t, &TopicsSearchResult{}, "{}")

	u := &TopicsSearchResult{
		Total:             Ptr(2),
		IncompleteResults: Ptr(false),
		Topics: []*TopicResult{
			{
				Name:             Ptr("t1"),
				DisplayName:      Ptr("tt"),
				ShortDescription: Ptr("t desc"),
				Description:      Ptr("desc"),
				CreatedBy:        Ptr("mi"),
				CreatedAt:        &Timestamp{referenceTime},
				UpdatedAt:        Ptr("2006-01-02T15:04:05Z"),
				Featured:         Ptr(true),
				Curated:          Ptr(true),
				Score:            Ptr(123.0),
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
	t.Parallel()
	testJSONMarshal(t, &LabelResult{}, "{}")

	u := &LabelResult{
		ID:          Ptr(int64(11)),
		URL:         Ptr("url"),
		Name:        Ptr("label"),
		Color:       Ptr("green"),
		Default:     Ptr(true),
		Description: Ptr("desc"),
		Score:       Ptr(123.0),
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

func TestIssuesSearchResult_Marshal(t *testing.T) {
	t.Parallel()
	testJSONMarshal(t, &IssuesSearchResult{}, "{}")

	u := &IssuesSearchResult{
		Total:             Ptr(48),
		IncompleteResults: Ptr(false),
		Issues: []*Issue{
			{
				ID:                Ptr(int64(1)),
				Number:            Ptr(1),
				State:             Ptr("s"),
				Locked:            Ptr(false),
				Title:             Ptr("title"),
				Body:              Ptr("body"),
				AuthorAssociation: Ptr("aa"),
				User:              &User{ID: Ptr(int64(1))},
				Labels:            []*Label{{ID: Ptr(int64(1))}},
				Assignee:          &User{ID: Ptr(int64(1))},
				Comments:          Ptr(1),
				ClosedAt:          &Timestamp{referenceTime},
				CreatedAt:         &Timestamp{referenceTime},
				UpdatedAt:         &Timestamp{referenceTime},
				ClosedBy:          &User{ID: Ptr(int64(1))},
				URL:               Ptr("url"),
				HTMLURL:           Ptr("hurl"),
				CommentsURL:       Ptr("curl"),
				EventsURL:         Ptr("eurl"),
				LabelsURL:         Ptr("lurl"),
				RepositoryURL:     Ptr("rurl"),
				Milestone:         &Milestone{ID: Ptr(int64(1))},
				PullRequestLinks:  &PullRequestLinks{URL: Ptr("url")},
				Repository:        &Repository{ID: Ptr(int64(1))},
				Reactions:         &Reactions{TotalCount: Ptr(1)},
				Assignees:         []*User{{ID: Ptr(int64(1))}},
				NodeID:            Ptr("nid"),
				TextMatches:       []*TextMatch{{ObjectURL: Ptr("ourl")}},
				ActiveLockReason:  Ptr("alr"),
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

func TestLabelsSearchResult_Marshal(t *testing.T) {
	t.Parallel()
	testJSONMarshal(t, &LabelsSearchResult{}, "{}")

	u := &LabelsSearchResult{
		Total:             Ptr(5),
		IncompleteResults: Ptr(false),
		Labels: []*LabelResult{
			{
				ID:          Ptr(int64(1)),
				URL:         Ptr("https://www.example.com"),
				Name:        Ptr("test name"),
				Color:       Ptr("green"),
				Default:     Ptr(true),
				Description: Ptr("testDescription"),
				Score:       Ptr(1.0),
			},
		},
	}

	want := `{
		"total_count": 5,
		"incomplete_results": false,
		"items": [
			{
				"id": 1,
				"url": "https://www.example.com",
				"name": "test name",
				"color": "green",
				"default": true,
				"description": "testDescription",
				"score": 1
			}
		]
	}`

	testJSONMarshal(t, u, want)
}

func TestCommitResult_Marshal(t *testing.T) {
	t.Parallel()
	testJSONMarshal(t, &CommitResult{}, "{}")

	c := &CommitResult{
		SHA:         Ptr("test"),
		HTMLURL:     Ptr("hurl"),
		CommentsURL: Ptr("curl"),
		URL:         Ptr("url"),
		Repository:  &Repository{ID: Ptr(int64(1))},
		Score:       Ptr(123.0),
		Commit:      &Commit{SHA: Ptr("test")},
		Author:      &User{ID: Ptr(int64(1))},
		Committer:   &User{ID: Ptr(int64(1))},
	}

	want := `{
		"sha": "test",
		"commit": {
		   "sha": "test"
		},
		"author": {
		   "id": 1
		},
		"committer": {
		   "id": 1
		},
		"html_url": "hurl",
		"url": "url",
		"comments_url": "curl",
		"repository": {
		   "id": 1
		},
		"score": 123
	 }`

	testJSONMarshal(t, c, want)
}

func TestUsersSearchResult_Marshal(t *testing.T) {
	t.Parallel()
	testJSONMarshal(t, &UsersSearchResult{}, "{}")

	u := &UsersSearchResult{
		Total:             Ptr(2),
		IncompleteResults: Ptr(false),
		Users: []*User{{
			Login:      Ptr("loginTest"),
			ID:         Ptr(int64(1)),
			NodeID:     Ptr("NodeTest"),
			AvatarURL:  Ptr("AvatarURLTest"),
			HTMLURL:    Ptr("Hurl"),
			GravatarID: Ptr("gravatarIDTest"),
			Name:       Ptr("nameTest"),
		}},
	}

	want := `{
		"total_count": 2,
		"incomplete_results": false,
		"items": [
		   {
			  "login": "loginTest",
			  "id": 1,
			  "node_id": "NodeTest",
			  "avatar_url": "AvatarURLTest",
			  "html_url": "Hurl",
			  "gravatar_id": "gravatarIDTest",
			  "name": "nameTest"
		   }
		]
	 }`

	testJSONMarshal(t, u, want)
}

func TestCodeSearchResult_Marshal(t *testing.T) {
	t.Parallel()
	testJSONMarshal(t, &CodeSearchResult{}, "{}")

	u := &CodeSearchResult{
		Total:             Ptr(4),
		IncompleteResults: Ptr(false),
		CodeResults:       []*CodeResult{{Name: Ptr("n")}},
	}

	want := `{
		"total_count" : 4,
		"incomplete_results" : false,
		"items" : [{"name": "n"}]
	}`

	testJSONMarshal(t, u, want)
}
