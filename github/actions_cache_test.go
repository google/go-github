// Copyright 2022 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package github

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestActionsService_ListCaches(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/repos/o/r/actions/caches", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testFormValues(t, r, values{"page": "2"})
		fmt.Fprint(w,
			`{
				"total_count":1,
				"actions_caches":[{"id":1}]
			}`,
		)
	})

	opts := &ActionsCacheListOptions{ListOptions: ListOptions{Page: 2}}
	ctx := context.Background()
	cacheList, _, err := client.Actions.ListCaches(ctx, "o", "r", opts)
	if err != nil {
		t.Errorf("Actions.ListCaches returned error: %v", err)
	}

	want := &ActionsCacheList{TotalCount: 1, ActionsCaches: []*ActionsCache{{ID: Ptr(int64(1))}}}
	if !cmp.Equal(cacheList, want) {
		t.Errorf("Actions.ListCaches returned %+v, want %+v", cacheList, want)
	}

	const methodName = "ListCaches"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Actions.ListCaches(ctx, "\n", "\n", opts)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Actions.ListCaches(ctx, "o", "r", opts)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestActionsService_ListCaches_invalidOwner(t *testing.T) {
	t.Parallel()
	client, _, _ := setup(t)

	ctx := context.Background()
	_, _, err := client.Actions.ListCaches(ctx, "%", "r", nil)
	testURLParseError(t, err)
}

func TestActionsService_ListCaches_invalidRepo(t *testing.T) {
	t.Parallel()
	client, _, _ := setup(t)

	ctx := context.Background()
	_, _, err := client.Actions.ListCaches(ctx, "o", "%", nil)
	testURLParseError(t, err)
}

func TestActionsService_ListCaches_notFound(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/repos/o/r/actions/caches", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		w.WriteHeader(http.StatusNotFound)
	})

	ctx := context.Background()
	caches, resp, err := client.Actions.ListCaches(ctx, "o", "r", nil)
	if err == nil {
		t.Errorf("Expected HTTP 404 response")
	}
	if got, want := resp.Response.StatusCode, http.StatusNotFound; got != want {
		t.Errorf("Actions.ListCaches return status %d, want %d", got, want)
	}
	if caches != nil {
		t.Errorf("Actions.ListCaches return %+v, want nil", caches)
	}
}

func TestActionsService_DeleteCachesByKey(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/repos/o/r/actions/caches", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
		testFormValues(t, r, values{"key": "1", "ref": "main"})
	})

	ctx := context.Background()
	_, err := client.Actions.DeleteCachesByKey(ctx, "o", "r", "1", Ptr("main"))
	if err != nil {
		t.Errorf("Actions.DeleteCachesByKey return error: %v", err)
	}

	const methodName = "DeleteCachesByKey"
	testBadOptions(t, methodName, func() (err error) {
		_, err = client.Actions.DeleteCachesByKey(ctx, "\n", "\n", "\n", Ptr("\n"))
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		return client.Actions.DeleteCachesByKey(ctx, "o", "r", "1", Ptr("main"))
	})
}

func TestActionsService_DeleteCachesByKey_invalidOwner(t *testing.T) {
	t.Parallel()
	client, _, _ := setup(t)

	ctx := context.Background()
	_, err := client.Actions.DeleteCachesByKey(ctx, "%", "r", "1", Ptr("main"))
	testURLParseError(t, err)
}

func TestActionsService_DeleteCachesByKey_invalidRepo(t *testing.T) {
	t.Parallel()
	client, _, _ := setup(t)

	ctx := context.Background()
	_, err := client.Actions.DeleteCachesByKey(ctx, "o", "%", "1", Ptr("main"))
	testURLParseError(t, err)
}
func TestActionsService_DeleteCachesByKey_notFound(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/repos/o/r/actions/artifacts/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
		w.WriteHeader(http.StatusNotFound)
	})

	ctx := context.Background()
	resp, err := client.Actions.DeleteCachesByKey(ctx, "o", "r", "1", Ptr("main"))
	if err == nil {
		t.Errorf("Expected HTTP 404 response")
	}
	if got, want := resp.Response.StatusCode, http.StatusNotFound; got != want {
		t.Errorf("Actions.DeleteCachesByKey return status %d, want %d", got, want)
	}
}

func TestActionsService_DeleteCachesByID(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/repos/o/r/actions/caches/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
	})

	ctx := context.Background()
	_, err := client.Actions.DeleteCachesByID(ctx, "o", "r", 1)
	if err != nil {
		t.Errorf("Actions.DeleteCachesByID return error: %v", err)
	}

	const methodName = "DeleteCachesByID"
	testBadOptions(t, methodName, func() (err error) {
		_, err = client.Actions.DeleteCachesByID(ctx, "\n", "\n", 0)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		return client.Actions.DeleteCachesByID(ctx, "o", "r", 1)
	})
}

func TestActionsService_DeleteCachesByID_invalidOwner(t *testing.T) {
	t.Parallel()
	client, _, _ := setup(t)

	ctx := context.Background()
	_, err := client.Actions.DeleteCachesByID(ctx, "%", "r", 1)
	testURLParseError(t, err)
}

func TestActionsService_DeleteCachesByID_invalidRepo(t *testing.T) {
	t.Parallel()
	client, _, _ := setup(t)

	ctx := context.Background()
	_, err := client.Actions.DeleteCachesByID(ctx, "o", "%", 1)
	testURLParseError(t, err)
}

func TestActionsService_DeleteCachesByID_notFound(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("repos/o/r/actions/caches/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
		w.WriteHeader(http.StatusNotFound)
	})

	ctx := context.Background()
	resp, err := client.Actions.DeleteCachesByID(ctx, "o", "r", 1)
	if err == nil {
		t.Errorf("Expected HTTP 404 response")
	}
	if got, want := resp.Response.StatusCode, http.StatusNotFound; got != want {
		t.Errorf("Actions.DeleteCachesByID return status %d, want %d", got, want)
	}
}

func TestActionsService_GetCacheUsageForRepo(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/repos/o/r/actions/cache/usage", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w,
			`{
				"full_name":"test-cache",
				"active_caches_size_in_bytes":1000,
				"active_caches_count":1
			}`,
		)
	})

	ctx := context.Background()
	cacheUse, _, err := client.Actions.GetCacheUsageForRepo(ctx, "o", "r")
	if err != nil {
		t.Errorf("Actions.GetCacheUsageForRepo returned error: %v", err)
	}

	want := &ActionsCacheUsage{FullName: "test-cache", ActiveCachesSizeInBytes: 1000, ActiveCachesCount: 1}
	if !cmp.Equal(cacheUse, want) {
		t.Errorf("Actions.GetCacheUsageForRepo returned %+v, want %+v", cacheUse, want)
	}

	const methodName = "GetCacheUsageForRepo"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Actions.GetCacheUsageForRepo(ctx, "\n", "\n")
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Actions.GetCacheUsageForRepo(ctx, "o", "r")
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestActionsService_GetCacheUsageForRepo_invalidOwner(t *testing.T) {
	t.Parallel()
	client, _, _ := setup(t)

	ctx := context.Background()
	_, _, err := client.Actions.GetCacheUsageForRepo(ctx, "%", "r")
	testURLParseError(t, err)
}

func TestActionsService_GetCacheUsageForRepo_invalidRepo(t *testing.T) {
	t.Parallel()
	client, _, _ := setup(t)

	ctx := context.Background()
	_, _, err := client.Actions.GetCacheUsageForRepo(ctx, "o", "%")
	testURLParseError(t, err)
}

func TestActionsService_GetCacheUsageForRepo_notFound(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/repos/o/r/actions/cache/usage", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		w.WriteHeader(http.StatusNotFound)
	})

	ctx := context.Background()
	caches, resp, err := client.Actions.GetCacheUsageForRepo(ctx, "o", "r")
	if err == nil {
		t.Errorf("Expected HTTP 404 response")
	}
	if got, want := resp.Response.StatusCode, http.StatusNotFound; got != want {
		t.Errorf("Actions.GetCacheUsageForRepo return status %d, want %d", got, want)
	}
	if caches != nil {
		t.Errorf("Actions.GetCacheUsageForRepo return %+v, want nil", caches)
	}
}

func TestActionsService_ListCacheUsageByRepoForOrg(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/orgs/o/actions/cache/usage-by-repository", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testFormValues(t, r, values{"page": "2", "per_page": "1"})
		fmt.Fprint(w,
			`{
				"total_count":1,
				"repository_cache_usages":[{"full_name":"test-cache","active_caches_size_in_bytes":1000,"active_caches_count":1}]
			}`,
		)
	})

	opts := &ListOptions{PerPage: 1, Page: 2}
	ctx := context.Background()
	cacheList, _, err := client.Actions.ListCacheUsageByRepoForOrg(ctx, "o", opts)
	if err != nil {
		t.Errorf("Actions.ListCacheUsageByRepoForOrg returned error: %v", err)
	}

	want := &ActionsCacheUsageList{TotalCount: 1, RepoCacheUsage: []*ActionsCacheUsage{{FullName: "test-cache", ActiveCachesSizeInBytes: 1000, ActiveCachesCount: 1}}}
	if !cmp.Equal(cacheList, want) {
		t.Errorf("Actions.ListCacheUsageByRepoForOrg returned %+v, want %+v", cacheList, want)
	}

	const methodName = "ListCacheUsageByRepoForOrg"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Actions.ListCacheUsageByRepoForOrg(ctx, "\n", opts)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Actions.ListCacheUsageByRepoForOrg(ctx, "o", opts)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestActionsService_ListCacheUsageByRepoForOrg_invalidOrganization(t *testing.T) {
	t.Parallel()
	client, _, _ := setup(t)

	ctx := context.Background()
	_, _, err := client.Actions.ListCacheUsageByRepoForOrg(ctx, "%", nil)
	testURLParseError(t, err)
}

func TestActionsService_ListCacheUsageByRepoForOrg_notFound(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/orgs/o/actions/cache/usage-by-repository", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		w.WriteHeader(http.StatusNotFound)
	})

	ctx := context.Background()
	caches, resp, err := client.Actions.ListCacheUsageByRepoForOrg(ctx, "o", nil)
	if err == nil {
		t.Errorf("Expected HTTP 404 response")
	}
	if got, want := resp.Response.StatusCode, http.StatusNotFound; got != want {
		t.Errorf("Actions.ListCacheUsageByRepoForOrg return status %d, want %d", got, want)
	}
	if caches != nil {
		t.Errorf("Actions.ListCacheUsageByRepoForOrg return %+v, want nil", caches)
	}
}

func TestActionsService_GetCacheUsageForOrg(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/orgs/o/actions/cache/usage", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w,
			`{
				"total_active_caches_size_in_bytes":1000,
				"total_active_caches_count":1
			}`,
		)
	})

	ctx := context.Background()
	cache, _, err := client.Actions.GetTotalCacheUsageForOrg(ctx, "o")
	if err != nil {
		t.Errorf("Actions.GetTotalCacheUsageForOrg returned error: %v", err)
	}

	want := &TotalCacheUsage{TotalActiveCachesUsageSizeInBytes: 1000, TotalActiveCachesCount: 1}
	if !cmp.Equal(cache, want) {
		t.Errorf("Actions.GetTotalCacheUsageForOrg returned %+v, want %+v", cache, want)
	}

	const methodName = "GetTotalCacheUsageForOrg"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Actions.GetTotalCacheUsageForOrg(ctx, "\n")
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Actions.GetTotalCacheUsageForOrg(ctx, "o")
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestActionsService_GetCacheUsageForOrg_invalidOrganization(t *testing.T) {
	t.Parallel()
	client, _, _ := setup(t)

	ctx := context.Background()
	_, _, err := client.Actions.GetTotalCacheUsageForOrg(ctx, "%")
	testURLParseError(t, err)
}

func TestActionsService_GetCacheUsageForOrg_notFound(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/orgs/o/actions/cache/usage", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		w.WriteHeader(http.StatusNotFound)
	})

	ctx := context.Background()
	caches, resp, err := client.Actions.GetTotalCacheUsageForOrg(ctx, "o")
	if err == nil {
		t.Errorf("Expected HTTP 404 response")
	}
	if got, want := resp.Response.StatusCode, http.StatusNotFound; got != want {
		t.Errorf("Actions.GetTotalCacheUsageForOrg return status %d, want %d", got, want)
	}
	if caches != nil {
		t.Errorf("Actions.GetTotalCacheUsageForOrg return %+v, want nil", caches)
	}
}

func TestActionsService_GetCacheUsageForEnterprise(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/enterprises/e/actions/cache/usage", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w,
			`{
				"total_active_caches_size_in_bytes":1000,
				"total_active_caches_count":1
			}`,
		)
	})

	ctx := context.Background()
	cache, _, err := client.Actions.GetTotalCacheUsageForEnterprise(ctx, "e")
	if err != nil {
		t.Errorf("Actions.GetTotalCacheUsageForEnterprise returned error: %v", err)
	}

	want := &TotalCacheUsage{TotalActiveCachesUsageSizeInBytes: 1000, TotalActiveCachesCount: 1}
	if !cmp.Equal(cache, want) {
		t.Errorf("Actions.GetTotalCacheUsageForEnterprise returned %+v, want %+v", cache, want)
	}

	const methodName = "GetTotalCacheUsageForEnterprise"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Actions.GetTotalCacheUsageForEnterprise(ctx, "\n")
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Actions.GetTotalCacheUsageForEnterprise(ctx, "e")
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestActionsService_GetCacheUsageForEnterprise_invalidEnterprise(t *testing.T) {
	t.Parallel()
	client, _, _ := setup(t)

	ctx := context.Background()
	_, _, err := client.Actions.GetTotalCacheUsageForEnterprise(ctx, "%")
	testURLParseError(t, err)
}

func TestActionsService_GetCacheUsageForEnterprise_notFound(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/enterprises/e/actions/cache/usage", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		w.WriteHeader(http.StatusNotFound)
	})

	ctx := context.Background()
	caches, resp, err := client.Actions.GetTotalCacheUsageForEnterprise(ctx, "o")
	if err == nil {
		t.Errorf("Expected HTTP 404 response")
	}
	if got, want := resp.Response.StatusCode, http.StatusNotFound; got != want {
		t.Errorf("Actions.GetTotalCacheUsageForEnterprise return status %d, want %d", got, want)
	}
	if caches != nil {
		t.Errorf("Actions.GetTotalCacheUsageForEnterprise return %+v, want nil", caches)
	}
}

func TestActionsCache_Marshal(t *testing.T) {
	t.Parallel()
	testJSONMarshal(t, &ActionsCache{}, "{}")

	u := &ActionsCache{
		ID:             Ptr(int64(1)),
		Ref:            Ptr("refAction"),
		Key:            Ptr("key1"),
		Version:        Ptr("alpha"),
		LastAccessedAt: &Timestamp{referenceTime},
		CreatedAt:      &Timestamp{referenceTime},
		SizeInBytes:    Ptr(int64(1)),
	}

	want := `{
		"id": 1,
		"ref": "refAction",
		"key": "key1",
		"version": "alpha",
		"last_accessed_at": ` + referenceTimeStr + `,
		"created_at": ` + referenceTimeStr + `,
		"size_in_bytes": 1
	}`

	testJSONMarshal(t, u, want)
}

func TestActionsCacheList_Marshal(t *testing.T) {
	t.Parallel()
	testJSONMarshal(t, &ActionsCacheList{}, "{}")

	u := &ActionsCacheList{
		TotalCount: 2,
		ActionsCaches: []*ActionsCache{
			{
				ID:             Ptr(int64(1)),
				Key:            Ptr("key1"),
				Version:        Ptr("alpha"),
				LastAccessedAt: &Timestamp{referenceTime},
				CreatedAt:      &Timestamp{referenceTime},
				SizeInBytes:    Ptr(int64(1)),
			},
			{
				ID:             Ptr(int64(2)),
				Ref:            Ptr("refAction"),
				LastAccessedAt: &Timestamp{referenceTime},
				CreatedAt:      &Timestamp{referenceTime},
				SizeInBytes:    Ptr(int64(1)),
			},
		},
	}
	want := `{
		"total_count": 2,
		"actions_caches": [{
				"id": 1,
				"key": "key1",
				"version": "alpha",
				"last_accessed_at": ` + referenceTimeStr + `,
				"created_at": ` + referenceTimeStr + `,
				"size_in_bytes": 1
			},
			{
				"id": 2,
				"ref": "refAction",
				"last_accessed_at": ` + referenceTimeStr + `,
				"created_at": ` + referenceTimeStr + `,
				"size_in_bytes": 1
		}]
	}`
	testJSONMarshal(t, u, want)
}

func TestActionsCacheUsage_Marshal(t *testing.T) {
	t.Parallel()
	testJSONMarshal(t, &ActionsCacheUsage{}, "{}")

	u := &ActionsCacheUsage{
		FullName:                "cache_usage1",
		ActiveCachesSizeInBytes: 2,
		ActiveCachesCount:       2,
	}

	want := `{
		"full_name": "cache_usage1",
		"active_caches_size_in_bytes": 2,
		"active_caches_count": 2
	}`

	testJSONMarshal(t, u, want)
}

func TestActionsCacheUsageList_Marshal(t *testing.T) {
	t.Parallel()
	testJSONMarshal(t, &ActionsCacheUsageList{}, "{}")

	u := &ActionsCacheUsageList{
		TotalCount: 1,
		RepoCacheUsage: []*ActionsCacheUsage{
			{
				FullName:                "cache_usage1",
				ActiveCachesSizeInBytes: 2,
				ActiveCachesCount:       2,
			},
		},
	}

	want := `{
		"total_count": 1,
		"repository_cache_usages": [{
			"full_name": "cache_usage1",
			"active_caches_size_in_bytes": 2,
			"active_caches_count": 2
		}]
	}`

	testJSONMarshal(t, u, want)
}

func TestTotalCacheUsage_Marshal(t *testing.T) {
	t.Parallel()
	testJSONMarshal(t, &TotalCacheUsage{}, "{}")

	u := &TotalCacheUsage{
		TotalActiveCachesUsageSizeInBytes: 2,
		TotalActiveCachesCount:            2,
	}

	want := `{
		"total_active_caches_size_in_bytes": 2,
		"total_active_caches_count": 2
	}`

	testJSONMarshal(t, u, want)
}
