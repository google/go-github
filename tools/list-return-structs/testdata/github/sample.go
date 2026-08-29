// Package github is a fixture used only by list-return-structs tests.
package github

import "context"

// Returned directly by Get and Create; never an input. Should be INCLUDED.
type Widget struct {
	Name string `json:"name,omitempty"`
	Spec WidgetSpec
}

// Field of Widget (a return struct); never an input. Should be INCLUDED
// transitively.
type WidgetSpec struct {
	Detail string `json:"detail"`
}

// Field of CreateWidgetRequest (an input struct); never returned and never
// directly an opts/body param. Should be EXCLUDED transitively.
type RequestMeta struct {
	Token string `json:"token,omitempty"`
}

// Used as an opts param. Should be EXCLUDED.
type ListWidgetsOptions struct {
	Page int `json:"page,omitempty"`
}

// Used as a body param. Should be EXCLUDED.
type CreateWidgetRequest struct {
	Title string `json:"title"`
	Meta  RequestMeta
}

// Returned by DualThing AND used as its body. Should be EXCLUDED.
type Dual struct {
	Val string `json:"val,omitempty"`
}

// Returned by ListEmbedded. Should be INCLUDED.
type Embedded struct {
	Inner Inner
}

// Field of Embedded (a return struct). Should be INCLUDED transitively.
type Inner struct {
	X int `json:"x,omitempty"`
}

// FakeService is a receiver for the fixture methods.
type FakeService struct{}

// Response is a placeholder return type (deliberately not a struct here).
type Response struct{}

func (s *FakeService) Get(ctx context.Context, opts *ListWidgetsOptions) (*Widget, *Response, error) {
	return nil, nil, nil
}

func (s *FakeService) Create(ctx context.Context, body CreateWidgetRequest) (*Widget, *Response, error) {
	return nil, nil, nil
}

func (s *FakeService) DualThing(ctx context.Context, body Dual) (*Dual, *Response, error) {
	return nil, nil, nil
}

func (s *FakeService) ListEmbedded(ctx context.Context) (*Embedded, *Response, error) {
	return nil, nil, nil
}
