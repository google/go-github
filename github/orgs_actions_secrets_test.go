package github

import (
	"context"
	"fmt"
	"net/http"
	"reflect"
	"testing"
)

func TestOrganizationsService_GetPublicKey(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/orgs/o/actions/secrets/public-key", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{"key_id":"012345678912345678","key":"2Sg8iYjAxxmI2LvUXpJjkYrMxURPc8r+dB7TJyvv1234"}`)
	})

	key, _, err := client.Organizations.GetPublicKey(context.Background(), "o")
	if err != nil {
		t.Errorf("OrgsActions.GetPublicKey returned error: %v", err)
	}

	want := &OrgsActionsPublicKey{KeyID: String("012345678912345678"), Key: String("2Sg8iYjAxxmI2LvUXpJjkYrMxURPc8r+dB7TJyvv1234")}
	if !reflect.DeepEqual(key, want) {
		t.Errorf("OrgsActions.GetPublicKey returned %+v, want %+v", key, want)
	}
}