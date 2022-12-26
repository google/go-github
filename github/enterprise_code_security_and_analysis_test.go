// Copyright 2022 The go-github AUTHORS. All rights reserved.
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

	"github.com/google/go-cmp/cmp"
)

func TestEnterpriseService_GetCodeSecurityAndAnalysis(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/enterprises/e/code_security_and_analysis", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")

		fmt.Fprint(w, `
		{
		  "advanced_security_enabled_for_new_repositories": true,
		  "secret_scanning_enabled_for_new_repositories": true,
		  "secret_scanning_push_protection_enabled_for_new_repositories": true,
		  "secret_scanning_push_protection_custom_link": "https://github.com/test-org/test-repo/blob/main/README.md"
		}`)
	})

	ctx := context.Background()

	const methodName = "GetCodeSecurityAndAnalysis"

	enterpriseSecurityAnalysisSettings, _, err := client.Enterprise.GetCodeSecurityAndAnalysis(ctx, "e")
	if err != nil {
		t.Errorf("Enterprise.%v returned error: %v", methodName, err)
	}
	want := &EnterpriseSecurityAnalysisSettings{
		AdvancedSecurityEnabledForNewRepositories:             true,
		SecretScanningEnabledForNewRepositories:               true,
		SecretScanningPushProtectionEnabledForNewRepositories: true,
		SecretScanningPushProtectionCustomLink:                "https://github.com/test-org/test-repo/blob/main/README.md",
	}

	if !cmp.Equal(enterpriseSecurityAnalysisSettings, want) {
		t.Errorf("Enterprise.%v return \ngot: %+v,\nwant:%+v", methodName, enterpriseSecurityAnalysisSettings, want)
	}

	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Enterprise.GetCodeSecurityAndAnalysis(ctx, "o")
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Enterprise.GetCodeSecurityAndAnalysis(ctx, "e")
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestEnterpriseService_UpdateCodeSecurityAndAnalysis(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	input := &EnterpriseSecurityAnalysisSettings{
		AdvancedSecurityEnabledForNewRepositories:             true,
		SecretScanningEnabledForNewRepositories:               true,
		SecretScanningPushProtectionEnabledForNewRepositories: true,
		SecretScanningPushProtectionCustomLink:                "https://github.com/test-org/test-repo/blob/main/README.md",
	}

	mux.HandleFunc("/enterprises/e/code_security_and_analysis", func(w http.ResponseWriter, r *http.Request) {
		v := new(EnterpriseSecurityAnalysisSettings)
		json.NewDecoder(r.Body).Decode(v)

		testMethod(t, r, "PATCH")
		if !cmp.Equal(v, input) {
			t.Errorf("Request body = %+v, want %+v", v, input)
		}
	})

	ctx := context.Background()

	const methodName = "UpdateCodeSecurityAndAnalysis"

	_, err := client.Enterprise.UpdateCodeSecurityAndAnalysis(ctx, "e", input)
	if err != nil {
		t.Errorf("Enterprise.%v returned error: %v", methodName, err)
	}

	testBadOptions(t, methodName, func() (err error) {
		_, err = client.Enterprise.UpdateCodeSecurityAndAnalysis(ctx, "o", input)
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		return client.Enterprise.UpdateCodeSecurityAndAnalysis(ctx, "e", input)
	})
}

func TestEnterpriseService_EnableAdvancedSecurity(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/enterprises/e/advanced_security/enable_all", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
	})

	ctx := context.Background()

	const methodName = "EnableAdvancedSecurity"

	_, err := client.Enterprise.EnableAdvancedSecurity(ctx, "e")
	if err != nil {
		t.Errorf("Enterprise.%v returned error: %v", methodName, err)
	}

	testBadOptions(t, methodName, func() (err error) {
		_, err = client.Enterprise.EnableAdvancedSecurity(ctx, "o")
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		return client.Enterprise.EnableAdvancedSecurity(ctx, "e")
	})
}

func TestEnterpriseService_DisableAdvancedSecurity(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/enterprises/e/advanced_security/disable_all", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
	})

	ctx := context.Background()

	const methodName = "DisableAdvancedSecurity"

	_, err := client.Enterprise.DisableAdvancedSecurity(ctx, "e")
	if err != nil {
		t.Errorf("Enterprise.%v returned error: %v", methodName, err)
	}

	testBadOptions(t, methodName, func() (err error) {
		_, err = client.Enterprise.DisableAdvancedSecurity(ctx, "o")
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		return client.Enterprise.DisableAdvancedSecurity(ctx, "e")
	})
}

func TestEnterpriseService_EnableSecretScanning(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/enterprises/e/secret_scanning/enable_all", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
	})

	ctx := context.Background()

	const methodName = "EnableSecretScanning"

	_, err := client.Enterprise.EnableSecretScanning(ctx, "e")
	if err != nil {
		t.Errorf("Enterprise.%v returned error: %v", methodName, err)
	}

	testBadOptions(t, methodName, func() (err error) {
		_, err = client.Enterprise.EnableSecretScanning(ctx, "o")
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		return client.Enterprise.EnableSecretScanning(ctx, "e")
	})
}

func TestEnterpriseService_DisableSecretScanning(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/enterprises/e/secret_scanning/disable_all", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
	})

	ctx := context.Background()

	const methodName = "DisableSecretScanning"

	_, err := client.Enterprise.DisableSecretScanning(ctx, "e")
	if err != nil {
		t.Errorf("Enterprise.%v returned error: %v", methodName, err)
	}

	testBadOptions(t, methodName, func() (err error) {
		_, err = client.Enterprise.DisableSecretScanning(ctx, "o")
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		return client.Enterprise.DisableSecretScanning(ctx, "e")
	})
}

func TestEnterpriseService_EnableSecretScanningPushProtection(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/enterprises/e/secret_scanning_push_protection/enable_all", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
	})

	ctx := context.Background()

	const methodName = "EnableSecretScanningPushProtection"

	_, err := client.Enterprise.EnableSecretScanningPushProtection(ctx, "e")
	if err != nil {
		t.Errorf("Enterprise.%v returned error: %v", methodName, err)
	}

	testBadOptions(t, methodName, func() (err error) {
		_, err = client.Enterprise.EnableSecretScanningPushProtection(ctx, "o")
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		return client.Enterprise.EnableSecretScanningPushProtection(ctx, "e")
	})
}

func TestEnterpriseService_DisableSecretScanningPushProtection(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/enterprises/e/secret_scanning_push_protection/disable_all", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
	})

	ctx := context.Background()

	const methodName = "DisableSecretScanningPushProtection"

	_, err := client.Enterprise.DisableSecretScanningPushProtection(ctx, "e")
	if err != nil {
		t.Errorf("Enterprise.%v returned error: %v", methodName, err)
	}

	testBadOptions(t, methodName, func() (err error) {
		_, err = client.Enterprise.DisableSecretScanningPushProtection(ctx, "o")
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		return client.Enterprise.DisableSecretScanningPushProtection(ctx, "e")
	})
}
