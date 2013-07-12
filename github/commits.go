package github

import (
	"time"
)

// Commit represents a GitHub commit
type Commit struct {
	Added     []string   `json:"added,omitempty"`
	Author    User       `json:"author,omitempty"`
	Committer User       `json:"commiter,omitempty"`
	Distinct  bool       `json:"distinct,omitempty"`
	ID        string     `json:"id,omitempty"`
	Message   string     `json:"message,omitempty"`
	Modified  []string   `json:"modified,omitempty"`
	Removed   []string   `json:"removed,omitempty"`
	Timestamp *time.Time `json:"timestamp,omitempty"`
	URL       string     `json:"url,omitempty"`
}
