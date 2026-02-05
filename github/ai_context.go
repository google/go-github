// github/ai_context.go

package github

import (
	"time"
)

// [OPERATION SILICON DIET]
// ToAgentContext transforms "Fat Structs" into high-signal context maps.
// Physics: Reduces payload entropy by >80% to maximize LLM context window efficiency.

// ---------------------------------------------------------------------
// 1. ISSUES (The Unit of Work)
// ---------------------------------------------------------------------

func (i *Issue) ToAgentContext() map[string]interface{} {
	if i == nil {
		return nil
	}

	ctx := map[string]interface{}{
		"number":     i.GetNumber(),
		"title":      i.GetTitle(),
		"state":      i.GetState(),
		"locked":     i.GetLocked(),
		"created_at": i.GetCreatedAt().Format(time.RFC3339),
		"updated_at": i.GetUpdatedAt().Format(time.RFC3339),
		"body":       trimBody(i.GetBody()),
		"html_url":   i.GetHTMLURL(),
	}

	if i.User != nil {
		ctx["author"] = i.User.GetLogin()
	}

	// Flatten Labels: []struct -> []string
	if len(i.Labels) > 0 {
		labels := make([]string, 0, len(i.Labels))
		for _, l := range i.Labels {
			if name := l.GetName(); name != "" {
				labels = append(labels, name)
			}
		}
		ctx["labels"] = labels
	}

	// Flatten Assignees
	if len(i.Assignees) > 0 {
		assignees := make([]string, 0, len(i.Assignees))
		for _, a := range i.Assignees {
			assignees = append(assignees, a.GetLogin())
		}
		ctx["assignees"] = assignees
	}

	// Contextual Hinting
	if i.IsPullRequest() {
		ctx["is_pull_request"] = true
	}

	return ctx
}

// ---------------------------------------------------------------------
// 2. PULL REQUESTS (The Code Vector)
// ---------------------------------------------------------------------

func (p *PullRequest) ToAgentContext() map[string]interface{} {
	if p == nil {
		return nil
	}

	ctx := map[string]interface{}{
		"number":        p.GetNumber(),
		"title":         p.GetTitle(),
		"state":         p.GetState(),
		"body":          trimBody(p.GetBody()),
		"html_url":      p.GetHTMLURL(),
		"draft":         p.GetDraft(),
		"merged":        p.GetMerged(),
		"mergeable":     p.GetMergeable(),
		"additions":     p.GetAdditions(),
		"deletions":     p.GetDeletions(),
		"changed_files": p.GetChangedFiles(),
		"created_at":    p.GetCreatedAt().Format(time.RFC3339),
	}

	if p.User != nil {
		ctx["author"] = p.User.GetLogin()
	}

	// Vector Definition (Source -> Target)
	if p.Head != nil {
		ctx["head_ref"] = p.Head.GetRef()
		ctx["head_sha"] = p.Head.GetSHA()
	}
	if p.Base != nil {
		ctx["base_ref"] = p.Base.GetRef()
	}

	return ctx
}

// ---------------------------------------------------------------------
// 3. COMMENTS (The Discussion Vector)
// ---------------------------------------------------------------------

func (c *IssueComment) ToAgentContext() map[string]interface{} {
	if c == nil {
		return nil
	}
	
	ctx := map[string]interface{}{
		"id":         c.GetID(),
		"body":       trimBody(c.GetBody()),
		"created_at": c.GetCreatedAt().Format(time.RFC3339),
		"html_url":   c.GetHTMLURL(),
	}

	if c.User != nil {
		ctx["author"] = c.User.GetLogin()
	}
	
	return ctx
}

// ---------------------------------------------------------------------
// 4. REPOSITORIES (The Battlefield)
// ---------------------------------------------------------------------

func (r *Repository) ToAgentContext() map[string]interface{} {
	if r == nil {
		return nil
	}

	ctx := map[string]interface{}{
		"name":           r.GetName(),
		"full_name":      r.GetFullName(),
		"description":    r.GetDescription(),
		"html_url":       r.GetHTMLURL(),
		"language":       r.GetLanguage(),
		"stars":          r.GetStargazersCount(),
		"forks":          r.GetForksCount(),
		"open_issues":    r.GetOpenIssuesCount(),
		"default_branch": r.GetDefaultBranch(),
		"private":        r.GetPrivate(),
		"archived":       r.GetArchived(),
	}

	if len(r.Topics) > 0 {
		ctx["topics"] = r.Topics
	}

	return ctx
}

// trimBody prevents token overflow from massive descriptions.
func trimBody(s string) string {
	const maxLen = 4000
	if len(s) > maxLen {
		return s[:maxLen] + "...[TRUNCATED]"
	}
	return s
}
