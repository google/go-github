// Copyright 2026 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package github

import "testing"

func TestRunnerGroupRequests_SelectedWorkflows(t *testing.T) {
	t.Parallel()

	for _, tt := range []struct {
		name      string
		workflows []string
		field     string
	}{
		{name: "omitted"},
		{name: "empty", workflows: []string{}, field: `,"selected_workflows":[]`},
		{name: "selected", workflows: []string{"o/r/.github/workflows/build.yml@refs/heads/main"}, field: `,"selected_workflows":["o/r/.github/workflows/build.yml@refs/heads/main"]`},
	} {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			for _, request := range []struct {
				name string
				body any
				tail string
			}{
				{"create organization", CreateRunnerGroupRequest{Name: "runner-group", RestrictedToWorkflows: new(true), SelectedWorkflows: tt.workflows}, ""},
				{"update organization", UpdateRunnerGroupRequest{Name: "runner-group", RestrictedToWorkflows: new(true), SelectedWorkflows: tt.workflows}, ""},
				{"create enterprise", CreateEnterpriseRunnerGroupRequest{Name: "runner-group", RestrictedToWorkflows: new(true), SelectedWorkflows: tt.workflows}, ""},
				{"update enterprise", UpdateEnterpriseRunnerGroupRequest{Name: new("runner-group"), RestrictedToWorkflows: new(true), SelectedWorkflows: tt.workflows}, ""},
				{"update organization and remove network", UpdateRunnerGroupRequest{Name: "runner-group", RestrictedToWorkflows: new(true), SelectedWorkflows: tt.workflows, RemoveNetworkConfiguration: true}, `,"network_configuration_id":null`},
				{"update enterprise and remove network", UpdateEnterpriseRunnerGroupRequest{Name: new("runner-group"), RestrictedToWorkflows: new(true), SelectedWorkflows: tt.workflows, RemoveNetworkConfiguration: true}, `,"network_configuration_id":null`},
			} {
				t.Run(request.name, func(t *testing.T) {
					t.Parallel()
					want := `{"name":"runner-group","restricted_to_workflows":true` + tt.field + request.tail + `}`
					testJSONMarshalOnly(t, request.body, want)
				})
			}
		})
	}
}
