package github

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"testing"

	"github.com/google/go-cmp/cmp"
)

// FixtureManager handles loading and updating test fixtures.
type FixtureManager struct {
	dir        string
	updateMode bool
	token      string
	client     *Client
}

// NewFixtureManager creates a new fixture manager.
func NewFixtureManager(t *testing.T, dir string) *FixtureManager {
	updateMode := os.Getenv("GITHUB_UPDATE_FIXTURE") != ""
	token := os.Getenv("GITHUB_AUTH_TOKEN")

	var client *Client
	if updateMode && token != "" {
		client = NewClient(nil).WithAuthToken(token)
	}

	return &FixtureManager{
		dir:        dir,
		updateMode: updateMode,
		token:      token,
		client:     client,
	}
}

// LoadOrFetch loads a fixture from file or fetches from GitHub API and saves it.
// The fetch function receives the client and should make a real API call.
func (fm *FixtureManager) LoadOrFetch(t *testing.T, name string, fetch func(*Client) (any, error)) []byte {
	path := filepath.Join(fm.dir, name+".json")

	if fm.updateMode {
		t.Logf("Updating fixture: %s from GitHub API", name)
		data, err := fetch(fm.client)
		if err != nil {
			t.Fatalf("Failed to fetch fixture data from GitHub API: %v", err)
		}

		if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
			t.Fatalf("Failed to create fixture directory: %v", err)
		}

		// Marshal to JSON
		prettyBytes, err := json.MarshalIndent(data, "", "  ")
		if err != nil {
			t.Fatalf("Failed to marshal fixture data: %v", err)
		}

		if err := os.WriteFile(path, prettyBytes, 0644); err != nil {
			t.Fatalf("Failed to write fixture file: %v", err)
		}

		return prettyBytes
	}

	bytes, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("Failed to load fixture file %s: %v", path, err)
	}

	return bytes
}

func TestPagination_Scan2_ListComments(t *testing.T) {
	t.Parallel()

	client, mux, _ := setup(t)
	fm := NewFixtureManager(t, filepath.Join("testdata", "pagination"))

	issue := 2618
	mux.HandleFunc(fmt.Sprintf("/repos/google/go-github/issues/%d/comments", issue), func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		page := r.URL.Query().Get("page")
		var fixture []byte
		loaderFn := func(page int) func(client *Client) (any, error) {
			opts := IssueListCommentsOptions{
				ListOptions: ListOptions{
					PerPage: 5,
					Page:    page,
				},
			}
			return func(client *Client) (any, error) {
				comments, resp, err := client.Issues.ListComments(t.Context(), "google", "go-github", issue, &opts)
				fmt.Printf("Link: %s\n\n", resp.Header["Link"][0])
				return comments, err
			}
		}
		switch page {
		case "", "1":
			w.Header().Set("Link", `<https://api.github.com/repositories/10270722/issues/2618/comments?page=2&per_page=5>; rel="next", <https://api.github.com/repositories/10270722/issues/2618/comments?page=5&per_page=5>; rel="last"`)
			fixture = fm.LoadOrFetch(t, "list_comments_page1", loaderFn(1))
			w.Write(fixture)
		case "2":
			w.Header().Set("Link", `<https://api.github.com/repositories/10270722/issues/2618/comments?page=1&per_page=5>; rel="prev", <https://api.github.com/repositories/10270722/issues/2618/comments?page=3&per_page=5>; rel="next", <https://api.github.com/repositories/10270722/issues/2618/comments?page=5&per_page=5>; rel="last", <https://api.github.com/repositories/10270722/issues/2618/comments?page=1&per_page=5>; rel="first"`)
			fixture = fm.LoadOrFetch(t, "list_comments_page2", loaderFn(2))
			w.Write(fixture)
		case "3":
			w.Header().Set("Link", `<https://api.github.com/repositories/10270722/issues/2618/comments?page=2&per_page=5>; rel="prev", <https://api.github.com/repositories/10270722/issues/2618/comments?page=4&per_page=5>; rel="next", <https://api.github.com/repositories/10270722/issues/2618/comments?page=5&per_page=5>; rel="last", <https://api.github.com/repositories/10270722/issues/2618/comments?page=1&per_page=5>; rel="first"`)
			fixture = fm.LoadOrFetch(t, "list_comments_page3", loaderFn(3))
			w.Write(fixture)
		case "4":
			w.Header().Set("Link", `<https://api.github.com/repositories/10270722/issues/2618/comments?page=3&per_page=5>; rel="prev", <https://api.github.com/repositories/10270722/issues/2618/comments?page=5&per_page=5>; rel="next", <https://api.github.com/repositories/10270722/issues/2618/comments?page=5&per_page=5>; rel="last", <https://api.github.com/repositories/10270722/issues/2618/comments?page=1&per_page=5>; rel="first"`)
			fixture = fm.LoadOrFetch(t, "list_comments_page4", loaderFn(4))
			w.Write(fixture)
		case "5":
			w.Header().Set("Link", `<https://api.github.com/repositories/10270722/issues/2618/comments?page=4&per_page=5>; rel="prev", <https://api.github.com/repositories/10270722/issues/2618/comments?page=1&per_page=5>; rel="first"`)
			fixture = fm.LoadOrFetch(t, "list_comments_page5", loaderFn(5))
			w.Write(fixture)
		}
	})

	ctx := t.Context()
	opts := &IssueListCommentsOptions{}

	var comments []*IssueComment
	for c, err := range Scan2(func(p PaginationOption) ([]*IssueComment, *Response, error) {
		return client.Issues.ListComments(ctx, "google", "go-github", issue, opts, p)
	}) {
		if err != nil {
			t.Fatalf("Scan2 iterator returned error: %v", err)
		}
		comments = append(comments, c)
	}

	wantCommentIDs := []int64{
		1372144817, 1386971275, 1396478534, 1446756890, 1656588501,
		1659541626, 1664125370, 1684770158, 1685120763, 1884764520,
		1912173472, 1912236179, 1912244287, 1912386258, 1919918050,
		1919936396, 1919948684, 1920002236, 1920009324, 1920084186,
		1979228975, 1994323766, 1994383750, 2708656405,
	}
	commentIDs := make([]int64, len(comments))
	for i, c := range comments {
		commentIDs[i] = c.GetID()
	}

	if !cmp.Equal(commentIDs, wantCommentIDs) {
		t.Errorf("Got %+v, want %+v", commentIDs, wantCommentIDs)
	}
}
