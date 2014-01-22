package github

import (
	"fmt"
	"net/url"
)

// FileOptions are shared parameters for CreateFile, UpdateFile, and DeleteFile.
type FileOptions struct {
	Path      string        `json:"path,omitempty"`
	Message   string        `json:"message,omitempty"`
	Branch    string        `json:"branch,omitempty"`
	Author    *CommitAuthor `json:"author,omitempty"`
	Committer *CommitAuthor `json:"committer,omitempty"`
}

// RepositoryContentResponse holds the parsed response from CreateFile, UpdateFile, and DeleteFile.
type RepositoryContentResponse struct {
	Content *RepositoryFile `json:"content,omitempty"`
	Commit  *CommitWithURL  `json:"commit,omitempty"`
}

// CommitWithURL contains commit information along with GitHub URLs for the commit.
type CommitWithURL struct {
	Commit
	URL     *string `json:"url,omitempty"`
	HTMLURL *string `json:"html_url,omitempty"`
}

// RepositoryFileLinks contains GitHub links to a file or folder in a repository.
type RepositoryFileLinks struct {
	Git  *string `json:"git,omitempty"`
	Self *string `json:"self,omitempty"`
	HTML *string `json:"html,omitempty"`
}

// RepositoryFile represents a file in a repository.
type RepositoryFile struct {
	Type            *string              `json:"type,omitempty"`
	Encoding        *string              `json:"encoding,omitempty"`
	Target          *string              `json:"target,omitempty"`
	Size            *int                 `json:"size,omitempty"`
	Name            *string              `json:"name,omitempty"`
	Path            *string              `json:"path,omitempty"`
	Content         *string              `json:"content,omitempty"`
	SHA             *string              `json:"sha,omitempty"`
	URL             *string              `json:"url,omitempty"`
	GitURL          *string              `json:"git_url,omitempty"`
	HTMLURL         *string              `json:"html_url,omitempty"`
	SubmoduleGitURL *string              `json:"submodule_git_url,omitempty"`
	Links           *RepositoryFileLinks `json:"_links,omitempty"`
}

// RefOption represents an optional ref parameter, which can be a SHA,
// branch, or tag
type RefOption struct {
	Ref string `url:"ref,omitempty"`
}

// GetReadme returns the README for a repository. Ref (optional; defaults to
// 'master') can be a SHA, branch, or tag.
//
// GitHub API docs: http://developer.github.com/v3/repos/contents/#get-the-readme
func (s *RepositoriesService) GetReadme(owner, repo string, opt *RefOption) (*RepositoryFile, *Response, error) {
	u := fmt.Sprintf("/repos/%s/%s/readme", owner, repo)
	u, err := addOptions(u, opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	readme := new(RepositoryFile)
	resp, err := s.client.Do(req, readme)
	if err != nil {
		return nil, resp, err
	}

	return readme, resp, err
}

// RepositoryContentsOptions specifies optional parameters for GetContents
type RepositoryContentsOptions struct {
	Ref  string `url:"ref,omitempty"`
	Path string `url:"path,omitempty"`
}

// GetContents returns the file(s) at the given path. Path and ref are optional.
//
// GitHub API docs: http://developer.github.com/v3/repos/contents/#get-contents
func (s *RepositoriesService) GetContents(owner, repo string, opt *RepositoryContentsOptions) ([]RepositoryFile, *Response, error) {
	u := fmt.Sprintf("/repos/%s/%s/contents/%s", owner, repo, opt.Path)
	u, err := addOptions(u, opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	file := new(RepositoryFile)
	files := new([]RepositoryFile)

	resp, i, err := s.client.PolymorphicDo(req, files, file)
	if i == 1 {
		files = &[]RepositoryFile{*file}
	}

	return *files, resp, err
}

// RepositoryCreateFileOptions specifies parameters for creating a file with CreateFile
// Path, Message, and Content are required fields
type RepositoryCreateFileOptions struct {
	FileOptions
	Content *[]byte `json:"content,omitempty"`
}

// CreateFile creates a new file at a path with the given content and
// returns the commit and file metadata. (Note that GitHub's API requires
// content to be base64-encoded, but go's encoding/json automatically does
// this for []byte fields.
//
// GitHub API docs: http://developer.github.com/v3/repos/contents/#create-a-file
func (s *RepositoriesService) CreateFile(owner, repo string, opt *RepositoryCreateFileOptions) (*RepositoryContentResponse, *Response, error) {
	u := fmt.Sprintf("/repos/%s/%s/contents/%s", owner, repo, opt.Path)
	req, err := s.client.NewRequest("PUT", u, opt)
	if err != nil {
		return nil, nil, err
	}

	createResp := new(RepositoryContentResponse)
	resp, err := s.client.Do(req, createResp)
	if err != nil {
		return nil, resp, err
	}

	return createResp, resp, err
}

// RepositoryUpdateFileOptions specifies parameters for updating a file with UpdateFile
// Path, Message, SHA, Content, and Branch are required fields
type RepositoryUpdateFileOptions struct {
	FileOptions
	Content *[]byte `json:"content,omitempty"`
	SHA     string  `json:"sha,omitempty"`
}

// UpdateFile updates a file at a path with the given content and
// returns the commit and file metadata. Requires the SHA to be updated,
// and will fail if the branch has advanced.
// (Note that GitHub's API requires content to be base64-encoded, but go's
// encoding/json automatically does this for []byte fields.
//
// GitHub API docs: http://developer.github.com/v3/repos/contents/#update-a-file
func (s *RepositoriesService) UpdateFile(owner, repo string, opt *RepositoryUpdateFileOptions) (*RepositoryContentResponse, *Response, error) {
	u := fmt.Sprintf("/repos/%s/%s/contents/%s", owner, repo, opt.Path)
	req, err := s.client.NewRequest("PUT", u, opt)
	if err != nil {
		return nil, nil, err
	}

	createResp := new(RepositoryContentResponse)
	resp, err := s.client.Do(req, createResp)
	if err != nil {
		return nil, resp, err
	}

	return createResp, resp, err
}

// RepositoryDeleteFileOptions specifies parameters for DeleteFile.
// Path, Message, SHA, and Branch are required fields
type RepositoryDeleteFileOptions struct {
	FileOptions
	SHA string `json:"sha,omitempty"`
}

// DeleteFile deletes the file at a path and returns the commit.
// Requires the SHA to be updated, and will fail if the branch has advanced.
//
// GitHub API docs: http://developer.github.com/v3/repos/contents/#delete-a-file
func (s *RepositoriesService) DeleteFile(owner, repo string, opt *RepositoryDeleteFileOptions) (*RepositoryContentResponse, *Response, error) {
	u := fmt.Sprintf("/repos/%s/%s/contents/%s", owner, repo, opt.Path)

	req, err := s.client.NewRequest("DELETE", u, opt)
	if err != nil {
		return nil, nil, err
	}

	createResp := new(RepositoryContentResponse)
	resp, err := s.client.Do(req, createResp)
	if err != nil {
		return nil, resp, err
	}

	return createResp, resp, err
}

// GetArchiveLinkOptions contains an optional ref parameter for GetArchiveLink.
type GetArchiveLinkOptions struct {
	Ref *string `json:"path,omitempty"`
}

// GetArchiveLink returns an URL to download a tarball or zipball archive for
// a repository. GitHub returns the URL using a 302 redirect, but this method
// does not follow the redirect. archive_format is either 'tarball' or 'zipball'.
//
// GitHub API docs: http://developer.github.com/v3/repos/contents/#get-archive-link
func (s *RepositoriesService) GetArchiveLink(owner, repo, archiveFormat string, opt *RefOption) (*url.URL, *Response, error) {
	u := fmt.Sprintf("/repos/%s/%s/%s/%s", owner, repo, archiveFormat, opt.Ref)

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	resp, err := s.noRedirectHTTPClient.Do(req)
	if resp.StatusCode != 302 {
		return nil, nil, err
	}

	parsedURL, err := url.Parse(resp.Header.Get("Location"))
	return parsedURL, newResponse(resp), err
}
