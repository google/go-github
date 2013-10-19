package github

import (
  "testing"
  "net/http"
  "reflect"
  "fmt"
)

func TestActivityService_ListWatchers(t *testing.T) {
  setup()
  defer teardown()

  mux.HandleFunc("/repos/google/go-github/subscribers", func(w http.ResponseWriter, r *http.Request) {
    testMethod(t, r, "GET")
    fmt.Fprint(w, `[{"id":1}]`)
  })

  watchers, _, err := client.Activity.ListWatchers("google", "go-github")
  if err != nil {
    t.Errorf("Activity.ListWatchers returned error: %v", err)
  }
  want := []User{{ID: Int(1)}}
  if !reflect.DeepEqual(watchers, want) {
    t.Errorf("Activity.ListWatchers returned %+v, want %+v", watchers, want)
  }
}
