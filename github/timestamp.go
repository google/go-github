// Copyright 2013 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package github

import (
	"strconv"
	"time"
)

// TimeStamp represents a time that can be imported in RFC3339 or Unix
// formatting. This is necessary for some fields since the GitHub API is
// inconsistent in how it represents times. All exported methods of time.Time
// can be called on TimeStamp.
type TimeStamp struct {
	time.Time
}

// UnmarshalJSON implements the json.Unmarshaler interface.
// Time is expected in RFC3339 or Unix format.
func (t *TimeStamp) UnmarshalJSON(data []byte) (err error) {
	str := string(data)
	i, err := strconv.ParseInt(str, 10, 64)
	if err == nil {
		(*t).Time = time.Unix(i, 0)
	} else {
		(*t).Time, err = time.Parse(`"`+time.RFC3339+`"`, str)
	}
	return
}

// Equal reports whether t and u are equal based on time.Equal
func (t TimeStamp) Equal(u TimeStamp) bool {
	return t.Time.Equal(u.Time)
}
