// Copyright 2013 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package github

import (
	"testing"
	"time"
)

func TestParameters_withOptionalParameter_withEmptyValue(t *testing.T) {
	// given
	baseUrl := "http://base"
	param := "parameter"
	var value *string = nil

	expected := "http://base"

	// when
	url := withOptionalParameter(baseUrl, param, value)

	//	then
	if url != expected {
		t.Errorf("Generated url '%v' did not equal: %v", url, expected)
	}
}

func TestParameters_withOptionalParameter_withGivenValue(t *testing.T) {
	// given
	url := "http://base"
	param := "parameter"
	value := "value"

	expected := "http://base?parameter=value"

	// when
	url = withOptionalParameter(url, param, &value)

	//	then
	if url != expected {
		t.Errorf("Generated url '%v' did not equal: %v", url, expected)
	}
}

func TestParameters_withOptionalParameter_withMultipleGivenValues(t *testing.T) {
	// given
	url := "http://base"
	param1 := "parameter1"
	value1 := "value1"
	param2 := "parameter2"
	value2 := "value2"

	expected := "http://base?parameter1=value1&parameter2=value2"

	// when
	url = withOptionalParameter(url, param1, &value1)
	url = withOptionalParameter(url, param2, &value2)

	//	then
	if url != expected {
		t.Errorf("Generated url '%v' did not equal: %v", url, expected)
	}
}

func TestParameters_withOptionalTimeParameter_withGivenValue(t *testing.T) {
	// given
	url := "http://base"
	param := "parameter"
	value := time.Date(2013, time.January, 1, 0, 0, 0, 0, time.UTC)
	expected := "http://base?parameter=2013-01-01T00%3A00%3A00Z" // todo verify this is the date format we need...

	// when
	url = withOptionalTimeParameter(url, param, value)

	//	then
	if url != expected {
		t.Errorf("Generated url '%v' did not equal: %v", url, expected)
	}
}
