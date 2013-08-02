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

		query, sort, order, page, perPage, err := searchOptsFromQueryString(r.URL.RawQuery)
		if err != nil {
			t.Errorf("Could not parse query string: %v", err)
		}
		if query == "foo" && sort == "" && order == "" {
			fmt.Fprint(w, `{
"total_count": 1,
"items": [{ "id": 1 }]
}`)
		} else if query == "bar" && sort == "stars" && order == "asc" {
			fmt.Fprintf(w, `{
"total_count": 2,
"items": [{ "id": 1 }, { "id": 2}]
}`)
		} else if query == "baz" && page == 2 && perPage == 1 {
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

func searchOptsFromQueryString(queryString string) (query string, sort string, order string, page int, perPage int, err error) {
	values, err := url.ParseQuery(queryString)
	if err != nil {
		return
	}

	if qs, in := values["q"]; in && len(qs) == 1 {
		query = qs[0]
	}
	if ss, in := values["sort"]; in && len(ss) == 1 {
		sort = ss[0]
	}
	if os, in := values["order"]; in && len(os) == 1 {
		order = os[0]
	}
	if ps, in := values["page"]; in && len(ps) == 1 {
		page, err = strconv.Atoi(ps[0])
		if err != nil {
			return
		}
	}
	if pps, in := values["per_page"]; in && len(pps) == 1 {
		perPage, err = strconv.Atoi(pps[0])
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
