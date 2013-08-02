package github

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"reflect"
	"strconv"
	"strings"
	"testing"
)

type searchTestCase struct {
	Query  string
	Opts   SearchOptions
	Result interface{}
}

type searchFunc func(query string, opt *SearchOptions) (interface{}, error)

func TestRepositories(t *testing.T) {
	testCases := []searchTestCase{
		{
			Query: "blah",
			Opts: SearchOptions{
				Sort:    "forks",
				Order:   Order_Desc,
				Page:    2,
				PerPage: 1,
			},
			Result: &RepositoriesSearchResult{
				Total: 2,
				Repos: []Repository{{ID: 0}, {ID: 1}},
			},
		},
		{
			Query: "foo",
			Result: &RepositoriesSearchResult{
				Total: 1,
				Repos: []Repository{{ID: 1}},
			},
		},
		{
			Query: "bar",
			Opts: SearchOptions{
				Sort:  "stars",
				Order: Order_Asc,
			},
			Result: &RepositoriesSearchResult{
				Total: 1,
				Repos: []Repository{{ID: 1}, {ID: 2}},
			},
		},
		{
			Query: "baz",
			Opts: SearchOptions{
				Page:    2,
				PerPage: 1,
			},
			Result: &RepositoriesSearchResult{
				Total: 10,
				Repos: []Repository{{ID: 2}},
			},
		},
	}
	searcher := func(query string, opts *SearchOptions) (interface{}, error) {
		result, err := client.Search.Repositories(query, opts)
		return result, err
	}
	testSearchOnType(t, Search_Repositories, searcher, testCases)
}

func TestIssues(t *testing.T) {
	testCases := []searchTestCase{
		{
			Query: "blah",
			Opts: SearchOptions{
				Sort:    "forks",
				Order:   Order_Desc,
				Page:    2,
				PerPage: 1,
			},
			Result: &IssuesSearchResult{
				Total:  2,
				Issues: []Issue{{Number: 0}, {Number: 1}},
			},
		},
		{
			Query: "foo",
			Result: &IssuesSearchResult{
				Total:  1,
				Issues: []Issue{{Number: 1}},
			},
		},
		{
			Query: "bar",
			Opts: SearchOptions{
				Sort:  "stars",
				Order: Order_Asc,
			},
			Result: &IssuesSearchResult{
				Total:  1,
				Issues: []Issue{{Number: 1}, {Number: 2}},
			},
		},
		{
			Query: "baz",
			Opts: SearchOptions{
				Page:    2,
				PerPage: 1,
			},
			Result: &IssuesSearchResult{
				Total:  10,
				Issues: []Issue{{Number: 2}},
			},
		},
	}
	searcher := func(query string, opts *SearchOptions) (interface{}, error) {
		result, err := client.Search.Issues(query, opts)
		return result, err
	}
	testSearchOnType(t, Search_Issues, searcher, testCases)
}

func TestUsers(t *testing.T) {
	testCases := []searchTestCase{
		{
			Query: "blah",
			Opts: SearchOptions{
				Sort:    "forks",
				Order:   Order_Desc,
				Page:    2,
				PerPage: 1,
			},
			Result: &UsersSearchResult{
				Total: 2,
				Users: []User{{ID: 0}, {ID: 1}},
			},
		},
		{
			Query: "foo",
			Result: &UsersSearchResult{
				Total: 1,
				Users: []User{{ID: 1}},
			},
		},
		{
			Query: "bar",
			Opts: SearchOptions{
				Sort:  "stars",
				Order: Order_Asc,
			},
			Result: &UsersSearchResult{
				Total: 1,
				Users: []User{{ID: 1}, {ID: 2}},
			},
		},
		{
			Query: "baz",
			Opts: SearchOptions{
				Page:    2,
				PerPage: 1,
			},
			Result: &UsersSearchResult{
				Total: 10,
				Users: []User{{ID: 2}},
			},
		},
	}
	searcher := func(query string, opts *SearchOptions) (interface{}, error) {
		result, err := client.Search.Users(query, opts)
		return result, err
	}
	testSearchOnType(t, Search_Users, searcher, testCases)
}

func TestCode(t *testing.T) {
	testCases := []searchTestCase{
		{
			Query: "blah",
			Opts: SearchOptions{
				Sort:    "forks",
				Order:   Order_Desc,
				Page:    2,
				PerPage: 1,
			},
			Result: &CodeSearchResult{
				Total:       2,
				CodeResults: []CodeResult{{Name: "0"}, {Name: "1"}},
			},
		},
		{
			Query: "foo",
			Result: &CodeSearchResult{
				Total:       1,
				CodeResults: []CodeResult{{Name: "1"}},
			},
		},
		{
			Query: "bar",
			Opts: SearchOptions{
				Sort:  "stars",
				Order: Order_Asc,
			},
			Result: &CodeSearchResult{
				Total:       1,
				CodeResults: []CodeResult{{Name: "1"}, {Name: "2"}},
			},
		},
		{
			Query: "baz",
			Opts: SearchOptions{
				Page:    2,
				PerPage: 1,
			},
			Result: &CodeSearchResult{
				Total:       10,
				CodeResults: []CodeResult{{Name: "2"}},
			},
		},
	}
	searcher := func(query string, opts *SearchOptions) (interface{}, error) {
		result, err := client.Search.Code(query, opts)
		return result, err
	}
	testSearchOnType(t, Search_Code, searcher, testCases)
}

// Helper function that runs test cases against GitHub search API for
// a specific type
func testSearchOnType(t *testing.T, searchType string, searcher searchFunc, testCases []searchTestCase) {
	setup()
	defer teardown()

	mux.HandleFunc(fmt.Sprintf("/search/%s", searchType), func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testHeaderExperimental(t, r)

		query, opts, err := searchOptsFromQueryString(r.URL.RawQuery)
		if err != nil {
			t.Errorf("Could not parse query string: %v", err)
		}

		for _, testCase := range testCases {
			if testCase.Query == query && reflect.DeepEqual(&testCase.Opts, opts) {
				json.NewEncoder(w).Encode(testCase.Result)
				return
			}
		}

		t.Errorf(`Search %s with query "%s" and opts %+v expected to return result, but none found`, searchType, query, *opts)
	})

	for _, testCase := range testCases {
		result, err := searcher(testCase.Query, &testCase.Opts)
		if err != nil {
			t.Errorf(`Search %s with query "%s" and opts %+v returned error: %v`, searchType, testCase.Query, testCase.Opts, err)
		}
		if !reflect.DeepEqual(testCase.Result, result) {
			t.Errorf(`Search %s with query "%s" and options %+v returned %+v, but expected %+v`,
				searchType, testCase.Query, testCase.Opts, result, testCase.Result)
		}
	}
}

// Parses search options from query string
func searchOptsFromQueryString(queryString string) (query string, opts *SearchOptions, err error) {
	values, err := url.ParseQuery(queryString)
	if err != nil {
		return
	}

	if q, in := values["q"]; in && len(q) == 1 {
		query = q[0]
	}

	opts = new(SearchOptions)
	if sort, in := values["sort"]; in && len(sort) == 1 {
		opts.Sort = sort[0]
	}
	if order, in := values["order"]; in && len(order) == 1 {
		opts.Order = order[0]
	}
	if page, in := values["page"]; in && len(page) == 1 {
		opts.Page, err = strconv.Atoi(page[0])
		if err != nil {
			return
		}
	}
	if perPage, in := values["per_page"]; in && len(perPage) == 1 {
		opts.PerPage, err = strconv.Atoi(perPage[0])
		if err != nil {
			return
		}
	}

	return
}

// Checks that HTTP request contains GitHub experimental media accept header
func testHeaderExperimental(t *testing.T, r *http.Request) {
	if !strings.Contains(r.Header.Get("Accept"), mimePreview) {
		t.Errorf("Header does not contain Accept:%s", mimePreview)
	}
}
