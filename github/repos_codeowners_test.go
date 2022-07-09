// Copyright 2022 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package github

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestRepositoriesService_GetCodeownersErrors(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/repos/o/r/codeowners/errors", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testHeader(t, r, "Accept", mediaTypeV3)
		fmt.Fprint(w, `{
		  "errors": [
			{
			  "line": 1,
			  "column": 1,
			  "kind": "Invalid pattern",
			  "source": "***/*.rb @monalisa",
			  "suggestion": "Did you mean **/*.rb?",
			  "message": "Invalid pattern on line 3: Did you mean **/*.rb?\n\n  ***/*.rb @monalisa\n  ^",
			  "path": ".github/CODEOWNERS"
			}
		  ]
		}
	`)
	})

	ctx := context.Background()
	codeownersErrors, _, err := client.Repositories.GetCodeownersErrors(ctx, "o", "r")
	if err != nil {
		t.Errorf("Repositories.GetCodeownersErrors returned error: %v", err)
	}

	want := &CodeownersErrors{
		Errors: []*CodeownersError{
			{
				Line:       1,
				Column:     1,
				Kind:       "Invalid pattern",
				Source:     "***/*.rb @monalisa",
				Suggestion: String("Did you mean **/*.rb?"),
				Message:    "Invalid pattern on line 3: Did you mean **/*.rb?\n\n  ***/*.rb @monalisa\n  ^",
				Path:       ".github/CODEOWNERS",
			},
		},
	}
	if !cmp.Equal(codeownersErrors, want) {
		t.Errorf("Repositories.GetCodeownersErrors returned %+v, want %+v", codeownersErrors, want)
	}

	const methodName = "GetCodeownersErrors"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Repositories.GetCodeownersErrors(ctx, "\n", "\n")
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Repositories.GetCodeownersErrors(ctx, "o", "r")
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestCodeownersErrors_Marshal(t *testing.T) {
	testJSONMarshal(t, &CodeownersErrors{}, "{}")

	u := &CodeownersErrors{
		Errors: []*CodeownersError{
			{
				Line:       1,
				Column:     1,
				Kind:       "Invalid pattern",
				Source:     "***/*.rb @monalisa",
				Suggestion: String("Did you mean **/*.rb?"),
				Message:    "Invalid pattern on line 3: Did you mean **/*.rb?\n\n  ***/*.rb @monalisa\n  ^",
				Path:       ".github/CODEOWNERS",
			},
		},
	}

	want := `{
	  "errors": [
		{
		  "line": 1,
		  "column": 1,
		  "kind": "Invalid pattern",
		  "source": "***/*.rb @monalisa",
		  "suggestion": "Did you mean **/*.rb?",
		  "message": "Invalid pattern on line 3: Did you mean **/*.rb?\n\n  ***/*.rb @monalisa\n  ^",
		  "path": ".github/CODEOWNERS"
		}
	  ]
	}
`
	testJSONMarshal(t, u, want)
}
