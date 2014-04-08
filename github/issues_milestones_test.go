// Copyright 2013 The go-github AUTHORS. All rights reserved.
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

func TestIssuesService_ListMilestones(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/repos/o/r/milestones", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testFormValues(t, r, values{
			"state":     "closed",
			"sort":      "due_date",
			"direction": "asc",
		})
		fmt.Fprint(w, `[{"number":1}]`)
	})

	opt := &MilestoneListOptions{"closed", "due_date", "asc"}
	milestones, _, err := client.Issues.ListMilestones("o", "r", opt)

	if err != nil {
		t.Errorf("IssuesService.ListMilestones returned error: %v", err)
	}

	want := []Milestone{{Number: Int(1)}}
	if !reflect.DeepEqual(milestones, want) {
		t.Errorf("IssuesService.ListMilestones returned %+v, want %+v", milestones, want)
	}
}

func TestIssuesService_ListMilestones_invalidOwner(t *testing.T) {
	_, _, err := client.Issues.ListMilestones("%", "r", nil)
	testURLParseError(t, err)
}

