// Copyright 2025 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package github

import "testing"

//TODO: implement new SecretScanningPatternConfigs method tests
func TestSecretScanningService_ListPatternConfigsForEnterprise(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)
	//TODO:
}

func TestSecretScanningService_ListPatternConfigsForOrg(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)
	//TODO:
}

func TestSecretScanningService_UpdatePatternConfigsForEnterprise(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)
	//TODO:
}

func TestSecretScanningService_UpdatePatternConfigsForOrg(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)
	//TODO:
}

func TestSecretScanningPatternConfigs_Marshal(t *testing.T) {
	t.Parallel()
	testJSONMarshal(t, &SecretScanningPatternConfigs{}, `{}`)
	//TODO:
}

func TestSecretScanningPatternOverride_Marshal(t *testing.T) {
	t.Parallel()
	testJSONMarshal(t, &SecretScanningPatternOverride{}, `{}`)
	//TODO:
}

func TestSecretScanningPatternConfigsUpdate_Marshal(t *testing.T) {
	t.Parallel()
	testJSONMarshal(t, &SecretScanningPatternConfigsUpdate{}, `{}`)
	//TODO:
}

func TestSecretScanningPatternConfigsUpdateOptions_Marshal(t *testing.T) {
	t.Parallel()
	testJSONMarshal(t, &SecretScanningPatternConfigsUpdateOptions{}, `{}`)
	//TODO:
}

func TestSecretScanningProviderPatternSetting_Marshal(t *testing.T) {
	t.Parallel()
	testJSONMarshal(t, &SecretScanningProviderPatternSetting{}, `{}`)
	//TODO:
}
