package github

import (
	"context"
	"fmt"
)

type OrgsActionsPublicKey struct {
	KeyID *string `json:"key_id"`
	Key   *string `json:"key"`
}

func (s *OrganizationsService) GetPublicKey(ctx context.Context, owner string) (*OrgsActionsPublicKey, *Response, error) {
	u := fmt.Sprintf("orgs/%v/actions/secrets/public-key", owner)
	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	pubKey := new(OrgsActionsPublicKey)
	resp, err := s.client.Do(ctx, req, pubKey)
	if err != nil {
		return nil, resp, err
	}

	return pubKey, resp, nil
}