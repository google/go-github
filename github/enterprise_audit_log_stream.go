// Copyright 2026 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package github

import (
	"context"
	"encoding/json"
	"fmt"
)

// AuditLogStream represents an audit log stream configuration for an enterprise.
type AuditLogStream struct {
	ID            *int64          `json:"id,omitempty"`
	StreamType    *string         `json:"stream_type,omitempty"`
	StreamDetails json.RawMessage `json:"stream_details,omitempty"`
	Enabled       *bool           `json:"enabled,omitempty"`
	CreatedAt     *Timestamp      `json:"created_at,omitempty"`
	UpdatedAt     *Timestamp      `json:"updated_at,omitempty"`
	PausedAt      *Timestamp      `json:"paused_at,omitempty"`
}

// AuditLogStreamConfig represents a configuration for creating/updating an audit log stream.
type AuditLogStreamConfig struct {
	Enabled        *bool                      `json:"enabled,omitempty"`
	StreamType     *string                    `json:"stream_type,omitempty"`
	VendorSpecific AuditLogStreamVendorConfig `json:"vendor_specific,omitempty"`
}

// AuditLogStreamVendorConfig is a marker interface for vendor-specific audit log
// stream configurations.
type AuditLogStreamVendorConfig interface {
	isAuditLogStreamVendorConfig()
}

// AzureBlobConfig represents vendor specific config for Azure Blob Storage.
type AzureBlobConfig struct {
	KeyID           *string `json:"key_id,omitempty"`
	EncryptedSASURL *string `json:"encrypted_sas_url,omitempty"`
	Container       *string `json:"container,omitempty"`
}

// AzureHubConfig represents vendor specific config for Azure Event Hubs.
type AzureHubConfig struct {
	Name                *string `json:"name,omitempty"`
	EncryptedConnString *string `json:"encrypted_connstring,omitempty"`
	KeyID               *string `json:"key_id,omitempty"`
}

// AmazonS3OIDCConfig represents vendor specific config for Amazon S3 with OIDC authentication.
type AmazonS3OIDCConfig struct {
	Bucket             *string `json:"bucket,omitempty"`
	Region             *string `json:"region,omitempty"`
	KeyID              *string `json:"key_id,omitempty"`
	AuthenticationType *string `json:"authentication_type,omitempty"`
	ARNRole            *string `json:"arn_role,omitempty"`
}

// AmazonS3AccessKeysConfig represents vendor specific config for Amazon S3 with access keys authentication.
type AmazonS3AccessKeysConfig struct {
	Bucket               *string `json:"bucket,omitempty"`
	Region               *string `json:"region,omitempty"`
	KeyID                *string `json:"key_id,omitempty"`
	AuthenticationType   *string `json:"authentication_type,omitempty"`
	EncryptedSecretKey   *string `json:"encrypted_secret_key,omitempty"`
	EncryptedAccessKeyID *string `json:"encrypted_access_key_id,omitempty"`
}

// SplunkConfig represents vendor specific config for Splunk.
type SplunkConfig struct {
	Domain         *string `json:"domain,omitempty"`
	Port           *uint16 `json:"port,omitempty"`
	KeyID          *string `json:"key_id,omitempty"`
	EncryptedToken *string `json:"encrypted_token,omitempty"`
	SSLVerify      *bool   `json:"ssl_verify,omitempty"`
}

// HecConfig represents vendor specific config for HTTPS Event Collector.
type HecConfig struct {
	Domain         *string `json:"domain,omitempty"`
	Port           *uint16 `json:"port,omitempty"`
	KeyID          *string `json:"key_id,omitempty"`
	EncryptedToken *string `json:"encrypted_token,omitempty"`
	Path           *string `json:"path,omitempty"`
	SSLVerify      *bool   `json:"ssl_verify,omitempty"`
}

// GoogleCloudConfig represents vendor specific config for Google Cloud Storage.
type GoogleCloudConfig struct {
	Bucket                   *string `json:"bucket,omitempty"`
	KeyID                    *string `json:"key_id,omitempty"`
	EncryptedJSONCredentials *string `json:"encrypted_json_credentials,omitempty"`
}

// DatadogConfig represents vendor specific config for Datadog.
type DatadogConfig struct {
	EncryptedToken *string `json:"encrypted_token,omitempty"`
	Site           *string `json:"site,omitempty"` // Can be one of: US, US3, US5, EU1, US1-FED, AP1
	KeyID          *string `json:"key_id,omitempty"`
}

// Implement the marker interface for all vendor config types.
func (*AzureBlobConfig) isAuditLogStreamVendorConfig()          {}
func (*AzureHubConfig) isAuditLogStreamVendorConfig()           {}
func (*AmazonS3OIDCConfig) isAuditLogStreamVendorConfig()       {}
func (*AmazonS3AccessKeysConfig) isAuditLogStreamVendorConfig() {}
func (*SplunkConfig) isAuditLogStreamVendorConfig()             {}
func (*HecConfig) isAuditLogStreamVendorConfig()                {}
func (*GoogleCloudConfig) isAuditLogStreamVendorConfig()        {}
func (*DatadogConfig) isAuditLogStreamVendorConfig()            {}

// Helper functions for constructing AuditLogStreamConfig per vendor type.

// NewAzureBlobStreamConfig returns an AuditLogStreamConfig for Azure Blob Storage.
func NewAzureBlobStreamConfig(enabled bool, cfg *AzureBlobConfig) *AuditLogStreamConfig {
	streamType := "AzureBlobStorage"
	return &AuditLogStreamConfig{Enabled: &enabled, StreamType: &streamType, VendorSpecific: cfg}
}

// NewAzureHubStreamConfig returns an AuditLogStreamConfig for Azure Event Hubs.
func NewAzureHubStreamConfig(enabled bool, cfg *AzureHubConfig) *AuditLogStreamConfig {
	streamType := "AzureHubs"
	return &AuditLogStreamConfig{Enabled: &enabled, StreamType: &streamType, VendorSpecific: cfg}
}

// NewAmazonS3OIDCStreamConfig returns an AuditLogStreamConfig for Amazon S3 with OIDC auth.
func NewAmazonS3OIDCStreamConfig(enabled bool, cfg *AmazonS3OIDCConfig) *AuditLogStreamConfig {
	streamType := "AmazonS3"
	return &AuditLogStreamConfig{Enabled: &enabled, StreamType: &streamType, VendorSpecific: cfg}
}

// NewAmazonS3AccessKeysStreamConfig returns an AuditLogStreamConfig for Amazon S3 with access keys auth.
func NewAmazonS3AccessKeysStreamConfig(enabled bool, cfg *AmazonS3AccessKeysConfig) *AuditLogStreamConfig {
	streamType := "AmazonS3"
	return &AuditLogStreamConfig{Enabled: &enabled, StreamType: &streamType, VendorSpecific: cfg}
}

// NewSplunkStreamConfig returns an AuditLogStreamConfig for Splunk.
func NewSplunkStreamConfig(enabled bool, cfg *SplunkConfig) *AuditLogStreamConfig {
	streamType := "Splunk"
	return &AuditLogStreamConfig{Enabled: &enabled, StreamType: &streamType, VendorSpecific: cfg}
}

// NewHecStreamConfig returns an AuditLogStreamConfig for HTTPS Event Collector.
func NewHecStreamConfig(enabled bool, cfg *HecConfig) *AuditLogStreamConfig {
	streamType := "Hec"
	return &AuditLogStreamConfig{Enabled: &enabled, StreamType: &streamType, VendorSpecific: cfg}
}

// NewGoogleCloudStreamConfig returns an AuditLogStreamConfig for Google Cloud Storage.
func NewGoogleCloudStreamConfig(enabled bool, cfg *GoogleCloudConfig) *AuditLogStreamConfig {
	streamType := "GoogleCloudStorage"
	return &AuditLogStreamConfig{Enabled: &enabled, StreamType: &streamType, VendorSpecific: cfg}
}

// NewDatadogStreamConfig returns an AuditLogStreamConfig for Datadog.
func NewDatadogStreamConfig(enabled bool, cfg *DatadogConfig) *AuditLogStreamConfig {
	streamType := "Datadog"
	return &AuditLogStreamConfig{Enabled: &enabled, StreamType: &streamType, VendorSpecific: cfg}
}

// ListAuditLogStreams lists the audit log stream configurations for an enterprise.
//
// GitHub API docs: https://docs.github.com/enterprise-cloud@latest/rest/enterprise-admin/audit-log#list-audit-log-stream-configurations-for-an-enterprise
//
//meta:operation GET /enterprises/{enterprise}/audit-log/streams
func (s *EnterpriseService) ListAuditLogStreams(ctx context.Context, enterprise string) ([]*AuditLogStream, *Response, error) {
	u := fmt.Sprintf("enterprises/%v/audit-log/streams", enterprise)

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	var streams []*AuditLogStream
	resp, err := s.client.Do(ctx, req, &streams)
	if err != nil {
		return nil, resp, err
	}

	return streams, resp, nil
}

// CreateAuditLogStream creates an audit log streaming configuration for an enterprise.
//
// GitHub API docs: https://docs.github.com/enterprise-cloud@latest/rest/enterprise-admin/audit-log#create-an-audit-log-streaming-configuration-for-an-enterprise
//
//meta:operation POST /enterprises/{enterprise}/audit-log/streams
func (s *EnterpriseService) CreateAuditLogStream(ctx context.Context, enterprise string, config *AuditLogStreamConfig) (*AuditLogStream, *Response, error) {
	u := fmt.Sprintf("enterprises/%v/audit-log/streams", enterprise)

	req, err := s.client.NewRequest("POST", u, config)
	if err != nil {
		return nil, nil, err
	}

	var stream *AuditLogStream
	resp, err := s.client.Do(ctx, req, &stream)
	if err != nil {
		return nil, resp, err
	}

	return stream, resp, nil
}

// GetAuditLogStream gets an audit log streaming configuration for an enterprise.
//
// GitHub API docs: https://docs.github.com/enterprise-cloud@latest/rest/enterprise-admin/audit-log#get-an-audit-log-streaming-configuration-for-an-enterprise
//
//meta:operation GET /enterprises/{enterprise}/audit-log/streams/{stream_id}
func (s *EnterpriseService) GetAuditLogStream(ctx context.Context, enterprise string, streamID int64) (*AuditLogStream, *Response, error) {
	u := fmt.Sprintf("enterprises/%v/audit-log/streams/%v", enterprise, streamID)

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	var stream *AuditLogStream
	resp, err := s.client.Do(ctx, req, &stream)
	if err != nil {
		return nil, resp, err
	}

	return stream, resp, nil
}

// UpdateAuditLogStream updates an audit log streaming configuration for an enterprise.
//
// GitHub API docs: https://docs.github.com/enterprise-cloud@latest/rest/enterprise-admin/audit-log#update-an-audit-log-streaming-configuration-for-an-enterprise
//
//meta:operation PATCH /enterprises/{enterprise}/audit-log/streams/{stream_id}
func (s *EnterpriseService) UpdateAuditLogStream(ctx context.Context, enterprise string, streamID int64, config *AuditLogStreamConfig) (*AuditLogStream, *Response, error) {
	u := fmt.Sprintf("enterprises/%v/audit-log/streams/%v", enterprise, streamID)

	req, err := s.client.NewRequest("PATCH", u, config)
	if err != nil {
		return nil, nil, err
	}

	var stream *AuditLogStream
	resp, err := s.client.Do(ctx, req, &stream)
	if err != nil {
		return nil, resp, err
	}

	return stream, resp, nil
}

// DeleteAuditLogStream deletes an audit log streaming configuration for an enterprise.
//
// GitHub API docs: https://docs.github.com/enterprise-cloud@latest/rest/enterprise-admin/audit-log#delete-an-audit-log-streaming-configuration-for-an-enterprise
//
//meta:operation DELETE /enterprises/{enterprise}/audit-log/streams/{stream_id}
func (s *EnterpriseService) DeleteAuditLogStream(ctx context.Context, enterprise string, streamID int64) (*Response, error) {
	u := fmt.Sprintf("enterprises/%v/audit-log/streams/%v", enterprise, streamID)

	req, err := s.client.NewRequest("DELETE", u, nil)
	if err != nil {
		return nil, err
	}

	return s.client.Do(ctx, req, nil)
}
