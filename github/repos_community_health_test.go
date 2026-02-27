// Copyright 2017 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package github

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
)

func TestRepositoriesService_GetCommunityHealthMetrics(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/repos/o/r/community/profile", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{
				"health_percentage": 100,
				"description": "My first repository on GitHub!",
				"documentation": null,
				"files": {
					"code_of_conduct": {
						"name": "Contributor Covenant",
						"key": "contributor_covenant",
						"url": null,
						"html_url": "https://github.com/octocat/Hello-World/blob/master/CODE_OF_CONDUCT.md"
					},
					"code_of_conduct_file": {
						"url": "https://api.github.com/repos/octocat/Hello-World/contents/CODE_OF_CONDUCT.md",
						"html_url": "https://github.com/octocat/Hello-World/blob/master/CODE_OF_CONDUCT.md"
					},
					"contributing": {
						"url": "https://api.github.com/repos/octocat/Hello-World/contents/CONTRIBUTING",
						"html_url": "https://github.com/octocat/Hello-World/blob/master/CONTRIBUTING"
					},
					"issue_template": {
						"url": "https://api.github.com/repos/octocat/Hello-World/contents/ISSUE_TEMPLATE",
						"html_url": "https://github.com/octocat/Hello-World/blob/master/ISSUE_TEMPLATE"
					},
					"pull_request_template": {
						"url": "https://api.github.com/repos/octocat/Hello-World/contents/PULL_REQUEST_TEMPLATE",
						"html_url": "https://github.com/octocat/Hello-World/blob/master/PULL_REQUEST_TEMPLATE"
					},
					"license": {
						"name": "MIT License",
						"key": "mit",
						"spdx_id": "MIT",
						"url": "https://api.github.com/licenses/mit",
						"html_url": "https://github.com/octocat/Hello-World/blob/master/LICENSE",
						"node_id": "MDc6TGljZW5zZW1pdA=="
					},
					"readme": {
						"url": "https://api.github.com/repos/octocat/Hello-World/contents/README.md",
						"html_url": "https://github.com/octocat/Hello-World/blob/master/README.md"
					}
				},
				"updated_at": "2017-02-28T00:00:00Z",
				"content_reports_enabled": true
			}`)
	})

	ctx := t.Context()
	got, _, err := client.Repositories.GetCommunityHealthMetrics(ctx, "o", "r")
	if err != nil {
		t.Errorf("Repositories.GetCommunityHealthMetrics returned error: %v", err)
	}

	updatedAt := time.Date(2017, time.February, 28, 0, 0, 0, 0, time.UTC)
	want := &CommunityHealthMetrics{
		HealthPercentage:      Ptr(100),
		Description:           Ptr("My first repository on GitHub!"),
		UpdatedAt:             &Timestamp{updatedAt},
		ContentReportsEnabled: Ptr(true),
		Files: &CommunityHealthFiles{
			CodeOfConduct: &Metric{
				Name:    Ptr("Contributor Covenant"),
				Key:     Ptr("contributor_covenant"),
				HTMLURL: Ptr("https://github.com/octocat/Hello-World/blob/master/CODE_OF_CONDUCT.md"),
			},
			CodeOfConductFile: &Metric{
				URL:     Ptr("https://api.github.com/repos/octocat/Hello-World/contents/CODE_OF_CONDUCT.md"),
				HTMLURL: Ptr("https://github.com/octocat/Hello-World/blob/master/CODE_OF_CONDUCT.md"),
			},
			Contributing: &Metric{
				URL:     Ptr("https://api.github.com/repos/octocat/Hello-World/contents/CONTRIBUTING"),
				HTMLURL: Ptr("https://github.com/octocat/Hello-World/blob/master/CONTRIBUTING"),
			},
			IssueTemplate: &Metric{
				URL:     Ptr("https://api.github.com/repos/octocat/Hello-World/contents/ISSUE_TEMPLATE"),
				HTMLURL: Ptr("https://github.com/octocat/Hello-World/blob/master/ISSUE_TEMPLATE"),
			},
			PullRequestTemplate: &Metric{
				URL:     Ptr("https://api.github.com/repos/octocat/Hello-World/contents/PULL_REQUEST_TEMPLATE"),
				HTMLURL: Ptr("https://github.com/octocat/Hello-World/blob/master/PULL_REQUEST_TEMPLATE"),
			},
			License: &Metric{
				Name:    Ptr("MIT License"),
				Key:     Ptr("mit"),
				SPDXID:  Ptr("MIT"),
				URL:     Ptr("https://api.github.com/licenses/mit"),
				HTMLURL: Ptr("https://github.com/octocat/Hello-World/blob/master/LICENSE"),
				NodeID:  Ptr("MDc6TGljZW5zZW1pdA=="),
			},
			Readme: &Metric{
				URL:     Ptr("https://api.github.com/repos/octocat/Hello-World/contents/README.md"),
				HTMLURL: Ptr("https://github.com/octocat/Hello-World/blob/master/README.md"),
			},
		},
	}
	if !cmp.Equal(got, want) {
		t.Errorf("Repositories.GetCommunityHealthMetrics:\ngot:\n%v\nwant:\n%v", Stringify(got), Stringify(want))
	}

	const methodName = "GetCommunityHealthMetrics"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Repositories.GetCommunityHealthMetrics(ctx, "\n", "\n")
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Repositories.GetCommunityHealthMetrics(ctx, "o", "r")
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestMetric_Marshal(t *testing.T) {
	t.Parallel()
	testJSONMarshal(t, &Metric{}, `{
		"name": null,
		"key": null,
		"spdx_id": null,
		"url": null,
		"html_url": null,
		"node_id": null
	}`)

	r := &Metric{
		Name:    Ptr("name"),
		Key:     Ptr("key"),
		SPDXID:  Ptr("spdx_id"),
		URL:     Ptr("url"),
		HTMLURL: Ptr("hurl"),
		NodeID:  Ptr("node_id"),
	}

	want := `{
		"name": "name",
		"key": "key",
		"spdx_id": "spdx_id",
		"url": "url",
		"html_url": "hurl",
		"node_id": "node_id"
	}`

	testJSONMarshal(t, r, want)
}

func TestCommunityHealthFiles_Marshal(t *testing.T) {
	t.Parallel()
	testJSONMarshal(t, &CommunityHealthFiles{}, `{
		"code_of_conduct": null,
		"code_of_conduct_file": null,
		"contributing": null,
		"issue_template": null,
		"pull_request_template": null,
		"license": null,
		"readme": null
	}`)

	r := &CommunityHealthFiles{
		CodeOfConduct: &Metric{
			Name:    Ptr("name"),
			Key:     Ptr("key"),
			URL:     Ptr("url"),
			HTMLURL: Ptr("hurl"),
		},
		CodeOfConductFile: &Metric{
			Name:    Ptr("name"),
			Key:     Ptr("key"),
			URL:     Ptr("url"),
			HTMLURL: Ptr("hurl"),
		},
		Contributing: &Metric{
			Name:    Ptr("name"),
			Key:     Ptr("key"),
			URL:     Ptr("url"),
			HTMLURL: Ptr("hurl"),
		},
		IssueTemplate: &Metric{
			Name:    Ptr("name"),
			Key:     Ptr("key"),
			URL:     Ptr("url"),
			HTMLURL: Ptr("hurl"),
		},
		PullRequestTemplate: &Metric{
			Name:    Ptr("name"),
			Key:     Ptr("key"),
			URL:     Ptr("url"),
			HTMLURL: Ptr("hurl"),
		},
		License: &Metric{
			Name:    Ptr("name"),
			Key:     Ptr("key"),
			SPDXID:  Ptr("spdx_id"),
			URL:     Ptr("url"),
			HTMLURL: Ptr("hurl"),
			NodeID:  Ptr("node_id"),
		},
		Readme: &Metric{
			Name:    Ptr("name"),
			Key:     Ptr("key"),
			URL:     Ptr("url"),
			HTMLURL: Ptr("hurl"),
		},
	}

	want := `{
		"code_of_conduct": {
			"name": "name",
			"key": "key",
			"url": "url",
			"html_url": "hurl",
			"node_id": null,
			"spdx_id": null
		},
		"code_of_conduct_file": {
			"name": "name",
			"key": "key",
			"url": "url",
			"html_url": "hurl",
			"node_id": null,
			"spdx_id": null
		},
		"contributing": {
			"name": "name",
			"key": "key",
			"url": "url",
			"html_url": "hurl",
			"node_id": null,
			"spdx_id": null
		},
		"issue_template": {
			"name": "name",
			"key": "key",
			"url": "url",
			"html_url": "hurl",
			"node_id": null,
			"spdx_id": null
		},
		"pull_request_template": {
			"name": "name",
			"key": "key",
			"url": "url",
			"html_url": "hurl",
			"node_id": null,
			"spdx_id": null
		},
		"license": {
			"name": "name",
			"key": "key",
			"spdx_id": "spdx_id",
			"url": "url",
			"html_url": "hurl",
			"node_id": "node_id"
		},
		"readme": {
			"name": "name",
			"key": "key",
			"url": "url",
			"html_url": "hurl",
			"node_id": null,
			"spdx_id": null
		}
	}`

	testJSONMarshal(t, r, want)
}

func TestCommunityHealthMetrics_Marshal(t *testing.T) {
	t.Parallel()
	testJSONMarshal(t, &CommunityHealthMetrics{}, `{
		"health_percentage": null,
		"description": null,
		"documentation": null,
		"files": null,
		"updated_at": null,
		"content_reports_enabled": null
	}`)

	r := &CommunityHealthMetrics{
		HealthPercentage: Ptr(1),
		Description:      Ptr("desc"),
		Documentation:    Ptr("docs"),
		Files: &CommunityHealthFiles{
			CodeOfConduct: &Metric{
				Name:    Ptr("name"),
				Key:     Ptr("key"),
				URL:     Ptr("url"),
				HTMLURL: Ptr("hurl"),
			},
			CodeOfConductFile: &Metric{
				Name:    Ptr("name"),
				Key:     Ptr("key"),
				URL:     Ptr("url"),
				HTMLURL: Ptr("hurl"),
			},
			Contributing: &Metric{
				Name:    Ptr("name"),
				Key:     Ptr("key"),
				URL:     Ptr("url"),
				HTMLURL: Ptr("hurl"),
			},
			IssueTemplate: &Metric{
				Name:    Ptr("name"),
				Key:     Ptr("key"),
				URL:     Ptr("url"),
				HTMLURL: Ptr("hurl"),
			},
			PullRequestTemplate: &Metric{
				Name:    Ptr("name"),
				Key:     Ptr("key"),
				URL:     Ptr("url"),
				HTMLURL: Ptr("hurl"),
			},
			License: &Metric{
				Name:    Ptr("name"),
				Key:     Ptr("key"),
				SPDXID:  Ptr("spdx_id"),
				URL:     Ptr("url"),
				HTMLURL: Ptr("hurl"),
				NodeID:  Ptr("node_id"),
			},
			Readme: &Metric{
				Name:    Ptr("name"),
				Key:     Ptr("key"),
				URL:     Ptr("url"),
				HTMLURL: Ptr("hurl"),
			},
		},
		UpdatedAt:             &Timestamp{referenceTime},
		ContentReportsEnabled: Ptr(true),
	}

	want := `{
		"health_percentage": 1,
		"description": "desc",
		"documentation": "docs",
		"files": {
			"code_of_conduct": {
				"name": "name",
				"key": "key",
				"url": "url",
				"html_url": "hurl",
				"node_id": null,
				"spdx_id": null
			},
			"code_of_conduct_file": {
				"name": "name",
				"key": "key",
				"url": "url",
				"html_url": "hurl",
				"node_id": null,
				"spdx_id": null
			},
			"contributing": {
				"name": "name",
				"key": "key",
				"url": "url",
				"html_url": "hurl",
				"node_id": null,
				"spdx_id": null
			},
			"issue_template": {
				"name": "name",
				"key": "key",
				"url": "url",
				"html_url": "hurl",
				"node_id": null,
				"spdx_id": null
			},
			"pull_request_template": {
				"name": "name",
				"key": "key",
				"url": "url",
				"html_url": "hurl",
				"node_id": null,
				"spdx_id": null
			},
			"license": {
				"name": "name",
				"key": "key",
				"spdx_id": "spdx_id",
				"url": "url",
				"html_url": "hurl",
				"node_id": "node_id"
			},
			"readme": {
				"name": "name",
				"key": "key",
				"url": "url",
				"html_url": "hurl",
				"node_id": null,
				"spdx_id": null
			}
		},
		"updated_at": ` + referenceTimeStr + `,
		"content_reports_enabled": true
	}`

	testJSONMarshal(t, r, want)
}
