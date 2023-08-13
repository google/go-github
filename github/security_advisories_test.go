// Copyright 2023 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package github

import (
	"context"
	"net/http"
	"testing"
)

func TestSecurityAdvisoriesService_RequestCVE(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/repos/o/r/security-advisories/ghsa_id_ok/cve", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		w.WriteHeader(http.StatusOK)
	})

	mux.HandleFunc("/repos/o/r/security-advisories/ghsa_id_accepted/cve", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		w.WriteHeader(http.StatusAccepted)
	})

	ctx := context.Background()
	_, err := client.SecurityAdvisories.RequestCVE(ctx, "o", "r", "ghsa_id_ok")
	if err != nil {
		t.Errorf("SecurityAdvisoriesService.RequestCVE returned error: %v", err)
	}

	_, err = client.SecurityAdvisories.RequestCVE(ctx, "o", "r", "ghsa_id_accepted")
	if err != nil {
		t.Errorf("SecurityAdvisoriesService.RequestCVE returned error: %v", err)
	}

	const methodName = "RequestCVE"
	testBadOptions(t, methodName, func() (err error) {
		_, err = client.SecurityAdvisories.RequestCVE(ctx, "\n", "\n", "\n")
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		resp, err := client.SecurityAdvisories.RequestCVE(ctx, "o", "r", "ghsa_id")
		if err == nil {
			t.Errorf("testNewRequestAndDoFailure %v should have return err", methodName)
		}
		return resp, err
	})
}
