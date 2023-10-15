package github

type AService struct{}

// Get gets an A
//
// GitHub API docs: https://docs.github.com/rest/a/a#overridden-get-a
func (s *AService) Get() {}

// MissingFromMetadata isn't mentioned in metadata.yaml
func (s *AService) MissingFromMetadata() {}

// Undocumented uses an undocumented operation
//
// Note: Undocumented uses the undocumented GitHub API endpoint "GET /undocumented/{undocumented_id}".
func (s *AService) Undocumented() {}

// OutdatedLinks has links that are outdated or wrong
//
// GitHub API docs: https://docs.github.com/rest/a/a#update-a
func (s *AService) OutdatedLinks() {}

// AcceptableLinks has links that resolve to the same link as the metadata
//
// GitHub API docs: https://docs.github.com/en/rest/a/a#update-a
// GitHub API docs: https://docs.github.com/en/rest/a/a#overridden-get-a
func (s *AService) AcceptableLinks() {}

// GitHub API docs: https://docs.github.com/rest/a/a#overridden-get-a
func (s *AService) Uncommented() {
	// has no pre-existing comment
}
