package github

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestOrganizationsService_GetAllOrganizationRepositoryRulesets(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/orgs/o/rulesets", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testFormValues(t, r, values{})
		fmt.Fprint(w, `[{
			"id": 26110,
			"name": "test ruleset",
			"target": "branch",
			"source_type": "Organization",
			"source": "my-org",
			"enforcement": "active",
			"bypass_mode": "none",
			"node_id": "nid",
			"_links": {
			  "self": {
				"href": "https://api.github.com/orgs/o/rulesets/26110"
			  }
			}
		}]`)
	})

	ctx := context.Background()
	rulesets, _, err := client.Organizations.GetAllOrganizationRepositoryRulesets(ctx, "o")
	if err != nil {
		t.Errorf("Organizations.GetAllOrganizationRepositoryRulesets returned error: %v", err)
	}

	want := []*Ruleset{{
		ID:          26110,
		Name:        "test ruleset",
		Target:      String("branch"),
		SourceType:  String("Organization"),
		Source:      "my-org",
		Enforcement: "active",
		BypassMode:  String("none"),
		NodeID:      String("nid"),
		Links: &RulesetLinks{
			Self: &RulesetLink{HRef: String("https://api.github.com/orgs/o/rulesets/26110")},
		},
	}}
	if !cmp.Equal(rulesets, want) {
		t.Errorf("Organizations.GetAllOrganizationRepositoryRulesets returned %+v, want %+v", rulesets, want)
	}

	const methodName = "GetAllOrganizationRepositoryRulesets"

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Organizations.GetAllOrganizationRepositoryRulesets(ctx, "o")
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}
