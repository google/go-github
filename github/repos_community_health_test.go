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
		testHeader(t, r, "Accept", mediaTypeRepositoryCommunityHealthMetricsPreview)
		fmt.Fprintf(w, `{
				"health_percentage": 100,
				"files": {
					"code_of_conduct": {
						"name": "Contributor Covenant",
						"key": "contributor_covenant",
						"url": null,
						"html_url": "https://github.com/octocat/Hello-World/blob/master/CODE_OF_CONDUCT.md"
					},
					"contributing": {
						"url": "https://api.github.com/repos/octocat/Hello-World/contents/CONTRIBUTING",
						"html_url": "https://github.com/octocat/Hello-World/blob/master/CONTRIBUTING"
					},
					"license": {
						"name": "MIT License",
						"key": "mit",
						"url": "https://api.github.com/licenses/mit",
						"html_url": "https://github.com/octocat/Hello-World/blob/master/LICENSE"
					},
					"readme": {
						"url": "https://api.github.com/repos/octocat/Hello-World/contents/README.md",
						"html_url": "https://github.com/octocat/Hello-World/blob/master/README.md"
					}
				},
				"updated_at": "2017-02-28T00:00:00Z"
			}`)
	})

	ctx := context.Background()
	got, _, err := client.Repositories.GetCommunityHealthMetrics(ctx, "o", "r")
	if err != nil {
		t.Errorf("Repositories.GetCommunityHealthMetrics returned error: %v", err)
	}

	updatedAt := time.Date(2017, time.February, 28, 0, 0, 0, 0, time.UTC)
	want := &CommunityHealthMetrics{
		HealthPercentage: Int(100),
		UpdatedAt:        &updatedAt,
		Files: &CommunityHealthFiles{
			CodeOfConduct: &Metric{
				Name:    String("Contributor Covenant"),
				Key:     String("contributor_covenant"),
				HTMLURL: String("https://github.com/octocat/Hello-World/blob/master/CODE_OF_CONDUCT.md"),
			},
			Contributing: &Metric{
				URL:     String("https://api.github.com/repos/octocat/Hello-World/contents/CONTRIBUTING"),
				HTMLURL: String("https://github.com/octocat/Hello-World/blob/master/CONTRIBUTING"),
			},
			License: &Metric{
				Name:    String("MIT License"),
				Key:     String("mit"),
				URL:     String("https://api.github.com/licenses/mit"),
				HTMLURL: String("https://github.com/octocat/Hello-World/blob/master/LICENSE"),
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
		URL:     String("url"),
		HTMLURL: String("hurl"),
	}

	want := `{
		"name": "name",
		"key": "key",
		"url": "url",
		"html_url": "hurl"
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
			URL:     String("url"),
			HTMLURL: String("hurl"),
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
			"url": "url",
			"html_url": "hurl"
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
		Files: &CommunityHealthFiles{
			CodeOfConduct: &Metric{
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
				URL:     String("url"),
				HTMLURL: String("hurl"),
			},
			Readme: &Metric{
				Name:    String("name"),
				Key:     String("key"),
				URL:     String("url"),
				HTMLURL: String("hurl"),
			},
		},
		UpdatedAt: &referenceTime,
	}

	want := `{
		"health_percentage": 1,
		"files": {
			"code_of_conduct": {
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
				"url": "url",
				"html_url": "hurl"
			},
			"readme": {
				"name": "name",
				"key": "key",
				"url": "url",
				"html_url": "hurl"
			}
		},
		"updated_at": ` + referenceTimeStr + `
	}`

	testJSONMarshal(t, r, want)
}
