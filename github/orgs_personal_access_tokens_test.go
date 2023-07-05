package github

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestOrganizationsService_ReviewPersonalAccessTokenRequest(t *testing.T) {
	type body struct {
		Action string `json:"action"`
		Reason string `json:"reason,omitempty"`
	}

	client, mux, _, teardown := setup()
	defer teardown()

	input := &body{
		Action: "a",
		Reason: "r",
	}

	mux.HandleFunc("/orgs/o/personal-access-token-requests/r", func(w http.ResponseWriter, r *http.Request) {
		v := new(body)
		json.NewDecoder(r.Body).Decode(v)

		testMethod(t, r, http.MethodPost)
		if !cmp.Equal(v, input) {
			t.Errorf("Request body = %+v, want %+v", v, input)
		}

		w.WriteHeader(http.StatusNoContent)
	})

	ctx := context.Background()
	res, err := client.Organizations.ReviewPersonalAccessTokenRequest(ctx, "o", "r", input.Action, input.Reason)
	if err != nil {
		t.Errorf("Organizations.ReviewPersonalAccessTokenRequest returned error: %v", err)
	}

	if res.StatusCode != http.StatusNoContent {
		t.Errorf("Organizations.ReviewPersonalAccessTokenRequest returned %v, want %v", res.StatusCode, http.StatusNoContent)
	}

	const methodName = "ReviewPersonalAccessTokenRequest"
	testBadOptions(t, methodName, func() (err error) {
		_, err = client.Organizations.ReviewPersonalAccessTokenRequest(ctx, "\n", "", "", "")
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		return client.Organizations.ReviewPersonalAccessTokenRequest(ctx, "o", "r", input.Action, input.Reason)
	})
}
