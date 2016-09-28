// Copyright 2016 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package github

// IntegrationsService provides access to the installation related functions
// in the GitHub API.
//
// GitHub API docs: https://developer.github.com/v3/integrations/
type IntegrationsService service

// ListInstallations lists the installations that the current integration has.
//
// GitHub API docs: https://developer.github.com/v3/integrations/#find-installations
func (s *IntegrationsService) ListInstallations(opt *ListOptions) ([]*Installation, *Response, error) {
	u, err := addOptions("integration/installations", opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	req.Header.Set("Accept", mediaTypeIntegrationPreview)

	i := new([]*Installation)
	resp, err := s.client.Do(req, &i)
	if err != nil {
		return nil, resp, err
	}

	return *i, resp, err
}
