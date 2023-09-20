// Copyright 2014 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package github

import (
	"bytes"
	"context"
	"fmt"
	"net/url"
)

// MarkdownOptions specifies optional parameters to the Markdown method.
type MarkdownOptions struct {
	// Mode identifies the rendering mode. Possible values are:
	//   markdown - render a document as plain Markdown, just like
	//   README files are rendered.
	//
	//   gfm - to render a document as user-content, e.g. like user
	//   comments or issues are rendered. In GFM mode, hard line breaks are
	//   always taken into account, and issue and user mentions are linked
	//   accordingly.
	//
	// Default is "markdown".
	Mode string

	// Context identifies the repository context. Only taken into account
	// when rendering as "gfm".
	Context string
}

type markdownRequest struct {
	Text    *string `json:"text,omitempty"`
	Mode    *string `json:"mode,omitempty"`
	Context *string `json:"context,omitempty"`
}

// MarkdownService provides access to markdown-related functions in the GitHub API.
type MarkdownService service

// Markdown renders an arbitrary Markdown document.
//
// GitHub API docs: https://docs.github.com/rest/markdown/markdown#render-a-markdown-document
func (s *MarkdownService) Markdown(ctx context.Context, text string, opts *MarkdownOptions) (string, *Response, error) {
	request := &markdownRequest{Text: String(text)}
	if opts != nil {
		if opts.Mode != "" {
			request.Mode = String(opts.Mode)
		}
		if opts.Context != "" {
			request.Context = String(opts.Context)
		}
	}

	req, err := s.client.NewRequest("POST", "markdown", request)
	if err != nil {
		return "", nil, err
	}

	buf := new(bytes.Buffer)
	resp, err := s.client.Do(ctx, req, buf)
	if err != nil {
		return "", resp, err
	}

	return buf.String(), resp, nil
}

// EmojisService provides access to emoji-related functions in the GitHub API.
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

// ListEmojis returns the emojis available to use on GitHub.
//
// Deprecated: Use EmojisService.ListEmojis instead
func (c *Client) ListEmojis(ctx context.Context) (map[string]string, *Response, error) {
	return c.Emojis.ListEmojis(ctx)
}

// CodeOfConduct represents a code of conduct.
type CodeOfConduct struct {
	Name *string `json:"name,omitempty"`
	Key  *string `json:"key,omitempty"`
	URL  *string `json:"url,omitempty"`
	Body *string `json:"body,omitempty"`
}

func (c *CodeOfConduct) String() string {
	return Stringify(c)
}

// CodesOfConductService provides access to code-of-conduct-related functions in the GitHub API.
type CodesOfConductService service

// ListCodesOfConduct returns all codes of conduct.
//
// GitHub API docs: https://docs.github.com/rest/codes-of-conduct/codes-of-conduct#get-all-codes-of-conduct
func (s *CodesOfConductService) ListCodesOfConduct(ctx context.Context) ([]*CodeOfConduct, *Response, error) {
	req, err := s.client.NewRequest("GET", "codes_of_conduct", nil)
	if err != nil {
		return nil, nil, err
	}

	// TODO: remove custom Accept header when this API fully launches.
	req.Header.Set("Accept", mediaTypeCodesOfConductPreview)

	var cs []*CodeOfConduct
	resp, err := s.client.Do(ctx, req, &cs)
	if err != nil {
		return nil, resp, err
	}

	return cs, resp, nil
}

// ListCodesOfConduct
// Deprecated: Use CodesOfConductService.ListCodesOfConduct instead
func (c *Client) ListCodesOfConduct(ctx context.Context) ([]*CodeOfConduct, *Response, error) {
	return c.CodesOfConduct.ListCodesOfConduct(ctx)
}

// GetCodeOfConduct returns an individual code of conduct.
//
// GitHub API docs: https://docs.github.com/rest/codes-of-conduct/codes-of-conduct#get-a-code-of-conduct
func (s *CodesOfConductService) GetCodeOfConduct(ctx context.Context, key string) (*CodeOfConduct, *Response, error) {
	u := fmt.Sprintf("codes_of_conduct/%s", key)
	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	// TODO: remove custom Accept header when this API fully launches.
	req.Header.Set("Accept", mediaTypeCodesOfConductPreview)

	coc := new(CodeOfConduct)
	resp, err := s.client.Do(ctx, req, coc)
	if err != nil {
		return nil, resp, err
	}

	return coc, resp, nil
}

// GetCodeOfConduct
// Deprecated: Use CodesOfConductService.GetCodeOfConduct instead
func (c *Client) GetCodeOfConduct(ctx context.Context, key string) (*CodeOfConduct, *Response, error) {
	return c.CodesOfConduct.GetCodeOfConduct(ctx, key)
}

// MetaService provides access to functions in the GitHub API that GitHub categorizes as "meta".
type MetaService service

// APIMeta represents metadata about the GitHub API.
type APIMeta struct {
	// An Array of IP addresses in CIDR format specifying the addresses
	// that incoming service hooks will originate from on GitHub.com.
	Hooks []string `json:"hooks,omitempty"`

	// An Array of IP addresses in CIDR format specifying the Git servers
	// for GitHub.com.
	Git []string `json:"git,omitempty"`

	// Whether authentication with username and password is supported.
	// (GitHub Enterprise instances using CAS or OAuth for authentication
	// will return false. Features like Basic Authentication with a
	// username and password, sudo mode, and two-factor authentication are
	// not supported on these servers.)
	VerifiablePasswordAuthentication *bool `json:"verifiable_password_authentication,omitempty"`

	// An array of IP addresses in CIDR format specifying the addresses
	// which serve GitHub Pages websites.
	Pages []string `json:"pages,omitempty"`

	// An Array of IP addresses specifying the addresses that source imports
	// will originate from on GitHub.com.
	Importer []string `json:"importer,omitempty"`

	// An array of IP addresses in CIDR format specifying the IP addresses
	// GitHub Actions will originate from.
	Actions []string `json:"actions,omitempty"`

	// An array of IP addresses in CIDR format specifying the IP addresses
	// Dependabot will originate from.
	Dependabot []string `json:"dependabot,omitempty"`

	// A map of algorithms to SSH key fingerprints.
	SSHKeyFingerprints map[string]string `json:"ssh_key_fingerprints,omitempty"`

	// An array of SSH keys.
	SSHKeys []string `json:"ssh_keys,omitempty"`

	// An array of IP addresses in CIDR format specifying the addresses
	// which serve GitHub websites.
	Web []string `json:"web,omitempty"`

	// An array of IP addresses in CIDR format specifying the addresses
	// which serve GitHub APIs.
	API []string `json:"api,omitempty"`
}

// APIMeta returns information about GitHub.com, the service. Or, if you access
// this endpoint on your organization’s GitHub Enterprise installation, this
// endpoint provides information about that installation.
//
// GitHub API docs: https://docs.github.com/rest/meta/meta#get-github-meta-information
func (s *MetaService) APIMeta(ctx context.Context) (*APIMeta, *Response, error) {
	req, err := s.client.NewRequest("GET", "meta", nil)
	if err != nil {
		return nil, nil, err
	}

	meta := new(APIMeta)
	resp, err := s.client.Do(ctx, req, meta)
	if err != nil {
		return nil, resp, err
	}

	return meta, resp, nil
}

// APIMeta
// Deprecated: Use MetaService.APIMeta instead.
func (c *Client) APIMeta(ctx context.Context) (*APIMeta, *Response, error) {
	return c.Meta.APIMeta(ctx)
}

// Octocat returns an ASCII art octocat with the specified message in a speech
// bubble. If message is empty, a random zen phrase is used.
//
// GitHub API docs: https://docs.github.com/rest/meta/meta#get-octocat
func (s *MetaService) Octocat(ctx context.Context, message string) (string, *Response, error) {
	u := "octocat"
	if message != "" {
		u = fmt.Sprintf("%s?s=%s", u, url.QueryEscape(message))
	}

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return "", nil, err
	}

	buf := new(bytes.Buffer)
	resp, err := s.client.Do(ctx, req, buf)
	if err != nil {
		return "", resp, err
	}

	return buf.String(), resp, nil
}

// Octocat
// Deprecated: Use MetaService.Octocat instead.
func (c *Client) Octocat(ctx context.Context, message string) (string, *Response, error) {
	return c.Meta.Octocat(ctx, message)
}

// Zen returns a random line from The Zen of GitHub.
//
// see also: http://warpspire.com/posts/taste/
//
// GitHub API docs: https://docs.github.com/rest/meta/meta#get-the-zen-of-github
func (s *MetaService) Zen(ctx context.Context) (string, *Response, error) {
	req, err := s.client.NewRequest("GET", "zen", nil)
	if err != nil {
		return "", nil, err
	}

	buf := new(bytes.Buffer)
	resp, err := s.client.Do(ctx, req, buf)
	if err != nil {
		return "", resp, err
	}

	return buf.String(), resp, nil
}

// Zen
// Deprecated: Use MetaService.Zen instead.
func (c *Client) Zen(ctx context.Context) (string, *Response, error) {
	return c.Meta.Zen(ctx)
}
