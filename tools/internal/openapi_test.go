package internal

import (
	"context"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/google/go-github/v55/github"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type dummyContentsClient struct {
	t                *testing.T
	downloadContents func(ctx context.Context, owner, repo, filepath string, opts *github.RepositoryContentGetOptions) (io.ReadCloser, *github.Response, error)
	getContents      func(ctx context.Context, owner, repo, path string, opts *github.RepositoryContentGetOptions) (*github.RepositoryContent, []*github.RepositoryContent, *github.Response, error)
}

func (d dummyContentsClient) DownloadContents(ctx context.Context, owner, repo, filepath string, opts *github.RepositoryContentGetOptions) (io.ReadCloser, *github.Response, error) {
	d.t.Helper()
	return d.downloadContents(ctx, owner, repo, filepath, opts)
}

func (d dummyContentsClient) GetContents(ctx context.Context, owner, repo, path string, opts *github.RepositoryContentGetOptions) (*github.RepositoryContent, []*github.RepositoryContent, *github.Response, error) {
	d.t.Helper()
	return d.getContents(ctx, owner, repo, path, opts)
}

func TestGetDescriptions(t *testing.T) {
	ctx := context.Background()
	client := dummyContentsClient{
		t: t,
		getContents: func(_ context.Context, owner, repo, path string, opts *github.RepositoryContentGetOptions) (*github.RepositoryContent, []*github.RepositoryContent, *github.Response, error) {
			t.Helper()
			assert.Equal(t, descriptionsOwnerName, owner)
			assert.Equal(t, descriptionsRepoName, repo)
			assert.Equal(t, "descriptions", path)
			assert.Equal(t, &github.RepositoryContentGetOptions{Ref: "main"}, opts)
			dir := []*github.RepositoryContent{
				{Name: github.String("ghes-3.1")},
				{Name: github.String("ghes-3.2")},
				{Name: github.String("ignore_me")},
				{Name: github.String("api.github.com")},
				{Name: github.String("ghec")},
			}
			return nil, dir, &github.Response{Response: &http.Response{StatusCode: 200}}, nil
		},
		downloadContents: func(_ context.Context, owner, repo, filepath string, opts *github.RepositoryContentGetOptions) (io.ReadCloser, *github.Response, error) {
			assert.Equal(t, descriptionsOwnerName, owner)
			assert.Equal(t, descriptionsRepoName, repo)
			assert.Equal(t, &github.RepositoryContentGetOptions{Ref: "main"}, opts)
			contents := `{"info": {"version": "1.2.3"}}`
			switch filepath {
			case "descriptions/api.github.com/api.github.com.json":
			case "descriptions/ghes-3.1/ghes-3.1.json":
			case "descriptions/ghes-3.2/ghes-3.2.json":
			case "descriptions/ghec/ghec.json":
			default:
				t.Errorf("unexpected filepath: %s", filepath)
			}
			return io.NopCloser(strings.NewReader(contents)), &github.Response{Response: &http.Response{StatusCode: 200}}, nil
		},
	}
	wantFiles := []string{
		"descriptions/api.github.com/api.github.com.json",
		"descriptions/ghec/ghec.json",
		"descriptions/ghes-3.2/ghes-3.2.json",
		"descriptions/ghes-3.1/ghes-3.1.json",
	}
	got, err := getDescriptions(ctx, client, "main")
	require.NoError(t, err)
	require.Equal(t, len(wantFiles), len(got))
	for i := range wantFiles {
		require.Equal(t, wantFiles[i], got[i].filename)
		require.Equal(t, "1.2.3", got[i].description.Info.Version)
	}
}
