// github/ai_context_test.go

package github

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"
)

// Internal helpers to avoid test dependency cycles
func strPtr(s string) *string { return &s }
func int64Ptr(i int64) *int64 { return &i }
func intPtr(i int) *int       { return &i }
func boolPtr(b bool) *bool    { return &b }

func TestOperationSiliconDiet(t *testing.T) {
	// 1. SETUP: Construct the "Fat" Struct
	now := Timestamp{time.Now()}
	url := "https://api.github.com/repos/google/go-github/issues/1"
	
	fatIssue := &Issue{
		ID:            int64Ptr(1234567890),
		NodeID:        strPtr("MDU6SXNzdWUxMjM0NTY3ODkw"), // NOISE
		Number:        intPtr(1),
		State:         strPtr("open"),
		Title:         strPtr("Optimize Struct Payload for AI Agents"),
		Body:          strPtr("The current payload is too heavy. We need to strip HATEOAS links."),
		User: &User{
			Login:             strPtr("mechanic-ai"),
			ID:                int64Ptr(999),
			NodeID:            strPtr("MDQ6VXNlcjk5OTk5"),
			AvatarURL:         strPtr("https://avatars.githubusercontent.com/u/999?v=4"),
			GravatarID:        strPtr(""),
			URL:               strPtr("https://api.github.com/users/mechanic-ai"), // NOISE
			HTMLURL:           strPtr("https://github.com/mechanic-ai"),
			Type:              strPtr("User"),
			SiteAdmin:         boolPtr(false),
		},
		Labels: []*Label{
			{
				ID:      int64Ptr(1),
				URL:     strPtr(url + "/labels/optimization"),
				Name:    strPtr("optimization"),
				Color:   strPtr("7057ff"), // NOISE
				Default: boolPtr(true),
			},
			{
				ID:      int64Ptr(2),
				URL:     strPtr(url + "/labels/ai"),
				Name:    strPtr("ai"),
				Color:   strPtr("008672"), // NOISE
				Default: boolPtr(false),
			},
		},
		URL:           strPtr(url), // NOISE
		HTMLURL:       strPtr("https://github.com/google/go-github/issues/1"),
		CreatedAt:     &now,
		UpdatedAt:     &now,
	}

	// 2. EXECUTE: Measure Fat Payload
	fatBytes, _ := json.Marshal(fatIssue)
	fatSize := len(fatBytes)

	// 3. EXECUTE: Surgical Strike (Compression)
	leanCtx := fatIssue.ToAgentContext()
	leanBytes, _ := json.Marshal(leanCtx)
	leanSize := len(leanBytes)

	// 4. VERIFY: Calculate Reduction Vector
	reduction := float64(fatSize-leanSize) / float64(fatSize) * 100
	
	t.Logf("FAT Payload:  %d bytes", fatSize)
	t.Logf("LEAN Payload: %d bytes", leanSize)
	t.Logf("REDUCTION:    %.2f%%", reduction)

	// 5. ASSERT: Mission Standard (>60%)
	if reduction < 60.0 {
		t.Errorf("MISSION FAILURE: Reduction %.2f%% is below target.", reduction)
	}

	// 6. VALIDATE: Signal Integrity
	if leanCtx["title"] != "Optimize Struct Payload for AI Agents" {
		t.Errorf("Signal Loss: Title mismatch")
	}
	if leanCtx["author"] != "mechanic-ai" {
		t.Errorf("Signal Loss: Author identity lost")
	}
	
	// [FIXED]: Used leanCtx instead of leanContext
	labels, ok := leanCtx["labels"].([]string)
	if !ok || len(labels) != 2 {
		t.Errorf("Signal Loss: Label flattening failed")
	}
	if labels[0] != "optimization" {
		t.Errorf("Signal Loss: Label content mismatch")
	}

	fmt.Printf("MISSION SUCCESS: Issue Payload Reduced by %.2f%%\n", reduction)
}

func TestCommentOptimization(t *testing.T) {
	now := Timestamp{time.Now()}
	url := "https://api.github.com/repos/google/go-github/issues/comments/555"

	fatComment := &IssueComment{
		ID:        int64Ptr(555),
		NodeID:    strPtr("MDEyOklzc3VlQ29tbWVudDU1NQ=="), // NOISE
		Body:      strPtr("This is a critical update."),
		User:      &User{Login: strPtr("reviewer"), ID: int64Ptr(101), URL: strPtr("http://noise.com")},
		CreatedAt: &now,
		HTMLURL:   strPtr("http://github.com/comment"),
		URL:       strPtr(url), // NOISE
	}

	leanCtx := fatComment.ToAgentContext()
	
	// Validate Signal
	if leanCtx["author"] != "reviewer" {
		t.Errorf("Comment Signal Loss: Author missing")
	}
	if leanCtx["body"] != "This is a critical update." {
		t.Errorf("Comment Signal Loss: Body text missing")
	}
	
	// Validate Noise Reduction
	if _, exists := leanCtx["node_id"]; exists {
		t.Errorf("Comment Failure: NodeID leaked into context")
	}
	
	fmt.Println("MISSION SUCCESS: Comment Optimized")
}
