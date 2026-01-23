package github

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestRepositoryRulesetRules_MarshalJSON_Empty(t *testing.T) {
	t.Parallel()

	got, err := json.Marshal(&RepositoryRulesetRules{})
	if err != nil {
		t.Fatalf("json.Marshal returned error: %v", err)
	}

	if string(got) != "[]" {
		t.Fatalf("expected empty array for no rules, got %s", got)
	}
}

func TestMarshalRepositoryRulesetRule_UpdateTypeValidation(t *testing.T) {
	t.Parallel()

	if _, err := marshalRepositoryRulesetRule(RulesetRuleTypeUpdate, &EmptyRuleParameters{}); err == nil {
		t.Fatal("expected type validation error, got nil")
	} else if !strings.Contains(err.Error(), "UpdateRuleParameters") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestMarshalRepositoryRulesetRule_UpdateNoParams(t *testing.T) {
	t.Parallel()

	bytes, err := marshalRepositoryRulesetRule(RulesetRuleTypeUpdate, (*UpdateRuleParameters)(nil))
	if err != nil {
		t.Fatalf("marshalRepositoryRulesetRule returned error: %v", err)
	}

	if string(bytes) != `{"type":"update"}` {
		t.Fatalf("marshalRepositoryRulesetRule expected type-only payload, got %s", bytes)
	}
}

func TestRepositoryRulesetRules_UnmarshalJSON(t *testing.T) {
	t.Parallel()

	payload := `[{"type":"creation"},{"type":"required_deployments","parameters":{"required_deployment_environments":["prod"]}}]`
	var rules RepositoryRulesetRules

	if err := json.Unmarshal([]byte(payload), &rules); err != nil {
		t.Fatalf("json.Unmarshal returned error: %v", err)
	}

	if rules.Creation == nil {
		t.Fatalf("Creation rule was not populated: %#v", rules.Creation)
	}

	if rules.RequiredDeployments == nil || len(rules.RequiredDeployments.RequiredDeploymentEnvironments) != 1 || rules.RequiredDeployments.RequiredDeploymentEnvironments[0] != "prod" {
		t.Fatalf("RequiredDeployments not populated as expected: %#v", rules.RequiredDeployments)
	}
}

func TestRepositoryRule_UnmarshalJSON(t *testing.T) {
	t.Parallel()

	var creation RepositoryRule
	if err := json.Unmarshal([]byte(`{"type":"creation"}`), &creation); err != nil {
		t.Fatalf("json.Unmarshal returned error: %v", err)
	}

	if creation.Type != RulesetRuleTypeCreation {
		t.Fatalf("Type mismatch, got %v", creation.Type)
	}

	if creation.Parameters != nil {
		t.Fatalf("creation rule should not carry parameters, got %#v", creation.Parameters)
	}

	var update RepositoryRule
	if err := json.Unmarshal([]byte(`{"type":"update","parameters":{"update_allows_fetch_and_merge":true}}`), &update); err != nil {
		t.Fatalf("json.Unmarshal returned error: %v", err)
	}

	params, ok := update.Parameters.(*UpdateRuleParameters)
	if !ok || params == nil {
		t.Fatalf("update parameters not decoded: %#v", update.Parameters)
	}

	if !params.UpdateAllowsFetchAndMerge {
		t.Fatalf("UpdateAllowsFetchAndMerge should be true, got %#v", params)
	}
}
