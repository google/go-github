package github

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestOrganizationsService_ListCustomRoles(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/organizations/o/custom_roles", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{"total_count": 1, "custom_roles": [{ "id": 1, "name": "Developer"}]}`)
	})

	ctx := context.Background()
	apps, _, err := client.Organizations.ListCustomRoles(ctx, "o")
	if err != nil {
		t.Errorf("Organizations.ListCustomRoles returned error: %v", err)
	}

	want := &OrginizationCustomRoles{TotalCount: Int(1), CustomRoles: []*CustomRoles{{ID: Int64(1), Name: String("Developer")}}}
	if !cmp.Equal(apps, want) {
		t.Errorf("Organizations.ListCustomRoles returned %+v, want %+v", apps, want)
	}

	const methodName = "ListCustomRoles"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Organizations.ListCustomRoles(ctx, "\no")
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Organizations.ListCustomRoles(ctx, "o")
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}
