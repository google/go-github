// Copyright 2013 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package github

import (
	"context"
	"net/http"
	"testing"
)

func TestRepositoriesService_EnableLfs(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/repos/o/r/lfs", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PUT")

		w.WriteHeader(http.StatusNoContent)
	})

	ctx := context.Background()
	if _, err := client.Repositories.EnableLfs(ctx, "o", "r"); err != nil {
		t.Errorf("Repositories.EnableLfs returned error: %v", err)
	}

	const methodName = "EnableLfs"
	testBadOptions(t, methodName, func() (err error) {
		_, err = client.Repositories.EnableLfs(ctx, "\n", "\n")
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		return client.Repositories.EnableLfs(ctx, "o", "r")
	})
}

func TestRepositoriesService_DisableLfs(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/repos/o/r/lfs", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")

		w.WriteHeader(http.StatusNoContent)
	})

	ctx := context.Background()
	if _, err := client.Repositories.DisableLfs(ctx, "o", "r"); err != nil {
		t.Errorf("Repositories.DisableLfs returned error: %v", err)
	}

	const methodName = "DisableLfs"
	testBadOptions(t, methodName, func() (err error) {
		_, err = client.Repositories.DisableLfs(ctx, "\n", "\n")
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		return client.Repositories.DisableLfs(ctx, "o", "r")
	})
}
