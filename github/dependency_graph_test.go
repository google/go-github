// Copyright 2023 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package github

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestDependencyGraphService_GetSBOM(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/repos/owner/repo/dependency-graph/sbom", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{
			"sbom": {
				"creationInfo": {
					"created": `+referenceTimeStr+`
				},
				"name": "owner/repo",
				"packages": [
					{
						"name": "rubygems:rails",
						"versionInfo": "1.0.0"
					}
				]
			}
		}`)
	})

	ctx := t.Context()
	sbom, _, err := client.DependencyGraph.GetSBOM(ctx, "owner", "repo")
	if err != nil {
		t.Errorf("DependencyGraph.GetSBOM returned error: %v", err)
	}

	want := &SBOM{
		&SBOMInfo{
			CreationInfo: &CreationInfo{
				Created: &referenceTimestamp,
			},
			Name: new("owner/repo"),
			Packages: []*RepoDependencies{
				{
					Name:        new("rubygems:rails"),
					VersionInfo: new("1.0.0"),
				},
			},
		},
	}

	if !cmp.Equal(sbom, want) {
		t.Errorf("DependencyGraph.GetSBOM returned %+v, want %+v", sbom, want)
	}

	const methodName = "GetSBOM"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.DependencyGraph.GetSBOM(ctx, "\n", "\n")
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.DependencyGraph.GetSBOM(ctx, "owner", "repo")
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestSBOMGeneration_UUID(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		gen  *SBOMGeneration
		want string
	}{
		{
			name: "reads the UUID off a generation response",
			gen: &SBOMGeneration{
				SBOMURL: new("https://api.github.com/repos/owner/repo/dependency-graph/sbom/fetch-report/c0ccba21-ccba-4292-9afd-a64781f7e98a"),
			},
			want: "c0ccba21-ccba-4292-9afd-a64781f7e98a",
		},
		{
			name: "nil SBOMGeneration",
			gen:  nil,
		},
		{
			name: "unset SBOMURL",
			gen:  &SBOMGeneration{},
		},
		{
			name: "SBOMURL has no path segments",
			gen:  &SBOMGeneration{SBOMURL: new("c0ccba21")},
		},
		{
			name: "SBOMURL ends in a separator",
			gen: &SBOMGeneration{
				SBOMURL: new("https://api.github.com/repos/owner/repo/dependency-graph/sbom/fetch-report/"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if got := tt.gen.UUID(); got != tt.want {
				t.Errorf("SBOMGeneration.UUID() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestDependencyGraphService_GenerateSBOM(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/repos/owner/repo/dependency-graph/sbom/generate-report", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		w.WriteHeader(http.StatusCreated)
		fmt.Fprint(w, `{"sbom_url":"https://api.github.com/repos/owner/repo/dependency-graph/sbom/fetch-report/1234"}`)
	})

	ctx := t.Context()
	generation, resp, err := client.DependencyGraph.GenerateSBOM(ctx, "owner", "repo")
	if err != nil {
		t.Errorf("DependencyGraph.GenerateSBOM returned error: %v", err)
	}

	if resp.StatusCode != http.StatusCreated {
		t.Errorf("DependencyGraph.GenerateSBOM returned status = %v, want %v", resp.StatusCode, http.StatusCreated)
	}

	want := &SBOMGeneration{
		SBOMURL: new("https://api.github.com/repos/owner/repo/dependency-graph/sbom/fetch-report/1234"),
	}
	if !cmp.Equal(generation, want) {
		t.Errorf("DependencyGraph.GenerateSBOM returned %+v, want %+v", generation, want)
	}

	const methodName = "GenerateSBOM"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.DependencyGraph.GenerateSBOM(ctx, "\n", "\n")
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.DependencyGraph.GenerateSBOM(ctx, "owner", "repo")
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestDependencyGraphService_FetchSBOM_Download(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		downloadStatus int
		downloadBody   string
		wantSBOM       *SBOM
		wantErr        error
	}{
		{
			name:           "success",
			downloadStatus: http.StatusOK,
			downloadBody:   `{"name":"owner/repo"}`,
			wantSBOM: &SBOM{
				SBOM: &SBOMInfo{Name: new("owner/repo")},
			},
		},
		{
			name:           "forbidden",
			downloadStatus: http.StatusForbidden,
			wantErr: &ErrorResponse{
				Response: &http.Response{StatusCode: http.StatusForbidden},
			},
		},
		{
			name:           "truncated JSON",
			downloadStatus: http.StatusOK,
			downloadBody:   `{"name":"owner/`,
			wantErr:        io.ErrUnexpectedEOF,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			client, mux, _ := setup(t)

			client, err := client.Clone(WithAuthToken("secret-token"))
			if err != nil {
				t.Fatalf("failed to clone client: %v", err)
			}

			downloadServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				testMethod(t, r, "GET")
				testHeader(t, r, "Authorization", "")
				w.WriteHeader(tt.downloadStatus)
				fmt.Fprint(w, tt.downloadBody)
			}))
			t.Cleanup(downloadServer.Close)

			mux.HandleFunc("/repos/owner/repo/dependency-graph/sbom/fetch-report/1234", func(w http.ResponseWriter, r *http.Request) {
				testMethod(t, r, "GET")
				testHeader(t, r, "Authorization", "Bearer secret-token")
				http.Redirect(w, r, downloadServer.URL, http.StatusFound)
			})

			sbom, redirectURL, resp, err := client.DependencyGraph.FetchSBOM(t.Context(), "owner", "repo", "1234", http.DefaultClient)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("DependencyGraph.FetchSBOM error = %#v, want %#v", err, tt.wantErr)
			}

			if resp == nil || resp.StatusCode != http.StatusFound {
				t.Errorf("DependencyGraph.FetchSBOM response = %v, want status 302", resp)
			}
			if diff := cmp.Diff(tt.wantSBOM, sbom); diff != "" {
				t.Errorf("DependencyGraph.FetchSBOM SBOM mismatch (-want +got):\n%v", diff)
			}
			if redirectURL != "" {
				t.Errorf("DependencyGraph.FetchSBOM redirectURL = %q, want empty", redirectURL)
			}
		})
	}
}

func TestDependencyGraphService_FetchSBOM_Redirect(t *testing.T) {
	t.Parallel()
	client, mux, serverURL := setup(t)
	downloadURL := serverURL + baseURLPath + "/download"

	mux.HandleFunc("/repos/owner/repo/dependency-graph/sbom/fetch-report/1234", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		http.Redirect(w, r, downloadURL, http.StatusFound)
	})
	mux.HandleFunc("/download", func(http.ResponseWriter, *http.Request) {
		t.Error("download URL was fetched with a nil followRedirectsClient")
	})

	sbom, redirectURL, resp, err := client.DependencyGraph.FetchSBOM(t.Context(), "owner", "repo", "1234", nil)
	if err != nil {
		t.Fatalf("DependencyGraph.FetchSBOM returned error: %v", err)
	}

	if sbom != nil {
		t.Errorf("DependencyGraph.FetchSBOM SBOM = %v, want nil", sbom)
	}
	if redirectURL != downloadURL {
		t.Errorf("DependencyGraph.FetchSBOM redirectURL = %q, want %q", redirectURL, downloadURL)
	}
	if resp == nil || resp.StatusCode != http.StatusFound {
		t.Errorf("DependencyGraph.FetchSBOM response = %v, want status 302", resp)
	}
}

func TestDependencyGraphService_FetchSBOM_DownloadTransportError(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	const downloadURL = "https://example.com/sbom"
	wantErr := errors.New("download failed")
	downloadClient := &http.Client{
		Transport: roundTripperFunc(func(r *http.Request) (*http.Response, error) {
			testMethod(t, r, "GET")
			if got := r.URL.String(); got != downloadURL {
				t.Errorf("download URL = %q, want %q", got, downloadURL)
			}
			return nil, wantErr
		}),
	}

	mux.HandleFunc("/repos/owner/repo/dependency-graph/sbom/fetch-report/1234", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		http.Redirect(w, r, downloadURL, http.StatusFound)
	})

	sbom, redirectURL, resp, err := client.DependencyGraph.FetchSBOM(t.Context(), "owner", "repo", "1234", downloadClient)
	if !errors.Is(err, wantErr) {
		t.Errorf("DependencyGraph.FetchSBOM error = %v, want %v", err, wantErr)
	}

	if sbom != nil || redirectURL != "" {
		t.Errorf("DependencyGraph.FetchSBOM returned (%v, %q), want (nil, empty)", sbom, redirectURL)
	}
	if resp == nil || resp.StatusCode != http.StatusFound {
		t.Errorf("DependencyGraph.FetchSBOM response = %v, want status 302", resp)
	}
}

func TestDependencyGraphService_FetchSBOM_AcceptedError(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/repos/owner/repo/dependency-graph/sbom/fetch-report/1234", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		w.WriteHeader(http.StatusAccepted)
	})

	sbom, redirectURL, resp, err := client.DependencyGraph.FetchSBOM(t.Context(), "owner", "repo", "1234", http.DefaultClient)

	if _, ok := errors.AsType[*AcceptedError](err); !ok {
		t.Errorf("DependencyGraph.FetchSBOM error = %v, want *AcceptedError", err)
	}

	if sbom != nil || redirectURL != "" {
		t.Errorf("DependencyGraph.FetchSBOM returned (%v, %q), want (nil, empty)", sbom, redirectURL)
	}
	if resp == nil || resp.StatusCode != http.StatusAccepted {
		t.Errorf("DependencyGraph.FetchSBOM response = %v, want status 202", resp)
	}
}

func TestDependencyGraphService_FetchSBOM_NoRedirect(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/repos/owner/repo/dependency-graph/sbom/fetch-report/1234", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		w.WriteHeader(http.StatusOK)
	})

	sbom, redirectURL, resp, err := client.DependencyGraph.FetchSBOM(t.Context(), "owner", "repo", "1234", http.DefaultClient)
	if err == nil {
		t.Error("DependencyGraph.FetchSBOM returned no error for a non-redirect response")
	}

	if sbom != nil || redirectURL != "" {
		t.Errorf("DependencyGraph.FetchSBOM returned (%v, %q), want (nil, empty)", sbom, redirectURL)
	}
	if resp == nil || resp.StatusCode != http.StatusOK {
		t.Errorf("DependencyGraph.FetchSBOM response = %v, want status 200", resp)
	}
}

func TestDependencyGraphService_FetchSBOM_requestErrors(t *testing.T) {
	t.Parallel()
	client, _, _ := setup(t)

	ctx := t.Context()

	const methodName = "FetchSBOM"
	testBadOptions(t, methodName, func() (err error) {
		_, _, _, err = client.DependencyGraph.FetchSBOM(ctx, "\n", "\n", "\n", http.DefaultClient)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		sbom, redirectURL, resp, err := client.DependencyGraph.FetchSBOM(ctx, "owner", "repo", "1234", http.DefaultClient)
		if sbom != nil || redirectURL != "" {
			t.Errorf("testNewRequestAndDoFailure %v returned (%v, %q), want (nil, empty)", methodName, sbom, redirectURL)
		}
		return resp, err
	})
}
