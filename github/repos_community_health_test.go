// Copyright 2017 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package github

import (
	"context"
	"fmt"
	"net/http"
	"reflect"
	"testing"
)

func TestRepositoriesService_GetCommunityHealthMetrics(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/repositories/o/r/community/profile", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testHeader(t, r, "Accept", mediaTypeRepositoryCommunityHealthMetricsPreview)
		fmt.Fprintf(w, `{"health_percentage":75}`)
	})

	got, _, err := client.Repositories.GetCommunityHealthMetrics(context.Background(), "o", "r")
	if err != nil {
		t.Errorf("Repositories.GetCommunityHealthMetrics returned error: %v", err)
	}

	want := &CommunityHealthMetrics{HealthPercentage: Int(75)}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Repositories.GetCommunityHealthMetrics = %+v, want %+v", got, want)
	}
}
