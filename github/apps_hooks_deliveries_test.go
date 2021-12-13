package github

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestAppsService_ListHookDeliveries(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/app/hook/deliveries", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testFormValues(t, r, values{"cursor": "v1_12077215967"})
		fmt.Fprint(w, `[{"id":1}, {"id":2}]`)
	})

	opts := &ListCursorOptions{Cursor: "v1_12077215967"}

	ctx := context.Background()

	deliveries, _, err := client.Apps.ListHookDeliveries(ctx, opts)
	if err != nil {
		t.Errorf("Apps.ListHookDeliveries returned error: %v", err)
	}

	want := []*HookDelivery{{ID: Int64(1)}, {ID: Int64(2)}}
	if d := cmp.Diff(deliveries, want); d != "" {
		t.Errorf("Apps.ListHooks want (-), got (+):\n%s", d)
	}
}
