// Copyright 2025 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package github

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
)

func TestPrivateRegistriesService_ListOrganizationPrivateRegistries(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/orgs/o/private-registries", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testFormValues(t, r, values{
			"page": "2",
		})
		fmt.Fprint(w, `{
  "total_count": 1,
  "configurations": [
    {
      "name": "MAVEN_REPOSITORY_SECRET",
      "registry_type": "maven_repository",
      "username": "monalisa",
      "created_at": "2019-08-10T14:59:22Z",
      "updated_at": "2020-01-10T14:59:22Z",
      "visibility": "selected"
    }
  ]
}`)
	})

	opts := &ListOptions{Page: 2}
	ctx := t.Context()
	privateRegistries, _, err := client.PrivateRegistries.ListOrganizationPrivateRegistries(ctx, "o", opts)
	if err != nil {
		t.Fatalf("PrivateRegistries.ListOrganizationPrivateRegistries returned error: %v", err)
	}

	want := &PrivateRegistries{
		TotalCount: Ptr(1),
		Configurations: []*PrivateRegistry{
			{
				Name:         Ptr("MAVEN_REPOSITORY_SECRET"),
				RegistryType: Ptr("maven_repository"),
				Username:     Ptr("monalisa"),
				CreatedAt:    &Timestamp{time.Date(2019, time.August, 10, 14, 59, 22, 0, time.UTC)},
				UpdatedAt:    &Timestamp{time.Date(2020, time.January, 10, 14, 59, 22, 0, time.UTC)},
				Visibility:   Ptr(PrivateRegistryVisibilitySelected),
			},
		},
	}
	if diff := cmp.Diff(want, privateRegistries); diff != "" {
		t.Errorf("PrivateRegistries.ListOrganizationPrivateRegistries mismatch (-want +got):\\n%v", diff)
	}

	const methodName = "ListOrganizationPrivateRegistries"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.PrivateRegistries.ListOrganizationPrivateRegistries(ctx, "\n", opts)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.PrivateRegistries.ListOrganizationPrivateRegistries(ctx, "o", opts)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})

	// still allow both set (no validation enforced) – ensure it does not error
	ctxBypass := context.WithValue(t.Context(), BypassRateLimitCheck, true)
	if _, _, err = client.PrivateRegistries.ListOrganizationPrivateRegistries(ctxBypass, "o", opts); err != nil {
		t.Fatalf("unexpected error when both before/after set: %v", err)
	}
}

func TestPrivateRegistriesService_CreateOrganizationPrivateRegistry(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	input := &CreateOrganizationPrivateRegistry{
		RegistryType:          "maven_repository",
		URL:                   "https://maven.pkg.github.com/OWNER/REPOSITORY",
		Username:              Ptr("monalisa"),
		EncryptedValue:        "encrypted_value",
		KeyID:                 "key_id",
		Visibility:            PrivateRegistryVisibilitySelected,
		SelectedRepositoryIDs: []int64{1, 2, 3},
	}

	mux.HandleFunc("/orgs/o/private-registries", func(w http.ResponseWriter, r *http.Request) {
		v := new(CreateOrganizationPrivateRegistry)
		assertNilError(t, json.NewDecoder(r.Body).Decode(v))

		testMethod(t, r, "POST")
		if !cmp.Equal(v, input) {
			t.Errorf("Request body = %+v, want %+v", v, input)
		}

		fmt.Fprint(w, `{
  "name": "MAVEN_REPOSITORY_SECRET",
  "registry_type": "maven_repository",
  "username": "monalisa",
  "created_at": "2019-08-10T14:59:22Z",
  "updated_at": "2020-01-10T14:59:22Z",
  "visibility": "selected"
}`)
	})

	ctx := t.Context()
	privateRegistry, _, err := client.PrivateRegistries.CreateOrganizationPrivateRegistry(ctx, "o", *input)
	if err != nil {
		t.Fatalf("PrivateRegistries.CreateOrganizationPrivateRegistries returned error: %v", err)
	}

	want := &PrivateRegistry{
		Name:         Ptr("MAVEN_REPOSITORY_SECRET"),
		RegistryType: Ptr("maven_repository"),
		Username:     Ptr("monalisa"),
		CreatedAt:    &Timestamp{time.Date(2019, time.August, 10, 14, 59, 22, 0, time.UTC)},
		UpdatedAt:    &Timestamp{time.Date(2020, time.January, 10, 14, 59, 22, 0, time.UTC)},
		Visibility:   Ptr(PrivateRegistryVisibilitySelected),
	}
	if diff := cmp.Diff(want, privateRegistry); diff != "" {
		t.Errorf("PrivateRegistries.CreateOrganizationPrivateRegistries mismatch (-want +got):\\n%v", diff)
	}

	const methodName = "CreateOrganizationPrivateRegistry"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.PrivateRegistries.CreateOrganizationPrivateRegistry(ctx, "\n", *input)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.PrivateRegistries.CreateOrganizationPrivateRegistry(ctx, "o", *input)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestPrivateRegistriesService_GetOrganizationPrivateRegistriesPublicKey(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/orgs/o/private-registries/public-key", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{
  "key_id": "0123456789",
  "key": "public_key"
}`)
	})
	ctx := t.Context()
	publicKey, _, err := client.PrivateRegistries.GetOrganizationPrivateRegistriesPublicKey(ctx, "o")
	if err != nil {
		t.Fatalf("PrivateRegistries.GetOrganizationPrivateRegistriesPublicKey returned error: %v", err)
	}

	want := &PublicKey{
		KeyID: Ptr("0123456789"),
		Key:   Ptr("public_key"),
	}
	if diff := cmp.Diff(want, publicKey); diff != "" {
		t.Errorf("PrivateRegistries.GetOrganizationPrivateRegistriesPublicKey mismatch (-want +got):\\n%v", diff)
	}

	const methodName = "GetOrganizationPrivateRegistriesPublicKey"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.PrivateRegistries.GetOrganizationPrivateRegistriesPublicKey(ctx, "\n")
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.PrivateRegistries.GetOrganizationPrivateRegistriesPublicKey(ctx, "o")
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestPrivateRegistriesService_GetOrganizationPrivateRegistry(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/orgs/o/private-registries/MAVEN_REPOSITORY_SECRET", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{
  "name": "MAVEN_REPOSITORY_SECRET",
  "registry_type": "maven_repository",
  "username": "monalisa",
  "created_at": "2019-08-10T14:59:22Z",
  "updated_at": "2020-01-10T14:59:22Z",
  "visibility": "selected"
}`)
	})
	ctx := t.Context()
	privateRegistry, _, err := client.PrivateRegistries.GetOrganizationPrivateRegistry(ctx, "o", "MAVEN_REPOSITORY_SECRET")
	if err != nil {
		t.Fatalf("PrivateRegistries.GetOrganizationPrivateRegistry returned error: %v", err)
	}

	want := &PrivateRegistry{
		Name:         Ptr("MAVEN_REPOSITORY_SECRET"),
		RegistryType: Ptr("maven_repository"),
		Username:     Ptr("monalisa"),
		CreatedAt:    &Timestamp{time.Date(2019, time.August, 10, 14, 59, 22, 0, time.UTC)},
		UpdatedAt:    &Timestamp{time.Date(2020, time.January, 10, 14, 59, 22, 0, time.UTC)},
		Visibility:   Ptr(PrivateRegistryVisibilitySelected),
	}
	if diff := cmp.Diff(want, privateRegistry); diff != "" {
		t.Errorf("PrivateRegistries.GetOrganizationPrivateRegistry mismatch (-want +got):\\n%v", diff)
	}

	const methodName = "GetOrganizationPrivateRegistry"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.PrivateRegistries.GetOrganizationPrivateRegistry(ctx, "\n", "MAVEN_REPOSITORY_SECRET")
		return err
	})
	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.PrivateRegistries.GetOrganizationPrivateRegistry(ctx, "o", "MAVEN_REPOSITORY_SECRET")
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestPrivateRegistries_UpdateOrganizationPrivateRegistry(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	input := &UpdateOrganizationPrivateRegistry{
		Username:       Ptr("monalisa"),
		EncryptedValue: Ptr("encrypted_value"),
		KeyID:          Ptr("key_id"),
		Visibility:     Ptr(PrivateRegistryVisibilitySelected),
	}

	mux.HandleFunc("/orgs/o/private-registries/MAVEN_REPOSITORY_SECRET", func(w http.ResponseWriter, r *http.Request) {
		v := new(UpdateOrganizationPrivateRegistry)
		assertNilError(t, json.NewDecoder(r.Body).Decode(v))

		testMethod(t, r, "PATCH")
		if !cmp.Equal(v, input) {
			t.Errorf("Request body = %+v, want %+v", v, input)
		}

		fmt.Fprint(w, `{
  "name": "MAVEN_REPOSITORY_SECRET",
  "registry_type": "maven_repository",
  "username": "monalisa",
  "created_at": "2019-08-10T14:59:22Z",
  "updated_at": "2020-01-10T14:59:22Z",
  "visibility": "selected"
}`)
	})

	ctx := t.Context()
	privateRegistry, _, err := client.PrivateRegistries.UpdateOrganizationPrivateRegistry(ctx, "o", "MAVEN_REPOSITORY_SECRET", *input)
	if err != nil {
		t.Fatalf("PrivateRegistries.UpdateOrganizationPrivateRegistry returned error: %v", err)
	}

	want := &PrivateRegistry{
		Name:         Ptr("MAVEN_REPOSITORY_SECRET"),
		RegistryType: Ptr("maven_repository"),
		Username:     Ptr("monalisa"),
		CreatedAt:    &Timestamp{time.Date(2019, time.August, 10, 14, 59, 22, 0, time.UTC)},
		UpdatedAt:    &Timestamp{time.Date(2020, time.January, 10, 14, 59, 22, 0, time.UTC)},
		Visibility:   Ptr(PrivateRegistryVisibilitySelected),
	}
	if diff := cmp.Diff(want, privateRegistry); diff != "" {
		t.Errorf("PrivateRegistries.UpdateOrganizationPrivateRegistry mismatch (-want +got):\\n%v", diff)
	}

	const methodName = "UpdateOrganizationPrivateRegistry"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.PrivateRegistries.UpdateOrganizationPrivateRegistry(ctx, "\n", "MAVEN_REPOSITORY_SECRET", *input)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.PrivateRegistries.UpdateOrganizationPrivateRegistry(ctx, "o", "MAVEN_REPOSITORY_SECRET", *input)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestPrivateRegistriesService_DeleteOrganizationPrivateRegistry(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/orgs/o/private-registries/MAVEN_REPOSITORY_SECRET", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
		w.WriteHeader(http.StatusNoContent)
	})
	ctx := t.Context()
	_, err := client.PrivateRegistries.DeleteOrganizationPrivateRegistry(ctx, "o", "MAVEN_REPOSITORY_SECRET")
	if err != nil {
		t.Fatalf("PrivateRegistries.DeleteOrganizationPrivateRegistry returned error: %v", err)
	}

	const methodName = "DeleteOrganizationPrivateRegistry"
	testBadOptions(t, methodName, func() (err error) {
		_, err = client.PrivateRegistries.DeleteOrganizationPrivateRegistry(ctx, "\n", "MAVEN_REPOSITORY_SECRET")
		return err
	})
	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		return client.PrivateRegistries.DeleteOrganizationPrivateRegistry(ctx, "o", "MAVEN_REPOSITORY_SECRET")
	})
}
