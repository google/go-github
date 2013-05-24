// Copyright 2013 Google. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file or at
// https://developers.google.com/open-source/licenses/bsd

package github

import (
  "encoding/json"
  "fmt"
  "net/http"
  "reflect"
  "testing"
)

func TestUsersService_Get_authenticatedUser(t *testing.T) {
  setup()
  defer teardown()

  mux.HandleFunc("/user", func(w http.ResponseWriter, r *http.Request) {
    if r.Method != "GET" {
      t.Errorf("Request method = %v, want %v", r.Method, "GET")
    }
    fmt.Fprint(w, `{"id":1}`)
  })

  user, err := client.Users.Get("")
  if err != nil {
    t.Errorf("Users.Get returned error: %v", err)
  }

  want := &User{ID: 1}
  if !reflect.DeepEqual(user, want) {
    t.Errorf("Users.Get returned %+v, want %+v", user, want)
  }
}

func TestUsersService_Get_specifiedUser(t *testing.T) {
  setup()
  defer teardown()

  mux.HandleFunc("/users/u", func(w http.ResponseWriter, r *http.Request) {
    if r.Method != "GET" {
      t.Errorf("Request method = %v, want %v", r.Method, "GET")
    }
    fmt.Fprint(w, `{"id":1}`)
  })

  user, err := client.Users.Get("u")
  if err != nil {
    t.Errorf("Users.Get returned error: %v", err)
  }

  want := &User{ID: 1}
  if !reflect.DeepEqual(user, want) {
    t.Errorf("Users.Get returned %+v, want %+v", user, want)
  }
}

func TestUsersService_Edit(t *testing.T) {
  setup()
  defer teardown()

  input := &User{Name: "n"}

  mux.HandleFunc("/user", func(w http.ResponseWriter, r *http.Request) {
    v := new(User)
    json.NewDecoder(r.Body).Decode(v)

    if r.Method != "PATCH" {
      t.Errorf("Request method = %v, want %v", r.Method, "GET")
    }
    if !reflect.DeepEqual(v, input) {
      t.Errorf("Request body = %+v, want %+v", v, input)
    }

    fmt.Fprint(w, `{"id":1}`)
  })

  user, err := client.Users.Edit(input)
  if err != nil {
    t.Errorf("Users.Edit returned error: %v", err)
  }

  want := &User{ID: 1}
  if !reflect.DeepEqual(user, want) {
    t.Errorf("Users.Edit returned %+v, want %+v", user, want)
  }
}

func TestUsersService_List(t *testing.T) {
  setup()
  defer teardown()

  mux.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {
    var v string
    if r.Method != "GET" {
      t.Errorf("Request method = %v, want %v", r.Method, "GET")
    }
    if v = r.FormValue("since"); v != "1" {
      t.Errorf("Request since parameter = %v, want %v", v, "1")
    }
    fmt.Fprint(w, `[{"id":2}]`)
  })

  opt := &UserListOptions{1}
  users, err := client.Users.List(opt)
  if err != nil {
    t.Errorf("Users.Get returned error: %v", err)
  }

  want := []User{User{ID: 2}}
  if !reflect.DeepEqual(users, want) {
    t.Errorf("Users.List returned %+v, want %+v", users, want)
  }
}
