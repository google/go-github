// Copyright 2023 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package github

import (
	"context"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
)

func TestCodespacesService_ListSecrets(t *testing.T) {
	t.Parallel()
	type test struct {
		name       string
		handleFunc func(*http.ServeMux)
		call       func(context.Context, *Client) (*Secrets, *Response, error)
		badCall    func(context.Context, *Client) (*Secrets, *Response, error)
		methodName string
	}
	opts := &ListOptions{Page: 2, PerPage: 2}
	tests := []*test{
		{
			name: "User",
			handleFunc: func(mux *http.ServeMux) {
				mux.HandleFunc("/user/codespaces/secrets", func(w http.ResponseWriter, r *http.Request) {
					testMethod(t, r, "GET")
					testFormValues(t, r, values{"per_page": "2", "page": "2"})
					fmt.Fprint(w, `{"total_count":4,"secrets":[{"name":"A","created_at":"2019-01-02T15:04:05Z","updated_at":"2020-01-02T15:04:05Z"},{"name":"B","created_at":"2019-01-02T15:04:05Z","updated_at":"2020-01-02T15:04:05Z"}]}`)
				})
			},
			call: func(ctx context.Context, client *Client) (*Secrets, *Response, error) {
				return client.Codespaces.ListUserSecrets(ctx, opts)
			},
			methodName: "ListUserSecrets",
		},
		{
			name: "Org",
			handleFunc: func(mux *http.ServeMux) {
				mux.HandleFunc("/orgs/o/codespaces/secrets", func(w http.ResponseWriter, r *http.Request) {
					testMethod(t, r, "GET")
					testFormValues(t, r, values{"per_page": "2", "page": "2"})
					fmt.Fprint(w, `{"total_count":4,"secrets":[{"name":"A","created_at":"2019-01-02T15:04:05Z","updated_at":"2020-01-02T15:04:05Z"},{"name":"B","created_at":"2019-01-02T15:04:05Z","updated_at":"2020-01-02T15:04:05Z"}]}`)
				})
			},
			call: func(ctx context.Context, client *Client) (*Secrets, *Response, error) {
				return client.Codespaces.ListOrgSecrets(ctx, "o", opts)
			},
			badCall: func(ctx context.Context, client *Client) (*Secrets, *Response, error) {
				return client.Codespaces.ListOrgSecrets(ctx, "\n", opts)
			},
			methodName: "ListOrgSecrets",
		},
		{
			name: "Repo",
			handleFunc: func(mux *http.ServeMux) {
				mux.HandleFunc("/repos/o/r/codespaces/secrets", func(w http.ResponseWriter, r *http.Request) {
					testMethod(t, r, "GET")
					testFormValues(t, r, values{"per_page": "2", "page": "2"})
					fmt.Fprint(w, `{"total_count":4,"secrets":[{"name":"A","created_at":"2019-01-02T15:04:05Z","updated_at":"2020-01-02T15:04:05Z"},{"name":"B","created_at":"2019-01-02T15:04:05Z","updated_at":"2020-01-02T15:04:05Z"}]}`)
				})
			},
			call: func(ctx context.Context, client *Client) (*Secrets, *Response, error) {
				return client.Codespaces.ListRepoSecrets(ctx, "o", "r", opts)
			},
			badCall: func(ctx context.Context, client *Client) (*Secrets, *Response, error) {
				return client.Codespaces.ListRepoSecrets(ctx, "\n", "\n", opts)
			},
			methodName: "ListRepoSecrets",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			client, mux, _ := setup(t)

			tt.handleFunc(mux)

			ctx := context.Background()
			secrets, _, err := tt.call(ctx, client)
			if err != nil {
				t.Errorf("Codespaces.%v returned error: %v", tt.methodName, err)
			}

			want := &Secrets{
				TotalCount: 4,
				Secrets: []*Secret{
					{Name: "A", CreatedAt: Timestamp{time.Date(2019, time.January, 02, 15, 04, 05, 0, time.UTC)}, UpdatedAt: Timestamp{time.Date(2020, time.January, 02, 15, 04, 05, 0, time.UTC)}},
					{Name: "B", CreatedAt: Timestamp{time.Date(2019, time.January, 02, 15, 04, 05, 0, time.UTC)}, UpdatedAt: Timestamp{time.Date(2020, time.January, 02, 15, 04, 05, 0, time.UTC)}},
				},
			}
			if !cmp.Equal(secrets, want) {
				t.Errorf("Codespaces.%v returned %+v, want %+v", tt.methodName, secrets, want)
			}

			if tt.badCall != nil {
				testBadOptions(t, tt.methodName, func() (err error) {
					_, _, err = tt.badCall(ctx, client)
					return err
				})
			}

			testNewRequestAndDoFailure(t, tt.methodName, client, func() (*Response, error) {
				got, resp, err := tt.call(ctx, client)
				if got != nil {
					t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", tt.methodName, got)
				}
				return resp, err
			})
		})
	}
}

func TestCodespacesService_GetSecret(t *testing.T) {
	t.Parallel()
	type test struct {
		name       string
		handleFunc func(*http.ServeMux)
		call       func(context.Context, *Client) (*Secret, *Response, error)
		badCall    func(context.Context, *Client) (*Secret, *Response, error)
		methodName string
	}
	tests := []*test{
		{
			name: "User",
			handleFunc: func(mux *http.ServeMux) {
				mux.HandleFunc("/user/codespaces/secrets/NAME", func(w http.ResponseWriter, r *http.Request) {
					testMethod(t, r, "GET")
					fmt.Fprint(w, `{"name":"A","created_at":"2019-01-02T15:04:05Z","updated_at":"2020-01-02T15:04:05Z"}`)
				})
			},
			call: func(ctx context.Context, client *Client) (*Secret, *Response, error) {
				return client.Codespaces.GetUserSecret(ctx, "NAME")
			},
			methodName: "GetUserSecret",
		},
		{
			name: "Org",
			handleFunc: func(mux *http.ServeMux) {
				mux.HandleFunc("/orgs/o/codespaces/secrets/NAME", func(w http.ResponseWriter, r *http.Request) {
					testMethod(t, r, "GET")
					fmt.Fprint(w, `{"name":"A","created_at":"2019-01-02T15:04:05Z","updated_at":"2020-01-02T15:04:05Z"}`)
				})
			},
			call: func(ctx context.Context, client *Client) (*Secret, *Response, error) {
				return client.Codespaces.GetOrgSecret(ctx, "o", "NAME")
			},
			badCall: func(ctx context.Context, client *Client) (*Secret, *Response, error) {
				return client.Codespaces.GetOrgSecret(ctx, "\n", "\n")
			},
			methodName: "GetOrgSecret",
		},
		{
			name: "Repo",
			handleFunc: func(mux *http.ServeMux) {
				mux.HandleFunc("/repos/o/r/codespaces/secrets/NAME", func(w http.ResponseWriter, r *http.Request) {
					testMethod(t, r, "GET")
					fmt.Fprint(w, `{"name":"A","created_at":"2019-01-02T15:04:05Z","updated_at":"2020-01-02T15:04:05Z"}`)
				})
			},
			call: func(ctx context.Context, client *Client) (*Secret, *Response, error) {
				return client.Codespaces.GetRepoSecret(ctx, "o", "r", "NAME")
			},
			badCall: func(ctx context.Context, client *Client) (*Secret, *Response, error) {
				return client.Codespaces.GetRepoSecret(ctx, "\n", "\n", "\n")
			},
			methodName: "GetRepoSecret",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			client, mux, _ := setup(t)

			tt.handleFunc(mux)

			ctx := context.Background()
			secret, _, err := tt.call(ctx, client)
			if err != nil {
				t.Errorf("Codespaces.%v returned error: %v", tt.methodName, err)
			}

			want := &Secret{Name: "A", CreatedAt: Timestamp{time.Date(2019, time.January, 02, 15, 04, 05, 0, time.UTC)}, UpdatedAt: Timestamp{time.Date(2020, time.January, 02, 15, 04, 05, 0, time.UTC)}}
			if !cmp.Equal(secret, want) {
				t.Errorf("Codespaces.%v returned %+v, want %+v", tt.methodName, secret, want)
			}

			if tt.badCall != nil {
				testBadOptions(t, tt.methodName, func() (err error) {
					_, _, err = tt.badCall(ctx, client)
					return err
				})
			}

			testNewRequestAndDoFailure(t, tt.methodName, client, func() (*Response, error) {
				got, resp, err := tt.call(ctx, client)
				if got != nil {
					t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", tt.methodName, got)
				}
				return resp, err
			})
		})
	}
}

func TestCodespacesService_CreateOrUpdateSecret(t *testing.T) {
	t.Parallel()
	type test struct {
		name       string
		handleFunc func(*http.ServeMux)
		call       func(context.Context, *Client, *EncryptedSecret) (*Response, error)
		badCall    func(context.Context, *Client, *EncryptedSecret) (*Response, error)
		methodName string
	}
	tests := []*test{
		{
			name: "User",
			handleFunc: func(mux *http.ServeMux) {
				mux.HandleFunc("/user/codespaces/secrets/NAME", func(w http.ResponseWriter, r *http.Request) {
					testMethod(t, r, "PUT")
					testHeader(t, r, "Content-Type", "application/json")
					testBody(t, r, `{"key_id":"1234","encrypted_value":"QIv="}`+"\n")
					w.WriteHeader(http.StatusCreated)
				})
			},
			call: func(ctx context.Context, client *Client, e *EncryptedSecret) (*Response, error) {
				return client.Codespaces.CreateOrUpdateUserSecret(ctx, e)
			},
			methodName: "CreateOrUpdateUserSecret",
		},
		{
			name: "Org",
			handleFunc: func(mux *http.ServeMux) {
				mux.HandleFunc("/orgs/o/codespaces/secrets/NAME", func(w http.ResponseWriter, r *http.Request) {
					testMethod(t, r, "PUT")
					testHeader(t, r, "Content-Type", "application/json")
					testBody(t, r, `{"key_id":"1234","encrypted_value":"QIv="}`+"\n")
					w.WriteHeader(http.StatusCreated)
				})
			},
			call: func(ctx context.Context, client *Client, e *EncryptedSecret) (*Response, error) {
				return client.Codespaces.CreateOrUpdateOrgSecret(ctx, "o", e)
			},
			badCall: func(ctx context.Context, client *Client, e *EncryptedSecret) (*Response, error) {
				return client.Codespaces.CreateOrUpdateOrgSecret(ctx, "\n", e)
			},
			methodName: "CreateOrUpdateOrgSecret",
		},
		{
			name: "Repo",
			handleFunc: func(mux *http.ServeMux) {
				mux.HandleFunc("/repos/o/r/codespaces/secrets/NAME", func(w http.ResponseWriter, r *http.Request) {
					testMethod(t, r, "PUT")
					testHeader(t, r, "Content-Type", "application/json")
					testBody(t, r, `{"key_id":"1234","encrypted_value":"QIv="}`+"\n")
					w.WriteHeader(http.StatusCreated)
				})
			},
			call: func(ctx context.Context, client *Client, e *EncryptedSecret) (*Response, error) {
				return client.Codespaces.CreateOrUpdateRepoSecret(ctx, "o", "r", e)
			},
			badCall: func(ctx context.Context, client *Client, e *EncryptedSecret) (*Response, error) {
				return client.Codespaces.CreateOrUpdateRepoSecret(ctx, "\n", "\n", e)
			},
			methodName: "CreateOrUpdateRepoSecret",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			client, mux, _ := setup(t)

			tt.handleFunc(mux)

			input := &EncryptedSecret{
				Name:           "NAME",
				EncryptedValue: "QIv=",
				KeyID:          "1234",
			}
			ctx := context.Background()
			_, err := tt.call(ctx, client, input)
			if err != nil {
				t.Errorf("Codespaces.%v returned error: %v", tt.methodName, err)
			}

			if tt.badCall != nil {
				testBadOptions(t, tt.methodName, func() (err error) {
					_, err = tt.badCall(ctx, client, input)
					return err
				})
			}

			testNewRequestAndDoFailure(t, tt.methodName, client, func() (*Response, error) {
				return tt.call(ctx, client, input)
			})
		})
	}
}

func TestCodespacesService_DeleteSecret(t *testing.T) {
	t.Parallel()
	type test struct {
		name       string
		handleFunc func(*http.ServeMux)
		call       func(context.Context, *Client) (*Response, error)
		badCall    func(context.Context, *Client) (*Response, error)
		methodName string
	}
	tests := []*test{
		{
			name: "User",
			handleFunc: func(mux *http.ServeMux) {
				mux.HandleFunc("/user/codespaces/secrets/NAME", func(w http.ResponseWriter, r *http.Request) {
					testMethod(t, r, "DELETE")
				})
			},
			call: func(ctx context.Context, client *Client) (*Response, error) {
				return client.Codespaces.DeleteUserSecret(ctx, "NAME")
			},
			methodName: "DeleteUserSecret",
		},
		{
			name: "Org",
			handleFunc: func(mux *http.ServeMux) {
				mux.HandleFunc("/orgs/o/codespaces/secrets/NAME", func(w http.ResponseWriter, r *http.Request) {
					testMethod(t, r, "DELETE")
				})
			},
			call: func(ctx context.Context, client *Client) (*Response, error) {
				return client.Codespaces.DeleteOrgSecret(ctx, "o", "NAME")
			},
			badCall: func(ctx context.Context, client *Client) (*Response, error) {
				return client.Codespaces.DeleteOrgSecret(ctx, "\n", "\n")
			},
			methodName: "DeleteOrgSecret",
		},
		{
			name: "Repo",
			handleFunc: func(mux *http.ServeMux) {
				mux.HandleFunc("/repos/o/r/codespaces/secrets/NAME", func(w http.ResponseWriter, r *http.Request) {
					testMethod(t, r, "DELETE")
				})
			},
			call: func(ctx context.Context, client *Client) (*Response, error) {
				return client.Codespaces.DeleteRepoSecret(ctx, "o", "r", "NAME")
			},
			badCall: func(ctx context.Context, client *Client) (*Response, error) {
				return client.Codespaces.DeleteRepoSecret(ctx, "\n", "\n", "\n")
			},
			methodName: "DeleteRepoSecret",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			client, mux, _ := setup(t)

			tt.handleFunc(mux)

			ctx := context.Background()
			_, err := tt.call(ctx, client)
			if err != nil {
				t.Errorf("Codespaces.%v returned error: %v", tt.methodName, err)
			}

			if tt.badCall != nil {
				testBadOptions(t, tt.methodName, func() (err error) {
					_, err = tt.badCall(ctx, client)
					return err
				})
			}

			testNewRequestAndDoFailure(t, tt.methodName, client, func() (*Response, error) {
				return tt.call(ctx, client)
			})
		})
	}
}

func TestCodespacesService_GetPublicKey(t *testing.T) {
	t.Parallel()
	type test struct {
		name       string
		handleFunc func(*http.ServeMux)
		call       func(context.Context, *Client) (*PublicKey, *Response, error)
		badCall    func(context.Context, *Client) (*PublicKey, *Response, error)
		methodName string
	}

	tests := []*test{
		{
			name: "User",
			handleFunc: func(mux *http.ServeMux) {
				mux.HandleFunc("/user/codespaces/secrets/public-key", func(w http.ResponseWriter, r *http.Request) {
					testMethod(t, r, "GET")
					fmt.Fprint(w, `{"key_id":"1234","key":"2Sg8iYjAxxmI2LvUXpJjkYrMxURPc8r+dB7TJyvv1234"}`)
				})
			},
			call: func(ctx context.Context, client *Client) (*PublicKey, *Response, error) {
				return client.Codespaces.GetUserPublicKey(ctx)
			},
			methodName: "GetUserPublicKey",
		},
		{
			name: "Org",
			handleFunc: func(mux *http.ServeMux) {
				mux.HandleFunc("/orgs/o/codespaces/secrets/public-key", func(w http.ResponseWriter, r *http.Request) {
					testMethod(t, r, "GET")
					fmt.Fprint(w, `{"key_id":"1234","key":"2Sg8iYjAxxmI2LvUXpJjkYrMxURPc8r+dB7TJyvv1234"}`)
				})
			},
			call: func(ctx context.Context, client *Client) (*PublicKey, *Response, error) {
				return client.Codespaces.GetOrgPublicKey(ctx, "o")
			},
			badCall: func(ctx context.Context, client *Client) (*PublicKey, *Response, error) {
				return client.Codespaces.GetOrgPublicKey(ctx, "\n")
			},
			methodName: "GetOrgPublicKey",
		},
		{
			name: "Repo",
			handleFunc: func(mux *http.ServeMux) {
				mux.HandleFunc("/repos/o/r/codespaces/secrets/public-key", func(w http.ResponseWriter, r *http.Request) {
					testMethod(t, r, "GET")
					fmt.Fprint(w, `{"key_id":"1234","key":"2Sg8iYjAxxmI2LvUXpJjkYrMxURPc8r+dB7TJyvv1234"}`)
				})
			},
			call: func(ctx context.Context, client *Client) (*PublicKey, *Response, error) {
				return client.Codespaces.GetRepoPublicKey(ctx, "o", "r")
			},
			badCall: func(ctx context.Context, client *Client) (*PublicKey, *Response, error) {
				return client.Codespaces.GetRepoPublicKey(ctx, "\n", "\n")
			},
			methodName: "GetRepoPublicKey",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			client, mux, _ := setup(t)

			tt.handleFunc(mux)

			ctx := context.Background()
			key, _, err := tt.call(ctx, client)
			if err != nil {
				t.Errorf("Codespaces.%v returned error: %v", tt.methodName, err)
			}

			want := &PublicKey{KeyID: Ptr("1234"), Key: Ptr("2Sg8iYjAxxmI2LvUXpJjkYrMxURPc8r+dB7TJyvv1234")}
			if !cmp.Equal(key, want) {
				t.Errorf("Codespaces.%v returned %+v, want %+v", tt.methodName, key, want)
			}

			if tt.badCall != nil {
				testBadOptions(t, tt.methodName, func() (err error) {
					_, _, err = tt.badCall(ctx, client)
					return err
				})
			}

			testNewRequestAndDoFailure(t, tt.methodName, client, func() (*Response, error) {
				got, resp, err := tt.call(ctx, client)
				if got != nil {
					t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", tt.methodName, got)
				}
				return resp, err
			})
		})
	}
}

func TestCodespacesService_ListSelectedReposForSecret(t *testing.T) {
	t.Parallel()
	type test struct {
		name       string
		handleFunc func(*http.ServeMux)
		call       func(context.Context, *Client) (*SelectedReposList, *Response, error)
		badCall    func(context.Context, *Client) (*SelectedReposList, *Response, error)
		methodName string
	}
	opts := &ListOptions{Page: 2, PerPage: 2}
	tests := []*test{
		{
			name: "User",
			handleFunc: func(mux *http.ServeMux) {
				mux.HandleFunc("/user/codespaces/secrets/NAME/repositories", func(w http.ResponseWriter, r *http.Request) {
					testMethod(t, r, "GET")
					fmt.Fprintf(w, `{"total_count":1,"repositories":[{"id":1}]}`)
				})
			},
			call: func(ctx context.Context, client *Client) (*SelectedReposList, *Response, error) {
				return client.Codespaces.ListSelectedReposForUserSecret(ctx, "NAME", opts)
			},
			methodName: "ListSelectedReposForUserSecret",
		},
		{
			name: "Org",
			handleFunc: func(mux *http.ServeMux) {
				mux.HandleFunc("/orgs/o/codespaces/secrets/NAME/repositories", func(w http.ResponseWriter, r *http.Request) {
					testMethod(t, r, "GET")
					fmt.Fprintf(w, `{"total_count":1,"repositories":[{"id":1}]}`)
				})
			},
			call: func(ctx context.Context, client *Client) (*SelectedReposList, *Response, error) {
				return client.Codespaces.ListSelectedReposForOrgSecret(ctx, "o", "NAME", opts)
			},
			badCall: func(ctx context.Context, client *Client) (*SelectedReposList, *Response, error) {
				return client.Codespaces.ListSelectedReposForOrgSecret(ctx, "\n", "\n", opts)
			},
			methodName: "ListSelectedReposForOrgSecret",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			client, mux, _ := setup(t)

			tt.handleFunc(mux)

			ctx := context.Background()
			repos, _, err := tt.call(ctx, client)
			if err != nil {
				t.Errorf("Codespaces.%v returned error: %v", tt.methodName, err)
			}

			want := &SelectedReposList{
				TotalCount: Ptr(1),
				Repositories: []*Repository{
					{ID: Ptr(int64(1))},
				},
			}

			if !cmp.Equal(repos, want) {
				t.Errorf("Codespaces.%v returned %+v, want %+v", tt.methodName, repos, want)
			}

			if tt.badCall != nil {
				testBadOptions(t, tt.methodName, func() (err error) {
					_, _, err = tt.badCall(ctx, client)
					return err
				})
			}

			testNewRequestAndDoFailure(t, tt.methodName, client, func() (*Response, error) {
				got, resp, err := tt.call(ctx, client)
				if got != nil {
					t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", tt.methodName, got)
				}
				return resp, err
			})
		})
	}
}

func TestCodespacesService_SetSelectedReposForSecret(t *testing.T) {
	t.Parallel()
	type test struct {
		name       string
		handleFunc func(*http.ServeMux)
		call       func(context.Context, *Client) (*Response, error)
		badCall    func(context.Context, *Client) (*Response, error)
		methodName string
	}
	ids := SelectedRepoIDs{64780797}
	tests := []*test{
		{
			name: "User",
			handleFunc: func(mux *http.ServeMux) {
				mux.HandleFunc("/user/codespaces/secrets/NAME/repositories", func(w http.ResponseWriter, r *http.Request) {
					testMethod(t, r, "PUT")
					testHeader(t, r, "Content-Type", "application/json")
					testBody(t, r, `{"selected_repository_ids":[64780797]}`+"\n")
				})
			},
			call: func(ctx context.Context, client *Client) (*Response, error) {
				return client.Codespaces.SetSelectedReposForUserSecret(ctx, "NAME", ids)
			},
			methodName: "SetSelectedReposForUserSecret",
		},
		{
			name: "Org",
			handleFunc: func(mux *http.ServeMux) {
				mux.HandleFunc("/orgs/o/codespaces/secrets/NAME/repositories", func(w http.ResponseWriter, r *http.Request) {
					testMethod(t, r, "PUT")
					testHeader(t, r, "Content-Type", "application/json")
					testBody(t, r, `{"selected_repository_ids":[64780797]}`+"\n")
				})
			},
			call: func(ctx context.Context, client *Client) (*Response, error) {
				return client.Codespaces.SetSelectedReposForOrgSecret(ctx, "o", "NAME", ids)
			},
			badCall: func(ctx context.Context, client *Client) (*Response, error) {
				return client.Codespaces.SetSelectedReposForOrgSecret(ctx, "\n", "\n", ids)
			},
			methodName: "SetSelectedReposForOrgSecret",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			client, mux, _ := setup(t)

			tt.handleFunc(mux)

			ctx := context.Background()
			_, err := tt.call(ctx, client)
			if err != nil {
				t.Errorf("Codespaces.%v returned error: %v", tt.methodName, err)
			}

			if tt.badCall != nil {
				testBadOptions(t, tt.methodName, func() (err error) {
					_, err = tt.badCall(ctx, client)
					return err
				})
			}

			testNewRequestAndDoFailure(t, tt.methodName, client, func() (*Response, error) {
				return tt.call(ctx, client)
			})
		})
	}
}

func TestCodespacesService_AddSelectedReposForSecret(t *testing.T) {
	t.Parallel()
	type test struct {
		name       string
		handleFunc func(*http.ServeMux)
		call       func(context.Context, *Client) (*Response, error)
		badCall    func(context.Context, *Client) (*Response, error)
		methodName string
	}
	repo := &Repository{ID: Ptr(int64(1234))}
	tests := []*test{
		{
			name: "User",
			handleFunc: func(mux *http.ServeMux) {
				mux.HandleFunc("/user/codespaces/secrets/NAME/repositories/1234", func(w http.ResponseWriter, r *http.Request) {
					testMethod(t, r, "PUT")
				})
			},
			call: func(ctx context.Context, client *Client) (*Response, error) {
				return client.Codespaces.AddSelectedRepoToUserSecret(ctx, "NAME", repo)
			},
			methodName: "AddSelectedRepoToUserSecret",
		},
		{
			name: "Org",
			handleFunc: func(mux *http.ServeMux) {
				mux.HandleFunc("/orgs/o/codespaces/secrets/NAME/repositories/1234", func(w http.ResponseWriter, r *http.Request) {
					testMethod(t, r, "PUT")
				})
			},
			call: func(ctx context.Context, client *Client) (*Response, error) {
				return client.Codespaces.AddSelectedRepoToOrgSecret(ctx, "o", "NAME", repo)
			},
			badCall: func(ctx context.Context, client *Client) (*Response, error) {
				return client.Codespaces.AddSelectedRepoToOrgSecret(ctx, "\n", "\n", repo)
			},
			methodName: "AddSelectedRepoToOrgSecret",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			client, mux, _ := setup(t)

			tt.handleFunc(mux)

			ctx := context.Background()
			_, err := tt.call(ctx, client)
			if err != nil {
				t.Errorf("Codespaces.%v returned error: %v", tt.methodName, err)
			}

			if tt.badCall != nil {
				testBadOptions(t, tt.methodName, func() (err error) {
					_, err = tt.badCall(ctx, client)
					return err
				})
			}

			testNewRequestAndDoFailure(t, tt.methodName, client, func() (*Response, error) {
				return tt.call(ctx, client)
			})
		})
	}
}

func TestCodespacesService_RemoveSelectedReposFromSecret(t *testing.T) {
	t.Parallel()
	type test struct {
		name       string
		handleFunc func(*http.ServeMux)
		call       func(context.Context, *Client) (*Response, error)
		badCall    func(context.Context, *Client) (*Response, error)
		methodName string
	}
	repo := &Repository{ID: Ptr(int64(1234))}
	tests := []*test{
		{
			name: "User",
			handleFunc: func(mux *http.ServeMux) {
				mux.HandleFunc("/user/codespaces/secrets/NAME/repositories/1234", func(w http.ResponseWriter, r *http.Request) {
					testMethod(t, r, "DELETE")
				})
			},
			call: func(ctx context.Context, client *Client) (*Response, error) {
				return client.Codespaces.RemoveSelectedRepoFromUserSecret(ctx, "NAME", repo)
			},
			methodName: "RemoveSelectedRepoFromUserSecret",
		},
		{
			name: "Org",
			handleFunc: func(mux *http.ServeMux) {
				mux.HandleFunc("/orgs/o/codespaces/secrets/NAME/repositories/1234", func(w http.ResponseWriter, r *http.Request) {
					testMethod(t, r, "DELETE")
				})
			},
			call: func(ctx context.Context, client *Client) (*Response, error) {
				return client.Codespaces.RemoveSelectedRepoFromOrgSecret(ctx, "o", "NAME", repo)
			},
			badCall: func(ctx context.Context, client *Client) (*Response, error) {
				return client.Codespaces.RemoveSelectedRepoFromOrgSecret(ctx, "\n", "\n", repo)
			},
			methodName: "RemoveSelectedRepoFromOrgSecret",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			client, mux, _ := setup(t)

			tt.handleFunc(mux)

			ctx := context.Background()
			_, err := tt.call(ctx, client)
			if err != nil {
				t.Errorf("Codespaces.%v returned error: %v", tt.methodName, err)
			}

			if tt.badCall != nil {
				testBadOptions(t, tt.methodName, func() (err error) {
					_, err = tt.badCall(ctx, client)
					return err
				})
			}

			testNewRequestAndDoFailure(t, tt.methodName, client, func() (*Response, error) {
				return tt.call(ctx, client)
			})
		})
	}
}
