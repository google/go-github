// Copyright 2013 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package github

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"testing"
)

func TestHooksService_Create(t *testing.T) {
	setup()
	defer teardown()

	input := &Hook{Name: "t"}

	mux.HandleFunc("/repos/o/r/hooks", func(w http.ResponseWriter, r *http.Request) {
		v := new(Hook)
		json.NewDecoder(r.Body).Decode(v)

		testMethod(t, r, "POST")
		if !reflect.DeepEqual(v, input) {
			t.Errorf("Request body = %+v, want %+v", v, input)
		}

		fmt.Fprint(w, `{"id":1}`)
	})

	hook, err := client.Hooks.Create("o", "r", input)
	if err != nil {
		t.Errorf("Hooks.Create returned error: %v", err)
	}

	want := &Hook{Id: 1}
	if !reflect.DeepEqual(hook, want) {
		t.Errorf("Hooks.Create returned %+v, want %+v", hook, want)
	}
}

func TestHooksService_List(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/repos/o/r/hooks", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `[{"id":1}, {"id":2}]`)
	})

	hooks, err := client.Hooks.List("o", "r")
	if err != nil {
		t.Errorf("Hooks.List returned error: %v", err)
	}

	want := []Hook{Hook{Id: 1}, Hook{Id: 2}}
	if !reflect.DeepEqual(hooks, want) {
		t.Errorf("Hooks.List returned %+v, want %+v", hooks, want)
	}
}

func TestHooksService_Get(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/repos/o/r/hooks/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{"id":1}`)
	})

	hook, err := client.Hooks.Get("o", "r", 1)
	if err != nil {
		t.Errorf("Hooks.Get returned error: %v", err)
	}

	want := &Hook{Id: 1}
	if !reflect.DeepEqual(hook, want) {
		t.Errorf("Hooks.Get returned %+v, want %+v", hook, want)
	}
}

func TestHooksService_Edit(t *testing.T) {
	setup()
	defer teardown()

	input := &Hook{Name: "t"}

	mux.HandleFunc("/repos/o/r/hooks/1", func(w http.ResponseWriter, r *http.Request) {
		v := new(Hook)
		json.NewDecoder(r.Body).Decode(v)

		testMethod(t, r, "PATCH")
		if !reflect.DeepEqual(v, input) {
			t.Errorf("Request body = %+v, want %+v", v, input)
		}

		fmt.Fprint(w, `{"id":1}`)
	})

	hook, err := client.Hooks.Edit("o", "r", 1, input)
	if err != nil {
		t.Errorf("Hooks.Edit returned error: %v", err)
	}

	want := &Hook{Id: 1}
	if !reflect.DeepEqual(hook, want) {
		t.Errorf("Hooks.Edit returned %+v, want %+v", hook, want)
	}
}

func TestHooksService_Delete(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/repos/o/r/hooks/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
	})

	err := client.Hooks.Delete("o", "r", 1)
	if err != nil {
		t.Errorf("Hooks.Delete returned error: %v", err)
	}
}

func TestHooksService_Test(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/repos/o/r/hooks/1/tests", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
	})

	err := client.Hooks.Test("o", "r", 1)
	if err != nil {
		t.Errorf("Hooks.Test returned error: %v", err)
	}
}
