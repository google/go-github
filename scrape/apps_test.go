package scrape

import (
	"net/http"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-github/v41/github"
)

func Test_AppRestrictionsEnabled(t *testing.T) {
	tests := []struct {
		description string
		testFile    string
		org         string
		want        bool
	}{
		{
			description: "return true for enabled orgs",
			testFile:    "access-restrictions-enabled.html",
			want:        true,
		},
		{
			description: "return false for disabled orgs",
			testFile:    "access-restrictions-disabled.html",
			want:        false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			client, mux, cleanup := setup()
			defer cleanup()

			mux.HandleFunc("/organizations/o/settings/oauth_application_policy", func(w http.ResponseWriter, r *http.Request) {
				copyTestFile(w, tt.testFile)
			})

			got, err := client.AppRestrictionsEnabled("o")
			if err != nil {
				t.Errorf("AppRestrictionsEnabled returned err: %v", err)
			}
			if want := tt.want; got != want {
				t.Errorf("AppRestrictionsEnabled returned %t, want %t", got, want)
			}
		})
	}
}

func Test_ListOAuthApps(t *testing.T) {
	client, mux, cleanup := setup()
	defer cleanup()

	mux.HandleFunc("/organizations/e/settings/oauth_application_policy", func(w http.ResponseWriter, r *http.Request) {
		copyTestFile(w, "access-restrictions-enabled.html")
	})

	got, err := client.ListOAuthApps("e")
	if err != nil {
		t.Errorf("ListOAuthApps(e) returned err: %v", err)
	}
	want := []OAuthApp{
		{
			ID:          22222,
			Name:        "Coveralls",
			Description: "Test coverage history and statistics.",
			State:       OAuthAppRequested,
			RequestedBy: "willnorris",
		},
		{
			ID:    530107,
			Name:  "Google Cloud Platform",
			State: OAuthAppApproved,
		},
		{
			ID:          231424,
			Name:        "GitKraken",
			Description: "An intuitive, cross-platform Git client that doesn't suck, built by @axosoft and made with @nodegit & @ElectronJS.",
			State:       OAuthAppDenied,
		},
	}
	if !cmp.Equal(got, want) {
		t.Errorf("ListOAuthApps(o) returned %v, want %v", got, want)
	}
}

func Test_CreateApp(t *testing.T) {
	client, mux, cleanup := setup()
	defer cleanup()

	mux.HandleFunc("/apps/settings/new", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusCreated)
	})

	if _, err := client.CreateApp(&AppManifest{
		URL: github.String("https://example.com"),
		HookAttributes: map[string]string{
			"url": "https://example.com/hook",
		},
	}); err != nil {
		t.Fatalf("CreateApp: %v", err)
	}
}
