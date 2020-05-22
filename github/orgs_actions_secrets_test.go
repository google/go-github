package github

import (
	"context"
	"fmt"
	"net/http"
	"reflect"
	"testing"
	"time"
)

func TestOrganizationsService_GetPublicKey(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/orgs/o/actions/secrets/public-key", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{"key_id":"012345678912345678","key":"2Sg8iYjAxxmI2LvUXpJjkYrMxURPc8r+dB7TJyvv1234"}`)
	})

	key, _, err := client.Organizations.GetPublicKey(context.Background(), "o")
	if err != nil {
		t.Errorf("Organizations.GetPublicKey returned error: %v", err)
	}

	want := &OrganizationPublicKey{KeyID: String("012345678912345678"), Key: String("2Sg8iYjAxxmI2LvUXpJjkYrMxURPc8r+dB7TJyvv1234")}
	if !reflect.DeepEqual(key, want) {
		t.Errorf("Organizations.GetPublicKey returned %+v, want %+v", key, want)
	}
}

func TestOrganizationsService_ListSecrets(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/orgs/o/actions/secrets", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testFormValues(t, r, values{"per_page": "2", "page": "2"})
		fmt.Fprint(w, `{"total_count":3,"secrets":[{"name":"GIST_ID","created_at":"2019-08-10T14:59:22Z","updated_at":"2020-01-10T14:59:22Z","visibility":"private"},{"name":"DEPLOY_TOKEN","created_at":"2019-08-10T14:59:22Z","updated_at":"2020-01-10T14:59:22Z","visibility":"all"},{"name":"GH_TOKEN","created_at":"2019-08-10T14:59:22Z","updated_at":"2020-01-10T14:59:22Z","visibility":"selected","selected_repositories_url":"https://api.github.com/orgs/octo-org/actions/secrets/SUPER_SECRET/repositories"}]}`)
	})

	opts := &ListOptions{Page: 2, PerPage: 2}
	secrets, _, err := client.Organizations.ListSecrets(context.Background(), "o", opts)
	if err != nil {
		t.Errorf("Organizations.ListSecrets returned error: %v", err)
	}

	want := &OrganizationSecrets{
		TotalCount: 3,
		Secrets: []*OrganizationSecret{
			{Name: "GIST_ID", CreatedAt: Timestamp{time.Date(2019, time.August, 10, 14, 59, 22, 0, time.UTC)}, UpdatedAt: Timestamp{time.Date(2020, time.January, 10, 14, 59, 22, 0, time.UTC)}, Visibility: "private"},
			{Name: "DEPLOY_TOKEN", CreatedAt: Timestamp{time.Date(2019, time.August, 10, 14, 59, 22, 0, time.UTC)}, UpdatedAt: Timestamp{time.Date(2020, time.January, 10, 14, 59, 22, 0, time.UTC)}, Visibility: "all"},
			{Name: "GH_TOKEN", CreatedAt: Timestamp{time.Date(2019, time.August, 10, 14, 59, 22, 0, time.UTC)}, UpdatedAt: Timestamp{time.Date(2020, time.January, 10, 14, 59, 22, 0, time.UTC)}, Visibility: "selected", SelectedRepositoriesUrl: "https://api.github.com/orgs/octo-org/actions/secrets/SUPER_SECRET/repositories"},
		},
	}
	if !reflect.DeepEqual(secrets, want) {
		t.Errorf("Organizations.ListSecrets returned %+v, want %+v", secrets, want)
	}
}

func TestOrganizationsService_GetSecret(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/orgs/o/actions/secrets/GH_TOKEN", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{"name":"GH_TOKEN","created_at":"2019-08-10T14:59:22Z","updated_at":"2020-01-10T14:59:22Z","visibility":"selected","selected_repositories_url":"https://api.github.com/orgs/octo-org/actions/secrets/SUPER_SECRET/repositories"}`)
	})

	secret, _, err := client.Organizations.GetSecret(context.Background(), "o", "GH_TOKEN")
	if err != nil {
		t.Errorf("Organizations.GetSecret returned error: %v", err)
	}

	want := &OrganizationSecret{
		Name:                    "GH_TOKEN",
		CreatedAt:               Timestamp{time.Date(2019, time.August, 10, 14, 59, 22, 0, time.UTC)},
		UpdatedAt:               Timestamp{time.Date(2020, time.January, 10, 14, 59, 22, 0, time.UTC)},
		Visibility:              "selected",
		SelectedRepositoriesUrl: "https://api.github.com/orgs/octo-org/actions/secrets/SUPER_SECRET/repositories",
	}

	if !reflect.DeepEqual(secret, want) {
		t.Errorf("Organizations.GetSecret returned %+v, want %+v", secret, want)
	}
}

func TestOrgnizationsService_CreateOrUpdateSecret(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/orgs/o/actions/secrets/NAME", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PUT")
		testHeader(t, r, "Content-Type", "application/json")
		testBody(t, r, `{"key_id":"1234","encrypted_value":"QIv=","visibility":"selected","selected_repository_ids":["A","B","C"]}`+"\n")
	})

	input := &OrganizationEncryptedSecret{
		Name:                  "NAME",
		EncryptedValue:        "QIv=",
		KeyID:                 "1234",
		Visibility:            "selected",
		SelectedRepositoryIDs: []string{"A", "B", "C"},
	}
	_, err := client.Organizations.CreateOrUpdateSecret(context.Background(), "o", input)
	if err != nil {
		t.Errorf("Organizations.CreateOrUpdateSecret returned error: %v", err)
	}
}

func TestOrgnizationsService_DeleteSecret(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/orgs/o/actions/secrets/NAME", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
	})

	_, err := client.Organizations.DeleteSecret(context.Background(), "o", "NAME")
	if err != nil {
		t.Errorf("Organizations.DeleteSecret returned error: %v", err)
	}
}

func TestOrganizationsService_ListSecretSelectedRepositories(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/orgs/o/actions/secrets/SECRET_NAME/repositories", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{"total_count":2,"repositories":[{"id":1},{"id":2}]}`)
	})

	secretSelectedRepositories, _, err := client.Organizations.ListSecretSelectedRepositories(context.Background(), "o", "SECRET_NAME")
	if err != nil {
		t.Errorf("Organizations.ListSecretSelectedRepositories returned error: %v", err)
	}

	want := &OrganizationSecretSelectedRepositories{
		TotalCount:   Int64(2),
		Repositories: []*Repository{{ID: Int64(1)}, {ID: Int64(2)}},
	}

	if !reflect.DeepEqual(secretSelectedRepositories, want) {
		t.Errorf("Organizations.ListSecretSelectedRepositories returned %+v, want %+v", secretSelectedRepositories, want)
	}
}
