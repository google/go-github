package github

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"testing"
)

func TestActivityService_ListWatchers(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/repos/o/r/subscribers", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testFormValues(t, r, values{
			"page": "2",
		})

		fmt.Fprint(w, `[{"id":1}]`)
	})

	watchers, _, err := client.Activity.ListWatchers("o", "r", &ListOptions{Page: 2})
	if err != nil {
		t.Errorf("Activity.ListWatchers returned error: %v", err)
	}

	want := []User{{ID: Int(1)}}
	if !reflect.DeepEqual(watchers, want) {
		t.Errorf("Activity.ListWatchers returned %+v, want %+v", watchers, want)
	}
}

func TestActivityService_ListWatchedRepositories_authenticatedUser(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/user/subscriptions", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `[{"id":1}]`)
	})

	watched, _, err := client.Activity.ListWatchedRepositories("")
	if err != nil {
		t.Errorf("Activity.ListWatchedRepositories returned error: %v", err)
	}

	want := []Repository{{ID: Int(1)}}
	if !reflect.DeepEqual(watched, want) {
		t.Errorf("Activity.ListWatchedRepositories returned %+v, want %+v", watched, want)
	}
}

func TestActivityService_ListWatchedRepositories_specifiedUser(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/users/u/subscriptions", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `[{"id":1}]`)
	})

	watched, _, err := client.Activity.ListWatchedRepositories("u")
	if err != nil {
		t.Errorf("Activity.ListWatchedRepositories returned error: %v", err)
	}

	want := []Repository{{ID: Int(1)}}
	if !reflect.DeepEqual(watched, want) {
		t.Errorf("Activity.ListWatchedRepositories returned %+v, want %+v", watched, want)
	}
}

func TestActivityService_GetRepositorySubscription_true(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/repos/o/r/subscription", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{"subscribed":true}`)
	})

	sub, _, err := client.Activity.GetRepositorySubscription("o", "r")
	if err != nil {
		t.Errorf("Activity.GetRepositorySubscription returned error: %v", err)
	}

	want := &RepositorySubscription{Subscribed: Bool(true)}
	if !reflect.DeepEqual(sub, want) {
		t.Errorf("Activity.GetRepositorySubscription returned %+v, want %+v", sub, want)
	}
}

func TestActivityService_GetRepositorySubscription_false(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/repos/o/r/subscription", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		w.WriteHeader(http.StatusNotFound)
	})

	sub, _, err := client.Activity.GetRepositorySubscription("o", "r")
	if err != nil {
		t.Errorf("Activity.GetRepositorySubscription returned error: %v", err)
	}

	var want *RepositorySubscription
	if !reflect.DeepEqual(sub, want) {
		t.Errorf("Activity.GetRepositorySubscription returned %+v, want %+v", sub, want)
	}
}

func TestActivityService_GetRepositorySubscription_error(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/repos/o/r/subscription", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		w.WriteHeader(http.StatusBadRequest)
	})

	_, _, err := client.Activity.GetRepositorySubscription("o", "r")
	if err == nil {
		t.Errorf("Expected HTTP 400 response")
	}
}

func TestActivityService_SetRepositorySubscription(t *testing.T) {
	setup()
	defer teardown()

	input := &RepositorySubscription{Subscribed: Bool(true)}

	mux.HandleFunc("/repos/o/r/subscription", func(w http.ResponseWriter, r *http.Request) {
		v := new(RepositorySubscription)
		json.NewDecoder(r.Body).Decode(v)

		testMethod(t, r, "PUT")
		if !reflect.DeepEqual(v, input) {
			t.Errorf("Request body = %+v, want %+v", v, input)
		}

		fmt.Fprint(w, `{"ignored":true}`)
	})

	sub, _, err := client.Activity.SetRepositorySubscription("o", "r", input)
	if err != nil {
		t.Errorf("Activity.SetRepositorySubscription returned error: %v", err)
	}

	want := &RepositorySubscription{Ignored: Bool(true)}
	if !reflect.DeepEqual(sub, want) {
		t.Errorf("Activity.SetRepositorySubscription returned %+v, want %+v", sub, want)
	}
}

func TestActivityService_DeleteRepositorySubscription(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/repos/o/r/subscription", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
	})

	_, err := client.Activity.DeleteRepositorySubscription("o", "r")
	if err != nil {
		t.Errorf("Activity.DeleteRepositorySubscription returned error: %v", err)
	}
}
