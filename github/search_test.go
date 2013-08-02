package github

import (
	"fmt"
	"net/http"
	"net/url"
	"reflect"
	"strconv"
	"strings"
	"testing"
)

func TestRepositories(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/search/repositories", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testHeaderExperimental(t, r)

		query, opts, err := searchOptsFromQueryString(r.URL.RawQuery)
		if err != nil {
			t.Errorf("Could not parse query string: %v", err)
		}
		if query == "foo" && opts.Sort == "" && opts.Order == "" {
			fmt.Fprint(w, `{
"total_count": 1,
"items": [{ "id": 1 }]
}`)
		} else if query == "bar" && opts.Sort == "stars" && opts.Order == "asc" {
			fmt.Fprintf(w, `{
"total_count": 2,
"items": [{ "id": 1 }, { "id": 2}]
}`)
		} else if query == "baz" && opts.Page == 2 && opts.PerPage == 1 {
			fmt.Fprintf(w, `{
"total_count": 10,
"items": [{ "id": 2}]
}`)
		}
	})

	// test case 1
	result, err := client.Search.Repositories("foo", nil)
	if err != nil {
		t.Errorf("Repositories search returned error: %v", err)
	}
	want := &RepositoriesSearchResult{
		TotalCount: 1,
		Repos:      []Repository{{ID: 1}},
	}
	if !reflect.DeepEqual(result, want) {
		t.Errorf("Search.Repoitories returned %+v, want %+v", result, want)
	}

	// test case 2
	result, err = client.Search.Repositories("bar", &SearchOptions{Sort: "stars", Order: SortOrder_Asc})
	if err != nil {
		t.Errorf("Repositories search returned error: %v", err)
	}
	want = &RepositoriesSearchResult{
		TotalCount: 2,
		Repos:      []Repository{{ID: 1}, {ID: 2}},
	}
	if !reflect.DeepEqual(result, want) {
		t.Errorf("Search.Repoitories returned %+v, want %+v", result, want)
	}

	// test case 3
	result, err = client.Search.Repositories("baz", &SearchOptions{Page: 2, PerPage: 1})
	if err != nil {
		t.Errorf("Repositories search returned error: %v", err)
	}
	want = &RepositoriesSearchResult{
		TotalCount: 10,
		Repos:      []Repository{{ID: 2}},
	}
	if !reflect.DeepEqual(result, want) {
		t.Errorf("Search.Repositories returned %+v, want %+v", result, want)
	}
}

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

func testHeaderExperimental(t *testing.T, r *http.Request) {
	if !strings.Contains(r.Header.Get("Accept"), mimePreview) {
		t.Errorf("Header does not contain Accept:%s", mimePreview)
	}
}
