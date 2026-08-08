package github

type DemoRequest struct {
	ID   int64   `json:"id"`
	Name string  `json:"name"`
	Note *string `json:"note,omitempty"`
}

func (s *svc) Create(ctx context.Context, body *DemoRequest) {
	s.client.NewRequest(ctx, "POST", "u", body)
}
