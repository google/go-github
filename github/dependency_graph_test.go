// Copyright 2023 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package github

import (
	"context"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
)

func TestDependencyGraphService_GetSBOM(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/repos/owner/repo/dependency-graph/sbom", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{
   "sbom":{
      "creationInfo":{
         "created":"2021-09-01T00:00:00Z"
      },
      "name":"owner/repo",
      "packages":[
                {
                "name":"rubygems:rails",
                "versionInfo":"1.0.0"
                }
            ]
        }
    }`)
	})

	ctx := context.Background()
	sbom, _, err := client.DependencyGraph.GetSBOM(ctx, "owner", "repo")
	if err != nil {
		t.Errorf("DependencyGraph.GetSBOM returned error: %v", err)
	}

	testTime := time.Date(2021, 9, 1, 0, 0, 0, 0, time.UTC)
	want := &SBOM{
		&SBOMInfo{
			CreationInfo: &CreationInfo{
				Created: &Timestamp{testTime},
			},
			Name: Ptr("owner/repo"),
			Packages: []*RepoDependencies{
				{
					Name:        Ptr("rubygems:rails"),
					VersionInfo: Ptr("1.0.0"),
				},
			},
		},
	}

	if !cmp.Equal(sbom, want) {
		t.Errorf("DependencyGraph.GetSBOM returned %+v, want %+v", sbom, want)
	}

	const methodName = "GetSBOM"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.DependencyGraph.GetSBOM(ctx, "\n", "\n")
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.DependencyGraph.GetSBOM(ctx, "owner", "repo")
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}
