// Copyright 2023 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"context"
	"fmt"
	"io"
	"regexp"
	"sort"
	"strconv"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/google/go-github/v69/github"
	"golang.org/x/sync/errgroup"
)

const (
	descriptionsOwnerName = "github"
	descriptionsRepoName  = "rest-api-description"
	descriptionsPath      = "descriptions"
)

type openapiFile struct {
	description  *openapi3.T
	filename     string
	plan         string
	planIdx      int
	releaseMajor int
	releaseMinor int
}

func getOpsFromGithub(ctx context.Context, client *github.Client, gitRef string) ([]*operation, error) {
	descs, err := getDescriptions(ctx, client, gitRef)
	if err != nil {
		return nil, err
	}
	var ops []*operation
	for _, desc := range descs {
		for p, pathItem := range desc.description.Paths.Map() {
			for method, op := range pathItem.Operations() {
				docURL := ""
				if op.ExternalDocs != nil {
					docURL = op.ExternalDocs.URL
				}
				ops = addOperation(ops, desc.filename, method+" "+p, docURL)
			}
		}
	}
	sortOperations(ops)
	return ops, nil
}

func (o *openapiFile) loadDescription(ctx context.Context, client *github.Client, gitRef string) error {
	contents, resp, err := client.Repositories.DownloadContents(
		ctx,
		descriptionsOwnerName,
		descriptionsRepoName,
		o.filename,
		&github.RepositoryContentGetOptions{Ref: gitRef},
	)
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		return fmt.Errorf("unexpected status code: %s", resp.Status)
	}
	b, err := io.ReadAll(contents)
	if err != nil {
		return err
	}
	err = contents.Close()
	if err != nil {
		return err
	}
	o.description, err = openapi3.NewLoader().LoadFromData(b)
	return err
}

// less sorts by the following rules:
//   - planIdx ascending
//   - releaseMajor descending
//   - releaseMinor descending
func (o *openapiFile) less(other *openapiFile) bool {
	if o.planIdx != other.planIdx {
		return o.planIdx < other.planIdx
	}
	if o.releaseMajor != other.releaseMajor {
		return o.releaseMajor > other.releaseMajor
	}
	return o.releaseMinor > other.releaseMinor
}

var dirPatterns = []*regexp.Regexp{
	regexp.MustCompile(`^(?P<plan>api\.github\.com)(-(?P<major>\d+)\.(?P<minor>\d+))?$`),
	regexp.MustCompile(`^(?P<plan>ghec)(-(?P<major>\d+)\.(?P<minor>\d+))?$`),
	regexp.MustCompile(`^(?P<plan>ghes)(-(?P<major>\d+)\.(?P<minor>\d+))?$`),
}

// getDescriptions loads OpenapiFiles for all the OpenAPI 3.0 description files in github/rest-api-description.
// This assumes that all directories in "descriptions/" contain OpenAPI 3.0 description files with the same
// name as the directory (plus the ".json" extension). For example, "descriptions/api.github.com/api.github.com.json".
// Results are sorted by these rules:
//   - Directories that don't match any of the patterns in dirPatterns are removed.
//   - Directories are sorted by the pattern that matched in the same order they appear in dirPatterns.
//   - Directories are then sorted by major and minor version in descending order.
func getDescriptions(ctx context.Context, client *github.Client, gitRef string) ([]*openapiFile, error) {
	_, dir, resp, err := client.Repositories.GetContents(
		ctx,
		descriptionsOwnerName,
		descriptionsRepoName,
		descriptionsPath,
		&github.RepositoryContentGetOptions{Ref: gitRef},
	)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("unexpected status code: %s", resp.Status)
	}
	files := make([]*openapiFile, 0, len(dir))
	for _, d := range dir {
		for i, pattern := range dirPatterns {
			m := pattern.FindStringSubmatch(d.GetName())
			if m == nil {
				continue
			}
			file := openapiFile{
				filename: fmt.Sprintf("descriptions/%s/%s.json", d.GetName(), d.GetName()),
				plan:     m[pattern.SubexpIndex("plan")],
				planIdx:  i,
			}
			rawMajor := m[pattern.SubexpIndex("major")]
			if rawMajor != "" {
				file.releaseMajor, err = strconv.Atoi(rawMajor)
				if err != nil {
					return nil, err
				}
			}
			rawMinor := m[pattern.SubexpIndex("minor")]
			if rawMinor != "" {
				file.releaseMinor, err = strconv.Atoi(rawMinor)
				if err != nil {
					return nil, err
				}
			}
			if file.plan == "ghes" && file.releaseMajor < 3 {
				continue
			}
			files = append(files, &file)
			break
		}
	}
	sort.Slice(files, func(i, j int) bool {
		return files[i].less(files[j])
	})
	g, ctx := errgroup.WithContext(ctx)
	for _, file := range files {
		f := file
		g.Go(func() error {
			return f.loadDescription(ctx, client, gitRef)
		})
	}
	err = g.Wait()
	if err != nil {
		return nil, err
	}
	return files, nil
}
