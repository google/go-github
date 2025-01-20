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
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestRepositoriesService_ListDeployments(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/repos/o/r/deployments", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testFormValues(t, r, values{"environment": "test"})
		fmt.Fprint(w, `[{"id":1}, {"id":2}]`)
	})

	opt := &DeploymentsListOptions{Environment: "test"}
	ctx := context.Background()
	deployments, _, err := client.Repositories.ListDeployments(ctx, "o", "r", opt)
	if err != nil {
		t.Errorf("Repositories.ListDeployments returned error: %v", err)
	}

	want := []*Deployment{{ID: Ptr(int64(1))}, {ID: Ptr(int64(2))}}
	if !cmp.Equal(deployments, want) {
		t.Errorf("Repositories.ListDeployments returned %+v, want %+v", deployments, want)
	}

	const methodName = "ListDeployments"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Repositories.ListDeployments(ctx, "\n", "\n", opt)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Repositories.ListDeployments(ctx, "o", "r", opt)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestRepositoriesService_GetDeployment(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/repos/o/r/deployments/3", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{"id":3}`)
	})

	ctx := context.Background()
	deployment, _, err := client.Repositories.GetDeployment(ctx, "o", "r", 3)
	if err != nil {
		t.Errorf("Repositories.GetDeployment returned error: %v", err)
	}

	want := &Deployment{ID: Ptr(int64(3))}

	if !cmp.Equal(deployment, want) {
		t.Errorf("Repositories.GetDeployment returned %+v, want %+v", deployment, want)
	}

	const methodName = "GetDeployment"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Repositories.GetDeployment(ctx, "\n", "\n", 3)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Repositories.GetDeployment(ctx, "o", "r", 3)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestRepositoriesService_CreateDeployment(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	input := &DeploymentRequest{Ref: Ptr("1111"), Task: Ptr("deploy"), TransientEnvironment: Ptr(true)}

	mux.HandleFunc("/repos/o/r/deployments", func(w http.ResponseWriter, r *http.Request) {
		v := new(DeploymentRequest)
		assertNilError(t, json.NewDecoder(r.Body).Decode(v))

		testMethod(t, r, "POST")
		wantAcceptHeaders := []string{mediaTypeDeploymentStatusPreview, mediaTypeExpandDeploymentStatusPreview}
		testHeader(t, r, "Accept", strings.Join(wantAcceptHeaders, ", "))
		if !cmp.Equal(v, input) {
			t.Errorf("Request body = %+v, want %+v", v, input)
		}

		fmt.Fprint(w, `{"ref": "1111", "task": "deploy"}`)
	})

	ctx := context.Background()
	deployment, _, err := client.Repositories.CreateDeployment(ctx, "o", "r", input)
	if err != nil {
		t.Errorf("Repositories.CreateDeployment returned error: %v", err)
	}

	want := &Deployment{Ref: Ptr("1111"), Task: Ptr("deploy")}
	if !cmp.Equal(deployment, want) {
		t.Errorf("Repositories.CreateDeployment returned %+v, want %+v", deployment, want)
	}

	const methodName = "CreateDeployment"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Repositories.CreateDeployment(ctx, "\n", "\n", input)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Repositories.CreateDeployment(ctx, "o", "r", input)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestRepositoriesService_DeleteDeployment(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/repos/o/r/deployments/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
		w.WriteHeader(http.StatusNoContent)
	})

	ctx := context.Background()
	resp, err := client.Repositories.DeleteDeployment(ctx, "o", "r", 1)
	if err != nil {
		t.Errorf("Repositories.DeleteDeployment returned error: %v", err)
	}
	if resp.StatusCode != http.StatusNoContent {
		t.Error("Repositories.DeleteDeployment should return a 204 status")
	}

	resp, err = client.Repositories.DeleteDeployment(ctx, "o", "r", 2)
	if err == nil {
		t.Error("Repositories.DeleteDeployment should return an error")
	}
	if resp.StatusCode != http.StatusNotFound {
		t.Error("Repositories.DeleteDeployment should return a 404 status")
	}

	const methodName = "DeleteDeployment"
	testBadOptions(t, methodName, func() (err error) {
		_, err = client.Repositories.DeleteDeployment(ctx, "\n", "\n", 1)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		return client.Repositories.DeleteDeployment(ctx, "o", "r", 1)
	})
}

func TestRepositoriesService_ListDeploymentStatuses(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	wantAcceptHeaders := []string{mediaTypeDeploymentStatusPreview, mediaTypeExpandDeploymentStatusPreview}
	mux.HandleFunc("/repos/o/r/deployments/1/statuses", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testHeader(t, r, "Accept", strings.Join(wantAcceptHeaders, ", "))
		testFormValues(t, r, values{"page": "2"})
		fmt.Fprint(w, `[{"id":1}, {"id":2}]`)
	})

	opt := &ListOptions{Page: 2}
	ctx := context.Background()
	statuses, _, err := client.Repositories.ListDeploymentStatuses(ctx, "o", "r", 1, opt)
	if err != nil {
		t.Errorf("Repositories.ListDeploymentStatuses returned error: %v", err)
	}

	want := []*DeploymentStatus{{ID: Ptr(int64(1))}, {ID: Ptr(int64(2))}}
	if !cmp.Equal(statuses, want) {
		t.Errorf("Repositories.ListDeploymentStatuses returned %+v, want %+v", statuses, want)
	}

	const methodName = "ListDeploymentStatuses"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Repositories.ListDeploymentStatuses(ctx, "\n", "\n", 1, opt)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Repositories.ListDeploymentStatuses(ctx, "o", "r", 1, opt)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestRepositoriesService_GetDeploymentStatus(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	wantAcceptHeaders := []string{mediaTypeDeploymentStatusPreview, mediaTypeExpandDeploymentStatusPreview}
	mux.HandleFunc("/repos/o/r/deployments/3/statuses/4", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testHeader(t, r, "Accept", strings.Join(wantAcceptHeaders, ", "))
		fmt.Fprint(w, `{"id":4}`)
	})

	ctx := context.Background()
	deploymentStatus, _, err := client.Repositories.GetDeploymentStatus(ctx, "o", "r", 3, 4)
	if err != nil {
		t.Errorf("Repositories.GetDeploymentStatus returned error: %v", err)
	}

	want := &DeploymentStatus{ID: Ptr(int64(4))}
	if !cmp.Equal(deploymentStatus, want) {
		t.Errorf("Repositories.GetDeploymentStatus returned %+v, want %+v", deploymentStatus, want)
	}

	const methodName = "GetDeploymentStatus"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Repositories.GetDeploymentStatus(ctx, "\n", "\n", 3, 4)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Repositories.GetDeploymentStatus(ctx, "o", "r", 3, 4)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestRepositoriesService_CreateDeploymentStatus(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	input := &DeploymentStatusRequest{State: Ptr("inactive"), Description: Ptr("deploy"), AutoInactive: Ptr(false)}

	mux.HandleFunc("/repos/o/r/deployments/1/statuses", func(w http.ResponseWriter, r *http.Request) {
		v := new(DeploymentStatusRequest)
		assertNilError(t, json.NewDecoder(r.Body).Decode(v))

		testMethod(t, r, "POST")
		wantAcceptHeaders := []string{mediaTypeDeploymentStatusPreview, mediaTypeExpandDeploymentStatusPreview}
		testHeader(t, r, "Accept", strings.Join(wantAcceptHeaders, ", "))
		if !cmp.Equal(v, input) {
			t.Errorf("Request body = %+v, want %+v", v, input)
		}

		fmt.Fprint(w, `{"state": "inactive", "description": "deploy"}`)
	})

	ctx := context.Background()
	deploymentStatus, _, err := client.Repositories.CreateDeploymentStatus(ctx, "o", "r", 1, input)
	if err != nil {
		t.Errorf("Repositories.CreateDeploymentStatus returned error: %v", err)
	}

	want := &DeploymentStatus{State: Ptr("inactive"), Description: Ptr("deploy")}
	if !cmp.Equal(deploymentStatus, want) {
		t.Errorf("Repositories.CreateDeploymentStatus returned %+v, want %+v", deploymentStatus, want)
	}

	const methodName = "CreateDeploymentStatus"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Repositories.CreateDeploymentStatus(ctx, "\n", "\n", 1, input)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Repositories.CreateDeploymentStatus(ctx, "o", "r", 1, input)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestDeploymentStatusRequest_Marshal(t *testing.T) {
	t.Parallel()
	testJSONMarshal(t, &DeploymentStatusRequest{}, "{}")

	r := &DeploymentStatusRequest{
		State:          Ptr("state"),
		LogURL:         Ptr("logurl"),
		Description:    Ptr("desc"),
		Environment:    Ptr("env"),
		EnvironmentURL: Ptr("eurl"),
		AutoInactive:   Ptr(false),
	}

	want := `{
		"state": "state",
		"log_url": "logurl",
		"description": "desc",
		"environment": "env",
		"environment_url": "eurl",
		"auto_inactive": false
	}`

	testJSONMarshal(t, r, want)
}

func TestDeploymentStatus_Marshal(t *testing.T) {
	t.Parallel()
	testJSONMarshal(t, &DeploymentStatus{}, "{}")

	r := &DeploymentStatus{
		ID:    Ptr(int64(1)),
		State: Ptr("state"),
		Creator: &User{
			Login:           Ptr("l"),
			ID:              Ptr(int64(1)),
			URL:             Ptr("u"),
			AvatarURL:       Ptr("a"),
			GravatarID:      Ptr("g"),
			Name:            Ptr("n"),
			Company:         Ptr("c"),
			Blog:            Ptr("b"),
			Location:        Ptr("l"),
			Email:           Ptr("e"),
			Hireable:        Ptr(true),
			Bio:             Ptr("b"),
			TwitterUsername: Ptr("t"),
			PublicRepos:     Ptr(1),
			Followers:       Ptr(1),
			Following:       Ptr(1),
			CreatedAt:       &Timestamp{referenceTime},
			SuspendedAt:     &Timestamp{referenceTime},
		},
		Description:    Ptr("desc"),
		Environment:    Ptr("env"),
		NodeID:         Ptr("nid"),
		CreatedAt:      &Timestamp{referenceTime},
		UpdatedAt:      &Timestamp{referenceTime},
		TargetURL:      Ptr("turl"),
		DeploymentURL:  Ptr("durl"),
		RepositoryURL:  Ptr("rurl"),
		EnvironmentURL: Ptr("eurl"),
		LogURL:         Ptr("lurl"),
		URL:            Ptr("url"),
	}

	want := `{
		"id": 1,
		"state": "state",
		"creator": {
			"login": "l",
			"id": 1,
			"avatar_url": "a",
			"gravatar_id": "g",
			"name": "n",
			"company": "c",
			"blog": "b",
			"location": "l",
			"email": "e",
			"hireable": true,
			"bio": "b",
			"twitter_username": "t",
			"public_repos": 1,
			"followers": 1,
			"following": 1,
			"created_at": ` + referenceTimeStr + `,
			"suspended_at": ` + referenceTimeStr + `,
			"url": "u"
		},
		"description": "desc",
		"environment": "env",
		"node_id": "nid",
		"created_at": ` + referenceTimeStr + `,
		"updated_at": ` + referenceTimeStr + `,
		"target_url": "turl",
		"deployment_url": "durl",
		"repository_url": "rurl",
		"environment_url": "eurl",
		"log_url": "lurl",
		"url": "url"
	}`

	testJSONMarshal(t, r, want)
}

func TestDeploymentRequest_Marshal(t *testing.T) {
	t.Parallel()
	testJSONMarshal(t, &DeploymentRequest{}, "{}")

	r := &DeploymentRequest{
		Ref:                   Ptr("ref"),
		Task:                  Ptr("task"),
		AutoMerge:             Ptr(false),
		RequiredContexts:      &[]string{"s"},
		Payload:               "payload",
		Environment:           Ptr("environment"),
		Description:           Ptr("description"),
		TransientEnvironment:  Ptr(false),
		ProductionEnvironment: Ptr(false),
	}

	want := `{
		"ref": "ref",
		"task": "task",
		"auto_merge": false,
		"required_contexts": ["s"],
		"payload": "payload",
		"environment": "environment",
		"description": "description",
		"transient_environment": false,
		"production_environment": false
	}`

	testJSONMarshal(t, r, want)
}

func TestDeployment_Marshal(t *testing.T) {
	t.Parallel()
	testJSONMarshal(t, &Deployment{}, "{}")

	str := "s"
	jsonMsg, _ := json.Marshal(str)

	r := &Deployment{
		URL:         Ptr("url"),
		ID:          Ptr(int64(1)),
		SHA:         Ptr("sha"),
		Ref:         Ptr("ref"),
		Task:        Ptr("task"),
		Payload:     jsonMsg,
		Environment: Ptr("env"),
		Description: Ptr("desc"),
		Creator: &User{
			Login:           Ptr("l"),
			ID:              Ptr(int64(1)),
			URL:             Ptr("u"),
			AvatarURL:       Ptr("a"),
			GravatarID:      Ptr("g"),
			Name:            Ptr("n"),
			Company:         Ptr("c"),
			Blog:            Ptr("b"),
			Location:        Ptr("l"),
			Email:           Ptr("e"),
			Hireable:        Ptr(true),
			Bio:             Ptr("b"),
			TwitterUsername: Ptr("t"),
			PublicRepos:     Ptr(1),
			Followers:       Ptr(1),
			Following:       Ptr(1),
			CreatedAt:       &Timestamp{referenceTime},
			SuspendedAt:     &Timestamp{referenceTime},
		},
		CreatedAt:     &Timestamp{referenceTime},
		UpdatedAt:     &Timestamp{referenceTime},
		StatusesURL:   Ptr("surl"),
		RepositoryURL: Ptr("rurl"),
		NodeID:        Ptr("nid"),
	}

	want := `{
		"url": "url",
		"id": 1,
		"sha": "sha",
		"ref": "ref",
		"task": "task",
		"payload": "s",
		"environment": "env",
		"description": "desc",
		"creator": {
			"login": "l",
			"id": 1,
			"avatar_url": "a",
			"gravatar_id": "g",
			"name": "n",
			"company": "c",
			"blog": "b",
			"location": "l",
			"email": "e",
			"hireable": true,
			"bio": "b",
			"twitter_username": "t",
			"public_repos": 1,
			"followers": 1,
			"following": 1,
			"created_at": ` + referenceTimeStr + `,
			"suspended_at": ` + referenceTimeStr + `,
			"url": "u"
		},
		"created_at": ` + referenceTimeStr + `,
		"updated_at": ` + referenceTimeStr + `,
		"statuses_url": "surl",
		"repository_url": "rurl",
		"node_id": "nid"
	}`

	testJSONMarshal(t, r, want)
}
