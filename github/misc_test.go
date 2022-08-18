// Copyright 2014 The go-github AUTHORS. All rights reserved.
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

func TestMarkdown(t *testing.T) {
	client, mux, _, teardown := setup()
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
		if !cmp.Equal(v, input) {
			t.Errorf("Request body = %+v, want %+v", v, input)
		}
		fmt.Fprint(w, `<h1>text</h1>`)
	})

	ctx := context.Background()
	md, _, err := client.Markdown(ctx, "# text #", &MarkdownOptions{
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
		got, resp, err := client.Markdown(ctx, "# text #", &MarkdownOptions{
			Mode:    "gfm",
			Context: "google/go-github",
		})
		if got != "" {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestListEmojis(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/emojis", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{"+1": "+1.png"}`)
	})

	ctx := context.Background()
	emoji, _, err := client.ListEmojis(ctx)
	if err != nil {
		t.Errorf("ListEmojis returned error: %v", err)
	}

	want := map[string]string{"+1": "+1.png"}
	if !cmp.Equal(want, emoji) {
		t.Errorf("ListEmojis returned %+v, want %+v", emoji, want)
	}

	const methodName = "ListEmojis"
	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.ListEmojis(ctx)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestListCodesOfConduct(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/codes_of_conduct", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testHeader(t, r, "Accept", mediaTypeCodesOfConductPreview)
		fmt.Fprint(w, `[{
						"key": "key",
						"name": "name",
						"url": "url"}
						]`)
	})

	ctx := context.Background()
	cs, _, err := client.ListCodesOfConduct(ctx)
	if err != nil {
		t.Errorf("ListCodesOfConduct returned error: %v", err)
	}

	want := []*CodeOfConduct{
		{
			Key:  String("key"),
			Name: String("name"),
			URL:  String("url"),
		}}
	if !cmp.Equal(want, cs) {
		t.Errorf("ListCodesOfConduct returned %+v, want %+v", cs, want)
	}

	const methodName = "ListCodesOfConduct"
	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.ListCodesOfConduct(ctx)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestGetCodeOfConduct(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/codes_of_conduct/k", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testHeader(t, r, "Accept", mediaTypeCodesOfConductPreview)
		fmt.Fprint(w, `{
						"key": "key",
						"name": "name",
						"url": "url",
						"body": "body"}`,
		)
	})

	ctx := context.Background()
	coc, _, err := client.GetCodeOfConduct(ctx, "k")
	if err != nil {
		t.Errorf("ListCodesOfConduct returned error: %v", err)
	}

	want := &CodeOfConduct{
		Key:  String("key"),
		Name: String("name"),
		URL:  String("url"),
		Body: String("body"),
	}
	if !cmp.Equal(want, coc) {
		t.Errorf("GetCodeOfConductByKey returned %+v, want %+v", coc, want)
	}

	const methodName = "GetCodeOfConduct"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.GetCodeOfConduct(ctx, "\n")
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.GetCodeOfConduct(ctx, "k")
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestAPIMeta_Marshal(t *testing.T) {
	testJSONMarshal(t, &APIMeta{}, "{}")

	a := &APIMeta{
		Hooks:                            []string{"h"},
		Git:                              []string{"g"},
		VerifiablePasswordAuthentication: Bool(true),
		Pages:                            []string{"p"},
		Importer:                         []string{"i"},
		Actions:                          []string{"a"},
		Dependabot:                       []string{"d"},
		SSHKeyFingerprints:               map[string]string{"a": "f"},
		SSHKeys:                          []string{"k"},
		API:                              []string{"a"},
		Web:                              []string{"w"},
	}
	want := `{
		"hooks":["h"],
		"git":["g"],
		"verifiable_password_authentication":true,
		"pages":["p"],
		"importer":["i"],
		"actions":["a"],
		"dependabot":["d"],
		"ssh_key_fingerprints":{"a":"f"},
		"ssh_keys":["k"],
		"api":["a"],
		"web":["w"]
	}`

	testJSONMarshal(t, a, want)
}

func TestAPIMeta(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/meta", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{"web":["w"],"api":["a"],"hooks":["h"], "git":["g"], "pages":["p"], "importer":["i"], "actions":["a"], "dependabot":["d"], "verifiable_password_authentication": true}`)
	})

	ctx := context.Background()
	meta, _, err := client.APIMeta(ctx)
	if err != nil {
		t.Errorf("APIMeta returned error: %v", err)
	}

	want := &APIMeta{
		Hooks:      []string{"h"},
		Git:        []string{"g"},
		Pages:      []string{"p"},
		Importer:   []string{"i"},
		Actions:    []string{"a"},
		Dependabot: []string{"d"},
		API:        []string{"a"},
		Web:        []string{"w"},

		VerifiablePasswordAuthentication: Bool(true),
	}
	if !cmp.Equal(want, meta) {
		t.Errorf("APIMeta returned %+v, want %+v", meta, want)
	}

	const methodName = "APIMeta"
	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.APIMeta(ctx)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestOctocat(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	input := "input"
	output := "sample text"

	mux.HandleFunc("/octocat", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testFormValues(t, r, values{"s": input})
		w.Header().Set("Content-Type", "application/octocat-stream")
		fmt.Fprint(w, output)
	})

	ctx := context.Background()
	got, _, err := client.Octocat(ctx, input)
	if err != nil {
		t.Errorf("Octocat returned error: %v", err)
	}

	if want := output; got != want {
		t.Errorf("Octocat returned %+v, want %+v", got, want)
	}

	const methodName = "Octocat"
	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Octocat(ctx, input)
		if got != "" {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestZen(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	output := "sample text"

	mux.HandleFunc("/zen", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		w.Header().Set("Content-Type", "text/plain;charset=utf-8")
		fmt.Fprint(w, output)
	})

	ctx := context.Background()
	got, _, err := client.Zen(ctx)
	if err != nil {
		t.Errorf("Zen returned error: %v", err)
	}

	if want := output; got != want {
		t.Errorf("Zen returned %+v, want %+v", got, want)
	}

	const methodName = "Zen"
	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Zen(ctx)
		if got != "" {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestListServiceHooks(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/hooks", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `[{
			"name":"n",
			"events":["e"],
			"supported_events":["s"],
			"schema":[
			  ["a", "b"]
			]
		}]`)
	})

	ctx := context.Background()
	hooks, _, err := client.ListServiceHooks(ctx)
	if err != nil {
		t.Errorf("ListServiceHooks returned error: %v", err)
	}

	want := []*ServiceHook{{
		Name:            String("n"),
		Events:          []string{"e"},
		SupportedEvents: []string{"s"},
		Schema:          [][]string{{"a", "b"}},
	}}
	if !cmp.Equal(hooks, want) {
		t.Errorf("ListServiceHooks returned %+v, want %+v", hooks, want)
	}

	const methodName = "ListServiceHooks"
	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.ListServiceHooks(ctx)
		if got != nil {
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

func TestCodeOfConduct_Marshal(t *testing.T) {
	testJSONMarshal(t, &CodeOfConduct{}, "{}")

	a := &CodeOfConduct{
		Name: String("name"),
		Key:  String("key"),
		URL:  String("url"),
		Body: String("body"),
	}

	want := `{
		"name": "name",
		"key": "key",
		"url": "url",
		"body": "body"
	}`

	testJSONMarshal(t, a, want)
}

func TestServiceHook_Marshal(t *testing.T) {
	testJSONMarshal(t, &ServiceHook{}, "{}")

	a := &ServiceHook{
		Name:            String("name"),
		Events:          []string{"e"},
		SupportedEvents: []string{"se"},
		Schema:          [][]string{{"g"}},
	}
	want := `{
		"name": "name",
		"events": [
			"e"
		],
		"supported_events": [
			"se"
		],
		"schema": [
			[
				"g"
			]
		]
	}`

	testJSONMarshal(t, a, want)
}
