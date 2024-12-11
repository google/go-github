// Copyright 2023 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

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
	t.Parallel()
	client, mux, _ := setup(t)

	input := &markdownRenderRequest{
		Text:    Ptr("# text #"),
		Mode:    Ptr("gfm"),
		Context: Ptr("google/go-github"),
	}
	mux.HandleFunc("/markdown", func(w http.ResponseWriter, r *http.Request) {
		v := new(markdownRenderRequest)
		assertNilError(t, json.NewDecoder(r.Body).Decode(v))

		testMethod(t, r, "POST")
		if !cmp.Equal(v, input) {
			t.Errorf("Request body = %+v, want %+v", v, input)
		}
		fmt.Fprint(w, `<h1>text</h1>`)
	})

	ctx := context.Background()
	md, _, err := client.Markdown.Render(ctx, "# text #", &MarkdownOptions{
		Mode:    "gfm",
		Context: "google/go-github",
	})
	if err != nil {
		t.Errorf("Render returned error: %v", err)
	}

	if want := "<h1>text</h1>"; want != md {
		t.Errorf("Render returned %+v, want %+v", md, want)
	}

	const methodName = "Render"
	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Markdown.Render(ctx, "# text #", &MarkdownOptions{
			Mode:    "gfm",
			Context: "google/go-github",
		})
		if got != "" {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestMarkdownRenderRequest_Marshal(t *testing.T) {
	t.Parallel()
	testJSONMarshal(t, &markdownRenderRequest{}, "{}")

	a := &markdownRenderRequest{
		Text:    Ptr("txt"),
		Mode:    Ptr("mode"),
		Context: Ptr("ctx"),
	}

	want := `{
		"text": "txt",
		"mode": "mode",
		"context": "ctx"
	}`

	testJSONMarshal(t, a, want)
}
