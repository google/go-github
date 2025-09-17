package github

import (
	"context"
	"fmt"
	"net/http"
	"testing"
)

func TestProjectsService_ListOrganizationProjects(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/orgs/o/projectsV2", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testHeader(t, r, "Accept", mediaTypeProjectsPreview)
		// Expect query params q, page, per_page when provided
		testFormValues(t, r, values{"q": "alpha", "page": "2", "per_page": "1"})
		fmt.Fprint(w, `[{"id":1,"title":"T1","created_at":"2011-01-02T15:04:05Z","updated_at":"2012-01-02T15:04:05Z"}]`)
	})

	opts := &ListProjectsOptions{Q: "alpha", ListOptions: ListOptions{Page: 2, PerPage: 1}}
	ctx := context.Background()
	projects, _, err := client.Projects.ListOrganizationProjects(ctx, "o", opts)
	if err != nil {
		t.Fatalf("Projects.ListOrganizationProjects returned error: %v", err)
	}
	if len(projects) != 1 || projects[0].GetID() != 1 || projects[0].GetTitle() != "T1" {
		t.Fatalf("Projects.ListOrganizationProjects returned %+v", projects)
	}

	const methodName = "ListOrganizationProjects"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Projects.ListOrganizationProjects(ctx, "\n", opts)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Projects.ListOrganizationProjects(ctx, "o", opts)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestProjectsService_GetOrganizationProject(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/orgs/o/projectsV2/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testHeader(t, r, "Accept", mediaTypeProjectsPreview)
		fmt.Fprint(w, `{"id":1,"title":"OrgProj","created_at":"2011-01-02T15:04:05Z","updated_at":"2012-01-02T15:04:05Z"}`)
	})

	ctx := context.Background()
	project, _, err := client.Projects.GetByOrg(ctx, "o", 1)
	if err != nil {
		t.Fatalf("Projects.GetByOrg returned error: %v", err)
	}
	if project.GetID() != 1 || project.GetTitle() != "OrgProj" {
		t.Fatalf("Projects.GetByOrg returned %+v", project)
	}

	const methodName = "GetByOrg"
	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Projects.GetByOrg(ctx, "o", 1)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestProjectsService_ListUserProjects(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/users/u/projectsV2", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testHeader(t, r, "Accept", mediaTypeProjectsPreview)
		testFormValues(t, r, values{"q": "beta", "page": "1", "per_page": "2"})
		fmt.Fprint(w, `[{"id":2,"title":"UProj","created_at":"2011-01-02T15:04:05Z","updated_at":"2012-01-02T15:04:05Z"}]`)
	})

	opts := &ListProjectsOptions{Q: "beta", ListOptions: ListOptions{Page: 1, PerPage: 2}}
	ctx := context.Background()
	projects, _, err := client.Projects.ListByUser(ctx, "u", opts)
	if err != nil {
		t.Fatalf("Projects.ListByUser returned error: %v", err)
	}
	if len(projects) != 1 || projects[0].GetID() != 2 || projects[0].GetTitle() != "UProj" {
		t.Fatalf("Projects.ListByUser returned %+v", projects)
	}

	const methodName = "ListByUser"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Projects.ListByUser(ctx, "\n", opts)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Projects.ListByUser(ctx, "u", opts)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestProjectsService_GetUserProject(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/users/u/projectsV2/2", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testHeader(t, r, "Accept", mediaTypeProjectsPreview)
		fmt.Fprint(w, `{"id":2,"title":"UProj","created_at":"2011-01-02T15:04:05Z","updated_at":"2012-01-02T15:04:05Z"}`)
	})

	ctx := context.Background()
	project, _, err := client.Projects.GetUserProject(ctx, "u", 2)
	if err != nil {
		t.Fatalf("Projects.GetUserProject returned error: %v", err)
	}
	if project.GetID() != 2 || project.GetTitle() != "UProj" {
		t.Fatalf("Projects.GetUserProject returned %+v", project)
	}

	const methodName = "GetUserProject"
	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Projects.GetUserProject(ctx, "u", 2)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

// Marshal test ensures V2 fields marshal correctly.
func TestProjectV2_Marshal(t *testing.T) {
	t.Parallel()
	testJSONMarshal(t, &ProjectV2{}, "{}")

	p := &ProjectV2{
		ID:          Ptr(int64(10)),
		Title:       Ptr("Title"),
		Description: Ptr("Desc"),
		Public:      Ptr(true),
		CreatedAt:   &Timestamp{referenceTime},
		UpdatedAt:   &Timestamp{referenceTime},
	}

	want := `{
        "id": 10,
        "title": "Title",
        "description": "Desc",
        "public": true,
        "created_at": ` + referenceTimeStr + `,
        "updated_at": ` + referenceTimeStr + `
    }`

	testJSONMarshal(t, p, want)
}
