// Copyright 2016 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package github

import (
	"fmt"
	"net/http"
	"reflect"
	"testing"
)

func TestIntegrationService_ListInstallations(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/integration/installations", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testHeader(t, r, "Accept", mediaTypeIntegrationPreview)
		testFormValues(t, r, values{
			"page":     "1",
			"per_page": "2",
		})
		fmt.Fprint(w, `[{"id":1}]`)
	})

	opt := &ListOptions{Page: 1, PerPage: 2}
	installations, _, err := client.Integrations.ListInstallations(opt)
	if err != nil {
		t.Errorf("Integration.ListInstallations returned error: %v", err)
	}

	want := []*Installation{{ID: Int(1)}}
	if !reflect.DeepEqual(installations, want) {
		t.Errorf("Integration.ListInstallations returned %+v, want %+v", installations, want)
	}
}
