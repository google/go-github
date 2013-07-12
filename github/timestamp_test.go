// Copyright 2013 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package github

import (
	"encoding/json"
	"testing"
	"time"
)

var (
	referenceTime = time.Unix(1136239445, 0)
	unixOrigin    = time.Unix(0, 0)
)

func TestMarshal(t *testing.T) {
	testCases := []struct {
		desc    string
		data    TimeStamp
		want    string
		wantErr bool
		equal   bool
	}{
		{"Reference", TimeStamp{referenceTime}, `"2006-01-02T17:04:05-05:00"`, false, true},
		{"Empty", TimeStamp{}, `"0001-01-01T00:00:00Z"`, false, true},
		{"Mismatch", TimeStamp{}, `"2006-01-02T17:04:05-05:00"`, false, false},
	}
	for _, tc := range testCases {
		out, err := json.Marshal(tc.data)
		if gotErr := err != nil; gotErr != tc.wantErr {
			t.Errorf("%s: gotErr=%v, wantErr=%v, err=%v", tc.desc, gotErr, tc.wantErr, err)
		}
		got := string(out)
		equal := got == tc.want
		if (got == tc.want) != tc.equal {
			t.Errorf("%s: got=%s, want=%s, equal=%v", tc.desc, got, tc.want, equal)
		}
	}
}

func TestUnmarshal(t *testing.T) {
	testCases := []struct {
		desc    string
		data    string
		want    TimeStamp
		wantErr bool
		equal   bool
	}{
		{"Reference", `"2006-01-02T17:04:05-05:00"`, TimeStamp{referenceTime}, false, true},
		{"ReferenceUnix", `1136239445`, TimeStamp{referenceTime}, false, true},
		{"Empty", `"0001-01-01T00:00:00Z"`, TimeStamp{}, false, true},
		{"EmptyUnix", `0`, TimeStamp{unixOrigin}, false, true},
		{"Mismatch", `"2006-01-02T17:04:05-05:00"`, TimeStamp{}, false, false},
		{"MismatchUnix", `0`, TimeStamp{}, false, false},
		{"Invalid", `"asdf"`, TimeStamp{referenceTime}, true, false},
	}
	for _, tc := range testCases {
		var got TimeStamp
		err := json.Unmarshal([]byte(tc.data), &got)
		if gotErr := err != nil; gotErr != tc.wantErr {
			t.Errorf("%s: gotErr=%v, wantErr=%v, err=%v", tc.desc, gotErr, tc.wantErr, err)
			continue
		}
		equal := got.Equal(tc.want)
		if equal != tc.equal {
			t.Errorf("%s: got=%#v, want=%#v, equal=%v", tc.desc, got, tc.want, equal)
		}
	}
}

func TestMarshalReflexivity(t *testing.T) {
	testCases := []struct {
		desc string
		data TimeStamp
	}{
		{"Reference", TimeStamp{referenceTime}},
		{"Empty", TimeStamp{}},
	}
	for _, tc := range testCases {
		data, err := json.Marshal(tc.data)
		if err != nil {
			t.Errorf("%s: Marshal err=%v", tc.desc, err)
		}
		var got TimeStamp
		err = json.Unmarshal(data, &got)
		if !got.Equal(tc.data) {
			t.Errorf("%s: %+v != %+v", got, data)
		}
	}
}

type WrappedTimeStamp struct {
	A    int
	Time TimeStamp
}

func TestWrappedMarshal(t *testing.T) {
	testCases := []struct {
		desc    string
		data    WrappedTimeStamp
		want    string
		wantErr bool
		equal   bool
	}{
		{"Reference", WrappedTimeStamp{0, TimeStamp{referenceTime}}, `{"A":0,"Time":"2006-01-02T17:04:05-05:00"}`, false, true},
		{"Empty", WrappedTimeStamp{}, `{"A":0,"Time":"0001-01-01T00:00:00Z"}`, false, true},
		{"Mismatch", WrappedTimeStamp{}, `{"A":0,"Time":"2006-01-02T17:04:05-05:00"}`, false, false},
	}
	for _, tc := range testCases {
		out, err := json.Marshal(tc.data)
		if gotErr := err != nil; gotErr != tc.wantErr {
			t.Errorf("%s: gotErr=%v, wantErr=%v, err=%v", tc.desc, gotErr, tc.wantErr, err)
		}
		got := string(out)
		equal := got == tc.want
		if equal != tc.equal {
			t.Errorf("%s: got=%s, want=%s, equal=%v", tc.desc, got, tc.want, equal)
		}
	}
}

func TestWrappedUnmarshal(t *testing.T) {
	testCases := []struct {
		desc    string
		data    string
		want    WrappedTimeStamp
		wantErr bool
		equal   bool
	}{
		{"Reference", `"2006-01-02T17:04:05-05:00"`, WrappedTimeStamp{0, TimeStamp{referenceTime}}, false, true},
		{"ReferenceUnix", `1136239445`, WrappedTimeStamp{0, TimeStamp{referenceTime}}, false, true},
		{"Empty", `"0001-01-01T00:00:00Z"`, WrappedTimeStamp{0, TimeStamp{}}, false, true},
		{"EmptyUnix", `0`, WrappedTimeStamp{0, TimeStamp{unixOrigin}}, false, true},
		{"Mismatch", `"2006-01-02T17:04:05-05:00"`, WrappedTimeStamp{0, TimeStamp{}}, false, false},
		{"MismatchUnix", `0`, WrappedTimeStamp{0, TimeStamp{}}, false, false},
		{"Invalid", `"asdf"`, WrappedTimeStamp{0, TimeStamp{referenceTime}}, true, false},
	}
	for _, tc := range testCases {
		var got TimeStamp
		err := json.Unmarshal([]byte(tc.data), &got)
		if gotErr := err != nil; gotErr != tc.wantErr {
			t.Errorf("%s: gotErr=%v, wantErr=%v, err=%v", tc.desc, gotErr, tc.wantErr, err)
			continue
		}
		equal := got.Time.Equal(tc.want.Time.Time)
		if equal != tc.equal {
			t.Errorf("%s: got=%#v, want=%#v, equal=%v", tc.desc, got, tc.want, equal)
		}
	}
}

func TestWrappedMarshalReflexivity(t *testing.T) {
	testCases := []struct {
		desc string
		data WrappedTimeStamp
	}{
		{"Reference", WrappedTimeStamp{0, TimeStamp{referenceTime}}},
		{"Empty", WrappedTimeStamp{0, TimeStamp{}}},
	}
	for _, tc := range testCases {
		bytes, err := json.Marshal(tc.data)
		if err != nil {
			t.Errorf("%s: Marshal err=%v", tc.desc, err)
		}
		var got WrappedTimeStamp
		err = json.Unmarshal(bytes, &got)
		if !got.Time.Equal(tc.data.Time) {
			t.Errorf("%s: %+v != %+v", tc.desc, got, tc.data)
		}
	}
}
