package github

type AService struct{}

// Get gets an A
//
//meta:operation GET /a/{non-canonical-id}
func (s *AService) Get() {}

// Undocumented uses an undocumented operation
//
//meta:operation GET /undocumented/{undocumented_id}
func (s *AService) Undocumented() {}

// OutdatedLinks has links that are outdated or wrong
//
// GitHub API docs: https://docs.github.com/rest/a/a#get-a
// GitHub API docs: https://example.com
// Note: Undocumented uses the undocumented GitHub API endpoint "GET /undocumented/{undocumented_id}".
//
//meta:operation post a/{a_id}
func (s *AService) OutdatedLinks() {}

//meta:operation GET /a/{a_id}
func (s *AService) Uncommented() {}

func (s *AService) unexported() {}

func NotAMethod() {}

type internalService struct{}

func (i *internalService) Get() {}
