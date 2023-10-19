package github

type AService struct{}

// Get gets an A
//
// GitHub API docs: https://docs.github.com/rest/a/a#overridden-get-a
//
//meta:operation GET /a/{a_id}
func (s *AService) Get() {}

// Undocumented uses an undocumented operation
//
// Note: Undocumented uses the undocumented GitHub API endpoint "GET /undocumented/{undocumented_id}".
//
//meta:operation GET /undocumented/{undocumented_id}
func (s *AService) Undocumented() {}

// OutdatedLinks has links that are outdated or wrong
//
// GitHub API docs: https://docs.github.com/rest/a/a#update-a
//
//meta:operation POST /a/{a_id}
func (s *AService) OutdatedLinks() {}

// GitHub API docs: https://docs.github.com/rest/a/a#overridden-get-a
//
//meta:operation GET /a/{a_id}
func (s *AService) Uncommented() {}

func (s *AService) unexported() {}

func NotAMethod() {}

type internalService struct{}

func (i *internalService) Get() {}
