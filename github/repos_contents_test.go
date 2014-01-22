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
	mux.HandleFunc("/repos/o/r/readme", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{
		  "type": "file",
		  "encoding": "base64",
		  "size": 5362,
		  "name": "README.md",
		  "path": "README.md"
		}`)
	})
	readme, _, err := client.Repositories.GetReadme("o", "r")
	if err != nil {
		t.Errorf("Repositories.GetReadme returned error: %v", err)
	}
	want := &RepositoryContent{Type: String("file"), Name: String("README.md"), Size: Int(5362), Encoding: String("base64"), Path: String("README.md")}
	if !reflect.DeepEqual(readme, want) {
		t.Errorf("Repositories.GetReadme returned %+v, want %+v", readme, want)
	}
}

func TestDecodeBadEncoding(t *testing.T) {
	setup()
	defer teardown()
	r := RepositoryContent{Encoding: String("bad")}
	_, err := r.Decode()
	if err == nil {
		t.Errorf("Should fail to decode non-base64")
	}
}

func TestDecode(t *testing.T) {
	setup()
	defer teardown()
	r := RepositoryContent{Encoding: String("base64"), Content: String("aGVsbG8=")}
	o, err := r.Decode()
	if err != nil {
		t.Errorf("Failed to decode content.")
	}
	want := "hello"
	if string(o) != want {
		t.Errorf("RepositoryContent.Decode returned %+v, want %+v", string(o), want)
	}
}
