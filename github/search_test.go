package github

import (
	"fmt"
	"net/http"
	"reflect"
	"strings"
	"testing"
)

func TestSearch_Repositories_opts(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/search/repositories", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testHeaderPreview(t, r)
		testFormValues(t, r, values{
			"q":        "blah",
			"sort":     "forks",
			"order":    "desc",
			"page":     "2",
			"per_page": "2",
		})

		fmt.Fprint(w, `{"total_count": 4, "items": [{"id":1},{"id":2}]}`)
	})

	opts := &SearchOptions{Sort: "forks", Order: "desc", Page: 2, PerPage: 2}
	result, _, err := client.Search.Repositories("blah", opts)
	if err != nil {
		t.Errorf("Search.Repositories returned error: %v", err)
	}

	want := &RepositoriesSearchResult{
		Total: 4,
		Repos: []Repository{{ID: 1}, {ID: 2}},
	}
	if !reflect.DeepEqual(result, want) {
		t.Errorf("Search.Repositories returned %+v, want %+v", result, want)
	}
}

func TestSearch_Repositories_noOptsNoResults(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/search/repositories", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testHeaderPreview(t, r)
		testFormValues(t, r, values{
			"q": "blah",
		})

		fmt.Fprint(w, `{"total_count": 0, "items": []}`)
	})

	result, _, err := client.Search.Repositories("blah", nil)
	if err != nil {
		t.Errorf("Search.Repositories returned error: %v", err)
	}

	want := &RepositoriesSearchResult{
		Total: 0,
		Repos: []Repository{},
	}
	if !reflect.DeepEqual(result, want) {
		t.Errorf("Search.Repositories returned %+v, want %+v", result, want)
	}
}

func TestSearch_Issues_opts(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/search/issues", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testHeaderPreview(t, r)
		testFormValues(t, r, values{
			"q":        "blah",
			"sort":     "forks",
			"order":    "desc",
			"page":     "2",
			"per_page": "2",
		})

		fmt.Fprint(w, `{"total_count": 4, "items": [{"number":1},{"number":2}]}`)
	})

	opts := &SearchOptions{Sort: "forks", Order: "desc", Page: 2, PerPage: 2}
	result, _, err := client.Search.Issues("blah", opts)
	if err != nil {
		t.Errorf("Search.Issues returned error: %v", err)
	}

	want := &IssuesSearchResult{
		Total:  4,
		Issues: []Issue{{Number: 1}, {Number: 2}},
	}
	if !reflect.DeepEqual(result, want) {
		t.Errorf("Search.Issues returned %+v, want %+v", result, want)
	}
}

func TestSearch_Issues_noOptsNoResults(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/search/issues", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testHeaderPreview(t, r)
		testFormValues(t, r, values{
			"q": "blah",
		})

		fmt.Fprint(w, `{"total_count": 0, "items": []}`)
	})

	result, _, err := client.Search.Issues("blah", nil)
	if err != nil {
		t.Errorf("Search.Issues returned error: %v", err)
	}

	want := &IssuesSearchResult{
		Total:  0,
		Issues: []Issue{},
	}
	if !reflect.DeepEqual(result, want) {
		t.Errorf("Search.Issues returned %+v, want %+v", result, want)
	}
}

func TestSearch_Users_opts(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/search/users", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testHeaderPreview(t, r)
		testFormValues(t, r, values{
			"q":        "blah",
			"sort":     "forks",
			"order":    "desc",
			"page":     "2",
			"per_page": "2",
		})

		fmt.Fprint(w, `{"total_count": 4, "items": [{"id":1},{"id":2}]}`)
	})

	opts := &SearchOptions{Sort: "forks", Order: "desc", Page: 2, PerPage: 2}
	result, _, err := client.Search.Users("blah", opts)
	if err != nil {
		t.Errorf("Search.Issues returned error: %v", err)
	}

	want := &UsersSearchResult{
		Total: 4,
		Users: []User{{ID: 1}, {ID: 2}},
	}
	if !reflect.DeepEqual(result, want) {
		t.Errorf("Search.Users returned %+v, want %+v", result, want)
	}
}

func TestSearch_Users_noOptsNoResults(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/search/users", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testHeaderPreview(t, r)
		testFormValues(t, r, values{
			"q": "blah",
		})

		fmt.Fprint(w, `{"total_count": 0, "items": []}`)
	})

	result, _, err := client.Search.Users("blah", nil)
	if err != nil {
		t.Errorf("Search.Users returned error: %v", err)
	}

	want := &UsersSearchResult{
		Total: 0,
		Users: []User{},
	}
	if !reflect.DeepEqual(result, want) {
		t.Errorf("Search.Users returned %+v, want %+v", result, want)
	}
}

func TestSearch_Code_opts(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/search/code", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testHeaderPreview(t, r)
		testFormValues(t, r, values{
			"q":        "blah",
			"sort":     "forks",
			"order":    "desc",
			"page":     "2",
			"per_page": "2",
		})

		fmt.Fprint(w, `{"total_count": 4, "items": [{"name":"1"},{"name":"2"}]}`)
	})

	opts := &SearchOptions{Sort: "forks", Order: "desc", Page: 2, PerPage: 2}
	result, _, err := client.Search.Code("blah", opts)
	if err != nil {
		t.Errorf("Search.Code returned error: %v", err)
	}

	want := &CodeSearchResult{
		Total:       4,
		CodeResults: []CodeResult{{Name: "1"}, {Name: "2"}},
	}
	if !reflect.DeepEqual(result, want) {
		t.Errorf("Search.Code returned %+v, want %+v", result, want)
	}
}

func TestSearch_Code_noOptsNoResults(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/search/code", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testHeaderPreview(t, r)
		testFormValues(t, r, values{
			"q": "blah",
		})

		fmt.Fprint(w, `{"total_count": 0, "items": []}`)
	})

	result, _, err := client.Search.Code("blah", nil)
	if err != nil {
		t.Errorf("Search.Code returned error: %v", err)
	}

	want := &CodeSearchResult{
		Total:       0,
		CodeResults: []CodeResult{},
	}
	if !reflect.DeepEqual(result, want) {
		t.Errorf("Search.Code returned %+v, want %+v", result, want)
	}
}

// Checks that HTTP request contains GitHub preview accept header
func testHeaderPreview(t *testing.T, r *http.Request) {
	if !strings.Contains(r.Header.Get("Accept"), mimePreview) {
		t.Errorf("Header does not contain Accept:%s", mimePreview)
	}
}
