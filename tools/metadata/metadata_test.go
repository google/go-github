// Copyright 2023 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"testing"
)

func Test_normalizedOpName(t *testing.T) {
	t.Parallel()
	for _, td := range []struct {
		name string
		want string
	}{
		{name: "", want: ""},
		{name: "get /foo/{id}", want: "GET /foo/*"},
		{name: "get foo", want: "GET /foo"},
	} {
		t.Run(td.name, func(t *testing.T) {
			t.Parallel()
			got := normalizedOpName(td.name)
			if got != td.want {
				t.Errorf("normalizedOpName() = %v, want %v", got, td.want)
			}
		})
	}
}

func Test_normalizeDocURL(t *testing.T) {
	t.Parallel()

	for _, td := range []struct {
		name   string
		docURL string
		want   string
	}{
		{
			name:   "invalid URL",
			docURL: "://bad",
			want:   "://bad",
		},
		{
			name:   "clean path and add api version",
			docURL: "https://docs.github.com//rest/repos/repos#get-a-repository",
			want:   "https://docs.github.com/rest/repos/repos?apiVersion=2022-11-28#get-a-repository",
		},
		{
			name:   "preserve query and add api version",
			docURL: "https://docs.github.com/rest/repos/repos?foo=bar",
			want:   "https://docs.github.com/rest/repos/repos?apiVersion=2022-11-28&foo=bar",
		},
		{
			name:   "replace existing api version",
			docURL: "https://docs.github.com/rest/repos/repos?apiVersion=2021-01-01&foo=bar",
			want:   "https://docs.github.com/rest/repos/repos?apiVersion=2022-11-28&foo=bar",
		},
		{
			name:   "enterprise cloud latest rest is normalized",
			docURL: "https://docs.github.com/enterprise-cloud@latest/rest/repos/repos#get-a-repository",
			want:   "https://docs.github.com/enterprise-cloud@latest/rest/repos/repos?apiVersion=2022-11-28#get-a-repository",
		},
		{
			name:   "non rest docs path unchanged",
			docURL: "https://gist.github.com/jonmagic/5282384165e0f86ef105",
			want:   "https://gist.github.com/jonmagic/5282384165e0f86ef105",
		},
		{
			name:   "non docs host unchanged",
			docURL: "https://example.com/rest/repos/repos",
			want:   "https://example.com/rest/repos/repos",
		},
	} {
		t.Run(td.name, func(t *testing.T) {
			t.Parallel()
			got := normalizeDocURL(td.docURL)
			if got != td.want {
				t.Errorf("normalizeDocURL() = %v, want %v", got, td.want)
			}
		})
	}
}
