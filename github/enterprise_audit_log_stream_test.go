// Copyright 2021 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package github

// func TestEnterpriseService_CreateAuditStream(t *testing.T) {
// 	t.Parallel()
// 	client, mux, _ := setup(t)

// 	mux.HandleFunc("/enterprises/e/audit-log/streams", func(w http.ResponseWriter, r *http.Request) {
// 		testMethod(t, r, "POST")

// 		fmt.Fprint(w, `{
// 			"id": 1,
// 			"stream_type": "Azure Blob Storage",
// 			"stream_details": "US",
// 			"enabled": true,
// 			"created_at": "2024-06-06T08:00:00Z",
// 			"updated_at": "2024-06-06T08:00:00Z",
// 			"paused_at": null
// 		}`)
// 	})

// 	ctx := context.Background()
// 	config := &AuditStreamConfig{
// 		Enabled:    true,
// 		StreamType: "Azure Blob Storage",
// 		VendorSpecific: AzureBlobConfig{
// 			KeyID:           "123",
// 			EncryptedSASURL: "base64-encrypted-sas-url",
// 		},
// 	}

// 	stream, resp, err := client.Enterprise.CreateAuditStream(ctx, "e", config)
// 	if err != nil {
// 		t.Fatalf("CreateAuditStream returned error: %v", err)
// 	}
// 	if resp == nil {
// 		t.Fatal("expected non-nil HTTP response")
// 	}
// 	if stream == nil {
// 		t.Fatal("expected non-nil AuditStreamEntry")
// 	}
// 	if stream.ID != 1 {
// 		t.Errorf("ID = %d, want 1", stream.ID)
// 	}
// 	if stream.StreamType != "Azure Blob Storage" {
// 		t.Errorf("StreamType = %q, want %q", stream.StreamType, "Azure Blob Storage")
// 	}
// 	if !stream.Enabled {
// 		t.Errorf("Enabled = %v, want true", stream.Enabled)
// 	}
// }

// func TestEnterpriseService_GetAuditStream(t *testing.T) {
// 	t.Parallel()
// 	client, mux, _ := setup(t)

// 	mux.HandleFunc("/enterprises/e/audit-log/streams/42", func(w http.ResponseWriter, r *http.Request) {
// 		testMethod(t, r, "GET")
// 		testHeader(t, r, "Accept", "application/vnd.github+json")
// 		fmt.Fprint(w, `{"id":42,"stream_type":"Azure Blob Storage","enabled":true}`)
// 	})

// 	got, resp, err := client.Enterprise.GetAuditStream(context.Background(), "e", 42)
// 	if err != nil {
// 		t.Fatalf("GetAuditStream returned error: %v", err)
// 	}
// 	if resp == nil {
// 		t.Fatal("expected non-nil response")
// 	}
// 	if got == nil || got.ID != 42 || !got.Enabled || got.StreamType != "Azure Blob Storage" {
// 		t.Fatalf("unexpected stream: %+v", got)
// 	}
// }

// func TestEnterpriseService_ListAuditStreams(t *testing.T) {
// 	t.Parallel()
// 	client, mux, _ := setup(t)

// 	mux.HandleFunc("/enterprises/e/audit-log/streams", func(w http.ResponseWriter, r *http.Request) {
// 		testMethod(t, r, "GET")
// 		fmt.Fprint(w, `[
// 			{
// 				"id": 1,
// 				"stream_type": "Azure Blob Storage",
// 				"stream_details": "US",
// 				"enabled": true,
// 				"created_at": "2024-06-06T08:00:00Z",
// 				"updated_at": "2024-06-06T08:00:00Z",
// 				"paused_at": null,
// 				"vendor_specific": {
// 					"key_id": "k1",
// 					"encrypted_sas_url": "enc1"
// 				}
// 			},
// 			{
// 				"id": 2,
// 				"stream_type": "Azure Blob Storage",
// 				"stream_details": "EU",
// 				"enabled": false,
// 				"created_at": "2024-06-07T08:00:00Z",
// 				"updated_at": "2024-06-07T08:00:00Z",
// 				"paused_at": null,
// 				"vendor_specific": {
// 					"key_id": "k2",
// 					"encrypted_sas_url": "enc2"
// 				}
// 			}
// 		]`)
// 	})

// 	ctx := context.Background()
// 	got, resp, err := client.Enterprise.ListAuditStreams(ctx, "e")
// 	if err != nil {
// 		t.Fatalf("ListAuditStreams returned error: %v", err)
// 	}
// 	if resp == nil {
// 		t.Fatal("expected non-nil response")
// 	}
// 	if len(got) != 2 {
// 		t.Fatalf("len(streams) = %d, want 2", len(got))
// 	}
// 	if got[0].ID != 1 {
// 		t.Errorf("streams[0].ID = %d, want 1", got[0].ID)
// 	}
// 	if got[0].StreamType != "Azure Blob Storage" {
// 		t.Errorf("streams[0].StreamType = %q, want %q", got[0].StreamType, "Azure Blob Storage")
// 	}
// 	if !got[0].Enabled {
// 		t.Errorf("streams[0].Enabled = %v, want true", got[0].Enabled)
// 	}
// 	if got[1].Enabled {
// 		t.Errorf("streams[1].Enabled = %v, want false", got[1].Enabled)
// 	}
// }

// func TestEnterpriseService_DeleteAuditStream(t *testing.T) {
// 	t.Parallel()
// 	client, mux, _ := setup(t)

// 	mux.HandleFunc("/enterprises/e/audit-log/streams/42", func(w http.ResponseWriter, r *http.Request) {
// 		testMethod(t, r, "DELETE")
// 		w.WriteHeader(http.StatusNoContent) // 204
// 	})

// 	ctx := context.Background()
// 	resp, err := client.Enterprise.DeleteAuditStream(ctx, "e", 42)
// 	if err != nil {
// 		t.Fatalf("DeleteAuditStream returned error: %v", err)
// 	}
// 	if resp == nil {
// 		t.Fatal("expected non-nil response")
// 	}
// 	if resp.Response.StatusCode != http.StatusNoContent {
// 		t.Errorf("StatusCode = %d, want %d", resp.Response.StatusCode, http.StatusNoContent)
// 	}
// }
