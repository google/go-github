package github

import (
	"fmt"
	"net/http"
	"reflect"
	"testing"
)

const README_CONTENT = `{
  "type": "file",
  "encoding": "base64",
  "size": 5362,
  "name": "README.md",
  "path": "README.md"
}`

func TestRepositoriesService_GetReadme(t *testing.T) {
	setup()
	defer teardown()
	mux.HandleFunc("/repos/o/r/readme", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, README_CONTENT)
	})
	readme, _, err := client.Repositories.GetReadme("o", "r")
	if err != nil {
		t.Errorf("Repositories.GetReadme returned error: %v", err)
	}
	want := RepositoryContent{Type: String("file"), Name: String("README.md"), Size: Int(5362), Encoding: String("base64"), Path: String("README.md")}
	if !reflect.DeepEqual(readme, want) {
		t.Errorf("Repositories.GetReadme returned %+v, want %+v", readme, want)
	}
}
