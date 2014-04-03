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
	want := &RepositoryContent{
		Type:     String("file"),
		Encoding: String("base64"),
		Size:     Int(5362),
		Name:     String("README.md"),
		Path:     String("README.md"),
	}
	if !reflect.DeepEqual(readme, want) {
		t.Errorf("Repositories.GetReadme returned %+v, want %+v", readme, want)
	}
}

func TestRepositoriesService_GetContents(t *testing.T) {
	setup()
	defer teardown()
	mux.HandleFunc("/repos/o/r/contents/package.json", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testFormValues(t, r, values{})
		fmt.Fprint(w, `{
		  "type": "file",
		  "encoding": "base64",
		  "size": 5362,
		  "name": "package.json",
		  "path": "package.json"
		}`)
	})
	content, _, err := client.Repositories.GetContents("o", "r", "package.json", "")
	if err != nil {
		t.Errorf("Repositories.GetContents returned error: %v", err)
	}
	want := &RepositoryContent{
		Type:     String("file"),
		Encoding: String("base64"),
		Size:     Int(5362),
		Name:     String("package.json"),
		Path:     String("package.json"),
	}
	if !reflect.DeepEqual(content, want) {
		t.Errorf("Repositories.GetContents returned %+v, want %+v", content, want)
	}
}

func TestRepositoriesService_GetContentsAtCustomRef(t *testing.T) {
	setup()
	defer teardown()
	mux.HandleFunc("/repos/o/r/contents/package.json", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testFormValues(t, r, values{"ref": "develop"})
		fmt.Fprint(w, `{
		  "type": "file",
		  "encoding": "base64",
		  "size": 5362,
		  "name": "package.json",
		  "path": "package.json"
		}`)
	})
	content, _, err := client.Repositories.GetContents("o", "r", "package.json", "develop")
	if err != nil {
		t.Errorf("Repositories.GetContents returned error: %v", err)
	}
	want := &RepositoryContent{
		Type:     String("file"),
		Encoding: String("base64"),
		Size:     Int(5362),
		Name:     String("package.json"),
		Path:     String("package.json"),
	}
	if !reflect.DeepEqual(content, want) {
		t.Errorf("Repositories.GetContents returned %+v, want %+v", content, want)
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
