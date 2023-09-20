package github

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestEmojisService_ListEmojis(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/emojis", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{"+1": "+1.png"}`)
	})

	ctx := context.Background()
	emoji, _, err := client.Emojis.ListEmojis(ctx)
	if err != nil {
		t.Errorf("ListEmojis returned error: %v", err)
	}

	want := map[string]string{"+1": "+1.png"}
	if !cmp.Equal(want, emoji) {
		t.Errorf("ListEmojis returned %+v, want %+v", emoji, want)
	}

	const methodName = "ListEmojis"
	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Emojis.ListEmojis(ctx)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}
