// Copyright 2013 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package github

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"testing"
)

func TestIssuesService_ListLabels(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/repos/o/r/labels", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testFormValues(t, r, values{"page": "2"})
		fmt.Fprint(w, `[{"name": "a"},{"name": "b"}]`)
	})

	opt := &ListOptions{Page: 2}
	ctx := context.Background()
	labels, _, err := client.Issues.ListLabels(ctx, "o", "r", opt)
	if err != nil {
		t.Errorf("Issues.ListLabels returned error: %v", err)
	}

	want := []*Label{{Name: String("a")}, {Name: String("b")}}
	if !reflect.DeepEqual(labels, want) {
		t.Errorf("Issues.ListLabels returned %+v, want %+v", labels, want)
	}
}

func TestIssuesService_ListLabels_invalidOwner(t *testing.T) {
	client, _, _, teardown := setup()
	defer teardown()

	ctx := context.Background()
	_, _, err := client.Issues.ListLabels(ctx, "%", "%", nil)
	testURLParseError(t, err)
}

func TestIssuesService_GetLabel(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/repos/o/r/labels/n", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{"url":"u", "name": "n", "color": "c", "description": "d"}`)
	})

	ctx := context.Background()
	label, _, err := client.Issues.GetLabel(ctx, "o", "r", "n")
	if err != nil {
		t.Errorf("Issues.GetLabel returned error: %v", err)
	}

	want := &Label{URL: String("u"), Name: String("n"), Color: String("c"), Description: String("d")}
	if !reflect.DeepEqual(label, want) {
		t.Errorf("Issues.GetLabel returned %+v, want %+v", label, want)
	}
}

func TestIssuesService_GetLabel_invalidOwner(t *testing.T) {
	client, _, _, teardown := setup()
	defer teardown()

	ctx := context.Background()
	_, _, err := client.Issues.GetLabel(ctx, "%", "%", "%")
	testURLParseError(t, err)
}

func TestIssuesService_CreateLabel(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	input := &Label{Name: String("n")}

	mux.HandleFunc("/repos/o/r/labels", func(w http.ResponseWriter, r *http.Request) {
		v := new(Label)
		json.NewDecoder(r.Body).Decode(v)

		testMethod(t, r, "POST")
		if !reflect.DeepEqual(v, input) {
			t.Errorf("Request body = %+v, want %+v", v, input)
		}

		fmt.Fprint(w, `{"url":"u"}`)
	})

	ctx := context.Background()
	label, _, err := client.Issues.CreateLabel(ctx, "o", "r", input)
	if err != nil {
		t.Errorf("Issues.CreateLabel returned error: %v", err)
	}

	want := &Label{URL: String("u")}
	if !reflect.DeepEqual(label, want) {
		t.Errorf("Issues.CreateLabel returned %+v, want %+v", label, want)
	}
}

func TestIssuesService_CreateLabel_invalidOwner(t *testing.T) {
	client, _, _, teardown := setup()
	defer teardown()

	ctx := context.Background()
	_, _, err := client.Issues.CreateLabel(ctx, "%", "%", nil)
	testURLParseError(t, err)
}

func TestIssuesService_EditLabel(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	input := &Label{Name: String("z")}

	mux.HandleFunc("/repos/o/r/labels/n", func(w http.ResponseWriter, r *http.Request) {
		v := new(Label)
		json.NewDecoder(r.Body).Decode(v)

		testMethod(t, r, "PATCH")
		if !reflect.DeepEqual(v, input) {
			t.Errorf("Request body = %+v, want %+v", v, input)
		}

		fmt.Fprint(w, `{"url":"u"}`)
	})

	ctx := context.Background()
	label, _, err := client.Issues.EditLabel(ctx, "o", "r", "n", input)
	if err != nil {
		t.Errorf("Issues.EditLabel returned error: %v", err)
	}

	want := &Label{URL: String("u")}
	if !reflect.DeepEqual(label, want) {
		t.Errorf("Issues.EditLabel returned %+v, want %+v", label, want)
	}
}

func TestIssuesService_EditLabel_invalidOwner(t *testing.T) {
	client, _, _, teardown := setup()
	defer teardown()

	ctx := context.Background()
	_, _, err := client.Issues.EditLabel(ctx, "%", "%", "%", nil)
	testURLParseError(t, err)
}

func TestIssuesService_DeleteLabel(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/repos/o/r/labels/n", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
	})

	ctx := context.Background()
	_, err := client.Issues.DeleteLabel(ctx, "o", "r", "n")
	if err != nil {
		t.Errorf("Issues.DeleteLabel returned error: %v", err)
	}
}

func TestIssuesService_DeleteLabel_invalidOwner(t *testing.T) {
	client, _, _, teardown := setup()
	defer teardown()

	ctx := context.Background()
	_, err := client.Issues.DeleteLabel(ctx, "%", "%", "%")
	testURLParseError(t, err)
}

func TestIssuesService_ListLabelsByIssue(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/repos/o/r/issues/1/labels", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testFormValues(t, r, values{"page": "2"})
		fmt.Fprint(w, `[{"name":"a","id":1},{"name":"b","id":2}]`)
	})

	opt := &ListOptions{Page: 2}
	ctx := context.Background()
	labels, _, err := client.Issues.ListLabelsByIssue(ctx, "o", "r", 1, opt)
	if err != nil {
		t.Errorf("Issues.ListLabelsByIssue returned error: %v", err)
	}

	want := []*Label{
		{Name: String("a"), ID: Int64(1)},
		{Name: String("b"), ID: Int64(2)},
	}
	if !reflect.DeepEqual(labels, want) {
		t.Errorf("Issues.ListLabelsByIssue returned %+v, want %+v", labels, want)
	}
}

func TestIssuesService_ListLabelsByIssue_invalidOwner(t *testing.T) {
	client, _, _, teardown := setup()
	defer teardown()

	ctx := context.Background()
	_, _, err := client.Issues.ListLabelsByIssue(ctx, "%", "%", 1, nil)
	testURLParseError(t, err)
}

func TestIssuesService_AddLabelsToIssue(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	input := []string{"a", "b"}

	mux.HandleFunc("/repos/o/r/issues/1/labels", func(w http.ResponseWriter, r *http.Request) {
		var v []string
		json.NewDecoder(r.Body).Decode(&v)

		testMethod(t, r, "POST")
		if !reflect.DeepEqual(v, input) {
			t.Errorf("Request body = %+v, want %+v", v, input)
		}

		fmt.Fprint(w, `[{"url":"u"}]`)
	})

	ctx := context.Background()
	labels, _, err := client.Issues.AddLabelsToIssue(ctx, "o", "r", 1, input)
	if err != nil {
		t.Errorf("Issues.AddLabelsToIssue returned error: %v", err)
	}

	want := []*Label{{URL: String("u")}}
	if !reflect.DeepEqual(labels, want) {
		t.Errorf("Issues.AddLabelsToIssue returned %+v, want %+v", labels, want)
	}
}

func TestIssuesService_AddLabelsToIssue_invalidOwner(t *testing.T) {
	client, _, _, teardown := setup()
	defer teardown()

	ctx := context.Background()
	_, _, err := client.Issues.AddLabelsToIssue(ctx, "%", "%", 1, nil)
	testURLParseError(t, err)
}

func TestIssuesService_RemoveLabelForIssue(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/repos/o/r/issues/1/labels/l", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
	})

	ctx := context.Background()
	_, err := client.Issues.RemoveLabelForIssue(ctx, "o", "r", 1, "l")
	if err != nil {
		t.Errorf("Issues.RemoveLabelForIssue returned error: %v", err)
	}
}

func TestIssuesService_RemoveLabelForIssue_invalidOwner(t *testing.T) {
	client, _, _, teardown := setup()
	defer teardown()

	ctx := context.Background()
	_, err := client.Issues.RemoveLabelForIssue(ctx, "%", "%", 1, "%")
	testURLParseError(t, err)
}

func TestIssuesService_ReplaceLabelsForIssue(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	input := []string{"a", "b"}

	mux.HandleFunc("/repos/o/r/issues/1/labels", func(w http.ResponseWriter, r *http.Request) {
		var v []string
		json.NewDecoder(r.Body).Decode(&v)

		testMethod(t, r, "PUT")
		if !reflect.DeepEqual(v, input) {
			t.Errorf("Request body = %+v, want %+v", v, input)
		}

		fmt.Fprint(w, `[{"url":"u"}]`)
	})

	ctx := context.Background()
	labels, _, err := client.Issues.ReplaceLabelsForIssue(ctx, "o", "r", 1, input)
	if err != nil {
		t.Errorf("Issues.ReplaceLabelsForIssue returned error: %v", err)
	}

	want := []*Label{{URL: String("u")}}
	if !reflect.DeepEqual(labels, want) {
		t.Errorf("Issues.ReplaceLabelsForIssue returned %+v, want %+v", labels, want)
	}
}

func TestIssuesService_ReplaceLabelsForIssue_invalidOwner(t *testing.T) {
	client, _, _, teardown := setup()
	defer teardown()

	ctx := context.Background()
	_, _, err := client.Issues.ReplaceLabelsForIssue(ctx, "%", "%", 1, nil)
	testURLParseError(t, err)
}

func TestIssuesService_RemoveLabelsForIssue(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/repos/o/r/issues/1/labels", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
	})

	ctx := context.Background()
	_, err := client.Issues.RemoveLabelsForIssue(ctx, "o", "r", 1)
	if err != nil {
		t.Errorf("Issues.RemoveLabelsForIssue returned error: %v", err)
	}
}

func TestIssuesService_RemoveLabelsForIssue_invalidOwner(t *testing.T) {
	client, _, _, teardown := setup()
	defer teardown()

	ctx := context.Background()
	_, err := client.Issues.RemoveLabelsForIssue(ctx, "%", "%", 1)
	testURLParseError(t, err)
}

func TestIssuesService_ListLabelsForMilestone(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/repos/o/r/milestones/1/labels", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testFormValues(t, r, values{"page": "2"})
		fmt.Fprint(w, `[{"name": "a"},{"name": "b"}]`)
	})

	opt := &ListOptions{Page: 2}
	ctx := context.Background()
	labels, _, err := client.Issues.ListLabelsForMilestone(ctx, "o", "r", 1, opt)
	if err != nil {
		t.Errorf("Issues.ListLabelsForMilestone returned error: %v", err)
	}

	want := []*Label{{Name: String("a")}, {Name: String("b")}}
	if !reflect.DeepEqual(labels, want) {
		t.Errorf("Issues.ListLabelsForMilestone returned %+v, want %+v", labels, want)
	}
}

func TestIssuesService_ListLabelsForMilestone_invalidOwner(t *testing.T) {
	client, _, _, teardown := setup()
	defer teardown()

	ctx := context.Background()
	_, _, err := client.Issues.ListLabelsForMilestone(ctx, "%", "%", 1, nil)
	testURLParseError(t, err)
}
