// Copyright 2014 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package github

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"testing"
)

func TestRepositoriesService_ListDeployments(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/repos/o/r/deployments", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testFormValues(t, r, values{"environment": "test"})
		fmt.Fprint(w, `[{"id":1}, {"id":2}]`)
	})

	opt := &DeploymentsListOptions{Environment: "test"}
	deployments, _, err := client.Repositories.ListDeployments("o", "r", opt)
	if err != nil {
		t.Errorf("Repositories.ListDeployments returned error: %v", err)
	}

	want := []RepositoryDeployment{{ID: Int(1)}, {ID: Int(2)}}
	if !reflect.DeepEqual(deployments, want) {
		t.Errorf("Repositories.ListDeployments returned %+v, want %+v", deployments, want)
	}
}

func TestRepositoriesService_CreateDeployment(t *testing.T) {
	setup()
	defer teardown()

	input := &RepositoryDeploymentRequest{Ref: String("1111"), Task: String("deploy")}

	mux.HandleFunc("/repos/o/r/deployments", func(w http.ResponseWriter, r *http.Request) {
		v := new(RepositoryDeploymentRequest)
		json.NewDecoder(r.Body).Decode(v)

		testMethod(t, r, "POST")
		if !reflect.DeepEqual(v, input) {
			t.Errorf("Request body = %+v, want %+v", v, input)
		}

		fmt.Fprint(w, `{"ref": "1111", "task": "deploy"}`)
	})

	deployment, _, err := client.Repositories.CreateDeployment("o", "r", input)
	if err != nil {
		t.Errorf("Repositories.CreateDeployment returned error: %v", err)
	}

	want := &RepositoryDeployment{Ref: String("1111"), Task: String("deploy")}
	if !reflect.DeepEqual(deployment, want) {
		t.Errorf("Repositories.CreateDeployment returned %+v, want %+v", deployment, want)
	}
}
