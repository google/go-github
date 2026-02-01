// Copyright 2026 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package github

import (
	"errors"
	"slices"
	"testing"
)

func TestScan(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		pages         [][]int
		responses     []*Response
		wantItems     []int
		wantErr       error
		callErrBefore bool // call error function before exhausting iterator
	}{
		{
			name: "single page",
			pages: [][]int{
				{1, 2, 3},
			},
			responses: []*Response{
				{NextPage: 0},
			},
			wantItems: []int{1, 2, 3},
		},
		{
			name: "multiple pages with offset pagination",
			pages: [][]int{
				{1, 2},
				{3, 4},
				{5},
			},
			responses: []*Response{
				{NextPage: 2},
				{NextPage: 3},
				{NextPage: 0},
			},
			wantItems: []int{1, 2, 3, 4, 5},
		},
		{
			name: "multiple pages with cursor pagination",
			pages: [][]int{
				{1, 2},
				{3, 4},
				{5},
			},
			responses: []*Response{
				{After: "cursor1"},
				{After: "cursor2"},
				{After: ""},
			},
			wantItems: []int{1, 2, 3, 4, 5},
		},
		{
			name: "error on first page",
			pages: [][]int{
				nil,
			},
			responses: []*Response{
				nil,
			},
			wantErr: errors.New("request failed"),
		},
		{
			name: "error on second page",
			pages: [][]int{
				{1, 2},
				nil,
			},
			responses: []*Response{
				{NextPage: 2},
				nil,
			},
			wantErr: errors.New("request failed"),
		},
		{
			name:          "error function called before iterator exhausted",
			callErrBefore: true,
		},
		{
			name: "empty pages",
			pages: [][]int{
				{},
			},
			responses: []*Response{
				{NextPage: 0},
			},
			wantItems: []int{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if tt.callErrBefore {
				defer func() {
					if r := recover(); r == nil {
						t.Error("expected panic but got none")
					}
				}()
			}

			pageIdx := 0
			f := func(PaginationOption) ([]int, *Response, error) {
				if pageIdx >= len(tt.pages) {
					t.Fatal("unexpected pagination call")
				}
				if tt.wantErr != nil && pageIdx == len(tt.pages)-1 {
					pageIdx++
					return nil, nil, tt.wantErr
				}
				page := tt.pages[pageIdx]
				resp := tt.responses[pageIdx]
				pageIdx++
				return page, resp, nil
			}

			it, hasErr := Scan(f)

			if tt.callErrBefore {
				_ = hasErr() // should panic
				return
			}

			got := slices.Collect(it)
			err := hasErr()

			if tt.wantErr != nil {
				if err == nil {
					t.Errorf("want error %v, got nil", tt.wantErr)
				}
				if err.Error() != tt.wantErr.Error() {
					t.Errorf("want error %v, got %v", tt.wantErr, err)
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			want := tt.wantItems
			if !slices.Equal(got, want) {
				t.Errorf("got %v, want %v", got, want)
			}
		})
	}
}

func TestScan2(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		pages     [][]string
		responses []*Response
		wantItems []string
		wantErr   error
	}{
		{
			name: "single page",
			pages: [][]string{
				{"a", "b", "c"},
			},
			responses: []*Response{
				{NextPage: 0},
			},
			wantItems: []string{"a", "b", "c"},
		},
		{
			name: "multiple pages with offset pagination",
			pages: [][]string{
				{"a", "b"},
				{"c", "d"},
				{"e"},
			},
			responses: []*Response{
				{NextPage: 2},
				{NextPage: 3},
				{NextPage: 0},
			},
			wantItems: []string{"a", "b", "c", "d", "e"},
		},
		{
			name: "error on first page",
			pages: [][]string{
				nil,
			},
			responses: []*Response{
				nil,
			},
			wantErr: errors.New("api error"),
		},
		{
			name: "error on second page",
			pages: [][]string{
				{"a", "b"},
				nil,
			},
			responses: []*Response{
				{NextPage: 2},
				nil,
			},
			wantErr: errors.New("api error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			pageIdx := 0
			f := func(PaginationOption) ([]string, *Response, error) {
				if pageIdx >= len(tt.pages) {
					t.Fatal("unexpected pagination call")
				}
				if tt.wantErr != nil && pageIdx == len(tt.pages)-1 {
					pageIdx++
					return nil, nil, tt.wantErr
				}
				page := tt.pages[pageIdx]
				resp := tt.responses[pageIdx]
				pageIdx++
				return page, resp, nil
			}

			var got []string
			var err error

			for item, itemErr := range Scan2(f) {
				if itemErr != nil {
					err = itemErr
					break
				}
				got = append(got, item)
			}

			if tt.wantErr != nil {
				if err == nil {
					t.Errorf("want error %v, got nil", tt.wantErr)
				}
				if err.Error() != tt.wantErr.Error() {
					t.Errorf("want error %v, got %v", tt.wantErr, err)
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			want := tt.wantItems
			if !slices.Equal(got, want) {
				t.Errorf("got %v, want %v", got, want)
			}
		})
	}
}

func TestMustIter(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		items     []int
		errorAt   int // position to error at, -1 for no error
		wantItems []int
		wantPanic bool
	}{
		{
			name:      "no error",
			items:     []int{1, 2, 3},
			errorAt:   -1,
			wantItems: []int{1, 2, 3},
		},
		{
			name:      "error on first item",
			items:     []int{1, 2, 3},
			errorAt:   0,
			wantPanic: true,
		},
		{
			name:      "error on second item",
			items:     []int{1, 2, 3},
			errorAt:   1,
			wantPanic: true,
		},
		{
			name:      "empty iterator",
			items:     []int{},
			errorAt:   -1,
			wantItems: []int{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if tt.wantPanic {
				defer func() {
					if r := recover(); r == nil {
						t.Error("expected panic but got none")
					}
				}()
			}

			// Create a Seq2 iterator that yields items with errors at specified position
			it := func(yield func(int, error) bool) {
				for i, item := range tt.items {
					if i == tt.errorAt {
						yield(item, errors.New("test error"))
						return
					}
					if !yield(item, nil) {
						return
					}
				}
			}

			var got []int
			for item := range MustIter(it) {
				got = append(got, item)
			}

			if tt.wantPanic {
				return
			}

			want := tt.wantItems
			if !slices.Equal(got, want) {
				t.Errorf("got %v, want %v", got, want)
			}
		})
	}
}

func TestScanAndCollect(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		pages     [][]float64
		responses []*Response
		wantItems []float64
		wantErr   error
	}{
		{
			name: "single page",
			pages: [][]float64{
				{1.5, 2.5, 3.5},
			},
			responses: []*Response{
				{NextPage: 0},
			},
			wantItems: []float64{1.5, 2.5, 3.5},
		},
		{
			name: "multiple pages",
			pages: [][]float64{
				{1.1, 2.2},
				{3.3, 4.4},
			},
			responses: []*Response{
				{NextPage: 2},
				{NextPage: 0},
			},
			wantItems: []float64{1.1, 2.2, 3.3, 4.4},
		},
		{
			name: "error",
			pages: [][]float64{
				{1.1},
				nil,
			},
			responses: []*Response{
				{NextPage: 2},
				nil,
			},
			wantErr: errors.New("collection failed"),
		},
		{
			name: "empty result",
			pages: [][]float64{
				{},
			},
			responses: []*Response{
				{NextPage: 0},
			},
			wantItems: []float64{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			pageIdx := 0
			f := func(PaginationOption) ([]float64, *Response, error) {
				if pageIdx >= len(tt.pages) {
					t.Fatal("unexpected pagination call")
				}
				if tt.wantErr != nil && pageIdx == len(tt.pages)-1 {
					pageIdx++
					return nil, nil, tt.wantErr
				}
				page := tt.pages[pageIdx]
				resp := tt.responses[pageIdx]
				pageIdx++
				return page, resp, nil
			}

			got, err := ScanAndCollect(f)

			if tt.wantErr != nil {
				if err == nil {
					t.Errorf("want error %v, got nil", tt.wantErr)
				}
				if err.Error() != tt.wantErr.Error() {
					t.Errorf("want error %v, got %v", tt.wantErr, err)
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			want := tt.wantItems
			if !slices.Equal(got, want) {
				t.Errorf("got %v, want %v", got, want)
			}
		})
	}
}
