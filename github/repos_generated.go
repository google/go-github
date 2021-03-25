// AUTO GENERATED FILE DO NOT EDIT

package github

import (
	"context"
)

// RepositoriesServiceInterface AUTO GENERATED DO NOT EDIT
type RepositoriesServiceInterface interface {
	// List the repositories for a user. Passing the empty string will list
	// repositories for the authenticated user.
	//
	// GitHub API docs: https://docs.github.com/en/free-pro-team@latest/rest/reference/repos/#list-repositories-for-the-authenticated-user
	// GitHub API docs: https://docs.github.com/en/free-pro-team@latest/rest/reference/repos/#list-repositories-for-a-user
	List(ctx context.Context, user string, opts *RepositoryListOptions) ([]*Repository, *Response, error)
	// ListByOrg lists the repositories for an organization.
	//
	// GitHub API docs: https://docs.github.com/en/free-pro-team@latest/rest/reference/repos/#list-organization-repositories
	ListByOrg(ctx context.Context, org string, opts *RepositoryListByOrgOptions) ([]*Repository, *Response, error)
	// ListAll lists all GitHub repositories in the order that they were created.
	//
	// GitHub API docs: https://docs.github.com/en/free-pro-team@latest/rest/reference/repos/#list-public-repositories
	ListAll(ctx context.Context, opts *RepositoryListAllOptions) ([]*Repository, *Response, error)
	// Create a new repository. If an organization is specified, the new
	// repository will be created under that org. If the empty string is
	// specified, it will be created for the authenticated user.
	//
	// Note that only a subset of the repo fields are used and repo must
	// not be nil.
	//
	// Also note that this method will return the response without actually
	// waiting for GitHub to finish creating the repository and letting the
	// changes propagate throughout its servers. You may set up a loop with
	// exponential back-off to verify repository's creation.
	//
	// GitHub API docs: https://docs.github.com/en/free-pro-team@latest/rest/reference/repos/#create-a-repository-for-the-authenticated-user
	// GitHub API docs: https://docs.github.com/en/free-pro-team@latest/rest/reference/repos/#create-an-organization-repository
	Create(ctx context.Context, org string, repo *Repository) (*Repository, *Response, error)
	// CreateFromTemplate generates a repository from a template.
	//
	// GitHub API docs: https://docs.github.com/en/free-pro-team@latest/rest/reference/repos/#create-a-repository-using-a-template
	CreateFromTemplate(ctx context.Context, templateOwner, templateRepo string, templateRepoReq *TemplateRepoRequest) (*Repository, *Response, error)
	// Get fetches a repository.
	//
	// GitHub API docs: https://docs.github.com/en/free-pro-team@latest/rest/reference/repos/#get-a-repository
	Get(ctx context.Context, owner, repo string) (*Repository, *Response, error)
	// GetCodeOfConduct gets the contents of a repository's code of conduct.
	//
	// GitHub API docs: https://docs.github.com/en/free-pro-team@latest/rest/reference/codes-of-conduct/#get-the-code-of-conduct-for-a-repository
	GetCodeOfConduct(ctx context.Context, owner, repo string) (*CodeOfConduct, *Response, error)
	// GetByID fetches a repository.
	//
	// Note: GetByID uses the undocumented GitHub API endpoint /repositories/:id.
	GetByID(ctx context.Context, id int64) (*Repository, *Response, error)
	// Edit updates a repository.
	//
	// GitHub API docs: https://docs.github.com/en/free-pro-team@latest/rest/reference/repos/#update-a-repository
	Edit(ctx context.Context, owner, repo string, repository *Repository) (*Repository, *Response, error)
	// Delete a repository.
	//
	// GitHub API docs: https://docs.github.com/en/free-pro-team@latest/rest/reference/repos/#delete-a-repository
	Delete(ctx context.Context, owner, repo string) (*Response, error)
	// GetVulnerabilityAlerts checks if vulnerability alerts are enabled for a repository.
	//
	// GitHub API docs: https://docs.github.com/en/free-pro-team@latest/rest/reference/repos/#check-if-vulnerability-alerts-are-enabled-for-a-repository
	GetVulnerabilityAlerts(ctx context.Context, owner, repository string) (bool, *Response, error)
	// EnableVulnerabilityAlerts enables vulnerability alerts and the dependency graph for a repository.
	//
	// GitHub API docs: https://docs.github.com/en/free-pro-team@latest/rest/reference/repos/#enable-vulnerability-alerts
	EnableVulnerabilityAlerts(ctx context.Context, owner, repository string) (*Response, error)
	// DisableVulnerabilityAlerts disables vulnerability alerts and the dependency graph for a repository.
	//
	// GitHub API docs: https://docs.github.com/en/free-pro-team@latest/rest/reference/repos/#disable-vulnerability-alerts
	DisableVulnerabilityAlerts(ctx context.Context, owner, repository string) (*Response, error)
	// EnableAutomatedSecurityFixes enables the automated security fixes for a repository.
	//
	// GitHub API docs: https://docs.github.com/en/free-pro-team@latest/rest/reference/repos/#enable-automated-security-fixes
	EnableAutomatedSecurityFixes(ctx context.Context, owner, repository string) (*Response, error)
	// DisableAutomatedSecurityFixes disables vulnerability alerts and the dependency graph for a repository.
	//
	// GitHub API docs: https://docs.github.com/en/free-pro-team@latest/rest/reference/repos/#disable-automated-security-fixes
	DisableAutomatedSecurityFixes(ctx context.Context, owner, repository string) (*Response, error)
	// ListContributors lists contributors for a repository.
	//
	// GitHub API docs: https://docs.github.com/en/free-pro-team@latest/rest/reference/repos/#list-repository-contributors
	ListContributors(ctx context.Context, owner string, repository string, opts *ListContributorsOptions) ([]*Contributor, *Response, error)
	// ListLanguages lists languages for the specified repository. The returned map
	// specifies the languages and the number of bytes of code written in that
	// language. For example:
	//
	//     {
	//       "C": 78769,
	//       "Python": 7769
	//     }
	//
	// GitHub API docs: https://docs.github.com/en/free-pro-team@latest/rest/reference/repos/#list-repository-languages
	ListLanguages(ctx context.Context, owner string, repo string) (map[string]int, *Response, error)
	// ListTeams lists the teams for the specified repository.
	//
	// GitHub API docs: https://docs.github.com/en/free-pro-team@latest/rest/reference/repos/#list-repository-teams
	ListTeams(ctx context.Context, owner string, repo string, opts *ListOptions) ([]*Team, *Response, error)
	// ListTags lists tags for the specified repository.
	//
	// GitHub API docs: https://docs.github.com/en/free-pro-team@latest/rest/reference/repos/#list-repository-tags
	ListTags(ctx context.Context, owner string, repo string, opts *ListOptions) ([]*RepositoryTag, *Response, error)
	// ListBranches lists branches for the specified repository.
	//
	// GitHub API docs: https://docs.github.com/en/free-pro-team@latest/rest/reference/repos/#list-branches
	ListBranches(ctx context.Context, owner string, repo string, opts *BranchListOptions) ([]*Branch, *Response, error)
	// GetBranch gets the specified branch for a repository.
	//
	// GitHub API docs: https://docs.github.com/en/free-pro-team@latest/rest/reference/repos/#get-a-branch
	GetBranch(ctx context.Context, owner, repo, branch string) (*Branch, *Response, error)
	// GetBranchProtection gets the protection of a given branch.
	//
	// GitHub API docs: https://docs.github.com/en/free-pro-team@latest/rest/reference/repos/#get-branch-protection
	GetBranchProtection(ctx context.Context, owner, repo, branch string) (*Protection, *Response, error)
	// GetRequiredStatusChecks gets the required status checks for a given protected branch.
	//
	// GitHub API docs: https://docs.github.com/en/free-pro-team@latest/rest/reference/repos/#get-status-checks-protection
	GetRequiredStatusChecks(ctx context.Context, owner, repo, branch string) (*RequiredStatusChecks, *Response, error)
	// ListRequiredStatusChecksContexts lists the required status checks contexts for a given protected branch.
	//
	// GitHub API docs: https://docs.github.com/en/free-pro-team@latest/rest/reference/repos/#get-all-status-check-contexts
	ListRequiredStatusChecksContexts(ctx context.Context, owner, repo, branch string) (contexts []string, resp *Response, err error)
	// UpdateBranchProtection updates the protection of a given branch.
	//
	// GitHub API docs: https://docs.github.com/en/free-pro-team@latest/rest/reference/repos/#update-branch-protection
	UpdateBranchProtection(ctx context.Context, owner, repo, branch string, preq *ProtectionRequest) (*Protection, *Response, error)
	// RemoveBranchProtection removes the protection of a given branch.
	//
	// GitHub API docs: https://docs.github.com/en/free-pro-team@latest/rest/reference/repos/#delete-branch-protection
	RemoveBranchProtection(ctx context.Context, owner, repo, branch string) (*Response, error)
	// GetSignaturesProtectedBranch gets required signatures of protected branch.
	//
	// GitHub API docs: https://docs.github.com/en/free-pro-team@latest/rest/reference/repos/#get-commit-signature-protection
	GetSignaturesProtectedBranch(ctx context.Context, owner, repo, branch string) (*SignaturesProtectedBranch, *Response, error)
	// RequireSignaturesOnProtectedBranch makes signed commits required on a protected branch.
	// It requires admin access and branch protection to be enabled.
	//
	// GitHub API docs: https://docs.github.com/en/free-pro-team@latest/rest/reference/repos/#create-commit-signature-protection
	RequireSignaturesOnProtectedBranch(ctx context.Context, owner, repo, branch string) (*SignaturesProtectedBranch, *Response, error)
	// OptionalSignaturesOnProtectedBranch removes required signed commits on a given branch.
	//
	// GitHub API docs: https://docs.github.com/en/free-pro-team@latest/rest/reference/repos/#delete-commit-signature-protection
	OptionalSignaturesOnProtectedBranch(ctx context.Context, owner, repo, branch string) (*Response, error)
	// UpdateRequiredStatusChecks updates the required status checks for a given protected branch.
	//
	// GitHub API docs: https://docs.github.com/en/free-pro-team@latest/rest/reference/repos/#update-status-check-protection
	UpdateRequiredStatusChecks(ctx context.Context, owner, repo, branch string, sreq *RequiredStatusChecksRequest) (*RequiredStatusChecks, *Response, error)
	// RemoveRequiredStatusChecks removes the required status checks for a given protected branch.
	//
	// GitHub API docs: https://docs.github.com/en/free-pro-team@latest/rest/reference/repos#remove-status-check-protection
	RemoveRequiredStatusChecks(ctx context.Context, owner, repo, branch string) (*Response, error)
	// License gets the contents of a repository's license if one is detected.
	//
	// GitHub API docs: https://docs.github.com/en/free-pro-team@latest/rest/reference/licenses/#get-the-license-for-a-repository
	License(ctx context.Context, owner, repo string) (*RepositoryLicense, *Response, error)
	// GetPullRequestReviewEnforcement gets pull request review enforcement of a protected branch.
	//
	// GitHub API docs: https://docs.github.com/en/free-pro-team@latest/rest/reference/repos/#get-pull-request-review-protection
	GetPullRequestReviewEnforcement(ctx context.Context, owner, repo, branch string) (*PullRequestReviewsEnforcement, *Response, error)
	// UpdatePullRequestReviewEnforcement patches pull request review enforcement of a protected branch.
	// It requires admin access and branch protection to be enabled.
	//
	// GitHub API docs: https://docs.github.com/en/free-pro-team@latest/rest/reference/repos/#update-pull-request-review-protection
	UpdatePullRequestReviewEnforcement(ctx context.Context, owner, repo, branch string, patch *PullRequestReviewsEnforcementUpdate) (*PullRequestReviewsEnforcement, *Response, error)
	// DisableDismissalRestrictions disables dismissal restrictions of a protected branch.
	// It requires admin access and branch protection to be enabled.
	//
	// GitHub API docs: https://docs.github.com/en/free-pro-team@latest/rest/reference/repos/#update-pull-request-review-protection
	DisableDismissalRestrictions(ctx context.Context, owner, repo, branch string) (*PullRequestReviewsEnforcement, *Response, error)
	// RemovePullRequestReviewEnforcement removes pull request enforcement of a protected branch.
	//
	// GitHub API docs: https://docs.github.com/en/free-pro-team@latest/rest/reference/repos/#delete-pull-request-review-protection
	RemovePullRequestReviewEnforcement(ctx context.Context, owner, repo, branch string) (*Response, error)
	// GetAdminEnforcement gets admin enforcement information of a protected branch.
	//
	// GitHub API docs: https://docs.github.com/en/free-pro-team@latest/rest/reference/repos/#get-admin-branch-protection
	GetAdminEnforcement(ctx context.Context, owner, repo, branch string) (*AdminEnforcement, *Response, error)
	// AddAdminEnforcement adds admin enforcement to a protected branch.
	// It requires admin access and branch protection to be enabled.
	//
	// GitHub API docs: https://docs.github.com/en/free-pro-team@latest/rest/reference/repos/#set-admin-branch-protection
	AddAdminEnforcement(ctx context.Context, owner, repo, branch string) (*AdminEnforcement, *Response, error)
	// RemoveAdminEnforcement removes admin enforcement from a protected branch.
	//
	// GitHub API docs: https://docs.github.com/en/free-pro-team@latest/rest/reference/repos/#delete-admin-branch-protection
	RemoveAdminEnforcement(ctx context.Context, owner, repo, branch string) (*Response, error)
	// ListAllTopics lists topics for a repository.
	//
	// GitHub API docs: https://docs.github.com/en/free-pro-team@latest/rest/reference/repos/#get-all-repository-topics
	ListAllTopics(ctx context.Context, owner, repo string) ([]string, *Response, error)
	// ReplaceAllTopics replaces topics for a repository.
	//
	// GitHub API docs: https://docs.github.com/en/free-pro-team@latest/rest/reference/repos/#replace-all-repository-topics
	ReplaceAllTopics(ctx context.Context, owner, repo string, topics []string) ([]string, *Response, error)
	// ListApps lists the GitHub apps that have push access to a given protected branch.
	// It requires the GitHub apps to have `write` access to the `content` permission.
	//
	// GitHub API docs: https://docs.github.com/en/free-pro-team@latest/rest/reference/repos/#get-apps-with-access-to-the-protected-branch
	ListApps(ctx context.Context, owner, repo, branch string) ([]*App, *Response, error)
	// ReplaceAppRestrictions replaces the apps that have push access to a given protected branch.
	// It removes all apps that previously had push access and grants push access to the new list of apps.
	// It requires the GitHub apps to have `write` access to the `content` permission.
	//
	// Note: The list of users, apps, and teams in total is limited to 100 items.
	//
	// GitHub API docs: https://docs.github.com/en/free-pro-team@latest/rest/reference/repos/#set-app-access-restrictions
	ReplaceAppRestrictions(ctx context.Context, owner, repo, branch string, slug []string) ([]*App, *Response, error)
	// AddAppRestrictions grants the specified apps push access to a given protected branch.
	// It requires the GitHub apps to have `write` access to the `content` permission.
	//
	// Note: The list of users, apps, and teams in total is limited to 100 items.
	//
	// GitHub API docs: https://docs.github.com/en/free-pro-team@latest/rest/reference/repos/#add-app-access-restrictions
	AddAppRestrictions(ctx context.Context, owner, repo, branch string, slug []string) ([]*App, *Response, error)
	// RemoveAppRestrictions removes the ability of an app to push to this branch.
	// It requires the GitHub apps to have `write` access to the `content` permission.
	//
	// Note: The list of users, apps, and teams in total is limited to 100 items.
	//
	// GitHub API docs: https://docs.github.com/en/free-pro-team@latest/rest/reference/repos/#remove-app-access-restrictions
	RemoveAppRestrictions(ctx context.Context, owner, repo, branch string, slug []string) ([]*App, *Response, error)
	// Transfer transfers a repository from one account or organization to another.
	//
	// This method might return an *AcceptedError and a status code of
	// 202. This is because this is the status that GitHub returns to signify that
	// it has now scheduled the transfer of the repository in a background task.
	// A follow up request, after a delay of a second or so, should result
	// in a successful request.
	//
	// GitHub API docs: https://docs.github.com/en/free-pro-team@latest/rest/reference/repos/#transfer-a-repository
	Transfer(ctx context.Context, owner, repo string, transfer TransferRequest) (*Repository, *Response, error)
	// Dispatch triggers a repository_dispatch event in a GitHub Actions workflow.
	//
	// GitHub API docs: https://docs.github.com/en/free-pro-team@latest/rest/reference/repos/#create-a-repository-dispatch-event
	Dispatch(ctx context.Context, owner, repo string, opts DispatchRequestOptions) (*Repository, *Response, error)
}
