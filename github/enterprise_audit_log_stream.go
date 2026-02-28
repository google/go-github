// Copyright 2026 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package github

import (
	"context"
	"fmt"
)

// AuditLogStream represents an audit log stream configuration for an enterprise.
type AuditLogStream struct {
	ID            *int64     `json:"id,omitempty"`
	StreamType    *string    `json:"stream_type,omitempty"`
	StreamDetails *string    `json:"stream_details,omitempty"`
	Enabled       *bool      `json:"enabled,omitempty"`
	CreatedAt     *Timestamp `json:"created_at,omitempty"`
	UpdatedAt     *Timestamp `json:"updated_at,omitempty"`
	PausedAt      *Timestamp `json:"paused_at,omitempty"`
}

// AuditLogStreamConfig represents a configuration for creating or updating an audit log stream.
type AuditLogStreamConfig struct {
	Enabled        bool                       `json:"enabled"`
	StreamType     string                     `json:"stream_type"`
	VendorSpecific AuditLogStreamVendorConfig `json:"vendor_specific"`
}

// AuditLogStreamVendorConfig is a sealed marker interface for vendor-specific audit log
// stream configurations. Only this package can define implementations.
type AuditLogStreamVendorConfig interface {
	isAuditLogStreamVendorConfig()
}

// AuditLogStreamKey represents the public key used to encrypt secrets for audit log streaming.
type AuditLogStreamKey struct {
	KeyID string `json:"key_id"`
	Key   string `json:"key"`
}

// AzureBlobConfig represents vendor-specific config for Azure Blob Storage.
type AzureBlobConfig struct {
	KeyID           *string `json:"key_id,omitempty"`
	EncryptedSasURL *string `json:"encrypted_sas_url,omitempty"`
	Container       *string `json:"container,omitempty"`
}

// AzureHubConfig represents vendor-specific config for Azure Event Hubs.
type AzureHubConfig struct {
	Name                *string `json:"name,omitempty"`
	EncryptedConnstring *string `json:"encrypted_connstring,omitempty"`
	KeyID               *string `json:"key_id,omitempty"`
}

// AmazonS3OIDCConfig represents vendor-specific config for Amazon S3 with OIDC authentication.
type AmazonS3OIDCConfig struct {
	Bucket             *string `json:"bucket,omitempty"`
	Region             *string `json:"region,omitempty"`
	KeyID              *string `json:"key_id,omitempty"`
	AuthenticationType *string `json:"authentication_type,omitempty"` // Value: "oidc"
	ArnRole            *string `json:"arn_role,omitempty"`
}

// AmazonS3AccessKeysConfig represents vendor-specific config for Amazon S3 with access key authentication.
type AmazonS3AccessKeysConfig struct {
	Bucket               *string `json:"bucket,omitempty"`
	Region               *string `json:"region,omitempty"`
	KeyID                *string `json:"key_id,omitempty"`
	AuthenticationType   *string `json:"authentication_type,omitempty"` // Value: "access_keys"
	EncryptedSecretKey   *string `json:"encrypted_secret_key,omitempty"`
	EncryptedAccessKeyID *string `json:"encrypted_access_key_id,omitempty"`
}

// SplunkConfig represents vendor-specific config for Splunk.
type SplunkConfig struct {
	Domain         *string `json:"domain,omitempty"`
	Port           *uint16 `json:"port,omitempty"`
	KeyID          *string `json:"key_id,omitempty"`
	EncryptedToken *string `json:"encrypted_token,omitempty"`
	SSLVerify      *bool   `json:"ssl_verify,omitempty"`
}

// HecConfig represents vendor-specific config for an HTTPS Event Collector (HEC) endpoint.
type HecConfig struct {
	Domain         *string `json:"domain,omitempty"`
	Port           *uint16 `json:"port,omitempty"`
	KeyID          *string `json:"key_id,omitempty"`
	EncryptedToken *string `json:"encrypted_token,omitempty"`
	Path           *string `json:"path,omitempty"`
	SSLVerify      *bool   `json:"ssl_verify,omitempty"`
}

// GoogleCloudConfig represents vendor-specific config for Google Cloud Storage.
type GoogleCloudConfig struct {
	Bucket                   *string `json:"bucket,omitempty"`
	KeyID                    *string `json:"key_id,omitempty"`
	EncryptedJSONCredentials *string `json:"encrypted_json_credentials,omitempty"`
}

// DatadogConfig represents vendor-specific config for Datadog.
type DatadogConfig struct {
	EncryptedToken *string `json:"encrypted_token,omitempty"`
	Site           *string `json:"site,omitempty"` // One of: US, US3, US5, EU1, US1-FED, AP1
	KeyID          *string `json:"key_id,omitempty"`
}

// Implement the sealed marker interface for all vendor config types.
func (*AzureBlobConfig) isAuditLogStreamVendorConfig()          {}
func (*AzureHubConfig) isAuditLogStreamVendorConfig()           {}
func (*AmazonS3OIDCConfig) isAuditLogStreamVendorConfig()       {}
func (*AmazonS3AccessKeysConfig) isAuditLogStreamVendorConfig() {}
func (*SplunkConfig) isAuditLogStreamVendorConfig()             {}
func (*HecConfig) isAuditLogStreamVendorConfig()                {}
func (*GoogleCloudConfig) isAuditLogStreamVendorConfig()        {}
func (*DatadogConfig) isAuditLogStreamVendorConfig()            {}

// Helper constructors for AuditLogStreamConfig.

// NewAzureBlobStreamConfig returns an AuditLogStreamConfig for Azure Blob Storage.
func NewAzureBlobStreamConfig(enabled bool, cfg *AzureBlobConfig) *AuditLogStreamConfig {
	streamType := "Azure Blob Storage"
	v := AuditLogStreamVendorConfig(cfg)
	return &AuditLogStreamConfig{Enabled: &enabled, StreamType: &streamType, VendorSpecific: &v}
}

// NewAzureHubStreamConfig returns an AuditLogStreamConfig for Azure Event Hubs.
func NewAzureHubStreamConfig(enabled bool, cfg *AzureHubConfig) *AuditLogStreamConfig {
	streamType := "Azure Event Hubs"
	v := AuditLogStreamVendorConfig(cfg)
	return &AuditLogStreamConfig{Enabled: &enabled, StreamType: &streamType, VendorSpecific: &v}
}

// NewAmazonS3OIDCStreamConfig returns an AuditLogStreamConfig for Amazon S3 with OIDC auth.
func NewAmazonS3OIDCStreamConfig(enabled bool, cfg *AmazonS3OIDCConfig) *AuditLogStreamConfig {
	streamType := "Amazon S3"
	v := AuditLogStreamVendorConfig(cfg)
	return &AuditLogStreamConfig{Enabled: &enabled, StreamType: &streamType, VendorSpecific: &v}
}

// NewAmazonS3AccessKeysStreamConfig returns an AuditLogStreamConfig for Amazon S3 with access key auth.
func NewAmazonS3AccessKeysStreamConfig(enabled bool, cfg *AmazonS3AccessKeysConfig) *AuditLogStreamConfig {
	streamType := "Amazon S3"
	v := AuditLogStreamVendorConfig(cfg)
	return &AuditLogStreamConfig{Enabled: &enabled, StreamType: &streamType, VendorSpecific: &v}
}

// NewSplunkStreamConfig returns an AuditLogStreamConfig for Splunk.
func NewSplunkStreamConfig(enabled bool, cfg *SplunkConfig) *AuditLogStreamConfig {
	streamType := "Splunk"
	v := AuditLogStreamVendorConfig(cfg)
	return &AuditLogStreamConfig{Enabled: &enabled, StreamType: &streamType, VendorSpecific: &v}
}

// NewHecStreamConfig returns an AuditLogStreamConfig for an HTTPS Event Collector endpoint.
func NewHecStreamConfig(enabled bool, cfg *HecConfig) *AuditLogStreamConfig {
	streamType := "HTTPS Event Collector"
	v := AuditLogStreamVendorConfig(cfg)
	return &AuditLogStreamConfig{Enabled: &enabled, StreamType: &streamType, VendorSpecific: &v}
}

// NewGoogleCloudStreamConfig returns an AuditLogStreamConfig for Google Cloud Storage.
func NewGoogleCloudStreamConfig(enabled bool, cfg *GoogleCloudConfig) *AuditLogStreamConfig {
	streamType := "Google Cloud Storage"
	v := AuditLogStreamVendorConfig(cfg)
	return &AuditLogStreamConfig{Enabled: &enabled, StreamType: &streamType, VendorSpecific: &v}
}

// NewDatadogStreamConfig returns an AuditLogStreamConfig for Datadog.
func NewDatadogStreamConfig(enabled bool, cfg *DatadogConfig) *AuditLogStreamConfig {
	streamType := "Datadog"
	v := AuditLogStreamVendorConfig(cfg)
	return &AuditLogStreamConfig{Enabled: &enabled, StreamType: &streamType, VendorSpecific: &v}
}

// GetAuditLogStreamKey retrieves the public key used to encrypt secrets for audit log streaming.
// Credentials must be encrypted with this key before being submitted via CreateAuditLogStream
// or UpdateAuditLogStream.
//
// GitHub API docs: https://docs.github.com/enterprise-cloud@latest/rest/enterprise-admin/audit-log#get-the-audit-log-stream-key-for-encrypting-secrets
//
//meta:operation GET /enterprises/{enterprise}/audit-log/stream-key
func (s *EnterpriseService) GetAuditLogStreamKey(ctx context.Context, enterprise string) (*AuditLogStreamKey, *Response, error) {
	u := fmt.Sprintf("enterprises/%v/audit-log/stream-key", enterprise)

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	var key *AuditLogStreamKey
	resp, err := s.client.Do(ctx, req, &key)
	if err != nil {
		return nil, resp, err
	}

	return key, resp, nil
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

// GetAuditLogStream gets a single audit log stream configuration for an enterprise.
//
// GitHub API docs: https://docs.github.com/enterprise-cloud@latest/rest/enterprise-admin/audit-log#list-one-audit-log-streaming-configuration-via-a-stream-id
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

// CreateAuditLogStream creates an audit log streaming configuration for an enterprise.
// Credentials in the config must be encrypted using the key returned by GetAuditLogStreamKey.
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

// UpdateAuditLogStream updates an existing audit log stream configuration for an enterprise.
// Credentials in the config must be encrypted using the key returned by GetAuditLogStreamKey.
//
// GitHub API docs: https://docs.github.com/enterprise-cloud@latest/rest/enterprise-admin/audit-log#update-an-existing-audit-log-stream-configuration
//
//meta:operation PUT /enterprises/{enterprise}/audit-log/streams/{stream_id}
func (s *EnterpriseService) UpdateAuditLogStream(ctx context.Context, enterprise string, streamID int64, config *AuditLogStreamConfig) (*AuditLogStream, *Response, error) {
	u := fmt.Sprintf("enterprises/%v/audit-log/streams/%v", enterprise, streamID)

	req, err := s.client.NewRequest("PUT", u, config)
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

// DeleteAuditLogStream deletes an audit log stream configuration for an enterprise.
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
