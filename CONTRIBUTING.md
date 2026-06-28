```go
// Repository represents a GitHub repository.
type Repository struct {
	ID     *int64  `json:"id,omitempty"`
	NodeID *string `json:"node_id,omitempty"`
	Owner  *User   `json:"owner,omitempty"`
	// ...
}
```