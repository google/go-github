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
	emptyTimeStr                     = `"0001-01-01T00:00:00Z"`
	referenceTimeStr                 = `"2006-01-02T15:04:05Z"`
	referenceTimeStrFractional       = `"2006-01-02T15:04:05.000Z"` // This format was returned by the Projects API before October 1, 2017.
	referenceUnixTimeStr             = `1136214245`
	referenceUnixTimeStrMilliSeconds = `1136214245000` // Millisecond-granular timestamps were introduced in the Audit log API.
)

var (
	referenceTime = time.Date(2006, time.January, 02, 15, 04, 05, 0, time.UTC)
	unixOrigin    = time.Unix(0, 0).In(time.UTC)
)

func TestTimestamp_Marshal(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		desc    string
		data    Timestamp
		want    string
		wantErr bool
		equal   bool
	}{
		{"Reference", Timestamp{referenceTime}, referenceTimeStr, false, true},
		{"Empty", Timestamp{}, emptyTimeStr, false, true},
		{"Mismatch", Timestamp{}, referenceTimeStr, false, false},
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

func TestTimestamp_Unmarshal(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		desc    string
		data    string
		want    Timestamp
		wantErr bool
		equal   bool
	}{
		{"Reference", referenceTimeStr, Timestamp{referenceTime}, false, true},
		{"ReferenceUnix", referenceUnixTimeStr, Timestamp{referenceTime}, false, true},
		{"ReferenceUnixMillisecond", referenceUnixTimeStrMilliSeconds, Timestamp{referenceTime}, false, true},
		{"ReferenceFractional", referenceTimeStrFractional, Timestamp{referenceTime}, false, true},
		{"Empty", emptyTimeStr, Timestamp{}, false, true},
		{"UnixStart", `0`, Timestamp{unixOrigin}, false, true},
		{"Mismatch", referenceTimeStr, Timestamp{}, false, false},
		{"MismatchUnix", `0`, Timestamp{}, false, false},
		{"Invalid", `"asdf"`, Timestamp{referenceTime}, true, false},
		{"OffByMillisecond", `1136214245001`, Timestamp{referenceTime}, false, false},
	}
	for _, tc := range testCases {
		var got Timestamp
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

func TestTimestamp_MarshalReflexivity(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		desc string
		data Timestamp
	}{
		{"Reference", Timestamp{referenceTime}},
		{"Empty", Timestamp{}},
	}
	for _, tc := range testCases {
		data, err := json.Marshal(tc.data)
		if err != nil {
			t.Errorf("%s: Marshal err=%v", tc.desc, err)
		}
		var got Timestamp
		err = json.Unmarshal(data, &got)
		if err != nil {
			t.Errorf("%s: Unmarshal err=%v", tc.desc, err)
		}
		if !got.Equal(tc.data) {
			t.Errorf("%s: %+v != %+v", tc.desc, got, data)
		}
	}
}

type WrappedTimestamp struct {
	A    int
	Time Timestamp
}

func TestWrappedTimestamp_Marshal(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		desc    string
		data    WrappedTimestamp
		want    string
		wantErr bool
		equal   bool
	}{
		{"Reference", WrappedTimestamp{0, Timestamp{referenceTime}}, fmt.Sprintf(`{"A":0,"Time":%s}`, referenceTimeStr), false, true},
		{"Empty", WrappedTimestamp{}, fmt.Sprintf(`{"A":0,"Time":%s}`, emptyTimeStr), false, true},
		{"Mismatch", WrappedTimestamp{}, fmt.Sprintf(`{"A":0,"Time":%s}`, referenceTimeStr), false, false},
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

func TestWrappedTimestamp_Unmarshal(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		desc    string
		data    string
		want    WrappedTimestamp
		wantErr bool
		equal   bool
	}{
		{"Reference", referenceTimeStr, WrappedTimestamp{0, Timestamp{referenceTime}}, false, true},
		{"ReferenceUnix", referenceUnixTimeStr, WrappedTimestamp{0, Timestamp{referenceTime}}, false, true},
		{"ReferenceUnixMillisecond", referenceUnixTimeStrMilliSeconds, WrappedTimestamp{0, Timestamp{referenceTime}}, false, true},
		{"Empty", emptyTimeStr, WrappedTimestamp{0, Timestamp{}}, false, true},
		{"UnixStart", `0`, WrappedTimestamp{0, Timestamp{unixOrigin}}, false, true},
		{"Mismatch", referenceTimeStr, WrappedTimestamp{0, Timestamp{}}, false, false},
		{"MismatchUnix", `0`, WrappedTimestamp{0, Timestamp{}}, false, false},
		{"Invalid", `"asdf"`, WrappedTimestamp{0, Timestamp{referenceTime}}, true, false},
		{"OffByMillisecond", `1136214245001`, WrappedTimestamp{0, Timestamp{referenceTime}}, false, false},
	}
	for _, tc := range testCases {
		var got Timestamp
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

func TestTimestamp_GetTime(t *testing.T) {
	t.Parallel()
	var t1 *Timestamp
	if t1.GetTime() != nil {
		t.Errorf("nil timestamp should return nil, got: %v", t1.GetTime())
	}
	t1 = &Timestamp{referenceTime}
	if !t1.GetTime().Equal(referenceTime) {
		t.Errorf("want reference time, got: %s", t1.GetTime().String())
	}
}

func TestWrappedTimestamp_MarshalReflexivity(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		desc string
		data WrappedTimestamp
	}{
		{"Reference", WrappedTimestamp{0, Timestamp{referenceTime}}},
		{"Empty", WrappedTimestamp{0, Timestamp{}}},
	}
	for _, tc := range testCases {
		bytes, err := json.Marshal(tc.data)
		if err != nil {
			t.Errorf("%s: Marshal err=%v", tc.desc, err)
		}
		var got WrappedTimestamp
		err = json.Unmarshal(bytes, &got)
		if err != nil {
			t.Errorf("%s: Unmarshal err=%v", tc.desc, err)
		}
		if !got.Time.Equal(tc.data.Time) {
			t.Errorf("%s: %+v != %+v", tc.desc, got, tc.data)
		}
	}
}
