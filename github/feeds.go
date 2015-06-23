package github

//FeedsService handles communication with the feeds related
//methods of the GitHub API.
//Note. client should provide basic basic auth.
//
//Github API docs: https://developer.github.com/v3/activity/feeds/
type FeedsService struct {
	client *Client
}

//Feeds represents the all the feeds available to the authenticated user
type Feeds struct {
	TimeLineURL                 *string      `json:"timeline_url,omitempty"`
	User                        *string      `json:"user_url,omitempty"`
	PublicURL                   *string      `json:"current_user_public_url,omitempty"`
	URL                         *string      `json:"current_user_url,omitempty"`
	CurrentUserActorURL         *string      `json:"current_user_actor_url,omitempty"`
	CurrentUserOrgannizationURL *string      `json:"current_user_organization_url,omitempty"`
	CurrentUserOrganizationList []string     `json:"current_user_organization_urls,omitempty"`
	Links                       *LinkService `json:"_links,omitempty"`
}

func (r Feeds) String() string {
	return Stringify(r)
}


//LinkService represents list of available links
type LinkService struct {
	Timeline                 *LinkData   `json:"timeline,omitempty"`
	User                     *LinkData   `json:"user,omitempty"`
	CurrentUserPublic        *LinkData   `json:"current_user_public,omitempty"`
	CurrentUser              *LinkData   `json:"current_user,omitempty"`
	CurrentUserActor         *LinkData   `json:"current_user_actor,omitempty"`
	CurrentUserOrganizations []*LinkData `json:"current_user_organizations,omitempty"`
}

//LinkData represents information from each link
type LinkData struct {
	Href string `json:"href, omitempty"`
	Type string `json:"type, omitempty"`
}

//ListFeeds Get feeds available to the authenticated user
//
//Github API Docs: https://developer.github.com/v3/activity/feeds/#list-feeds
func (s *FeedsService) ListFeeds() (*Feeds, *Response, error) {
	req, err := s.client.NewRequest("GET", "https://api.github.com/feeds", nil)
	if err != nil {
		return nil, nil, err
	}

	feeds := &Feeds{}
	resp, err := s.client.Do(req, feeds)
	if err != nil {
		return nil, resp, err
	}

	return feeds, resp, err
}
