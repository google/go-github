// Copyright 2017 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package github

import (
	"context"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
)

func TestRepositoriesService_GetCommunityHealthMetrics(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/repos/o/r/community/profile", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprintf(w, `{
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

	ctx := context.Background()
	got, _, err := client.Repositories.GetCommunityHealthMetrics(ctx, "o", "r")
	if err != nil {
		t.Errorf("Repositories.GetCommunityHealthMetrics returned error: %v", err)
	}

	updatedAt := time.Date(2017, time.February, 28, 0, 0, 0, 0, time.UTC)
	want := &CommunityHealthMetrics{
		HealthPercentage:      Int(100),
		Description:           String("My first repository on GitHub!"),
		UpdatedAt:             &Timestamp{updatedAt},
		ContentReportsEnabled: Bool(true),
		Files: &CommunityHealthFiles{
			CodeOfConduct: &Metric{
				Name:    String("Contributor Covenant"),
				Key:     String("contributor_covenant"),
				HTMLURL: String("https://github.com/octocat/Hello-World/blob/master/CODE_OF_CONDUCT.md"),
			},
			CodeOfConductFile: &Metric{
				URL:     String("https://api.github.com/repos/octocat/Hello-World/contents/CODE_OF_CONDUCT.md"),
				HTMLURL: String("https://github.com/octocat/Hello-World/blob/master/CODE_OF_CONDUCT.md"),
			},
			Contributing: &Metric{
				URL:     String("https://api.github.com/repos/octocat/Hello-World/contents/CONTRIBUTING"),
				HTMLURL: String("https://github.com/octocat/Hello-World/blob/master/CONTRIBUTING"),
			},
			IssueTemplate: &Metric{
				URL:     String("https://api.github.com/repos/octocat/Hello-World/contents/ISSUE_TEMPLATE"),
				HTMLURL: String("https://github.com/octocat/Hello-World/blob/master/ISSUE_TEMPLATE"),
			},
			PullRequestTemplate: &Metric{
				URL:     String("https://api.github.com/repos/octocat/Hello-World/contents/PULL_REQUEST_TEMPLATE"),
				HTMLURL: String("https://github.com/octocat/Hello-World/blob/master/PULL_REQUEST_TEMPLATE"),
			},
			License: &Metric{
				Name:    String("MIT License"),
				Key:     String("mit"),
				SPDXID:  String("MIT"),
				URL:     String("https://api.github.com/licenses/mit"),
				HTMLURL: String("https://github.com/octocat/Hello-World/blob/master/LICENSE"),
				NodeID:  String("MDc6TGljZW5zZW1pdA=="),
			},
			Readme: &Metric{
				URL:     String("https://api.github.com/repos/octocat/Hello-World/contents/README.md"),
				HTMLURL: String("https://github.com/octocat/Hello-World/blob/master/README.md"),
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
	testJSONMarshal(t, &Metric{}, "{}")

	r := &Metric{
		Name:    String("name"),
		Key:     String("key"),
		SPDXID:  String("spdx_id"),
		URL:     String("url"),
		HTMLURL: String("hurl"),
		NodeID:  String("node_id"),
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
	testJSONMarshal(t, &CommunityHealthFiles{}, "{}")

	r := &CommunityHealthFiles{
		CodeOfConduct: &Metric{
			Name:    String("name"),
			Key:     String("key"),
			URL:     String("url"),
			HTMLURL: String("hurl"),
		},
		CodeOfConductFile: &Metric{
			Name:    String("name"),
			Key:     String("key"),
			URL:     String("url"),
			HTMLURL: String("hurl"),
		},
		Contributing: &Metric{
			Name:    String("name"),
			Key:     String("key"),
			URL:     String("url"),
			HTMLURL: String("hurl"),
		},
		IssueTemplate: &Metric{
			Name:    String("name"),
			Key:     String("key"),
			URL:     String("url"),
			HTMLURL: String("hurl"),
		},
		PullRequestTemplate: &Metric{
			Name:    String("name"),
			Key:     String("key"),
			URL:     String("url"),
			HTMLURL: String("hurl"),
		},
		License: &Metric{
			Name:    String("name"),
			Key:     String("key"),
			SPDXID:  String("spdx_id"),
			URL:     String("url"),
			HTMLURL: String("hurl"),
			NodeID:  String("node_id"),
		},
		Readme: &Metric{
			Name:    String("name"),
			Key:     String("key"),
			URL:     String("url"),
			HTMLURL: String("hurl"),
		},
	}

	want := `{
		"code_of_conduct": {
			"name": "name",
			"key": "key",
			"url": "url",
			"html_url": "hurl"
		},
		"code_of_conduct_file": {
			"name": "name",
			"key": "key",
                        "url": "url",
                        "html_url": "hurl"
                },
		"contributing": {
			"name": "name",
			"key": "key",
			"url": "url",
			"html_url": "hurl"
		},
		"issue_template": {
			"name": "name",
			"key": "key",
			"url": "url",
			"html_url": "hurl"
		},
		"pull_request_template": {
			"name": "name",
			"key": "key",
			"url": "url",
			"html_url": "hurl"
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
			"html_url": "hurl"
		}
	}`

	testJSONMarshal(t, r, want)
}

func TestCommunityHealthMetrics_Marshal(t *testing.T) {
	testJSONMarshal(t, &CommunityHealthMetrics{}, "{}")

	r := &CommunityHealthMetrics{
		HealthPercentage: Int(1),
		Description:      String("desc"),
		Documentation:    String("docs"),
		Files: &CommunityHealthFiles{
			CodeOfConduct: &Metric{
				Name:    String("name"),
				Key:     String("key"),
				URL:     String("url"),
				HTMLURL: String("hurl"),
			},
			CodeOfConductFile: &Metric{
				Name:    String("name"),
				Key:     String("key"),
				URL:     String("url"),
				HTMLURL: String("hurl"),
			},
			Contributing: &Metric{
				Name:    String("name"),
				Key:     String("key"),
				URL:     String("url"),
				HTMLURL: String("hurl"),
			},
			IssueTemplate: &Metric{
				Name:    String("name"),
				Key:     String("key"),
				URL:     String("url"),
				HTMLURL: String("hurl"),
			},
			PullRequestTemplate: &Metric{
				Name:    String("name"),
				Key:     String("key"),
				URL:     String("url"),
				HTMLURL: String("hurl"),
			},
			License: &Metric{
				Name:    String("name"),
				Key:     String("key"),
				SPDXID:  String("spdx_id"),
				URL:     String("url"),
				HTMLURL: String("hurl"),
				NodeID:  String("node_id"),
			},
			Readme: &Metric{
				Name:    String("name"),
				Key:     String("key"),
				URL:     String("url"),
				HTMLURL: String("hurl"),
			},
		},
		UpdatedAt:             &Timestamp{referenceTime},
		ContentReportsEnabled: Bool(true),
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
				"html_url": "hurl"
			},
			"code_of_conduct_file": {
                                "name": "name",
                                "key": "key",
                                "url": "url",
                                "html_url": "hurl"
                        },
			"contributing": {
				"name": "name",
				"key": "key",
				"url": "url",
				"html_url": "hurl"
			},
			"issue_template": {
				"name": "name",
				"key": "key",
				"url": "url",
				"html_url": "hurl"
			},
			"pull_request_template": {
				"name": "name",
				"key": "key",
				"url": "url",
				"html_url": "hurl"
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
				"html_url": "hurl"
			}
		},
		"updated_at": ` + referenceTimeStr + `,
		"content_reports_enabled": true
	}`

	testJSONMarshal(t, r, want)
}
