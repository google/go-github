// Copyright 2014 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package github

import (
	"bytes"
)

// MarkdownOptions specifies optional parameters to the Markdown method.
type MarkdownOptions struct {
	// Mode identifies the rendering mode.  Possible values are:
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

	// Context identifies the repository context.  Only taken into account
	// when rendering as "gfm".
	Context string
}

type markdownRequest struct {
	Text    *string `json:"text,omitempty"`
	Mode    *string `json:"mode,omitempty"`
	Context *string `json:"context,omitempty"`
}

// Markdown renders an arbitrary Markdown document.
//
// GitHub API docs: https://developer.github.com/v3/markdown/
func (c *Client) Markdown(text string, opt *MarkdownOptions) (string, *Response, error) {
	request := &markdownRequest{Text: String(text)}
	if opt != nil {
		if opt.Mode != "" {
			request.Mode = String(opt.Mode)
		}
		if opt.Context != "" {
			request.Context = String(opt.Context)
		}
	}

	req, err := c.NewRequest("POST", "/markdown", request)
	if err != nil {
		return "", nil, err
	}

	buf := new(bytes.Buffer)
	resp, err := c.Do(req, buf)
	if err != nil {
		return "", resp, nil
	}

	return buf.String(), resp, nil
}
