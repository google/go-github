package github

import (
	"bytes"
	"context"
)

type MarkdownService service

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
