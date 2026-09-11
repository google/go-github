package github

// DemoRequest is a demo request body.
//
//meta:schema request POST /demo
type DemoRequest struct {
	ID   int64   `json:"id"`
	Name string  `json:"name"`
	Note *string `json:"note,omitempty"`
}
