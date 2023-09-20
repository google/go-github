package github

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestMarkdownService_Markdown(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	input := &markdownRequest{
		Text:    String("# text #"),
		Mode:    String("gfm"),
		Context: String("google/go-github"),
	}
	mux.HandleFunc("/markdown", func(w http.ResponseWriter, r *http.Request) {
		v := new(markdownRequest)
		assertNilError(t, json.NewDecoder(r.Body).Decode(v))

		testMethod(t, r, "POST")
		if !cmp.Equal(v, input) {
			t.Errorf("Request body = %+v, want %+v", v, input)
		}
		fmt.Fprint(w, `<h1>text</h1>`)
	})

	ctx := context.Background()
	md, _, err := client.Markdown.Markdown(ctx, "# text #", &MarkdownOptions{
		Mode:    "gfm",
		Context: "google/go-github",
	})
	if err != nil {
		t.Errorf("Markdown returned error: %v", err)
	}

	if want := "<h1>text</h1>"; want != md {
		t.Errorf("Markdown returned %+v, want %+v", md, want)
	}

	const methodName = "Markdown"
	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Markdown.Markdown(ctx, "# text #", &MarkdownOptions{
			Mode:    "gfm",
			Context: "google/go-github",
		})
		if got != "" {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestMarkdownRequest_Marshal(t *testing.T) {
	testJSONMarshal(t, &markdownRequest{}, "{}")

	a := &markdownRequest{
		Text:    String("txt"),
		Mode:    String("mode"),
		Context: String("ctx"),
	}

	want := `{
		"text": "txt",
		"mode": "mode",
		"context": "ctx"
	}`

	testJSONMarshal(t, a, want)
}