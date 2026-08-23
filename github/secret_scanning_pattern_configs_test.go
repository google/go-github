// Copyright 2025 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package github

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestSecretScanningService_ListPatternConfigsForEnterprise(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/enterprises/e/secret-scanning/pattern-configurations", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")

		fmt.Fprint(w, `{
			"pattern_config_version": "0ujsswThIGTUYm2K8FjOOfXtY1K",
			"provider_pattern_overrides": [
			  {
			    "token_type": "GITHUB_PERSONAL_ACCESS_TOKEN",
			    "slug": "github_personal_access_token_legacy_v2",
			    "display_name": "GitHub Personal Access Token (Legacy v2)",
			    "alert_total": 15,
			    "alert_total_percentage": 36,
			    "false_positives": 2,
			    "false_positive_rate": 13,
			    "bypass_rate": 13,
			    "default_setting": "enabled",
			    "setting": "enabled",
			    "enterprise_setting": "enabled"
			  }
			],
			"custom_pattern_overrides": [
			  {
			    "token_type": "cp_2",
			    "custom_pattern_version": "0ujsswThIGTUYm2K8FjOOfXtY1K",
			    "slug": "custom-api-key",
			    "display_name": "Custom API Key",
			    "alert_total": 15,
			    "alert_total_percentage": 36,
			    "false_positives": 3,
			    "false_positive_rate": 20,
			    "bypass_rate": 20,
			    "default_setting": "disabled",
			    "setting": "enabled"
			  }
			]
		}`)
	})

	ctx := t.Context()

	patternConfigs, _, err := client.SecretScanning.ListPatternConfigsForEnterprise(ctx, "e")
	if err != nil {
		t.Errorf("SecretScanning.ListPatternConfigsForEnterprise returned error: %v", err)
	}

	want := &SecretScanningPatternConfigs{
		PatternConfigVersion: new("0ujsswThIGTUYm2K8FjOOfXtY1K"),
		ProviderPatternOverrides: []*SecretScanningPatternOverride{
			{
				TokenType:            new("GITHUB_PERSONAL_ACCESS_TOKEN"),
				CustomPatternVersion: nil,
				Slug:                 new("github_personal_access_token_legacy_v2"),
				DisplayName:          new("GitHub Personal Access Token (Legacy v2)"),
				AlertTotal:           new(15),
				AlertTotalPercentage: new(36),
				FalsePositives:       new(2),
				FalsePositiveRate:    new(13),
				Bypassrate:           new(13),
				DefaultSetting:       new("enabled"),
				EnterpriseSetting:    new("enabled"),
				Setting:              new("enabled"),
			},
		},
		CustomPatternOverrides: []*SecretScanningPatternOverride{
			{
				TokenType:            new("cp_2"),
				CustomPatternVersion: new("0ujsswThIGTUYm2K8FjOOfXtY1K"),
				Slug:                 new("custom-api-key"),
				DisplayName:          new("Custom API Key"),
				AlertTotal:           new(15),
				AlertTotalPercentage: new(36),
				FalsePositives:       new(3),
				FalsePositiveRate:    new(20),
				Bypassrate:           new(20),
				DefaultSetting:       new("disabled"),
				EnterpriseSetting:    nil,
				Setting:              new("enabled"),
			},
		},
	}

	if !cmp.Equal(patternConfigs, want) {
		t.Errorf("SecretScanning.ListPatternConfigsForEnterprise returned %+v, want %+v", patternConfigs, want)
	}

	const methodName = "ListPatternConfigsForEnterprise"

	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.SecretScanning.ListPatternConfigsForEnterprise(ctx, "\n")
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		_, resp, err := client.SecretScanning.ListPatternConfigsForEnterprise(ctx, "e")
		return resp, err
	})
}

func TestSecretScanningService_ListPatternConfigsForOrg(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/orgs/o/secret-scanning/pattern-configurations", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")

		fmt.Fprint(w, `{
			"pattern_config_version": "0ujsswThIGTUYm2K8FjOOfXtY1K",
			"provider_pattern_overrides": [
			  {
			    "token_type": "GITHUB_PERSONAL_ACCESS_TOKEN",
			    "slug": "github_personal_access_token_legacy_v2",
			    "display_name": "GitHub Personal Access Token (Legacy v2)",
			    "alert_total": 15,
			    "alert_total_percentage": 36,
			    "false_positives": 2,
			    "false_positive_rate": 13,
			    "bypass_rate": 13,
			    "default_setting": "enabled",
			    "setting": "enabled",
			    "enterprise_setting": "enabled"
			  }
			],
			"custom_pattern_overrides": [
			  {
			    "token_type": "cp_2",
			    "custom_pattern_version": "0ujsswThIGTUYm2K8FjOOfXtY1K",
			    "slug": "custom-api-key",
			    "display_name": "Custom API Key",
			    "alert_total": 15,
			    "alert_total_percentage": 36,
			    "false_positives": 3,
			    "false_positive_rate": 20,
			    "bypass_rate": 20,
			    "default_setting": "disabled",
			    "setting": "enabled"
			  }
			]
		}`)
	})
	ctx := t.Context()

	patternConfigs, _, err := client.SecretScanning.ListPatternConfigsForOrg(ctx, "o")
	if err != nil {
		t.Errorf("SecretScanning.ListPatternConfigsForOrg returned error: %v", err)
	}

	want := &SecretScanningPatternConfigs{
		PatternConfigVersion: new("0ujsswThIGTUYm2K8FjOOfXtY1K"),
		ProviderPatternOverrides: []*SecretScanningPatternOverride{
			{
				TokenType:            new("GITHUB_PERSONAL_ACCESS_TOKEN"),
				CustomPatternVersion: nil,
				Slug:                 new("github_personal_access_token_legacy_v2"),
				DisplayName:          new("GitHub Personal Access Token (Legacy v2)"),
				AlertTotal:           new(15),
				AlertTotalPercentage: new(36),
				FalsePositives:       new(2),
				FalsePositiveRate:    new(13),
				Bypassrate:           new(13),
				DefaultSetting:       new("enabled"),
				EnterpriseSetting:    new("enabled"),
				Setting:              new("enabled"),
			},
		},
		CustomPatternOverrides: []*SecretScanningPatternOverride{
			{
				TokenType:            new("cp_2"),
				CustomPatternVersion: new("0ujsswThIGTUYm2K8FjOOfXtY1K"),
				Slug:                 new("custom-api-key"),
				DisplayName:          new("Custom API Key"),
				AlertTotal:           new(15),
				AlertTotalPercentage: new(36),
				FalsePositives:       new(3),
				FalsePositiveRate:    new(20),
				Bypassrate:           new(20),
				DefaultSetting:       new("disabled"),
				EnterpriseSetting:    nil,
				Setting:              new("enabled"),
			},
		},
	}

	if !cmp.Equal(patternConfigs, want) {
		t.Errorf("SecretScanning.ListPatternConfigsForOrg returned %+v, want %+v", patternConfigs, want)
	}

	const methodName = "ListPatternConfigsForOrg"

	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.SecretScanning.ListPatternConfigsForOrg(ctx, "\n")
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		_, resp, err := client.SecretScanning.ListPatternConfigsForOrg(ctx, "o")
		return resp, err
	})
}

func TestSecretScanningService_UpdatePatternConfigsForEnterprise(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/enterprises/e/secret-scanning/pattern-configurations", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PATCH")

		fmt.Fprint(w, `{
			"pattern_config_version": "0ujsswThIGTUYm2K8FjOOfXtY1K"
		}`)
	})

	ctx := t.Context()

	opts := &SecretScanningPatternConfigsUpdateOptions{
		PatternConfigVersion: new("0ujsswThIGTUYm2K8FjOOfXtY1K"),
		ProviderPatternSettings: []*SecretScanningProviderPatternSetting{
			{
				TokenType:             "GITHUB_PERSONAL_ACCESS_TOKEN",
				PushProtectionSetting: "enabled",
			},
		},
		CustomPatternSettings: []*SecretScanningCustomPatternSetting{
			{
				TokenType:             "cp_2",
				CustomPatternVersion:  new("0ujsswThIGTUYm2K8FjOOfXtY1K"),
				PushProtectionSetting: "enabled",
			},
		},
	}

	configsUpdate, _, err := client.SecretScanning.UpdatePatternConfigsForEnterprise(ctx, "e", opts)
	if err != nil {
		t.Errorf("SecretScanning.UpdatePatternConfigsForEnterprise returned error: %v", err)
	}

	want := &SecretScanningPatternConfigsUpdate{
		PatternConfigVersion: new("0ujsswThIGTUYm2K8FjOOfXtY1K"),
	}

	if !cmp.Equal(configsUpdate, want) {
		t.Errorf("SecretScanning.UpdatePatternConfigsForEnterprise returned %+v, want %+v", configsUpdate, want)
	}

	const methodName = "UpdatePatternConfigsForEnterprise"

	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.SecretScanning.UpdatePatternConfigsForEnterprise(ctx, "\n", opts)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		_, resp, err := client.SecretScanning.UpdatePatternConfigsForEnterprise(ctx, "o", opts)
		return resp, err
	})
}

func TestSecretScanningService_UpdatePatternConfigsForOrg(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/orgs/o/secret-scanning/pattern-configurations", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PATCH")

		fmt.Fprint(w, `{
			"pattern_config_version": "0ujsswThIGTUYm2K8FjOOfXtY1K"
		}`)
	})

	ctx := t.Context()

	opts := &SecretScanningPatternConfigsUpdateOptions{
		PatternConfigVersion: new("0ujsswThIGTUYm2K8FjOOfXtY1K"),
		ProviderPatternSettings: []*SecretScanningProviderPatternSetting{
			{
				TokenType:             "GITHUB_PERSONAL_ACCESS_TOKEN",
				PushProtectionSetting: "enabled",
			},
		},
		CustomPatternSettings: []*SecretScanningCustomPatternSetting{
			{
				TokenType:             "cp_2",
				CustomPatternVersion:  new("0ujsswThIGTUYm2K8FjOOfXtY1K"),
				PushProtectionSetting: "enabled",
			},
		},
	}

	configsUpdate, _, err := client.SecretScanning.UpdatePatternConfigsForOrg(ctx, "o", opts)
	if err != nil {
		t.Errorf("SecretScanning.UpdatePatternConfigsForOrg returned err: %v", err)
	}

	want := &SecretScanningPatternConfigsUpdate{
		PatternConfigVersion: new("0ujsswThIGTUYm2K8FjOOfXtY1K"),
	}

	if !cmp.Equal(configsUpdate, want) {
		t.Errorf("SecretScanning.UpdatePatternConfigsForOrg returned %+v, want %+v", configsUpdate, want)
	}

	const methodName = "UpdatePatternConfigsForOrg"

	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.SecretScanning.UpdatePatternConfigsForOrg(ctx, "\n", opts)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		_, resp, err := client.SecretScanning.UpdatePatternConfigsForOrg(ctx, "o", opts)
		return resp, err
	})
}
