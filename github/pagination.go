// Copyright 2026 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package github

import (
	"fmt"
	"iter"
	"slices"
)

// PaginationOption represents a pagination option for requests.
type PaginationOption = RequestOption

// Scan scans all pages for the given request function f and returns individual items in an iterator.
// If an error happens during pagination, the iterator stops immediately.
// The caller must consume the returned error function to retrieve potential errors.
func Scan[T any](f func(PaginationOption) ([]T, *Response, error)) (iter.Seq[T], func() error) {
	exhausted := false
	var e error
	it := func(yield func(T) bool) {
		defer func() {
			exhausted = true
		}()
		for t, err := range Scan2(f) {
			if err != nil {
				e = err
				return
			}

			if !yield(t) {
				return
			}
		}
	}
	hasErr := func() error {
		if !exhausted {
			panic("called error function of Scan iterator before iterator was exhausted")
		}
		return e
	}
	return it, hasErr
}

// Scan2 scans all pages for the given request function f and returns individual items and potential errors in an iterator.
// The caller must consume the error element of the iterator during each iteration
// to ensure that no errors happened.
func Scan2[T any](f func(PaginationOption) ([]T, *Response, error)) iter.Seq2[T, error] {
	return func(yield func(T, error) bool) {
		var nextOpt PaginationOption

	Pagination:
		for {
			ts, resp, err := f(nextOpt)
			if err != nil {
				var t T
				yield(t, err)
				return
			}

			for _, t := range ts {
				if !yield(t, nil) {
					return
				}
			}

			// the f request function was either configured for offset- or cursor-based pagination.
			switch {
			case resp.NextPage != 0:
				nextOpt = WithOffsetPagination(resp.NextPage)
			case resp.After != "":
				nextOpt = WithAfterPagination(resp.After)
			default:
				// no more pages
				break Pagination
			}
		}
	}
}

// MustIter provides a single item iterator for the provided two item iterator and panics if an error happens.
func MustIter[T any](it iter.Seq2[T, error]) iter.Seq[T] {
	return func(yield func(T) bool) {
		for x, err := range it {
			if err != nil {
				panic(fmt.Errorf("iterator produced an error: %w", err))
			}

			if !yield(x) {
				return
			}
		}
	}
}

// ScanAndCollect is a convenience function that collects all results and returns them as slice as well as an error if one happens.
func ScanAndCollect[T any](f func(p PaginationOption) ([]T, *Response, error)) ([]T, error) {
	it, hasErr := Scan(f)
	allItems := slices.Collect(it)
	if err := hasErr(); err != nil {
		return nil, err
	}
	return allItems, nil
}
