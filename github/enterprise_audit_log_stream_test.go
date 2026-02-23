// Copyright 2026 The go-github AUTHORS. All rights reserved.
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

func TestEnterpriseService_GetAuditLogStreamKey(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/enterprises/e/audit-log/stream-key", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{"key_id":"1234","key":"2Sg8iYjAxxmI2LvUXpJjkYrMxURPc8r+dB7TJyvv1234"}`)
	})

	ctx := context.Background()
	key, _, err := client.Enterprise.GetAuditLogStreamKey(ctx, "e")
	if err != nil {
		t.Errorf("Enterprise.GetAuditLogStreamKey returned error: %v", err)
	}

	want := &AuditLogStreamKey{
		KeyID:     Ptr("1234"),
		PublicKey: Ptr("2Sg8iYjAxxmI2LvUXpJjkYrMxURPc8r+dB7TJyvv1234"),
	}
	if !cmp.Equal(key, want) {
		t.Errorf("Enterprise.GetAuditLogStreamKey returned %+v, want %+v", key, want)
	}

	const methodName = "GetAuditLogStreamKey"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Enterprise.GetAuditLogStreamKey(ctx, "\n")
		return err
	})
	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Enterprise.GetAuditLogStreamKey(ctx, "e")
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestEnterpriseService_ListAuditLogStreams(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/enterprises/e/audit-log/streams", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `[{"id":1,"stream_type":"Splunk","stream_details":"US","enabled":true}]`)
	})

	ctx := context.Background()
	streams, _, err := client.Enterprise.ListAuditLogStreams(ctx, "e")
	if err != nil {
		t.Errorf("Enterprise.ListAuditLogStreams returned error: %v", err)
	}

	want := []*AuditLogStream{
		{
			ID:            Ptr(int64(1)),
			StreamType:    Ptr("Splunk"),
			StreamDetails: Ptr("US"),
			Enabled:       Ptr(true),
		},
	}
	if !cmp.Equal(streams, want) {
		t.Errorf("Enterprise.ListAuditLogStreams returned %+v, want %+v", streams, want)
	}

	const methodName = "ListAuditLogStreams"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Enterprise.ListAuditLogStreams(ctx, "\n")
		return err
	})
	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Enterprise.ListAuditLogStreams(ctx, "e")
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestEnterpriseService_GetAuditLogStream(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/enterprises/e/audit-log/streams/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{"id":1,"stream_type":"Datadog","stream_details":"US","enabled":true}`)
	})

	ctx := context.Background()
	stream, _, err := client.Enterprise.GetAuditLogStream(ctx, "e", 1)
	if err != nil {
		t.Errorf("Enterprise.GetAuditLogStream returned error: %v", err)
	}

	want := &AuditLogStream{
		ID:            Ptr(int64(1)),
		StreamType:    Ptr("Datadog"),
		StreamDetails: Ptr("US"),
		Enabled:       Ptr(true),
	}
	if !cmp.Equal(stream, want) {
		t.Errorf("Enterprise.GetAuditLogStream returned %+v, want %+v", stream, want)
	}

	const methodName = "GetAuditLogStream"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Enterprise.GetAuditLogStream(ctx, "\n", 1)
		return err
	})
	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Enterprise.GetAuditLogStream(ctx, "e", 1)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestEnterpriseService_CreateAuditLogStream(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/enterprises/e/audit-log/streams", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		fmt.Fprint(w, `{"id":2,"stream_type":"Datadog","stream_details":"US3","enabled":false}`)
	})

	input := NewDatadogStreamConfig(false, &DatadogConfig{
		EncryptedToken: Ptr("ENCRYPTED"),
		Site:           Ptr("US3"),
		KeyID:          Ptr("v1"),
	})

	ctx := context.Background()
	stream, _, err := client.Enterprise.CreateAuditLogStream(ctx, "e", input)
	if err != nil {
		t.Errorf("Enterprise.CreateAuditLogStream returned error: %v", err)
	}

	want := &AuditLogStream{
		ID:            Ptr(int64(2)),
		StreamType:    Ptr("Datadog"),
		StreamDetails: Ptr("US3"),
		Enabled:       Ptr(false),
	}
	if !cmp.Equal(stream, want) {
		t.Errorf("Enterprise.CreateAuditLogStream returned %+v, want %+v", stream, want)
	}

	const methodName = "CreateAuditLogStream"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Enterprise.CreateAuditLogStream(ctx, "\n", input)
		return err
	})
	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Enterprise.CreateAuditLogStream(ctx, "e", input)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestEnterpriseService_UpdateAuditLogStream(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/enterprises/e/audit-log/streams/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PATCH")
		fmt.Fprint(w, `{"id":1,"stream_type":"Splunk","stream_details":"splunk.example.com","enabled":true}`)
	})

	input := NewSplunkStreamConfig(true, &SplunkConfig{
		Domain:         Ptr("splunk.example.com"),
		Port:           Ptr(uint16(8089)),
		KeyID:          Ptr("v1"),
		EncryptedToken: Ptr("ENCRYPTED"),
		SSLVerify:      Ptr(true),
	})

	ctx := context.Background()
	stream, _, err := client.Enterprise.UpdateAuditLogStream(ctx, "e", 1, input)
	if err != nil {
		t.Errorf("Enterprise.UpdateAuditLogStream returned error: %v", err)
	}

	want := &AuditLogStream{
		ID:            Ptr(int64(1)),
		StreamType:    Ptr("Splunk"),
		StreamDetails: Ptr("splunk.example.com"),
		Enabled:       Ptr(true),
	}
	if !cmp.Equal(stream, want) {
		t.Errorf("Enterprise.UpdateAuditLogStream returned %+v, want %+v", stream, want)
	}

	const methodName = "UpdateAuditLogStream"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Enterprise.UpdateAuditLogStream(ctx, "\n", 1, input)
		return err
	})
	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Enterprise.UpdateAuditLogStream(ctx, "e", 1, input)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestEnterpriseService_DeleteAuditLogStream(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/enterprises/e/audit-log/streams/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
	})

	ctx := context.Background()
	_, err := client.Enterprise.DeleteAuditLogStream(ctx, "e", 1)
	if err != nil {
		t.Errorf("Enterprise.DeleteAuditLogStream returned error: %v", err)
	}

	const methodName = "DeleteAuditLogStream"
	testBadOptions(t, methodName, func() (err error) {
		_, err = client.Enterprise.DeleteAuditLogStream(ctx, "\n", 1)
		return err
	})
	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		return client.Enterprise.DeleteAuditLogStream(ctx, "e", 1)
	})
}

func TestNewAzureBlobStreamConfig(t *testing.T) {
	t.Parallel()
	cfg := &AzureBlobConfig{
		KeyID:           Ptr("v1"),
		EncryptedSASURL: Ptr("ENCRYPTED"),
		Container:       Ptr("my-container"),
	}
	got := NewAzureBlobStreamConfig(true, cfg)
	if got.StreamType == nil || *got.StreamType != "Azure Blob Storage" {
		t.Errorf("NewAzureBlobStreamConfig StreamType = %v, want Azure Blob Storage", got.StreamType)
	}
	if got.Enabled == nil || !*got.Enabled {
		t.Errorf("NewAzureBlobStreamConfig Enabled = %v, want true", got.Enabled)
	}
	if got.VendorSpecific != cfg {
		t.Errorf("NewAzureBlobStreamConfig VendorSpecific = %v, want %v", got.VendorSpecific, cfg)
	}
}

func TestNewDatadogStreamConfig(t *testing.T) {
	t.Parallel()
	cfg := &DatadogConfig{
		EncryptedToken: Ptr("ENCRYPTED"),
		Site:           Ptr("US"),
		KeyID:          Ptr("v1"),
	}
	got := NewDatadogStreamConfig(false, cfg)
	if got.StreamType == nil || *got.StreamType != "Datadog" {
		t.Errorf("NewDatadogStreamConfig StreamType = %v, want Datadog", got.StreamType)
	}
	if got.Enabled == nil || *got.Enabled {
		t.Errorf("NewDatadogStreamConfig Enabled = %v, want false", got.Enabled)
	}
}
