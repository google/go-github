package github

import "fmt"

// ListWatchers lists watchers of a particular repo.
//
// GitHub API Docs: http://developer.github.com/v3/activity/watching/#list-watchers
func (s *ActivityService) ListWatchers(owner, repo string) ([]User, *Response, error) {
	url := fmt.Sprintf("repos/%s/%s/subscribers", owner, repo)

	req, err := s.client.NewRequest("GET", url, nil)
	if err != nil {
		return nil, nil, err
	}

	watchers := new([]User)
	resp, err := s.client.Do(req, watchers)
	if err != nil {
		return nil, resp, err
	}

	return *watchers, resp, err
}
