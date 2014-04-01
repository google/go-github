// Copyright 2014 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package github

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"testing"
)

func TestMarkdown(t *testing.T) {
	setup()
	defer teardown()

	input := &markdownRequest{
		Text:    String("# text #"),
		Mode:    String("gfm"),
		Context: String("google/go-github"),
	}
	mux.HandleFunc("/markdown", func(w http.ResponseWriter, r *http.Request) {
		v := new(markdownRequest)
		json.NewDecoder(r.Body).Decode(v)

		testMethod(t, r, "POST")
		if !reflect.DeepEqual(v, input) {
			t.Errorf("Request body = %+v, want %+v", v, input)
		}
		fmt.Fprint(w, `<h1>text</h1>`)
	})

	md, _, err := client.Markdown("# text #", &MarkdownOptions{
		Mode:    "gfm",
		Context: "google/go-github",
	})
	if err != nil {
		t.Errorf("Markdown returned error: %v", err)
	}

	if want := "<h1>text</h1>"; want != md {
		t.Errorf("Markdown returned %+v, want %+v", md, want)
	}
}
