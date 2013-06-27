// Copyright 2013 The go-github AUTHORS. All rights reserved.
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
	"time"
)

func TestIssuesService_List_all(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/issues", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testFormValues(t, r, values{
			"filter":    "all",
			"state":     "closed",
			"labels":    "a,b",
			"sort":      "updated",
			"direction": "asc",
			"since":     "2002-02-10T15:30:00Z",
		})
		fmt.Fprint(w, `[{"number":1}]`)
	})

	opt := &IssueListOptions{
		"all", "closed", []string{"a", "b"}, "updated", "asc",
		time.Date(2002, time.February, 10, 15, 30, 0, 0, time.UTC),
	}
	issues, err := client.Issues.List(true, opt)

	if err != nil {
		t.Errorf("Issues.List returned error: %v", err)
	}

	want := []Issue{Issue{Number: 1}}
	if !reflect.DeepEqual(issues, want) {
		t.Errorf("Issues.List returned %+v, want %+v", issues, want)
	}
}

func TestIssuesService_List_owned(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/user/issues", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `[{"number":1}]`)
	})

	issues, err := client.Issues.List(false, nil)
	if err != nil {
		t.Errorf("Issues.List returned error: %v", err)
	}

	want := []Issue{Issue{Number: 1}}
	if !reflect.DeepEqual(issues, want) {
		t.Errorf("Issues.List returned %+v, want %+v", issues, want)
	}
}

func TestIssuesService_ListByOrg(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/orgs/o/issues", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `[{"number":1}]`)
	})

	issues, err := client.Issues.ListByOrg("o", nil)
	if err != nil {
		t.Errorf("Issues.ListByOrg returned error: %v", err)
	}

	want := []Issue{Issue{Number: 1}}
	if !reflect.DeepEqual(issues, want) {
		t.Errorf("Issues.List returned %+v, want %+v", issues, want)
	}
}

func TestIssuesService_ListByOrg_invalidOrg(t *testing.T) {
	_, err := client.Issues.ListByOrg("%", nil)
	testURLParseError(t, err)
}

func TestIssuesService_ListByRepo(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/repos/o/r/issues", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testFormValues(t, r, values{
			"milestone": "*",
			"state":     "closed",
			"assignee":  "a",
			"creator":   "c",
			"mentioned": "m",
			"labels":    "a,b",
			"sort":      "updated",
			"direction": "asc",
			"since":     "2002-02-10T15:30:00Z",
		})
		fmt.Fprint(w, `[{"number":1}]`)
	})

	opt := &IssueListByRepoOptions{
		"*", "closed", "a", "c", "m", []string{"a", "b"}, "updated", "asc",
		time.Date(2002, time.February, 10, 15, 30, 0, 0, time.UTC),
	}
	issues, err := client.Issues.ListByRepo("o", "r", opt)
	if err != nil {
		t.Errorf("Issues.ListByOrg returned error: %v", err)
	}

	want := []Issue{Issue{Number: 1}}
	if !reflect.DeepEqual(issues, want) {
		t.Errorf("Issues.List returned %+v, want %+v", issues, want)
	}
}

func TestIssuesService_ListByRepo_invalidOwner(t *testing.T) {
	_, err := client.Issues.ListByRepo("%", "r", nil)
	testURLParseError(t, err)
}

func TestIssuesService_Get(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/repos/o/r/issues/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{"number":1}`)
	})

	issue, err := client.Issues.Get("o", "r", 1)
	if err != nil {
		t.Errorf("Issues.Get returned error: %v", err)
	}

	want := &Issue{Number: 1}
	if !reflect.DeepEqual(issue, want) {
		t.Errorf("Issues.Get returned %+v, want %+v", issue, want)
	}
}

func TestIssuesService_Get_invalidOwner(t *testing.T) {
	_, err := client.Issues.Get("%", "r", 1)
	testURLParseError(t, err)
}

func TestIssuesService_Create(t *testing.T) {
	setup()
	defer teardown()

	input := &Issue{Title: "t"}

	mux.HandleFunc("/repos/o/r/issues", func(w http.ResponseWriter, r *http.Request) {
		v := new(Issue)
		json.NewDecoder(r.Body).Decode(v)

		testMethod(t, r, "POST")
		if !reflect.DeepEqual(v, input) {
			t.Errorf("Request body = %+v, want %+v", v, input)
		}

		fmt.Fprint(w, `{"number":1}`)
	})

	issue, err := client.Issues.Create("o", "r", input)
	if err != nil {
		t.Errorf("Issues.Create returned error: %v", err)
	}

	want := &Issue{Number: 1}
	if !reflect.DeepEqual(issue, want) {
		t.Errorf("Issues.Create returned %+v, want %+v", issue, want)
	}
}

func TestIssuesService_Create_invalidOwner(t *testing.T) {
	_, err := client.Issues.Create("%", "r", nil)
	testURLParseError(t, err)
}

func TestIssuesService_Edit(t *testing.T) {
	setup()
	defer teardown()

	input := &Issue{Title: "t"}

	mux.HandleFunc("/repos/o/r/issues/1", func(w http.ResponseWriter, r *http.Request) {
		v := new(Issue)
		json.NewDecoder(r.Body).Decode(v)

		testMethod(t, r, "PATCH")
		if !reflect.DeepEqual(v, input) {
			t.Errorf("Request body = %+v, want %+v", v, input)
		}

		fmt.Fprint(w, `{"number":1}`)
	})

	issue, err := client.Issues.Edit("o", "r", 1, input)
	if err != nil {
		t.Errorf("Issues.Edit returned error: %v", err)
	}

	want := &Issue{Number: 1}
	if !reflect.DeepEqual(issue, want) {
		t.Errorf("Issues.Edit returned %+v, want %+v", issue, want)
	}
}

func TestIssuesService_Edit_invalidOwner(t *testing.T) {
	_, err := client.Issues.Edit("%", "r", 1, nil)
	testURLParseError(t, err)
}

func TestIssuesService_ListAssignees(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/repos/o/r/assignees", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `[{"id":1}]`)
	})

	assignees, err := client.Issues.ListAssignees("o", "r")
	if err != nil {
		t.Errorf("Issues.List returned error: %v", err)
	}

	want := []User{User{ID: 1}}
	if !reflect.DeepEqual(assignees, want) {
		t.Errorf("Issues.ListAssignees returned %+v, want %+v", assignees, want)
	}
}

func TestIssuesService_ListAssignees_invalidOwner(t *testing.T) {
	_, err := client.Issues.ListAssignees("%", "r")
	testURLParseError(t, err)
}

func TestIssuesService_CheckAssignee_true(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/repos/o/r/assignees/u", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
	})

	assignee, err := client.Issues.CheckAssignee("o", "r", "u")
	if err != nil {
		t.Errorf("Issues.CheckAssignee returned error: %v", err)
	}
	if want := true; assignee != want {
		t.Errorf("Issues.CheckAssignee returned %+v, want %+v", assignee, want)
	}
}

func TestIssuesService_CheckAssignee_false(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/repos/o/r/assignees/u", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		w.WriteHeader(http.StatusNotFound)
	})

	assignee, err := client.Issues.CheckAssignee("o", "r", "u")
	if err != nil {
		t.Errorf("Issues.CheckAssignee returned error: %v", err)
	}
	if want := false; assignee != want {
		t.Errorf("Issues.CheckAssignee returned %+v, want %+v", assignee, want)
	}
}

func TestIssuesService_CheckAssignee_error(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/repos/o/r/assignees/u", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		http.Error(w, "BadRequest", http.StatusBadRequest)
	})

	assignee, err := client.Issues.CheckAssignee("o", "r", "u")
	if err == nil {
		t.Errorf("Expected HTTP 400 response")
	}
	if want := false; assignee != want {
		t.Errorf("Issues.CheckAssignee returned %+v, want %+v", assignee, want)
	}
}

func TestIssuesService_CheckAssignee_invalidOwner(t *testing.T) {
	_, err := client.Issues.CheckAssignee("%", "r", "u")
	testURLParseError(t, err)
}

func TestIssuesService_ListComments_allIssues(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/repos/o/r/issues/comments", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testFormValues(t, r, values{
			"sort":      "updated",
			"direction": "desc",
			"since":     "2002-02-10T15:30:00Z",
		})
		fmt.Fprint(w, `[{"id":1}]`)
	})

	opt := &IssueListCommentsOptions{"updated", "desc",
		time.Date(2002, time.February, 10, 15, 30, 0, 0, time.UTC),
	}
	comments, err := client.Issues.ListComments("o", "r", 0, opt)
	if err != nil {
		t.Errorf("Issues.ListComments returned error: %v", err)
	}

	want := []IssueComment{IssueComment{ID: 1}}
	if !reflect.DeepEqual(comments, want) {
		t.Errorf("Issues.ListComments returned %+v, want %+v", comments, want)
	}
}

func TestIssuesService_ListComments_specificIssue(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/repos/o/r/issues/1/comments", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `[{"id":1}]`)
	})

	comments, err := client.Issues.ListComments("o", "r", 1, nil)
	if err != nil {
		t.Errorf("Issues.ListComments returned error: %v", err)
	}

	want := []IssueComment{IssueComment{ID: 1}}
	if !reflect.DeepEqual(comments, want) {
		t.Errorf("Issues.ListComments returned %+v, want %+v", comments, want)
	}
}

func TestIssuesService_ListComments_invalidOwner(t *testing.T) {
	_, err := client.Issues.ListComments("%", "r", 1, nil)
	testURLParseError(t, err)
}

func TestIssuesService_GetComment(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/repos/o/r/issues/comments/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{"id":1}`)
	})

	comment, err := client.Issues.GetComment("o", "r", 1)
	if err != nil {
		t.Errorf("Issues.GetComment returned error: %v", err)
	}

	want := &IssueComment{ID: 1}
	if !reflect.DeepEqual(comment, want) {
		t.Errorf("Issues.GetComment returned %+v, want %+v", comment, want)
	}
}

func TestIssuesService_GetComment_invalidOrg(t *testing.T) {
	_, err := client.Issues.GetComment("%", "r", 1)
	testURLParseError(t, err)
}

func TestIssuesService_CreateComment(t *testing.T) {
	setup()
	defer teardown()

	input := &IssueComment{Body: "b"}

	mux.HandleFunc("/repos/o/r/issues/1/comments", func(w http.ResponseWriter, r *http.Request) {
		v := new(IssueComment)
		json.NewDecoder(r.Body).Decode(v)

		testMethod(t, r, "POST")
		if !reflect.DeepEqual(v, input) {
			t.Errorf("Request body = %+v, want %+v", v, input)
		}

		fmt.Fprint(w, `{"id":1}`)
	})

	comment, err := client.Issues.CreateComment("o", "r", 1, input)
	if err != nil {
		t.Errorf("Issues.CreateComment returned error: %v", err)
	}

	want := &IssueComment{ID: 1}
	if !reflect.DeepEqual(comment, want) {
		t.Errorf("Issues.CreateComment returned %+v, want %+v", comment, want)
	}
}

func TestIssuesService_CreateComment_invalidOrg(t *testing.T) {
	_, err := client.Issues.CreateComment("%", "r", 1, nil)
	testURLParseError(t, err)
}

func TestIssuesService_EditComment(t *testing.T) {
	setup()
	defer teardown()

	input := &IssueComment{Body: "b"}

	mux.HandleFunc("/repos/o/r/issues/comments/1", func(w http.ResponseWriter, r *http.Request) {
		v := new(IssueComment)
		json.NewDecoder(r.Body).Decode(v)

		testMethod(t, r, "PATCH")
		if !reflect.DeepEqual(v, input) {
			t.Errorf("Request body = %+v, want %+v", v, input)
		}

		fmt.Fprint(w, `{"id":1}`)
	})

	comment, err := client.Issues.EditComment("o", "r", 1, input)
	if err != nil {
		t.Errorf("Issues.EditComment returned error: %v", err)
	}

	want := &IssueComment{ID: 1}
	if !reflect.DeepEqual(comment, want) {
		t.Errorf("Issues.EditComment returned %+v, want %+v", comment, want)
	}
}

func TestIssuesService_EditComment_invalidOwner(t *testing.T) {
	_, err := client.Issues.EditComment("%", "r", 1, nil)
	testURLParseError(t, err)
}

func TestIssuesService_DeleteComment(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/repos/o/r/issues/comments/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
	})

	err := client.Issues.DeleteComment("o", "r", 1)
	if err != nil {
		t.Errorf("Issues.DeleteComments returned error: %v", err)
	}
}

func TestIssuesService_DeleteComment_invalidOwner(t *testing.T) {
	err := client.Issues.DeleteComment("%", "r", 1)
	testURLParseError(t, err)
}
