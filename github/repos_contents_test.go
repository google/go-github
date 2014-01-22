package github

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
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
		fmt.Fprint(w, `{"name":"n", "_links": {"self": "s"}}`)
	})

	comments, _, err := client.Repositories.GetReadme("o", "r", &RefOption{})
	if err != nil {
		t.Errorf("Repositories.GetReadme returned error: %v", err)
	}

	want := &RepositoryFile{
		Name:  String("n"),
		Links: &RepositoryFileLinks{Self: String("s")},
	}
	if !reflect.DeepEqual(comments, want) {
		t.Errorf("Repositories.GetReadme returned %+v, want %+v", comments, want)
	}
}

func TestRepositoriesService_GetReadme_withRef(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/repos/o/r/readme", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testFormValues(t, r, values{
			"ref": "r",
		})
		fmt.Fprint(w, `{"name":"n", "_links": {"self": "s"}}`)
	})

	comments, _, err := client.Repositories.GetReadme("o", "r", &RefOption{"r"})
	if err != nil {
		t.Errorf("Repositories.GetReadme_withRef returned error: %v", err)
	}

	want := &RepositoryFile{
		Name:  String("n"),
		Links: &RepositoryFileLinks{Self: String("s")},
	}
	if !reflect.DeepEqual(comments, want) {
		t.Errorf("Repositories.GetReadme_withRef returned %+v, want %+v", comments, want)
	}
}

func TestRepositoriesService_GetContents_singleFile(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/repos/o/r/contents/", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{"name":"n", "_links": {"self": "s"}}`)
	})

	comments, _, err := client.Repositories.GetContents("o", "r", &RepositoryContentsOptions{})
	if err != nil {
		t.Errorf("Repositories.GetContents returned error: %v", err)
	}

	want := []RepositoryFile{{
		Name:  String("n"),
		Links: &RepositoryFileLinks{Self: String("s")},
	}}
	if !reflect.DeepEqual(comments, want) {
		t.Errorf("Repositories.GetContents returned %+v, want %+v", comments, want)
	}
}

func TestRepositoriesService_GetContents_multipleFile(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/repos/o/r/contents/", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `[{"name":"n", "_links": {"self": "s"}}]`)
	})

	contents, _, err := client.Repositories.GetContents("o", "r", &RepositoryContentsOptions{})
	if err != nil {
		t.Errorf("Repositories.GetContents_multipleFile returned error: %v", err)
	}

	want := []RepositoryFile{{
		Name:  String("n"),
		Links: &RepositoryFileLinks{Self: String("s")},
	}}
	if !reflect.DeepEqual(contents, want) {
		t.Errorf("Repositories.GetContents_multipleFile returned %+v, want %+v", contents, want)
	}
}

func TestRepositoriesService_GetContents_singleFileWithPath(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/repos/o/r/contents/p", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{"name":"n", "size": 5, "_links": {"self": "s"}}`)
	})

	contents, _, err := client.Repositories.GetContents("o", "r", &RepositoryContentsOptions{Path: "p"})
	if err != nil {
		t.Errorf("Repositories.GetContents returned error: %v", err)
	}

	want := []RepositoryFile{{
		Name:  String("n"),
		Size:  Int(5),
		Links: &RepositoryFileLinks{Self: String("s")},
	}}
	if !reflect.DeepEqual(contents, want) {
		t.Errorf("Repositories.GetContents returned %+v, want %+v", contents, want)
	}
}

func TestRepositoriesService_CreateFile(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/repos/o/r/contents/p", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PUT")
		buf := new(bytes.Buffer)
		buf.ReadFrom(r.Body)
		req := struct{ Content string }{}
		json.Unmarshal(buf.Bytes(), &req)

		want := base64.StdEncoding.EncodeToString([]byte("c"))
		if req.Content != want {
			t.Errorf("Repositories.CreateFile contents base64-encoded as %v, want %v", req.Content, want)
		}
		fmt.Fprint(w, `{"content":{"name":"n"},"commit":{"sha":"a","url":"u"}}`)
	})

	content := []byte("c")
	repositoryCreateFileOptions := &RepositoryCreateFileOptions{
		FileOptions: FileOptions{
			Path:   "p",
			Author: &CommitAuthor{Name: String("n"), Email: String("e")},
		},
		Content: &content,
	}

	createResp, _, err := client.Repositories.CreateFile("o", "r", repositoryCreateFileOptions)
	if err != nil {
		t.Errorf("Repositories.CreateFile returned error: %v", err)
	}

	want := &RepositoryContentResponse{
		Content: &RepositoryFile{Name: String("n")},
		Commit: &CommitWithURL{
			URL:    String("u"),
			Commit: Commit{SHA: String("a")}},
	}
	if !reflect.DeepEqual(createResp, want) {
		t.Errorf("Repositories.CreateFile returned %+v, want %+v", createResp, want)
	}
}

func TestRepositoriesService_UpdateFile(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/repos/o/r/contents/p", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PUT")
		buf := new(bytes.Buffer)
		buf.ReadFrom(r.Body)
		req := struct{ Content string }{}
		json.Unmarshal(buf.Bytes(), &req)

		want := base64.StdEncoding.EncodeToString([]byte("c"))
		if req.Content != want {
			t.Errorf("Repositories.CreateFile contents base64-encoded as %v, want %v", req.Content, want)
		}
		fmt.Fprint(w, `{"content":{"name":"n"},"commit":{"sha":"a","url":"u"}}`)
	})

	content := []byte("c")
	repositoryUpdateFileOptions := &RepositoryUpdateFileOptions{
		FileOptions: FileOptions{
			Path:   "p",
			Author: &CommitAuthor{Name: String("n"), Email: String("e")},
		},
		SHA:     "s",
		Content: &content,
	}

	updateResp, _, err := client.Repositories.UpdateFile("o", "r", repositoryUpdateFileOptions)
	if err != nil {
		t.Errorf("Repositories.UpdateFile returned error: %v", err)
	}

	want := &RepositoryContentResponse{
		Content: &RepositoryFile{Name: String("n")},
		Commit: &CommitWithURL{
			URL:    String("u"),
			Commit: Commit{SHA: String("a")}},
	}
	if !reflect.DeepEqual(updateResp, want) {
		t.Errorf("Repositories.CreateFile returned %+v, want %+v", updateResp, want)
	}
}

func TestRepositoriesService_DeleteFile(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/repos/o/r/contents/p", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
		fmt.Fprint(w, `{"content":null,"commit":{"sha":"a","url":"u"}}`)
	})

	repositoryDeleteFileOptions := &RepositoryDeleteFileOptions{
		FileOptions: FileOptions{
			Path:   "p",
			Author: &CommitAuthor{Name: String("n"), Email: String("e")},
		},
		SHA: "s",
	}

	deleteResp, _, err := client.Repositories.DeleteFile("o", "r", repositoryDeleteFileOptions)
	if err != nil {
		t.Errorf("Repositories.DeleteFile returned error: %v", err)
	}

	want := &RepositoryContentResponse{
		Content: nil,
		Commit: &CommitWithURL{
			URL:    String("u"),
			Commit: Commit{SHA: String("a")}},
	}
	if !reflect.DeepEqual(deleteResp, want) {
		t.Errorf("Repositories.DeleteFile returned %+v, want %+v", deleteResp, want)
	}
}

func TestRepositoriesService_GetArchiveLink(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/repos/o/r/z/", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		http.Redirect(w, r, "http://github.com/a", http.StatusFound)
	})

	url, resp, err := client.Repositories.GetArchiveLink("o", "r", "z", &RefOption{})
	_, _, _ = url, resp, err

	if err != nil {
		t.Errorf("Repositories.GetArchiveLink returned error: %v", err)
	}

	if resp.StatusCode != http.StatusFound {
		t.Errorf("Repositories.GetArchiveLink returned status: %d, want %d", resp.StatusCode, http.StatusFound)
	}

	want := "http://github.com/a"
	if url.String() != want {
		t.Errorf("Repositories.GetArchiveLink returned %+v, want %+v", url.String(), want)
	}
}
