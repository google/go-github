package github

import (
	"context"
)

type EmojisService service

// ListEmojis returns the emojis available to use on GitHub.
//
// GitHub API docs: https://docs.github.com/rest/emojis/emojis#get-emojis
func (s *EmojisService) ListEmojis(ctx context.Context) (map[string]string, *Response, error) {
	req, err := s.client.NewRequest("GET", "emojis", nil)
	if err != nil {
		return nil, nil, err
	}

	var emoji map[string]string
	resp, err := s.client.Do(ctx, req, &emoji)
	if err != nil {
		return nil, resp, err
	}

	return emoji, resp, nil
}

// ListEmojis
// Deprecated: Use EmojisService.ListEmojis instead
func (c *Client) ListEmojis(ctx context.Context) (map[string]string, *Response, error) {
	return c.Emojis.ListEmojis(ctx)
}
