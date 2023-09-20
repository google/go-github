package github

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestCodesOfConductService_ListCodesOfConduct(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/codes_of_conduct", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testHeader(t, r, "Accept", mediaTypeCodesOfConductPreview)
		fmt.Fprint(w, `[{
						"key": "key",
						"name": "name",
						"url": "url"}
						]`)
	})

	ctx := context.Background()
	cs, _, err := client.CodesOfConduct.ListCodesOfConduct(ctx)
	if err != nil {
		t.Errorf("ListCodesOfConduct returned error: %v", err)
	}

	want := []*CodeOfConduct{
		{
			Key:  String("key"),
			Name: String("name"),
			URL:  String("url"),
		}}
	if !cmp.Equal(want, cs) {
		t.Errorf("ListCodesOfConduct returned %+v, want %+v", cs, want)
	}

	const methodName = "ListCodesOfConduct"
	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.CodesOfConduct.ListCodesOfConduct(ctx)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestCodesOfConductService_GetCodeOfConduct(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/codes_of_conduct/k", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testHeader(t, r, "Accept", mediaTypeCodesOfConductPreview)
		fmt.Fprint(w, `{
						"key": "key",
						"name": "name",
						"url": "url",
						"body": "body"}`,
		)
	})

	ctx := context.Background()
	coc, _, err := client.CodesOfConduct.GetCodeOfConduct(ctx, "k")
	if err != nil {
		t.Errorf("ListCodesOfConduct returned error: %v", err)
	}

	want := &CodeOfConduct{
		Key:  String("key"),
		Name: String("name"),
		URL:  String("url"),
		Body: String("body"),
	}
	if !cmp.Equal(want, coc) {
		t.Errorf("GetCodeOfConductByKey returned %+v, want %+v", coc, want)
	}

	const methodName = "GetCodeOfConduct"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.CodesOfConduct.GetCodeOfConduct(ctx, "\n")
		return err
	})

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.CodesOfConduct.GetCodeOfConduct(ctx, "k")
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestCodeOfConduct_Marshal(t *testing.T) {
	testJSONMarshal(t, &CodeOfConduct{}, "{}")

	a := &CodeOfConduct{
		Name: String("name"),
		Key:  String("key"),
		URL:  String("url"),
		Body: String("body"),
	}

	want := `{
		"name": "name",
		"key": "key",
		"url": "url",
		"body": "body"
	}`

	testJSONMarshal(t, a, want)
}
