// Copyright 2013 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package scrape

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"golang.org/x/net/html"
)

func Test_ParseForms(t *testing.T) {
	t.Parallel()
	tests := []struct {
		description string
		html        string
		forms       []*htmlForm
	}{
		{"no forms", `<html></html>`, nil},
		{"empty form", `<html><form></form></html>`, []*htmlForm{{Values: url.Values{}}}},
		{
			"single form with one value",
			`<html><form action="a" method="m"><input name="n1" value="v1"></form></html>`,
			[]*htmlForm{{Action: "a", Method: "m", Values: url.Values{"n1": {"v1"}}}},
		},
		{
			"two forms",
			`<html>
			  <form action="a1" method="m1"><input name="n1" value="v1"></form>
			  <form action="a2" method="m2"><input name="n2" value="v2"></form>
			</html>`,
			[]*htmlForm{
				{Action: "a1", Method: "m1", Values: url.Values{"n1": {"v1"}}},
				{Action: "a2", Method: "m2", Values: url.Values{"n2": {"v2"}}},
			},
		},
		{
			"form with radio buttons (none checked)",
			`<html><form>
			   <input type="radio" name="n1" value="v1">
			   <input type="radio" name="n1" value="v2">
			   <input type="radio" name="n1" value="v3">
			</form></html>`,
			[]*htmlForm{{Values: url.Values{}}},
		},
		{
			"form with radio buttons",
			`<html><form>
			   <input type="radio" name="n1" value="v1">
			   <input type="radio" name="n1" value="v2">
			   <input type="radio" name="n1" value="v3" checked>
			</form></html>`,
			[]*htmlForm{{Values: url.Values{"n1": {"v3"}}}},
		},
		{
			"form with checkboxes",
			`<html><form>
			   <input type="checkbox" name="n1" value="v1" checked>
			   <input type="checkbox" name="n2" value="v2">
			   <input type="checkbox" name="n3" value="v3" checked>
			</form></html>`,
			[]*htmlForm{{Values: url.Values{"n1": {"v1"}, "n3": {"v3"}}}},
		},
		{
			"single form with textarea",
			`<html><form><textarea name="n1">v1</textarea></form></html>`,
			[]*htmlForm{{Values: url.Values{"n1": {"v1"}}}},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.description, func(t *testing.T) {
			t.Parallel()
			node, err := html.Parse(strings.NewReader(tt.html))
			if err != nil {
				t.Fatalf("error parsing html: %v", err)
			}
			if got, want := parseForms(node), tt.forms; !cmp.Equal(got, want) {
				t.Errorf("parseForms(%q) returned %+v, want %+v", tt.html, got, want)
			}
		})
	}
}

func Test_FetchAndSumbitForm(t *testing.T) {
	t.Parallel()
	client, mux := setup(t)
	var submitted bool

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `<html><form action="/submit">
		  <input type="hidden" name="hidden" value="h">
		  <input type="text" name="name">
		</form></html>`)
	})
	mux.HandleFunc("/submit", func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		if err != nil {
			t.Fatalf("error parsing form: %v", err)
		}
		want := url.Values{"hidden": {"h"}, "name": {"n"}}
		if got := r.Form; !cmp.Equal(got, want) {
			t.Errorf("submitted form contained values %v, want %v", got, want)
		}
		submitted = true
	})

	setValues := func(values url.Values) { values.Set("name", "n") }
	_, err := fetchAndSubmitForm(client.Client, client.baseURL.String()+"/", setValues)
	if err != nil {
		t.Fatalf("fetchAndSubmitForm returned err: %v", err)
	}
	if !submitted {
		t.Error("form was never submitted")
	}
}
