// Copyright 2025 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package github

import "testing"

//TODO: implement new SecretScanningPatternConfigs method tests
func TestSecretScanningService_ListPatternConfigsForEnterprise(t *testing.T) {
	t.Parallel()
	// client, mux, _ := setup(t)
	//TODO:
}

func TestSecretScanningService_ListPatternConfigsForOrg(t *testing.T) {
	t.Parallel()
	// client, mux, _ := setup(t)
	//TODO:
}

func TestSecretScanningService_UpdatePatternConfigsForEnterprise(t *testing.T) {
	t.Parallel()
	// client, mux, _ := setup(t)
	//TODO:
}

func TestSecretScanningService_UpdatePatternConfigsForOrg(t *testing.T) {
	t.Parallel()
	// client, mux, _ := setup(t)
	//TODO:
}

func TestSecretScanningPatternConfigs_Marshal(t *testing.T) {
	t.Parallel()
	testJSONMarshal(t, &SecretScanningPatternConfigs{}, `{}`)

	v := &SecretScanningPatternConfigs{
		PatternConfigVersion: Ptr("0ujsswThIGTUYm2K8FjOOfXtY1K"),
		ProviderPatternOverrides: []*SecretScanningPatternOverride{
			{
				TokenType:            Ptr("GITHUB_PERSONAL_ACCESS_TOKEN"),
				CustomPatternVersion: nil,
				Slug:                 Ptr("github_personal_access_token_legacy_v2"),
				DisplayName:          Ptr("GitHub Personal Access Token (Legacy v2)"),
				AlertTotal:           Ptr(15),
				AlertTotalPercentage: Ptr(36),
				FalsePositives:       Ptr(2),
				FalsePositiveRate:    Ptr(13),
				Bypassrate:           Ptr(13),
				DefaultSetting:       Ptr("enabled"),
				EnterpriseSetting:    Ptr("enabled"),
				Setting:              Ptr("enabled"),
			},
		},
		CustomPatternOverrides: []*SecretScanningPatternOverride{
			{
				TokenType:            Ptr("cp_2"),
				CustomPatternVersion: Ptr("0ujsswThIGTUYm2K8FjOOfXtY1K"),
				Slug:                 Ptr("custom-api-key"),
				DisplayName:          Ptr("Custom API Key"),
				AlertTotal:           Ptr(15),
				AlertTotalPercentage: Ptr(36),
				FalsePositives:       Ptr(3),
				FalsePositiveRate:    Ptr(20),
				Bypassrate:           Ptr(20),
				DefaultSetting:       Ptr("disabled"),
				EnterpriseSetting:    nil,
				Setting:              Ptr("enabled"),
			},
		},
	}

	want := `{
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
	}`

	testJSONMarshal(t, v, want)
}

func TestSecretScanningPatternOverride_Marshal(t *testing.T) {
	t.Parallel()
	testJSONMarshal(t, &SecretScanningPatternOverride{}, `{}`)

	v := &SecretScanningPatternOverride{
		TokenType:            Ptr("GITHUB_PERSONAL_ACCESS_TOKEN"),
		CustomPatternVersion: nil,
		Slug:                 Ptr("github_personal_access_token_legacy_v2"),
		DisplayName:          Ptr("GitHub Personal Access Token (Legacy v2)"),
		AlertTotal:           Ptr(15),
		AlertTotalPercentage: Ptr(36),
		FalsePositives:       Ptr(2),
		FalsePositiveRate:    Ptr(13),
		Bypassrate:           Ptr(13),
		DefaultSetting:       Ptr("enabled"),
		EnterpriseSetting:    Ptr("enabled"),
		Setting:              Ptr("enabled"),
	}

	want := `{
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
    }`

	testJSONMarshal(t, v, want)

	v = &SecretScanningPatternOverride{
		TokenType:            Ptr("cp_2"),
		CustomPatternVersion: Ptr("0ujsswThIGTUYm2K8FjOOfXtY1K"),
		Slug:                 Ptr("custom-api-key"),
		DisplayName:          Ptr("Custom API Key"),
		AlertTotal:           Ptr(15),
		AlertTotalPercentage: Ptr(36),
		FalsePositives:       Ptr(3),
		FalsePositiveRate:    Ptr(20),
		Bypassrate:           Ptr(20),
		DefaultSetting:       Ptr("disabled"),
		EnterpriseSetting:    nil,
		Setting:              Ptr("enabled"),
	}

	want = `{
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
    }`

	testJSONMarshal(t, v, want)
}

func TestSecretScanningPatternConfigsUpdate_Marshal(t *testing.T) {
	t.Parallel()
	testJSONMarshal(t, &SecretScanningPatternConfigsUpdate{}, `{}`)

	v := &SecretScanningPatternConfigsUpdate{
		PatternConfigVersion: Ptr("0ujsswThIGTUYm2K8FjOOfXtY1K"),
	}

	want := `{
		"pattern_config_version": "0ujsswThIGTUYm2K8FjOOfXtY1K"
	}`

	testJSONMarshal(t, v, want)
}

func TestSecretScanningPatternConfigsUpdateOptions_Marshal(t *testing.T) {
	t.Parallel()
	testJSONMarshal(t, &SecretScanningPatternConfigsUpdateOptions{}, `{}`)

	v := &SecretScanningPatternConfigsUpdateOptions{
		PatternConfigVersion: Ptr("0ujsswThIGTUYm2K8FjOOfXtY1K"),
		ProviderPatternSettings: []*SecretScanningProviderPatternSetting{
			{
				TokenType:             "GITHUB_PERSONAL_ACCESS_TOKEN",
				PushProtectionSetting: "enabled",
			},
		},
		CustomPatternSettings: []*SecretScanningCustomPatternSetting{
			{
				TokenType:             "cp_2",
				CustomPatternVersion:  Ptr("0ujsswThIGTUYm2K8FjOOfXtY1K"),
				PushProtectionSetting: "enabled",
			},
		},
	}

	want := `{
  		"pattern_config_version": "0ujsswThIGTUYm2K8FjOOfXtY1K",
  		"provider_pattern_settings": [
    		{
		      "token_type": "GITHUB_PERSONAL_ACCESS_TOKEN",
		      "push_protection_setting": "enabled"
		    }
		  ],
		  "custom_pattern_settings": [
		    {
		      "token_type": "cp_2",
		      "custom_pattern_version": "0ujsswThIGTUYm2K8FjOOfXtY1K",
		      "push_protection_setting": "enabled"
		    }
		  ]
		}`

	testJSONMarshal(t, v, want)
}

func TestSecretScanningProviderPatternSetting_Marshal(t *testing.T) {
	t.Parallel()
	testJSONMarshal(t, &SecretScanningProviderPatternSetting{}, `{}`)

	v := SecretScanningProviderPatternSetting{
		TokenType:             "GITHUB_PERSONAL_ACCESS_TOKEN",
		PushProtectionSetting: "enabled",
	}

	want := `{
    	"token_type": "GITHUB_PERSONAL_ACCESS_TOKEN",
    	"push_protection_setting": "enabled"
    }`

	testJSONMarshal(t, v, want)
}

func TestSecretScanninCustomPatternSetting_Marshal(t *testing.T) {
	t.Parallel()
	testJSONMarshal(t, &SecretScanningCustomPatternSetting{}, `{}`)

	v := SecretScanningCustomPatternSetting{
		TokenType:             "cp_2",
		CustomPatternVersion:  Ptr("0ujsswThIGTUYm2K8FjOOfXtY1K"),
		PushProtectionSetting: "enabled",
	}

	want := `{
    	"token_type": "cp_2",
    	"custom_pattern_version": "0ujsswThIGTUYm2K8FjOOfXtY1K",
    	"push_protection_setting": "enabled"
    }`

	testJSONMarshal(t, v, want)
}
