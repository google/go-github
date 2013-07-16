// Copyright 2013 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package github

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"
)

const (
	referenceStr     = `"2006-01-02T15:04:05Z"`
	emptyStr         = `"0001-01-01T00:00:00Z"`
	referenceUnixStr = `1136214245`
)

var (
	referenceTime = time.Date(2006, 01, 02, 15, 04, 05, 0, time.UTC)
	unixOrigin    = time.Unix(0, 0).In(time.UTC)
)

func TestMarshal(t *testing.T) {
	testCases := []struct {
		desc    string
		data    TimeStamp
		want    string
		wantErr bool
		equal   bool
	}{
		{"Reference", TimeStamp{referenceTime}, referenceStr, false, true},
		{"Empty", TimeStamp{}, emptyStr, false, true},
		{"Mismatch", TimeStamp{}, referenceStr, false, false},
	}
	for _, tc := range testCases {
		out, err := json.Marshal(tc.data)
		if gotErr := err != nil; gotErr != tc.wantErr {
			t.Errorf("%s: gotErr=%v, wantErr=%v, err=%v", tc.desc, gotErr, tc.wantErr, err)
		}
		got := string(out)
		equal := got == tc.want
		if (got == tc.want) != tc.equal {
			t.Errorf("%s: got=%s, want=%s, equal=%v, want=%v", tc.desc, got, tc.want, equal, tc.equal)
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
		{"Reference", referenceStr, TimeStamp{referenceTime}, false, true},
		{"ReferenceUnix", `1136214245`, TimeStamp{referenceTime}, false, true},
		{"Empty", emptyStr, TimeStamp{}, false, true},
		{"UnixStart", `0`, TimeStamp{unixOrigin}, false, true},
		{"Mismatch", referenceStr, TimeStamp{}, false, false},
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
			t.Errorf("%s: got=%#v, want=%#v, equal=%v, want=%v", tc.desc, got, tc.want, equal, tc.equal)
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
		{"Reference", WrappedTimeStamp{0, TimeStamp{referenceTime}}, fmt.Sprintf(`{"A":0,"Time":%s}`, referenceStr), false, true},
		{"Empty", WrappedTimeStamp{}, fmt.Sprintf(`{"A":0,"Time":%s}`, emptyStr), false, true},
		{"Mismatch", WrappedTimeStamp{}, fmt.Sprintf(`{"A":0,"Time":%s}`, referenceStr), false, false},
	}
	for _, tc := range testCases {
		out, err := json.Marshal(tc.data)
		if gotErr := err != nil; gotErr != tc.wantErr {
			t.Errorf("%s: gotErr=%v, wantErr=%v, err=%v", tc.desc, gotErr, tc.wantErr, err)
		}
		got := string(out)
		equal := got == tc.want
		if equal != tc.equal {
			t.Errorf("%s: got=%s, want=%s, equal=%v, want=%v", tc.desc, got, tc.want, equal, tc.equal)
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
		{"Reference", referenceStr, WrappedTimeStamp{0, TimeStamp{referenceTime}}, false, true},
		{"ReferenceUnix", referenceUnixStr, WrappedTimeStamp{0, TimeStamp{referenceTime}}, false, true},
		{"Empty", emptyStr, WrappedTimeStamp{0, TimeStamp{}}, false, true},
		{"UnixStart", `0`, WrappedTimeStamp{0, TimeStamp{unixOrigin}}, false, true},
		{"Mismatch", referenceStr, WrappedTimeStamp{0, TimeStamp{}}, false, false},
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
			t.Errorf("%s: got=%#v, want=%#v, equal=%v, want=%v", tc.desc, got, tc.want, equal, tc.equal)
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
