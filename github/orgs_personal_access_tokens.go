package github

import (
	"context"
	"fmt"
	"net/http"
)

// Approves or denies a pending request to access organization resources via a fine-grained personal access token.
// Only GitHub Apps can call this API, using the `organization_personal_access_token_requests: write` permission.
// `action` can be one of `approve` or `deny`.
//
// GitHub API docs: https://docs.github.com/en/rest/orgs/personal-access-tokens?apiVersion=2022-11-28#review-a-request-to-access-organization-resources-with-a-fine-grained-personal-access-token
func (s *OrganizationsService) ReviewPersonalAccessTokenRequest(ctx context.Context, org, requestID, action, reason string) (*Response, error) {
	u := fmt.Sprintf("orgs/%v/personal-access-token-requests/%v", org, requestID)
	body := struct {
		Action string `json:"action"`
		Reason string `json:"reason,omitempty"`
	}{
		Action: action,
		Reason: reason,
	}

	req, err := s.client.NewRequest(http.MethodPost, u, &body)
	if err != nil {
		return nil, err
	}

	return s.client.Do(ctx, req, nil)
}
