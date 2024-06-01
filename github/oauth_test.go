package github

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"golang.org/x/oauth2"
)

const fakePrivateKey = `-----BEGIN RSA PRIVATE KEY-----
MIICXAIBAAKBgQCBEJcaraGrNm68LstWKMbwalt+n0kaz1pYkBiz+vuFU2K2BkO8
qrk2VXsJQKtLRJYBsEb5CzbgEPRQFoHbqo7lGlJlT8iLVpd0G0g5A9AfABtMnSCP
B5onUg3qAy+G16EtsG+5xOddq3CGs3Au17l9xPCinrpjUts+yYDnUdgKhwIDAQAB
AoGALJT8fRypEak1yw8m8dYYEgfLHwwKhpZFkP4FanYx17YcDOBRGaSnYZtZarLJ
+K/yWRb5DSHQjMmOevOrW6Oow/twOxgt/Qrdmtf7zJCBDuKeIdb5n2Eo1iCeD3jD
LgdM376i0zfgGOeUQFSjGARmAdXPYuITorXPur+fwqQQosECQQD1dzzLtW63GQSa
kO4zm8szLQQKtt22PiTMnqIEizb9GV8Ut96yRaqGXOPPW+WVK402OLyGxg8PxGQp
Pi2+aplxAkEAhpqJMsE/Zedy5mTeyrZ4twDqIlRzkKOOEsPjyvR6zC5kEN3xtE32
oxt2dh9oAINUabQRl4HsnYdvLkQZ/lKndwJBAIYktnXA3hNzuZ9aisJrZn2+IRJV
4w4gIe4s1u/SwKdKgTqKTUpxZgJtqxV77Bg8J7y0/tpMTJjaQ76CO29vRBECQBZR
u5w55PxNw4SmhzbPyPZ3ZVtp63u5Uw6TgssdfNssehU96B1ArDvhiuQNUx56qF2a
qSwZNBXu5iWizziXlgMCQFodovtv4Wk5WI7BYADWXiHeTCwxuZ4Mp1Z3jD9E7Le6
vvShKgb038c8wUyjsYYEmSI5SsdLAKrSqXsFKFMTKWg=
-----END RSA PRIVATE KEY-----`

func Test_applicationTokenSource_Token(t *testing.T) {
	type fields struct {
		id         string
		privateKey []byte
		expiration time.Duration
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "should error when applicationID is empty",
			fields: fields{
				id: "",
			},
			wantErr: true,
		},
		{
			name: "should error when private key is invalid",
			fields: fields{
				id:         "github-app-id",
				privateKey: []byte("invalid"),
			},
			wantErr: true,
		},
		{
			name: "should return a token",
			fields: fields{
				id:         "github-app-id",
				privateKey: []byte(fakePrivateKey),
			},
		},
		{
			name: "should return a token with custom expiration",
			fields: fields{
				id:         "github-app-id",
				privateKey: []byte(fakePrivateKey),
				expiration: 5 * time.Minute,
			},
		},
		{
			name: "should use default expiration when expiration is greater than DefaultApplicationTokenExpiration",
			fields: fields{
				id:         "github-app-id",
				privateKey: []byte(fakePrivateKey),
				expiration: 15 * time.Minute,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var opts []ApplicationTokenOpt
			if tt.fields.expiration != 0 {
				opts = append(opts, WithApplicationTokenExpiration(tt.fields.expiration))
			}

			tr, err := NewApplicationTokenSource(tt.fields.id, tt.fields.privateKey, opts...)

			if (err != nil) != tt.wantErr {
				t.Errorf("NewApplicationTokenSource() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tr == nil {
				return
			}

			got, err := tr.Token()
			if (err != nil) != tt.wantErr {
				t.Errorf("applicationTokenSource.Token() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if got.AccessToken == "" {
				t.Errorf("applicationTokenSource.Token() = %v, want non-empty AccessToken", got)
			}
			if got.TokenType != "Bearer" {
				t.Errorf("applicationTokenSource.Token() = %v, want TokenType to be Bearer", got)
			}
			if got.Expiry.IsZero() {
				t.Errorf("applicationTokenSource.Token() = %v, want non-zero Expiry", got)
			}
		})
	}
}

func Test_installationTokenSource_Token(t *testing.T) {
	tr, err := NewApplicationTokenSource("fake_app_id", []byte(fakePrivateKey))
	if err != nil {
		t.Fatalf("NewApplicationTokenSource() error = %v", err)
	}

	type fields struct {
		id   int64
		src  oauth2.TokenSource
		opts *InstallationTokenOptions
	}
	tests := []struct {
		name               string
		fields             fields
		accessTokenHandler func(http.ResponseWriter, *http.Request)
		want               *oauth2.Token
		wantErr            bool
	}{
		{
			name: "should error when create installation token fails",
			fields: fields{
				id:  2324,
				src: tr,
			},
			accessTokenHandler: func(w http.ResponseWriter, r *http.Request) {
				http.Error(w, "error", http.StatusInternalServerError)
			},
			wantErr: true,
		},
		{
			name: "should return a token",
			fields: fields{
				id:  2324,
				src: tr,
			},
			accessTokenHandler: func(w http.ResponseWriter, r *http.Request) {
				fmt.Fprint(w, `{"token":"t"}`)
			},
			want: &oauth2.Token{
				AccessToken: "t",
				TokenType:   "Bearer",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, mux, _, teardown := setup()
			defer teardown()

			mux.HandleFunc(fmt.Sprintf("/app/installations/%d/access_tokens", tt.fields.id), tt.accessTokenHandler)

			tr := NewInstallationTokenSource(tt.fields.id, tt.fields.src, WithInstallationTokenOptions(tt.fields.opts))

			tr.(*installationTokenSource).apps.client = client

			token, err := tr.Token()
			if (err != nil) != tt.wantErr {
				t.Errorf("installationTokenSource.Token() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !cmp.Equal(token, tt.want, cmpopts.IgnoreUnexported(oauth2.Token{})) {
				t.Errorf("installationTokenSource.Token() returned %+v, want %+v", token, tt.want)
			}
		})
	}
}
