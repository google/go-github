package github

import (
	"fmt"
	"net/http"
	"reflect"
	"testing"
)

func TestRepositoriesService_GetReadme(t *testing.T) {
	setup()
	defer teardown()
	mux.HandleFunc("repos/o/r/readme", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{"type": "file", "name": "readme.md"}`)
		readme, _, err := client.Repositories.GetReadme("o", "r")
		if err != nil {
			t.Errorf("Repositories.GetReadme returned error: %v", err)
		}
		want := RepositoryContent{Type: String("file"), Name: String("readme.md")}
		if !reflect.DeepEqual(readme, want) {
			t.Errorf("Repositories.GetReadme returned %+v, want %+v", readme, want)
		}
	})
}
