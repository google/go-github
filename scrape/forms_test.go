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
	tests := []struct {
		description string
		html        string
		forms       []htmlForm
	}{
		{"no forms", `<html></html>`, nil},
		{"empty form", `<html><form></form></html>`, []htmlForm{{Values: url.Values{}}}},
		{
			"single form with one value",
			`<html><form action="a" method="m"><input name="n1" value="v1"></form></html>`,
			[]htmlForm{{Action: "a", Method: "m", Values: url.Values{"n1": {"v1"}}}},
		},
		{
			"two forms",
			`<html>
			  <form action="a1" method="m1"><input name="n1" value="v1"></form>
			  <form action="a2" method="m2"><input name="n2" value="v2"></form>
			</html>`,
			[]htmlForm{
				{Action: "a1", Method: "m1", Values: url.Values{"n1": {"v1"}}},
				{Action: "a2", Method: "m2", Values: url.Values{"n2": {"v2"}}},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			node, err := html.Parse(strings.NewReader(tt.html))
			if err != nil {
				t.Errorf("error parsing html: %v", err)
			}
			if got, want := parseForms(node), tt.forms; !cmp.Equal(got, want) {
				t.Errorf("parseForms(%q) returned %+v, want %+v", tt.html, got, want)
			}
		})
	}
}

func Test_FetchAndSumbitForm(t *testing.T) {
	client, mux, cleanup := setup()
	defer cleanup()
	var submitted bool

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `<html><form action="/submit">
		  <input type="hidden" name="hidden" value="h">
		  <input type="text" name="name">
		</form></html>`)
	})
	mux.HandleFunc("/submit", func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		want := url.Values{"hidden": {"h"}, "name": {"n"}}
		if got := r.Form; !cmp.Equal(got, want) {
			t.Errorf("submitted form contained values %v, want %v", got, want)
		}
		submitted = true
	})

	setValues := func(values url.Values) { values.Set("name", "n") }
	fetchAndSubmitForm(client.Client, client.baseURL.String()+"/", setValues)
	if !submitted {
		t.Error("form was never submitted")
	}
}
