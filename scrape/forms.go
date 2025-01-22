// Copyright 2013 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// forms.go contains logic for parsing and submitting HTML forms.  None of this
// is specific to go-github in any way, and could easily be pulled out into a
// general purpose scraping library in the future.

package scrape

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"golang.org/x/net/html"
)

// htmlForm represents the basic elements of an HTML Form.
type htmlForm struct {
	// Action is the URL where the form will be submitted
	Action string
	// Method is the HTTP method to use when submitting the form
	Method string
	// Values contains form values to be submitted
	Values url.Values
}

// parseForms parses and returns all form elements beneath node.  Form values
// include all input and textarea elements within the form. The values of radio
// and checkbox inputs are included only if they are checked.
//
// In the future, we might want to allow a custom selector to be passed in to
// further restrict what forms will be returned.
func parseForms(node *html.Node) (forms []*htmlForm) {
	if node == nil {
		return nil
	}

	doc := goquery.NewDocumentFromNode(node)
	doc.Find("form").Each(func(_ int, s *goquery.Selection) {
		form := htmlForm{Values: url.Values{}}
		form.Action, _ = s.Attr("action")
		form.Method, _ = s.Attr("method")

		s.Find("input").Each(func(_ int, s *goquery.Selection) {
			name, _ := s.Attr("name")
			if name == "" {
				return
			}

			typ, _ := s.Attr("type")
			typ = strings.ToLower(typ)
			_, checked := s.Attr("checked")
			if (typ == "radio" || typ == "checkbox") && !checked {
				return
			}

			value, _ := s.Attr("value")
			form.Values.Add(name, value)
		})
		s.Find("textarea").Each(func(_ int, s *goquery.Selection) {
			name, _ := s.Attr("name")
			if name == "" {
				return
			}

			value := s.Text()
			form.Values.Add(name, value)
		})
		forms = append(forms, &form)
	})

	return forms
}

// fetchAndSubmitForm will fetch the page at urlStr, then parse and submit the first form found.
// setValues will be called with the parsed form values, allowing the caller to set any custom
// form values. Form submission will always use the POST method, regardless of the value of the
// method attribute in the form.  The response from submitting the parsed form is returned.
func fetchAndSubmitForm(client *http.Client, urlStr string, setValues func(url.Values)) (*http.Response, error) {
	resp, err := client.Get(urlStr)
	if err != nil {
		return nil, fmt.Errorf("error fetching url %q: %v", urlStr, err)
	}

	defer resp.Body.Close()
	root, err := html.Parse(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error parsing response: %v", err)
	}

	forms := parseForms(root)
	if len(forms) == 0 {
		return nil, fmt.Errorf("no forms found at %q", urlStr)
	}
	form := forms[0]

	actionURL, err := url.Parse(form.Action)
	if err != nil {
		return nil, fmt.Errorf("error parsing form action URL %q: %v", form.Action, err)
	}
	actionURL = resp.Request.URL.ResolveReference(actionURL)

	// allow caller to fill out the form
	if setValues != nil {
		setValues(form.Values)
	}

	resp, err = client.PostForm(actionURL.String(), form.Values)
	if err != nil {
		return nil, fmt.Errorf("error posting form: %v", err)
	}

	return resp, nil
}
