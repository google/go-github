package internal

import (
	"context"
	"fmt"
	"io"
	"regexp"
	"sort"
	"strconv"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/google/go-github/v55/github"
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

type contentsClient interface {
	DownloadContents(ctx context.Context, owner, repo, filepath string, opts *github.RepositoryContentGetOptions) (io.ReadCloser, *github.Response, error)
	GetContents(ctx context.Context, owner, repo, path string, opts *github.RepositoryContentGetOptions) (*github.RepositoryContent, []*github.RepositoryContent, *github.Response, error)
}

func (o *openapiFile) loadDescription(ctx context.Context, client contentsClient, gitRef string) error {
	contents, resp, err := client.DownloadContents(
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
func getDescriptions(ctx context.Context, client contentsClient, gitRef string) ([]*openapiFile, error) {
	_, dir, resp, err := client.GetContents(
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
			plan := m[pattern.SubexpIndex("plan")]
			major, _ := strconv.Atoi(m[pattern.SubexpIndex("major")])
			minor, _ := strconv.Atoi(m[pattern.SubexpIndex("minor")])
			if plan == "ghes" && major < 3 {
				continue
			}
			filename := fmt.Sprintf("descriptions/%s/%s.json", d.GetName(), d.GetName())
			files = append(files, &openapiFile{
				filename:     filename,
				plan:         plan,
				planIdx:      i,
				releaseMajor: major,
				releaseMinor: minor,
			})
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
