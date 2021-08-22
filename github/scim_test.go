package github

import (
	"context"
	"net/http"
	"testing"
)

func TestSCIMService_GetSCIMProvisioningInfoForUser(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/scim/v2/organizations/o/Users/123", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
	})

	ctx := context.Background()
	_, err := client.SCIM.GetSCIMProvisioningInfoForUser(ctx, "o", "123")
	if err != nil {
		t.Errorf("SCIM.GetSCIMProvisioningInfoForUser returned error: %v", err)
	}

	const methodName = "GetSCIMProvisioningInfoForUser"
	testBadOptions(t, methodName, func() error {
		_, err := client.SCIM.GetSCIMProvisioningInfoForUser(ctx, "\n", "123")
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		return client.SCIM.GetSCIMProvisioningInfoForUser(ctx, "o", "123")
	})
}
